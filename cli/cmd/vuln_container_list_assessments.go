package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/lacework/go-sdk/lwtime"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// vulContainerListAssessmentsCmd represents the list-assessments sub-command inside the container
	// vulnerability command
	vulContainerListAssessmentsCmd = &cobra.Command{
		Use:     "list-assessments",
		Aliases: []string{"list", "ls"},
		Short:   "List container vulnerability assessments (default last 24 hours)",
		Long: `List all container vulnerability assessments for the last 24 hours by default.

To customize the time range use use '--start', '--end', or '--range'.

The start and end times can be specified in one of the following formats:

    A. A relative time specifier
    B. RFC3339 date and time
    C. Epoch time in milliseconds

Or use a natural time range like.

    lacework vuln container list --range yesterday

The natural time range of 'yesterday' would represent a relative start time of '-1d@d'
and a relative end time of '@d'.

You can also pass '--fixable' to filter on containers with vulnerabilities that have
fixes available, or '--active' to filter on container images actively running in your
environment.`,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if vulCmdState.Csv {
				cli.EnableCSVOutput()
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			var (
				// the cache key changes depending on some filters that
				// will affect the data returned from our API's
				cacheKey    = generateContainerVulnListCacheKey()
				assessments []vulnerabilityAssessmentSummary
				filter      api.SearchFilter
				start       time.Time
				end         time.Time
				err         error
			)

			expired := cli.ReadCachedAsset(cacheKey, &assessments)
			if expired {
				if vulCmdState.Range != "" {
					cli.Log.Debugw("retrieving natural time range", "range", vulCmdState.Range)
					start, end, err = lwtime.ParseNatural(vulCmdState.Range)
					if err != nil {
						return errors.Wrap(err, "unable to parse natural time range")
					}

				} else {
					cli.Log.Debugw("parsing start time", "start", vulCmdState.Start)
					start, err = parseQueryTime(vulCmdState.Start)
					if err != nil {
						return errors.Wrap(err, "unable to parse start time")
					}

					cli.Log.Debugw("parsing end time", "end", vulCmdState.End)
					end, err = parseQueryTime(vulCmdState.End)
					if err != nil {
						return errors.Wrap(err, "unable to parse end time")
					}
				}

				// search for all active containers
				cli.Log.Infow("using filter with", "start_time", start, "end_time", end)
				filter.TimeFilter = &api.TimeFilter{
					StartTime: &start,
					EndTime:   &end,
				}

				timeRangeMsg := fmt.Sprintf(" in time range (%s to %s)",
					start.Format(time.RFC3339), end.Format(time.RFC3339))

				cli.StartProgress(fmt.Sprintf("Searching for active containers%s...", timeRangeMsg))
				activeContainers, err := cli.LwApi.V2.Entities.ListAllContainersWithFilters(
					api.SearchFilter{
						TimeFilter: filter.TimeFilter,
						Returns:    []string{"mid", "imageId", "startTime"},
					})
				cli.StopProgress()
				if err != nil {
					return errors.Wrap(err, "unable to search for active containers")
				}

				cli.Log.Infow("active containers found",
					"active_count", activeContainers.Total(),
					"entities_count", len(activeContainers.Data),
				)

				// get all container vulnerability assessments
				cli.Log.Infow("requesting list of assessments", "start_time", start, "end_time", end)
				cli.StartProgress(fmt.Sprintf("Fetching assessments%s...", timeRangeMsg))
				assessments, err = listVulnCtrAssessments(activeContainers, &filter)
				cli.StopProgress()
				if err != nil {
					return err
				}

				// write to cache if the request was successful
				cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), assessments)
			} else {
				cli.Log.Infow("assessments loaded from cache", "count", len(assessments))
			}

			// apply vuln ctr list-assessment filters (--active, --registries, --repositories, --fixable)
			if vulnCtrListAssessmentFiltersEnabled() {
				assessments = applyVulnCtrFilters(assessments)
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(assessments)
			}

			if len(assessments) == 0 {
				cli.OutputHuman(buildContainerAssessmentsError())
				return nil
			}

			// Build table output
			assessmentOutput := assessmentSummaryToOutputFormat(assessments)
			rows := vulAssessmentsToTable(assessmentOutput)
			headers := []string{"Registry", "Repository", "Last Scan",
				"Status", "Containers", "Vulnerabilities", "Image Digest"}
			switch {
			case len(rows) == 0:
				cli.OutputHuman(buildContainerAssessmentsError())
			case cli.CSVOutput():
				if err := cli.OutputCSV(headers, rows); err != nil {
					return errors.Wrap(err, "failed to create csv output")
				}
			default:
				cli.OutputHuman(renderSimpleTable(headers, rows))
				if !vulCmdState.Active {
					cli.OutputHuman(
						"\nTry adding '--active' to only show assessments of active containers.\n",
					)
				} else if !vulCmdState.Fixable {
					cli.OutputHuman(
						"\nTry adding '--fixable' to only show assessments with fixable vulnerabilities.\n",
					)
				}
			}

			return nil
		},
	}
)

func vulnCtrListAssessmentFiltersEnabled() bool {
	return len(vulCmdState.Repositories) > 0 ||
		len(vulCmdState.Registries) > 0 ||
		vulCmdState.Fixable ||
		vulCmdState.Active
}

func applyVulnCtrFilters(assessments []vulnerabilityAssessmentSummary) (filtered []vulnerabilityAssessmentSummary) {
	for _, a := range assessments {
		switch {
		case len(vulCmdState.Repositories) > 0:
			if !array.ContainsStr(vulCmdState.Repositories, a.Repository) {
				continue
			}
		case len(vulCmdState.Registries) > 0:
			if !array.ContainsStr(vulCmdState.Registries, a.Registry) {
				continue
			}
		case vulCmdState.Active:
			if a.ActiveContainers == 0 {
				continue
			}
		case vulCmdState.Fixable:
			var vulns []vulnerabilityCtrSummary
			for _, v := range a.Cves {
				if v.Fixable != 0 {
					vulns = append(vulns, v)
				}
				a.Cves = vulns
			}
		}

		filtered = append(filtered, a)
	}
	return
}

// The process to get the list of container assessments is
//
// 1) Check if the user provided a list of registries and repositories,
//    if so, use those filters instead of fetching the entire data from
//    all registries, repositories, local scanners, etc. (This is a memory
//    utilization improvement)
// 2) If no filter by registries and/or repos, then fetch all data from all
//    registries and all local scanners, we purposely split them in two search
//    requests since there could be so much data that we get to the 500,000 rows
//    if data and we could potentially miss some information
// 3) Either 1) or 2) will generate a tree of unique container vulnerability
//    assessments (see the `treeCtrVuln` type), with this tree we will generate
//    one last API request to unique evaluations per image (This is a memory
//    utilization improvement)
// 4) Finally, if we get information from the queried assessments, we build a
//    summary that will ultimately get stored in the cache for subsequent commands
//
func listVulnCtrAssessments(
	activeContainers api.ContainersEntityResponse, filter *api.SearchFilter,
) (assessments []vulnerabilityAssessmentSummary, err error) {

	// Collect only the image ID and the start time to build a tree of
	// images, the time they were evaluated, and the evaluation GUID.
	// This will tell us all images and their latest evaluation
	filter.Returns = []string{"imageId", "startTime", "evalGuid"}
	filter.Filters = []api.Filter{}
	treeOfContainerVuln := treeCtrVuln{}

	// if the user wants to only list assessments from a subset of registries,
	// use that filter instead of fetching data from all registries
	if len(vulCmdState.Registries) != 0 {
		filter.Filters = append(filter.Filters,
			api.Filter{
				Expression: "in",
				Field:      "evalCtx.image_info.registry",
				Values:     vulCmdState.Registries,
			})
	}

	// if the user wants to only list assessments from a subset of repositories,
	// use that filter instead of fetching data from all repositories
	if len(vulCmdState.Repositories) != 0 {
		filter.Filters = append(filter.Filters,
			api.Filter{
				Expression: "in",
				Field:      "evalCtx.image_info.repo",
				Values:     vulCmdState.Repositories,
			})
	}

	if len(filter.Filters) == 0 {
		// if not, then we need to fetch information from 1) all
		// container registries and 2) local scanners in two separate
		// searches since platform scanners might have way too much
		// data which may cause loosing the local scanners data
		//
		// find all container registries
		// cli.StartProgress("Fetching container registries...")
		registries, err := getContainerRegistries()
		// cli.StopProgress()
		if err != nil {
			return nil, err
		}
		cli.Log.Infow("container registries found", "count", len(registries))

		if len(registries) != 0 {
			// 1) search for all assessments from configured container registries
			filter.Filters = []api.Filter{
				{
					Expression: "in",
					Field:      "evalCtx.image_info.registry",
					Values:     registries,
				},
			}
			response, err := cli.LwApi.V2.Vulnerabilities.Containers.SearchAllPages(*filter)
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
			response, err := cli.LwApi.V2.Vulnerabilities.Containers.SearchAllPages(*filter)
			if err != nil {
				return assessments, errors.Wrap(err, "unable to search for container assessments")
			}

			treeOfContainerVuln.ParseData(response.Data)
		}
	} else {
		response, err := cli.LwApi.V2.Vulnerabilities.Containers.SearchAllPages(*filter)
		if err != nil {
			return assessments, errors.Wrap(err, "unable to search for container assessments")
		}

		treeOfContainerVuln.ParseData(response.Data)
	}

	if len(treeOfContainerVuln.ListEvalGuid()) != 0 {
		// Update the filter with the list of evaluation GUIDs and remove the "returns"
		filter.Returns = nil
		filter.Filters = []api.Filter{
			{
				Expression: "in",
				Field:      "evalGuid",
				Values:     treeOfContainerVuln.ListEvalGuid(),
			},
		}

		response, err := cli.LwApi.V2.Vulnerabilities.Containers.SearchAllPages(*filter)
		if err != nil {
			return assessments, errors.Wrap(err, "unable to search for container assessments")
		}

		assessments = buildVulnCtrAssessmentSummary(response.Data, activeContainers)
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
	vulnerabilities  []string
	StatusList       []string `json:"-"`
	fixableCount     int
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
				summary.Cves = append(imageMap[i].Cves, vulnerabilityCtrSummary{a.VulnID, a.FeatureKey.Name, a.FixInfo.FixAvailable, a.Severity})
				if a.FixInfo.FixAvailable != 0 {
					summary.fixableCount++
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
			fixableCount:     fixableCount,
		}
	}

	// Loop over image map and build result
	for _, v := range imageMap {
		uniqueAssessments = append(uniqueAssessments, v)
	}
	return
}

func buildContainerAssessmentsError() string {
	msg := "There are no"
	if vulCmdState.Active {
		msg = fmt.Sprintf("%s active containers", msg)
	} else {
		msg = fmt.Sprintf("%s assessments", msg)
	}

	if len(vulCmdState.Repositories) != 0 {
		msg = fmt.Sprintf("%s for the specified", msg)
		if len(vulCmdState.Repositories) == 1 {
			msg = fmt.Sprintf("%s Repository", msg)
		} else {
			msg = fmt.Sprintf("%s repositories", msg)
		}
	}

	if len(vulCmdState.Registries) != 0 {
		msg = fmt.Sprintf("%s for the specified", msg)
		if len(vulCmdState.Registries) == 1 {
			msg = fmt.Sprintf("%s registry", msg)
		} else {
			msg = fmt.Sprintf("%s registries", msg)
		}
	}

	if vulCmdState.Fixable {
		msg = fmt.Sprintf("%s with fixable vulnerabilities", msg)
	}

	return fmt.Sprintf("%s in your environment.\n", msg)
}

// assessmentSummaryToOutputFormat builds assessmentOutput from the raw response
func assessmentSummaryToOutputFormat(assessments []vulnerabilityAssessmentSummary) []assessmentOutput {
	var out []assessmentOutput

	// sort by active containers
	sort.Slice(assessments, func(i, j int) bool {
		return assessments[i].ActiveContainers > assessments[j].ActiveContainers
	})

	for _, ctr := range assessments {
		severities := []string{}
		fixableCount := 0
		for _, cve := range ctr.Cves {
			severities = append(severities, cve.Severity)
			if cve.Fixable != 0 {
				fixableCount++
			}
		}

		summaryString := severityCtrSummary(severities, ctr.fixableCount)

		out = append(out, assessmentOutput{
			imageRegistry:     ctr.Registry,
			imageRepo:         ctr.Repository,
			startTime:         ctr.ScanTime.UTC().Format(time.RFC3339),
			imageScanStatus:   ctr.ScanStatus,
			ndvContainers:     fmt.Sprintf("%d", ctr.ActiveContainers),
			assessmentSummary: summaryString,
			imageDigest:       ctr.Digest,
		})
	}
	return out
}

func severityCtrSummary(severities []string, fixable int) string {
	summary := &strings.Builder{}
	sevSummaries := make(map[string]int)
	for _, s := range severities {
		switch s {
		case "Critical":
			if v, ok := sevSummaries["Critical"]; ok {
				sevSummaries["Critical"] = v + 1
				continue
			}
			sevSummaries["Critical"] = 1
		case "High":
			if v, ok := sevSummaries["High"]; ok {
				sevSummaries["High"] = v + 1
				continue
			}
			sevSummaries["High"] = 1
		case "Medium":
			if v, ok := sevSummaries["Medium"]; ok {
				sevSummaries["Medium"] = v + 1
				continue
			}
			sevSummaries["Medium"] = 1
		case "Low":
			if v, ok := sevSummaries["Low"]; ok {
				sevSummaries["Low"] = v + 1
				continue
			}
			sevSummaries["Low"] = 1
		case "Info":
			if v, ok := sevSummaries["Info"]; ok {
				sevSummaries["Info"] = v + 1
				continue
			}
			sevSummaries["Info"] = 1
		}
	}

	var keys []string
	for k := range sevSummaries {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return api.SeverityOrder(keys[i]) < api.SeverityOrder(keys[j])
	})

	// only output the highest severity
	if len(keys) != 0 {
		v := sevSummaries[keys[0]]
		summary.WriteString(fmt.Sprintf("%d %s", v, keys[0]))
	}

	if fixable != 0 {
		summary.WriteString(fmt.Sprintf(" %d Fixable", fixable))
	}

	if len(keys) == 0 && fixable == 0 {
		summary.WriteString(fmt.Sprintf("None! Time for %s", randomEmoji()))
	}
	return summary.String()
}

// vulAssessmentsToTable returns assessments in format compatible with table output
func vulAssessmentsToTable(assessments []assessmentOutput) [][]string {
	var out [][]string
	for _, assessment := range assessments {
		out = append(out, []string{
			assessment.imageRegistry,
			assessment.imageRepo,
			assessment.startTime,
			assessment.imageScanStatus,
			assessment.ndvContainers,
			assessment.assessmentSummary,
			assessment.imageDigest,
		})
	}
	return out
}

type assessmentOutput struct {
	imageRegistry     string
	imageRepo         string
	startTime         string
	imageScanStatus   string
	ndvContainers     string
	assessmentSummary string
	imageDigest       string
}

// generateContainerVulnListCacheKey returns the cache key where the CLI will store vulnerability
// assessments for the next few minutes so that consecutive commands run faster and we avoid sending
// duplicate requests to our APIs.
//
// The criteria to generate this cache key is to use the prefix 'vulnerability/container/v2_{HASH}'
// with a hash value at the end of the prefix. The hash is generated based of filters that, if changed
// by the user, will change the data returned from our APIs
func generateContainerVulnListCacheKey() string {
	var cacheFiltersHash = cacheFiltersToBuildVulnContainerHash{
		Start:        vulCmdState.Start,        // mapped to '--start'
		End:          vulCmdState.End,          // mapped to '--end'
		Range:        vulCmdState.Range,        // mapped to '--range'
		Repositories: vulCmdState.Repositories, // mapped to '--repository' (multi-flag)
		Registries:   vulCmdState.Registries,   // mapped to '--registry'   (multi-flag)
	}
	return fmt.Sprintf("vulnerability/container/v2_%d", hash(cacheFiltersHash))
}

// struct that defines the filters we care about, that is, filters that
// when changed, will generate a different hash
type cacheFiltersToBuildVulnContainerHash struct {
	Start        string
	End          string
	Range        string
	Repositories []string
	Registries   []string
}
