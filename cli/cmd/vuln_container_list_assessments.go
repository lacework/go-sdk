package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/lacework/go-sdk/v2/api"
	"github.com/lacework/go-sdk/v2/lwtime"
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

				assessments, err = listVulnCtrAssessments(&api.SearchFilter{
					Filters: []api.Filter{
						{
							Expression: "gt",
							Field:      "lastScanTime",
							Value:      start.Format(time.RFC3339),
						},
						{
							Expression: "lt",
							Field:      "lastScanTime",
							Value:      end.Format(time.RFC3339),
						},
					},
				})
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
	return vulCmdState.Fixable ||
		vulCmdState.Active
}

func applyVulnCtrFilters(assessments []vulnerabilityAssessmentSummary) (filtered []vulnerabilityAssessmentSummary) {
	for _, a := range assessments {
		if vulCmdState.Active {
			if a.ActiveContainers == 0 {
				continue
			}
		}
		if vulCmdState.Fixable {
			if !a.HasFixableVulns() {
				continue
			}
		}

		filtered = append(filtered, a)
	}
	return
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
func listVulnCtrAssessments(
	filter *api.SearchFilter,
) (assessments []vulnerabilityAssessmentSummary, err error) {

	// if the user wants to only list assessments from a subset of registries,
	// use that filter instead of fetching data from all registries
	if len(vulCmdState.Registries) != 0 {
		filter.Filters = append(filter.Filters,
			api.Filter{
				Expression: "in",
				Field:      "registry",
				Values:     vulCmdState.Registries,
			})
	}
	// if the user wants to only list assessments from a subset of repositories,
	// use that filter instead of fetching data from all repositories
	if len(vulCmdState.Repositories) != 0 {
		filter.Filters = append(filter.Filters,
			api.Filter{
				Expression: "in",
				Field:      "repository",
				Values:     vulCmdState.Repositories,
			})
	}
	response, err := cli.LwApi.V2.VulnerabilityObservations.ImageSummary.SearchAllPages(*filter)
	if err != nil {
		return assessments, errors.Wrap(err, "unable to search for container assessments")
	}

	assessments = buildVulnCtrAssessmentSummary(response.Data)

	return
}

type vulnerabilityAssessmentSummary struct {
	ImageID          string         `json:"image_id"`
	Repository       string         `json:"repository"`
	Registry         string         `json:"registry"`
	Digest           string         `json:"digest"`
	ScanTime         time.Time      `json:"scan_time"`
	ScanStatus       string         `json:"scan_status"`
	ActiveContainers int            `json:"active_containers"`
	FixableCount     int            `json:"fixable_count"`
	VulnCount        map[string]int `json:"vuln_count"`
}

func (v vulnerabilityAssessmentSummary) HasFixableVulns() bool {
	return v.FixableCount != 0
}

func buildVulnCtrAssessmentSummary(
	assessments []api.VulnerabilityObservationsImageSummary,
) (uniqueAssessments []vulnerabilityAssessmentSummary) {

	imageMap := map[string]vulnerabilityAssessmentSummary{}

	for _, a := range assessments {
		i := fmt.Sprintf("%s-%s-%s", a.Registry, a.Repository, a.ImageId)
		scanTime, err := time.Parse(time.RFC3339, a.LastScanTime)
		if err != nil {
			fmt.Println("Error parsing last scan time: ", err)
		}

		vulnCount := make(map[string]int)
		vulnCount["Critical"] = a.VulnCountCritical
		vulnCount["High"] = a.VulnCountHigh
		vulnCount["Medium"] = a.VulnCountMedium
		vulnCount["Low"] = a.VulnCountLow
		vulnCount["Info"] = a.VulnCountInfo

		imageMap[i] = vulnerabilityAssessmentSummary{
			ImageID:          a.ImageId,
			Repository:       a.Repository,
			Registry:         a.Registry,
			Digest:           a.Digest,
			ScanTime:         scanTime,
			ScanStatus:       a.ScanStatus,
			ActiveContainers: a.ContainerCount,
			FixableCount: a.VulnCountCriticalFixable + a.VulnCountHighFixable + a.VulnCountMediumFixable +
				a.VulnCountLowFixable + a.VulnCountInfoFixable,
			VulnCount: vulnCount,
		}
	}
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
			msg = fmt.Sprintf("%s repository", msg)
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

		summaryString := severityCtrSummary(ctr.VulnCount, ctr.FixableCount)

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

func severityCtrSummary(vulnCount map[string]int, fixable int) string {
	summary := &strings.Builder{}

	severityOrder := []string{"Critical", "High", "Medium", "Low", "Info"}

	highestSeverity := ""
	highestSeverityVulnCount := 0

	//Find the highest-severity non-zero vuln count
	for _, severity := range severityOrder {
		if count, exists := vulnCount[severity]; exists && count > 0 {
			highestSeverity = severity
			highestSeverityVulnCount = count
			break
		}
	}

	// only output the highest severity
	if highestSeverityVulnCount > 0 {
		summary.WriteString(fmt.Sprintf("%d %s", highestSeverityVulnCount, highestSeverity))
	}

	if fixable != 0 {
		summary.WriteString(fmt.Sprintf(" %d Fixable", fixable))
	}

	if highestSeverityVulnCount == 0 && fixable == 0 {
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
