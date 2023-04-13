package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiV2(),
	)
	if err != nil {
		log.Fatal(err)
	}

	treeOfContainerVuln := treeCtrVuln{}
	now := time.Now().UTC()
	start := now.AddDate(0, 0, -1)
	filter := api.SearchFilter{
		TimeFilter: &api.TimeFilter{
			StartTime: &start,
			EndTime:   &now,
		},
		Filters: []api.Filter{
			api.Filter{
				Expression: "in",
				Field:      "evalCtx.image_info.repo",
				Values:     []string{"YOUR_REPOS"},
			},
		},
		Returns: []string{"imageId", "startTime", "evalGuid"},
	}

	response, err := lacework.V2.Vulnerabilities.Containers.SearchAllPages(filter)
	if err != nil {
		log.Fatal(err)
	}

	treeOfContainerVuln.ParseData(response.Data)
	if len(treeOfContainerVuln.ListEvalGuid()) == 0 {
		fmt.Println("no evaluations")
		os.Exit(0)
	}

	// Update the filter with the list of evaluation GUIDs and remove the "returns"
	filter.Returns = nil
	filter.Filters = []api.Filter{
		{
			Expression: "in",
			Field:      "evalGuid",
			Values:     treeOfContainerVuln.ListEvalGuid(),
		},
	}

	response, err = lacework.V2.Vulnerabilities.Containers.SearchAllPages(filter)
	if err != nil {
		log.Fatal(err)
	}

	activeContainers, err := lacework.V2.Entities.ListAllContainersWithFilters(
		api.SearchFilter{
			TimeFilter: &api.TimeFilter{
				StartTime: &start,
				EndTime:   &now,
			},
			Returns: []string{"mid", "imageId", "startTime"},
		})
	if err != nil {
		log.Fatal(err)
	}

	assessments := buildVulnCtrAssessmentSummary(response.Data, activeContainers)

	fmt.Printf("Total number of assessments: %s", len(assessments))
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
				latestContainer.StartTime = ctr.StartTime
				latestContainer.EvalGUID = ctr.EvalGUID
			}
		} else {
			// @afiune this is NOT thread safe!! But it is also not used in parallel executions
			*v = append(*v, ctrVuln{ctr.EvalGUID, ctr.ImageID, ctr.StartTime})
		}
	}
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

func buildVulnCtrAssessmentSummary(
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
				summary.Cves = append(imageMap[i].Cves, vulnerabilityCtrSummary{
					a.VulnID, a.FeatureKey.Name, a.FixInfo.FixAvailable, a.Severity})
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
