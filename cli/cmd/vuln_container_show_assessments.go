package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const regexAllTabs = "(\\t){1,}"

var (
	// vulContainerShowAssessmentCmd represents the show-assessment sub-command inside the container
	// vulnerability command
	vulContainerShowAssessmentCmd = &cobra.Command{
		Use:     "show-assessment <sha256:hash>",
		Aliases: []string{"show"},
		Short:   "Show results of a container vulnerability assessment",
		Long: `Show the results from a vulnerability assessment of a specified container.

Arguments:
    <sha256:hash> a sha256 hash of a container image (format: sha256:1ee...1d3b)

By default, this command expects a sha256 image digest or tag. To lookup an
assessment by its image id, use the flag '--image_id' followed by the sha256
image id.

To request an on-demand vulnerability scan:

    lacework vulnerability container scan <registry> <repository> <tag|digest>`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if vulCmdState.Csv {
				cli.EnableCSVOutput()

				// If --details or --packages is not passed, csv outputs nothing; defaulting to --details
				if !vulCmdState.Details && !vulCmdState.Packages {
					vulCmdState.Details = true
				}
			}

			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}
			err := showContainerAssessmentsWithSha256(args[0], api.SearchFilter{})
			var e *vulnerabilityPolicyError
			if errors.As(err, &e) {
				c.SilenceUsage = true
			}

			return err
		},
	}
)

func showContainerAssessmentsWithSha256(sha string, filter api.SearchFilter) error {
	var (
		assessment     api.VulnerabilitiesContainersResponse
		vulnerabilites []api.VulnerabilityContainer
		now            = time.Now().UTC()
		before         = now.AddDate(0, 0, -7) // 7 days from ago
		searchField    string
		err            error
	)

	if len(filter.Filters) > 0 {
		cli.Log.Debugw("retrieve assessment with filters ", filter.Filters)

		cli.StartProgress("Fetching assessment...")
		assessment, err = cli.LwApi.V2.Vulnerabilities.Containers.Search(filter)
		cli.StopProgress()

		if len(assessment.Data) == 0 {
			cli.OutputHuman("No active containers found.\n")
			return nil
		}
	}

	if vulCmdState.ImageID {
		searchField = "evalCtx.image_info.id"
		cli.Log.Debugw("retrieve image assessment", searchField, sha)

		cli.StartProgress("Fetching assessment...")
		assessment, err = cli.LwApi.V2.Vulnerabilities.Containers.Search(api.SearchFilter{
			TimeFilter: &api.TimeFilter{
				StartTime: &before,
				EndTime:   &now,
			},
			Filters: []api.Filter{{
				Expression: "eq",
				Field:      searchField,
				Value:      sha,
			}},
		})
		cli.StopProgress()

		if len(assessment.Data) == 0 {
			cli.OutputHuman("No active containers found.\n")
			return nil
		}

	}

	if sha != "" && !vulCmdState.ImageID {
		searchField = "evalCtx.image_info.digest"
		cli.StartProgress("Fetching assessment...")
		assessment, err = cli.LwApi.V2.Vulnerabilities.Containers.Search(api.SearchFilter{
			TimeFilter: &api.TimeFilter{
				StartTime: &before,
				EndTime:   &now,
			},
			Filters: []api.Filter{{
				Expression: "eq",
				Field:      searchField,
				Value:      sha,
			}},
		})
		cli.StopProgress()

		if len(assessment.Data) == 0 {
			cli.OutputHuman("No active containers found.\n")
			return nil
		}
	}

	for _, a := range assessment.Data {
		if a.Status == "VULNERABLE" {
			vulnerabilites = append(vulnerabilites, a)
		}
	}

	assessment.Data = vulnerabilites

	if err != nil {
		return errors.Wrap(err, "unable to show vulnerability assessment")
	}

	cli.Log.Debugw("image assessment", "details", assessment)

	if err := buildVulnContainerAssessmentReports(assessment); err != nil {
		return err
	}

	if vulFailureFlagsEnabled() {
		cli.Log.Infow("failure flags enabled",
			"fail_on_severity", vulCmdState.FailOnSeverity,
			"fail_on_fixable", vulCmdState.FailOnFixable,
		)
		vulnPolicy := NewVulnerabilityPolicyErrorV2(
			assessment,
			vulCmdState.FailOnSeverity,
			vulCmdState.FailOnFixable,
		)
		if vulnPolicy.NonCompliant() {
			return vulnPolicy
		}
	}

	return nil
}

// Build the cli output for vuln ctr 'show-assessment' command
func buildVulnContainerAssessmentReports(response api.VulnerabilitiesContainersResponse) error {
	assessment := response.Data
	if len(assessment) == 0 {
		if cli.JSONOutput() {
			// if no assessments are found return empty array
			return cli.OutputJSON([]any{})
		}
		cli.OutputHuman("Great news! This container image has no vulnerabilities... (time for %s)\n", randomEmoji())
		return nil
	}

	var details vulnerabilityDetailsReport
	details.VulnerabilityDetails = filterVulnerabilityContainer(assessment)
	response.Data = details.VulnerabilityDetails.Filtered
	details.Packages = filterVulnContainerImagePackages(details.VulnerabilityDetails.Filtered)

	switch {
	case cli.JSONOutput():
		filteredAssessment := assessment
		if err := cli.OutputJSON(filteredAssessment); err != nil {
			return err
		}
	case cli.CSVOutput():
		if err := cli.OutputCSV(buildVulnerabilityDetailsReportCSV(details)); err != nil {
			return err
		}
	default:
		if len(response.Data) == 0 {
			if vulCmdState.Severity != "" {
				cli.OutputHuman("There ano vulnerabilties found for this severity")
			}

			cli.OutputHuman("Great news! This container image has no vulnerabilities... (time for %s)\n", randomEmoji())
			return nil
		}
		summaryReport := buildVulnerabilitySummaryReportTable(response)
		detailsReport := buildVulnerabilityDetailsReportTable(details)
		cli.OutputHuman(buildVulnContainerAssessmentReportTable(summaryReport, detailsReport))
		if vulCmdState.Html {
			if err := generateVulnAssessmentHTML(response); err != nil {
				return err
			}
		}
	}

	return nil
}

func buildVulnerabilityDetailsReportCSV(details vulnerabilityDetailsReport) ([]string, [][]string) {
	if !(vulCmdState.Details || vulCmdState.Packages || vulFiltersEnabled()) {
		return nil, nil
	}

	if vulCmdState.Packages {
		return []string{"CVE Count", "Severity", "Package", "Current Version", "Fix Version"},
			vulContainerImagePackagesToTable(details.Packages)
	}

	return []string{"CVE ID", "Severity", "CVSSv2", "CVSSv3", "Package", "Current Version",
		"Fix Version", "Introduced in Layer"}, vulContainerImageLayersToCSV(details.VulnerabilityDetails)
}

func buildVulnerabilityDetailsReportTable(details vulnerabilityDetailsReport) string {
	report := &strings.Builder{}

	if vulCmdState.Details || vulCmdState.Packages || vulFiltersEnabled() {
		if vulCmdState.Packages {
			vulnPackagesTable := vulContainerImagePackagesToTable(details.Packages)

			report.WriteString(
				renderSimpleTable(
					[]string{"CVE Count", "Severity", "Package", "Current Version", "Fix Version"},
					vulnPackagesTable,
				),
			)

			if vulFiltersEnabled() {
				filteredOutput := fmt.Sprintf("%v of %v packages showing \n", details.Packages.totalPackages, details.Packages.totalUnfiltered)
				report.WriteString(filteredOutput)
			}
		} else {
			vulnImageTable := vulContainerImageLayersToTable(details.VulnerabilityDetails)

			report.WriteString(
				renderCustomTable(
					[]string{"CVE ID", "Severity", "Package", "Current Version",
						"Fix Version", "Introduced in Layer", "Status"},
					vulnImageTable,
					tableFunc(func(t *tablewriter.Table) {
						t.SetBorder(false)
						t.SetRowLine(true)
						t.SetColumnSeparator(" ")
						t.SetAlignment(tablewriter.ALIGN_LEFT)
					}),
				),
			)

			if vulFiltersEnabled() {
				filteredOutput := fmt.Sprintf("%v of %v vulnerabilities showing \n",
					details.VulnerabilityDetails.TotalVulnerabilitiesShowing, details.VulnerabilityDetails.TotalVulnerabilities)
				report.WriteString(filteredOutput)
			}
			if !vulCmdState.Html {
				report.WriteString("\nTry adding '--packages' to show a list of packages with CVE count.\n")
			}
		}
	}

	return report.String()
}

func buildVulnerabilitySummaryReportTable(response api.VulnerabilitiesContainersResponse) string {
	assessment := response.Data
	mainReport := &strings.Builder{}
	mainReport.WriteString(
		renderCustomTable(
			[]string{
				"Container Image Details",
				"Vulnerabilities",
			},
			[][]string{{
				renderCustomTable([]string{},
					vulContainerImageToTable(assessment[0].EvalCtx.ImageInfo),
					tableFunc(func(t *tablewriter.Table) {
						t.SetBorder(false)
						t.SetColumnSeparator("")
						t.SetAlignment(tablewriter.ALIGN_LEFT)
					}),
				),
				renderCustomTable([]string{"Severity", "Count", "Fixable"},
					vulContainerAssessmentToCountsTable(response),
					tableFunc(func(t *tablewriter.Table) {
						t.SetBorder(false)
						t.SetColumnSeparator(" ")
					}),
				),
			}},
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetAutoWrapText(false)
				t.SetColumnSeparator(" ")
			}),
		),
	)

	return mainReport.String()
}

func filterVulnContainerImagePackages(image []api.VulnerabilityContainer) filteredPackageTable {
	var filteredPackages []packageTable
	var aggregatedPackages []packageTable

	for _, i := range image {
		pack := packageTable{
			cveCount:       1,
			severity:       cases.Title(language.English).String(i.Severity),
			packageName:    i.FeatureKey.Name,
			currentVersion: i.FeatureKey.Version,
			fixVersion:     i.FixInfo.FixedVersion,
		}

		// filter fixable
		if vulCmdState.Fixable && i.FixInfo.FixedVersion == "" {
			filteredPackages = aggregatePackages(filteredPackages, pack)
			continue
		}

		//filter severity
		if vulCmdState.Severity != "" {
			if filterSeverity(i.Severity, vulCmdState.Severity) {
				filteredPackages = aggregatePackages(filteredPackages, pack)
				continue
			}
		}
		aggregatedPackages = aggregatePackages(aggregatedPackages, pack)
	}

	totalUnfiltered := len(filteredPackages) + len(aggregatedPackages)
	return filteredPackageTable{packages: aggregatedPackages, totalPackages: len(aggregatedPackages), totalUnfiltered: totalUnfiltered}
}

func vulContainerImagePackagesToTable(packageTable filteredPackageTable) [][]string {
	var out [][]string

	for _, p := range packageTable.packages {
		out = append(out, []string{
			strconv.Itoa(p.cveCount),
			p.severity,
			p.packageName,
			p.currentVersion,
			p.fixVersion,
		})
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return api.SeverityOrder(out[i][1]) < api.SeverityOrder(out[j][1])
	})

	return out
}

func filterVulnerabilityContainer(image []api.VulnerabilityContainer) filteredImageTable {
	var (
		vulns      = make(map[string]vulnTable)
		vulnIDs    []string
		vulnsCount int
		vulnList   []vulnTable
		filtered   []api.VulnerabilityContainer
	)

	for _, i := range image {
		vulnIDs = append(vulnIDs, fmt.Sprintf("%s-%s", i.VulnID, i.FeatureKey.Name))
		// filter: severity
		if vulCmdState.Severity != "" {
			if filterSeverity(i.Severity, vulCmdState.Severity) {
				continue
			}
		}
		// filter: fixable
		if vulCmdState.Fixable && i.FixInfo.FixedVersion == "" {
			continue
		}

		// Format IntroducedIn Field. In v2 response this field is not formatted with new lines.
		regex := regexp.MustCompile(regexAllTabs)
		introducedIn := regex.ReplaceAllString(i.FeatureProps.IntroducedIn, "\n")

		if _, ok := vulns[fmt.Sprintf("%s-%s", i.VulnID, i.FeatureKey.Name)]; !ok {
			vulns[fmt.Sprintf("%s-%s", i.VulnID, i.FeatureKey.Name)] = vulnTable{
				Name:           i.VulnID,
				Severity:       i.Severity,
				PackageName:    i.FeatureKey.Name,
				CurrentVersion: i.FeatureKey.Version,
				FixVersion:     i.FixInfo.FixedVersion,
				CreatedBy:      introducedIn,
				// Todo(v2): CVSSv3Score is missing from V2
				CVSSv3Score: 0,
				// Todo(v2): CVSSv2Score is missing from V2
				CVSSv2Score: 0,
				Status:      i.Status,
			}
			filtered = append(filtered, i)
		}
	}

	var uniqueIDs []string = array.Unique(vulnIDs)
	vulnsCount = len(uniqueIDs)

	for _, v := range vulns {
		vulnList = append(vulnList, v)
	}

	return filteredImageTable{
		Vulnerabilities:             vulnList,
		TotalVulnerabilitiesShowing: len(vulns),
		TotalVulnerabilities:        vulnsCount,
		Filtered:                    filtered,
	}
}

func vulContainerImageLayersToCSV(imageTable filteredImageTable) [][]string {
	var out [][]string
	for _, vuln := range imageTable.Vulnerabilities {
		out = append(out, []string{
			vuln.Name,
			vuln.Severity,
			strconv.FormatFloat(vuln.CVSSv2Score, 'f', 1, 64),
			strconv.FormatFloat(vuln.CVSSv3Score, 'f', 1, 64),
			vuln.PackageName,
			vuln.CurrentVersion,
			vuln.FixVersion,
			vuln.CreatedBy,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return api.SeverityOrder(out[i][1]) < api.SeverityOrder(out[j][1])
	})

	return out
}

func vulContainerImageLayersToTable(imageTable filteredImageTable) [][]string {
	var out [][]string
	for _, vuln := range imageTable.Vulnerabilities {
		out = append(out, []string{
			vuln.Name,
			vuln.Severity,
			vuln.PackageName,
			vuln.CurrentVersion,
			vuln.FixVersion,
			vuln.CreatedBy,
			vuln.Status,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return api.SeverityOrder(out[i][1]) < api.SeverityOrder(out[j][1])
	})

	return out
}

func vulContainerAssessmentToCountsTable(assessment api.VulnerabilitiesContainersResponse) [][]string {
	return [][]string{
		{"Critical", fmt.Sprint(assessment.CriticalVulnerabilities()),
			fmt.Sprint(assessment.VulnFixableCount("critical"))},
		{"High", fmt.Sprint(assessment.HighVulnerabilities()),
			fmt.Sprint(assessment.VulnFixableCount("high"))},
		{"Medium", fmt.Sprint(assessment.MediumVulnerabilities()),
			fmt.Sprint(assessment.VulnFixableCount("medium"))},
		{"Low", fmt.Sprint(assessment.LowVulnerabilities()),
			fmt.Sprint(assessment.VulnFixableCount("low"))},
		{"Info", fmt.Sprint(assessment.InfoVulnerabilities()),
			fmt.Sprint(assessment.VulnFixableCount("info"))},
	}
}

func vulContainerImageToTable(image api.ImageInfo) [][]string {
	return [][]string{
		{"ID", image.ID},
		{"Digest", image.Digest},
		{"Registry", image.Registry},
		{"Repository", image.Repo},
		{"Size", byteCountBinary(image.Size)},
		{"Created At", time.UnixMilli(image.CreatedTime).UTC().Format(time.RFC3339)},
		{"Tags", strings.Join(image.Tags, ",")},
	}
}

func aggregatePackages(slice []packageTable, s packageTable) []packageTable {
	for i, item := range slice {
		if item.packageName == s.packageName &&
			item.currentVersion == s.currentVersion &&
			item.severity == s.severity &&
			item.fixVersion == s.fixVersion {
			slice[i].cveCount++
			return slice
		}
	}
	return append(slice, s)
}

type vulnerabilityDetailsReport struct {
	VulnerabilityDetails filteredImageTable
	Packages             filteredPackageTable
}

type filteredImageTable struct {
	Vulnerabilities             []vulnTable
	TotalVulnerabilitiesShowing int
	TotalVulnerabilities        int
	Filtered                    []api.VulnerabilityContainer
}

type vulnTable struct {
	Name           string
	Severity       string
	PackageName    string
	CurrentVersion string
	FixVersion     string
	CreatedBy      string
	CVSSv2Score    float64
	CVSSv3Score    float64
	Status         string
}
