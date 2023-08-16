package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/pkg/errors"
)

type state struct {
	client *api.Client
}

// The process to get the list of container assessments is
//
//  1. Check if the user provided a list of registries and repositories,
//     if so, use those filters instead of fetching the entire data from
//     all registries, repositories, local scanners, etc. (This is a memory
//     utilization improvement)
//  2. If no filter by registries and/or repos, then fetch all data from all
//     registries and all local scanners, we purposely split them in two search
//     requests since there could be so much data that we get to the 500,000 rows
//     if data and we could potentially miss some information
//  3. Either 1) or 2) will generate a tree of unique container vulnerability
//     assessments (see the `treeCtrVuln` type), with this tree we will generate
//     one last API request to unique evaluations per image (This is a memory
//     utilization improvement)
//  4. Finally, if we get information from the queried assessments, we build a
//     summary that will ultimately get stored in the cache for subsequent commands
func (s state) listVulnCtrAssessments(
	activeContainers api.ContainersEntityResponse, filter *api.SearchFilter,
) (assessments []vulnerabilityAssessmentSummary, err error) {

	// Collect only the image ID and the start time to build a tree of
	// images, the time they were evaluated, and the evaluation GUID.
	// This will tell us all images and their latest evaluation
	filter.Returns = []string{"imageId", "startTime", "evalGuid"}
	filter.Filters = []api.Filter{}
	treeOfContainerVuln := treeCtrVuln{}

	// We need to fetch information from 1) all
	// container registries and 2) local scanners in two separate
	// searches since platform scanners might have way too much
	// data which may cause losing the local scanners data
	//
	// find all container registries
	registries, err := s.getContainerRegistries()
	if err != nil {
		return nil, err
	}

	if len(registries) != 0 {
		// 1) search for all assessments from configured container registries
		filter.Filters = []api.Filter{
			{
				Expression: "in",
				Field:      "evalCtx.image_info.registry",
				Values:     registries,
			},
		}
		response, err := s.client.V2.Vulnerabilities.Containers.SearchAllPages(*filter)
		if err != nil {
			return assessments, errors.Wrap(err, "unable to search for container assessments")
		}

		treeOfContainerVuln.ParseData(response.Data)

		// 2) search for assessments from local scanners, that is, non container registries
		filter.Filters = []api.Filter{
			{
				Expression: "not_in",
				Field:      "evalCtx.image_info.registry",
				Values:     registries,
			},
		}
	} else {
		response, err := s.client.V2.Vulnerabilities.Containers.SearchAllPages(*filter)
		if err != nil {
			return assessments, errors.Wrap(err, "unable to search for container assessments")
		}

		treeOfContainerVuln.ParseData(response.Data)
	}

	evalGuids := treeOfContainerVuln.ListEvalGuid()

	if len(evalGuids) != 0 {
		// Update the filter with the list of evaluation GUIDs and remove the "returns"
		filter.Returns = nil
		filter.Filters = []api.Filter{
			{
				Expression: "in",
				Field:      "evalGuid",
				Values:     evalGuids,
			},
		}

		response, err := s.client.V2.Vulnerabilities.Containers.SearchAllPages(*filter)
		if err != nil {
			return assessments, errors.Wrap(err, "unable to search for container assessments")
		}

		assessments = s.buildVulnCtrAssessmentSummary(response.Data, activeContainers)
	}
	return
}

// treeCtrVuln and ctrVuln are types that help us generate an tree of container
// vulnerability assessments that are unique per image ID, that is, there will
// never be duplicates of the same image with different evaluation guids (evalGuid)
type treeCtrVuln []ctrVuln
type ctrVuln struct {
	EvalGUID  string
	ImageID   string
	StartTime time.Time
}

func (v treeCtrVuln) Len() int {
	return len(v)
}
func (v treeCtrVuln) Get(imageID string) (*ctrVuln, bool) {
	for _, ctr := range v {
		if ctr.ImageID == imageID {
			return &ctr, true
		}
	}
	return nil, false
}

func (v treeCtrVuln) ListEvalGuid() (guids []string) {
	for _, ctr := range v {
		guids = append(guids, ctr.EvalGUID)
	}
	return
}
func (v treeCtrVuln) ListImageIDs() (ids []string) {
	for _, ctr := range v {
		ids = append(ids, ctr.ImageID)
	}
	return
}

func (v *treeCtrVuln) ParseData(data []api.VulnerabilityContainer) {
	for _, ctr := range data {
		latestContainer, exist := v.Get(ctr.ImageID)
		if exist {
			if latestContainer.EvalGUID == ctr.EvalGUID {
				continue
			}
			if ctr.StartTime.After(latestContainer.StartTime) {
				v.Replace(*latestContainer, ctrVuln{
					EvalGUID:  ctr.EvalGUID,
					ImageID:   ctr.ImageID,
					StartTime: ctr.StartTime,
				})
			}
		} else {
			// @afiune this is NOT thread safe!! But it is also not used in parallel executions
			*v = append(*v, ctrVuln{ctr.EvalGUID, ctr.ImageID, ctr.StartTime})
		}
	}
}

// Replace updates an existing ctrVuln in the treeCtrVuln slice with a new ctrVuln
func (v treeCtrVuln) Replace(old ctrVuln, new ctrVuln) {
	for i, ctrVuln := range v {
		if ctrVuln.EvalGUID == old.EvalGUID {
			v[i] = new
			break
		}
	}
}

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
	)
	if err != nil {
		log.Fatal(err)
	}

	app := state{lacework}

	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -7) // 7 days from ago
		filter = api.SearchFilter{
			TimeFilter: &api.TimeFilter{
				StartTime: &before,
				EndTime:   &now,
			},
			Returns: []string{"mid", "imageId", "startTime"},
		}
	)

	// search for all active containers
	activeContainers, err := lacework.V2.Entities.ListAllContainersWithFilters(filter)
	if err != nil {
		log.Fatal(err)
	}

	// get all container vulnerability assessments
	assessments, err := app.listVulnCtrAssessments(activeContainers, &filter)
	if err != nil {
		log.Fatal(err)
	}

	if len(assessments) == 0 {
		fmt.Println("There are no active containers in your environment")
	} else {
		for _, assessment := range assessments {
			fmt.Printf("Registry: %s\nRepo: %s\nActive Containers: %d\n",
				assessment.Registry, assessment.Repository, assessment.ActiveContainers,
			)
		}
		// Output:
		//
		// Registry: my-registry
		// Repo: foo
		// Active Containers: 5
		// Registry: my-registry
		// Repo: bar
		// Active Containers: 2
	}
}

func (s state) getContainerRegistries() ([]string, error) {
	var (
		registries            = make([]string, 0)
		regsIntegrations, err = s.client.V2.ContainerRegistries.List()
	)
	if err != nil {
		return registries, errors.Wrap(err, "unable to get container registry integrations")
	}

	for _, i := range regsIntegrations.Data {
		// avoid adding empty registries coming from the new local_scanner and avoid adding duplicate registries
		if i.ContainerRegistryDomain() == "" || array.ContainsStr(registries, i.ContainerRegistryDomain()) {
			continue
		}

		registries = append(registries, i.ContainerRegistryDomain())
	}

	return registries, nil
}

func (s state) buildVulnCtrAssessmentSummary(
	assessments []api.VulnerabilityContainer, activeContainers api.ContainersEntityResponse,
) (uniqueAssessments []vulnerabilityAssessmentSummary) {

	imageMap := map[string]vulnerabilityAssessmentSummary{}

	// build a map for our assessments per image
	for _, a := range assessments {
		i := fmt.Sprintf("%s-%s-%s", a.EvalCtx.ImageInfo.Registry, a.EvalCtx.ImageInfo.Repo, a.EvalCtx.ImageInfo.ID)
		if _, ok := imageMap[i]; ok {
			// if the image id assessment has already been added, then append the vulnerabilities
			summary := imageMap[i]
			summary.StatusList = append(imageMap[i].StatusList, a.Status)

			// check duplicate cves
			vulnKey := fmt.Sprintf("%s-%s", a.VulnID, a.FeatureKey.Name)
			if !array.ContainsStr(imageMap[i].vulnerabilities, vulnKey) && a.VulnID != "" {
				summary.vulnerabilities = append(imageMap[i].vulnerabilities, vulnKey)
				summary.Cves = append(imageMap[i].Cves,
					vulnerabilityCtrSummary{a.VulnID, a.FeatureKey.Name, a.FixInfo.FixAvailable, a.Severity},
				)
				if a.FixInfo.FixAvailable != 0 {
					summary.FixableCount++
				}
			}

			imageMap[i] = summary
			continue
		}

		fixableCount := 0
		if a.FixInfo.FixAvailable != 0 {
			fixableCount = 1
		}

		// search for active containers
		imageMap[i] = vulnerabilityAssessmentSummary{
			ImageID:          a.ImageID,
			Repository:       a.EvalCtx.ImageInfo.Repo,
			Registry:         a.EvalCtx.ImageInfo.Registry,
			Digest:           a.EvalCtx.ImageInfo.Digest,
			ScanTime:         a.StartTime,
			Cves:             []vulnerabilityCtrSummary{{a.VulnID, a.FeatureKey.Name, a.FixInfo.FixAvailable, a.Severity}},
			ScanStatus:       a.EvalCtx.ImageInfo.Status,
			ActiveContainers: activeContainers.Count(a.ImageID),
			vulnerabilities:  []string{fmt.Sprintf("%s-%s", a.VulnID, a.FeatureKey.Name)},
			StatusList:       []string{a.Status},
			FixableCount:     fixableCount,
		}
	}

	// Loop over image map and build result
	for _, v := range imageMap {
		uniqueAssessments = append(uniqueAssessments, v)
	}
	return
}

type vulnerabilityAssessmentSummary struct {
	ImageID          string                    `json:"image_id"`
	Repository       string                    `json:"repository"`
	Registry         string                    `json:"registry"`
	Digest           string                    `json:"digest"`
	ScanTime         time.Time                 `json:"scan_time"`
	Cves             []vulnerabilityCtrSummary `json:"cves"`
	ScanStatus       string                    `json:"scan_status"`
	ActiveContainers int                       `json:"active_containers"`
	FixableCount     int                       `json:"fixable_count"`
	vulnerabilities  []string
	StatusList       []string `json:"-"`
}

type vulnerabilityCtrSummary struct {
	Id       string `json:"vuln_id"`
	Pkg      string `json:"package_name"`
	Fixable  int    `json:"fixable"`
	Severity string `json:"severity"`
}

func (v vulnerabilityAssessmentSummary) Status() string {
	if array.ContainsStr(v.StatusList, "VULNERABLE") {
		return "VULNERABLE"
	}
	return "GOOD"
}

func (v vulnerabilityAssessmentSummary) HasFixableVulns() bool {
	for _, c := range v.Cves {
		if c.Fixable != 0 {
			return true
		}
	}
	return false
}
