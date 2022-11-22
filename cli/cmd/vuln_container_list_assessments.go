package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// vulContainerListAssessmentsCmd represents the list-assessments sub-command inside the container
	// vulnerability command
	vulContainerListAssessmentsCmd = &cobra.Command{
		Use:     "list-assessments",
		Aliases: []string{"list", "ls"},
		Short:   "List container vulnerability assessments (default last 7 days)",
		Long: `List all container vulnerability assessments for the last 7 days by default, or
pass --start and --end to specify a custom time range. You can pass --fixable to
filter on containers with vulnerabilities that have fixes available.`,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if vulCmdState.Csv {
				cli.EnableCSVOutput()
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			var (
				assessments      []vulnerabilityAssessmentSummary
				cacheKey         = "vulnerability/container/v2"
				filter           api.SearchFilter
				now              = time.Now().UTC()
				before           = now.AddDate(0, 0, -7) // 7 days from ago
				msg              = "Fetching assessments from the last 7 days"
				partialResultMsg string
			)

			expired := cli.ReadCachedAsset(cacheKey, &assessments)
			if expired {
				// before starting the search find all ctr reg
				cli.StartProgress("Fetching container registries...")
				registries, err := getContainerRegistries()
				cli.StopProgress()

				if err != nil {
					return err
				}

				if vulCmdState.Start != "" || vulCmdState.End != "" {
					start, end, errT := parseStartAndEndTime(vulCmdState.Start, vulCmdState.End)
					if errT != nil {
						return errors.Wrap(errT, "unable to parse time range")
					}

					cli.Log.Infow("requesting list of assessments from custom time range",
						"start_time", start, "end_time", end,
					)
					filter.TimeFilter = &api.TimeFilter{
						StartTime: &before,
						EndTime:   &now,
					}
					msg = fmt.Sprintf("Fetching assessments in date range %s - %s", start, end)
				}
				cli.StartProgress(msg)
				assessments, err = listVulnCtrAssessments(registries, filter)
				cli.StopProgress()
				cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), assessments)

				if err != nil {
					return err
				}
			}

			if len(assessments) == 0 {
				cli.OutputHuman("There are no container assessments for this environment.\n")
				return nil
			}

			// apply vuln ctr list-assessment filters (--registries, --repositories, --fixable)
			if vulnCtrListAssessmentFiltersEnabled() {
				assessments = applyVulnCtrFilters(assessments)
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(assessments)
			}

			// Build table output
			assessmentOutput := assessmentSummaryToOutputFormat(assessments)
			rows := vulAssessmentsToTable(assessmentOutput)
			headers := []string{"Registry", "Repository", "Last Scan", "Status", "Vulnerabilities", "Image Digest"}
			switch {
			// if the user wants to show only assessments of containers running
			// and we don't have any, show a friendly message
			case len(rows) == 0:
				cli.OutputHuman(buildContainerAssessmentsError())
			case cli.CSVOutput():
				if err := cli.OutputCSV(headers, rows); err != nil {
					return errors.Wrap(err, "failed to create csv output")
				}
			case partialResultMsg != "":
				cli.OutputHuman(partialResultMsg)
			default:
				cli.OutputHuman(renderSimpleTable(headers, rows))
				if !vulCmdState.Fixable {
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
	return len(vulCmdState.Repositories) > 0 || len(vulCmdState.Registries) > 0 || vulCmdState.Fixable
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

func listVulnCtrAssessments(registries []string, filter api.SearchFilter) (assessments []vulnerabilityAssessmentSummary, err error) {
	var ctrMap = map[string][]api.VulnerabilityContainer{}
	// for each ctr registry perform a search
	for _, registry := range registries {
		filter.Filters = []api.Filter{{
			Expression: "eq",
			Field:      "evalCtx.image_info.registry",
			Value:      registry,
		}}

		response, err := cli.LwApi.V2.Vulnerabilities.Containers.SearchAllPages(filter)
		if err != nil {
			return assessments, errors.Wrap(err, "unable to get assessments")
		}
		if len(response.Data) != 0 {
			ctrMap[registry] = response.Data
		}
	}

	for _, v := range ctrMap {
		assessments = append(assessments, buildVulnCtrAssessmentSummary(v)...)
	}

	return
}

type vulnerabilityAssessmentSummary struct {
	ImageID         string                    `json:"image_id"`
	Repository      string                    `json:"repository"`
	Registry        string                    `json:"registry"`
	Digest          string                    `json:"digest"`
	ScanTime        time.Time                 `json:"scan_time"`
	Cves            []vulnerabilityCtrSummary `json:"cves"`
	vulnerabilities []string
	StatusList      []string `json:"status_list"`
	fixableCount    int
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

func buildVulnCtrAssessmentSummary(assessments []api.VulnerabilityContainer) (uniqueAssessments []vulnerabilityAssessmentSummary) {
	var (
		imageMap = map[string]vulnerabilityAssessmentSummary{}
	)

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
		imageMap[i] = vulnerabilityAssessmentSummary{
			a.ImageID,
			a.EvalCtx.ImageInfo.Repo,
			a.EvalCtx.ImageInfo.Registry,
			a.EvalCtx.ImageInfo.Digest,
			a.StartTime,
			[]vulnerabilityCtrSummary{{a.VulnID, a.FeatureKey.Name, a.FixInfo.FixAvailable, a.Severity}},
			[]string{fmt.Sprintf("%s-%s", a.VulnID, a.FeatureKey.Name)},
			[]string{a.Status},
			fixableCount,
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

	sort.Slice(assessments, func(i, j int) bool {
		return assessments[i].Repository > assessments[j].Repository
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
			imageRegistry:   ctr.Registry,
			imageRepo:       ctr.Repository,
			startTime:       ctr.ScanTime.UTC().Format(time.RFC3339),
			imageScanStatus: ctr.Status(),
			//todo(v2): adding active containers blocked by RAIN-43538
			ndvContainers:     "1",
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
			//todo(v2): active containers blocked by RAIN-43538
			//assessment.ndvContainers,
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
