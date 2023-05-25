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
		Use:     "show-assessment <sha256:hash> [cve_id]",
		Aliases: []string{"show"},
		Short:   "Show results of a container vulnerability assessment",
		Long: `Show the vulnerability assessment results of the specified container.

Arguments:
    <sha256:hash> a sha256 hash of a container image (format: sha256:1ee...1d3b)

Note that the provided SHA is treated first as the image digest, but if no results
are found, this commands tries to use the SHA as the image id.

To request an on-demand vulnerability scan:

    lacework vulnerability container scan <registry> <repository> <tag|digest>

To see details for a single cve result in an assessment:

    lacework vulnerability show-assessment <sha256:hash> [cve_id]
`,
		Args: cobra.RangeArgs(1, 2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if vulCmdState.Csv {
				cli.EnableCSVOutput()

				// If --details or --packages is not passed, csv outputs nothing; defaulting to --details
				if !vulCmdState.Details && !vulCmdState.Packages {
					vulCmdState.Details = true
				}
			}

			if len(args) > 1 {
				vulCmdState.Cve = args[1]
			}

			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}
			err := showContainerAssessmentsWithSha256(args[0])
			var e *vulnerabilityPolicyError
			if errors.As(err, &e) {
				c.SilenceUsage = true
			}

			return err
		},
	}
)

func searchLastestEvaluationGuid(sha string) (string, error) {
	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -7) // 7 days from ago
		filter = api.SearchFilter{
			TimeFilter: &api.TimeFilter{
				StartTime: &before,
				EndTime:   &now,
			},
			Returns: []string{"evalGuid", "startTime"},
		}
	)

	// By default, we display the image digest in the command 'list-assessments',
	// so we start by fetching the image using the digest
	cli.Log.Infow("retrieve image assessment", "image_digest", sha)
	assessment, err := cli.LwApi.V2.Vulnerabilities.Containers.SearchAllPages(api.SearchFilter{
		Returns:    filter.Returns,
		TimeFilter: filter.TimeFilter,
		Filters: []api.Filter{{
			Expression: "eq",
			Field:      "evalCtx.image_info.digest",
			Value:      sha,
		}},
	})
	if err != nil {
		return "", err
	}

	if len(assessment.Data) == 0 {
		// provided sha was not an image digest, try using it as an image id instead
		cli.Log.Infow("retrieve image assessment", "image_id", sha)
		assessment, err = cli.LwApi.V2.Vulnerabilities.Containers.SearchAllPages(api.SearchFilter{
			TimeFilter: &api.TimeFilter{
				StartTime: &before,
				EndTime:   &now,
			},
			Filters: []api.Filter{{
				Expression: "eq",
				Field:      "evalCtx.image_info.id",
				Value:      sha,
			}},
		})
		if err != nil {
			return "", err
		}

		if len(assessment.Data) == 0 {
			return "", errors.New("no data found")
		}
	}

	return getUniqueEvalGUID(assessment), nil
}

func getUniqueEvalGUID(assessment api.VulnerabilitiesContainersResponse) string {
	var (
		guid      string
		startTime time.Time
	)
	for _, ctr := range assessment.Data {
		if ctr.EvalGUID != guid {
			if ctr.StartTime.After(startTime) {
				startTime = ctr.StartTime
				guid = ctr.EvalGUID
			}
		}
	}
	return guid
}

func showContainerAssessmentsWithSha256(sha string) error {
	var (
		cacheKey   = fmt.Sprintf("vulnerability/container/%s", sha)
		assessment api.VulnerabilitiesContainersResponse
	)
	expired := cli.ReadCachedAsset(cacheKey, &assessment)
	if expired {
		// search for the latest evaluation guid
		cli.StartProgress("Searching for latest container evaluation...")
		evalGUID, err := searchLastestEvaluationGuid(sha)
		cli.StopProgress()
		if err != nil {
			return errors.Wrapf(err, "unable to find assessment information of image %s", sha)
		}

		cli.Log.Infow("latest assessment found", "eval_guid", evalGUID)

		var (
			now    = time.Now().UTC()
			before = now.AddDate(0, 0, -7) // 7 days from ago
		)

		cli.StartProgress(
			fmt.Sprintf("Fetching assessment from container evaluation '%s'...", evalGUID),
		)
		assessment, err = cli.LwApi.V2.Vulnerabilities.Containers.SearchAllPages(api.SearchFilter{
			TimeFilter: &api.TimeFilter{
				StartTime: &before,
				EndTime:   &now,
			},
			Filters: []api.Filter{{
				Expression: "eq",
				Field:      "evalGuid",
				Value:      evalGUID,
			}},
		})
		cli.StopProgress()
		if err != nil {
			return errors.Wrap(err, "unable to fetch assessment data")
		}

		// write to cache if the request was successful
		cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), assessment)
	} else {
		cli.Log.Infow("assessment loaded from cache", "data_points", len(assessment.Data))
	}

	filterContainerAssessmentByVulnerable(&assessment)
	cli.Log.Infow("filtered assessment (status = vulnerable)", "data_points", len(assessment.Data))

	if err := outputContainerVulnerabilityAssessment(assessment); err != nil {
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

func buildVulnContainerAssessmentCve(assessment api.VulnerabilitiesContainersResponse) error {
	assessment.FilterSingleVulnIDData(vulCmdState.Cve)

	if cli.JSONOutput() {
		return cli.OutputJSON(assessment.Data)
	}

	if len(assessment.Data) == 0 {
		cli.OutputHuman("Unable to find results for '%s'\n", vulCmdState.Cve)
		return nil
	}

	var details vulnerabilityDetailsReport
	details.VulnerabilityDetails = filterVulnerabilityContainer(assessment.Data)

	cli.OutputHuman(buildVulnerabilitySingleCveReportTable(details))
	return nil
}

func filterContainerAssessmentByVulnerable(assessment *api.VulnerabilitiesContainersResponse) {
	var vulnerabilities []api.VulnerabilityContainer
	for _, a := range assessment.Data {
		if a.Status == "VULNERABLE" {
			vulnerabilities = append(vulnerabilities, a)
		}
	}
	assessment.Data = vulnerabilities
}

func outputContainerVulnerabilityAssessment(assessment api.VulnerabilitiesContainersResponse) error {
	if len(assessment.Data) == 0 {
		if cli.JSONOutput() {
			// if no assessments are found return empty array
			return cli.OutputJSON([]any{})
		}
		cli.OutputHuman(
			"Great news! This container image has no vulnerabilities... (time for %s)\n",
			randomEmoji(),
		)
		return nil
	}

	if vulCmdState.Cve != "" {
		if err := buildVulnContainerAssessmentCve(assessment); err != nil {
			return err
		}
	} else {
		if err := buildVulnContainerAssessmentReports(assessment); err != nil {
			return err
		}
	}

	return nil
}

// Build the cli output for vuln ctr 'show-assessment' command
func buildVulnContainerAssessmentReports(assessment api.VulnerabilitiesContainersResponse) error {
	var (
		filteredAssessment = assessment
		summaryReport      = buildVulnerabilitySummaryReportTable(assessment)
		details            vulnerabilityDetailsReport
	)

	details.VulnerabilityDetails = filterVulnerabilityContainer(assessment.Data)
	filteredAssessment.Data = details.VulnerabilityDetails.Filtered
	details.Packages = filterVulnContainerImagePackages(details.VulnerabilityDetails.Filtered)
	details.Packages.totalUnfiltered = countVulnContainerImagePackages(assessment.Data)

	switch {
	case cli.JSONOutput():
		return cli.OutputJSON(filteredAssessment.Data)
	case cli.CSVOutput():
		if !(vulCmdState.Details || vulCmdState.Packages || vulFiltersEnabled()) {
			return nil
		}

		if vulCmdState.Packages {
			return cli.OutputCSV([]string{
				"CVE Count", "Severity", "Package", "Current Version", "Fix Version"},
				vulContainerImagePackagesToTable(details.Packages))
		}

		return cli.OutputCSV([]string{
			"CVE ID", "Severity", "CVSSv2", "CVSSv3", "Package", "Current Version",
			"Fix Version", "Version Format", "Feed", "Src", "Start Time", "Status",
			"Namespace", "Image Digest", "Image ID", "Image Repo", "Image Registry",
			"Image Size", "Introduced in Layer"},
			vulContainerImageLayersToCSV(filteredAssessment.Data))
	default:
		if len(filteredAssessment.Data) == 0 {
			if vulCmdState.Severity != "" {
				cli.OutputHuman("There are no vulnerabilities found for this severity\n")
				return nil
			}

			cli.OutputHuman(
				"Great news! This container image has no vulnerabilities... (time for %s)\n",
				randomEmoji(),
			)
			return nil
		}

		detailsReport := buildVulnerabilityDetailsReportTable(details)
		cli.OutputHuman(buildVulnContainerAssessmentReportTable(summaryReport, detailsReport))
		if vulCmdState.Html {
			if err := generateVulnAssessmentHTML(filteredAssessment); err != nil {
				return err
			}
		}
	}

	return nil
}

func buildVulnerabilitySingleCveReportTable(details vulnerabilityDetailsReport) string {
	var singleCveTable = strings.Builder{}
	var singleCveTableContent = strings.Builder{}

	for _, cveData := range details.VulnerabilityDetails.Vulnerabilities {
		detailsTable := renderCustomTable([]string{},
			[][]string{
				{"PACKAGE", cveData.PackageName},
				{"CURRENT VERSION", cveData.CurrentVersion},
				{"NAMESPACE", cveData.Namespace},
				{"SEVERITY", cveData.Severity},
				{"FEED", cveData.Feed},
				{"SRC", cveData.Src},
				{"START TIME", cveData.StartTime},
				{"VERSION FORMAT", cveData.VersionFormat},
				{"FIX VERSION", cveData.FixVersion},
				{"STATUS", cveData.Status},
			},
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
				t.SetAutoWrapText(false)
				t.SetColWidth(150)
				t.SetAlignment(tablewriter.ALIGN_LEFT)
			}),
		)

		layerTable := renderCustomTable([]string{"INTRODUCED IN LAYERS"},
			introducedInLayerToTable(cveData),
			tableFunc(func(t *tablewriter.Table) {
				t.SetRowLine(true)
				t.SetBorder(false)
				t.SetAutoWrapText(true)
				t.SetAlignment(tablewriter.ALIGN_LEFT)
				t.SetColumnSeparator(" ")
			}),
		)

		singleCveTableContent.WriteString(fmt.Sprintf("%s\n", detailsTable))
		singleCveTableContent.WriteString(layerTable)
	}

	cveSummaryTable := renderOneLineCustomTable(details.VulnerabilityDetails.Vulnerabilities[0].Name,
		singleCveTableContent.String(),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		}),
	)

	singleCveTable.WriteString(cveSummaryTable)
	return singleCveTable.String()
}

func introducedInLayerToTable(vuln vulnTable) (resourceTable [][]string) {
	for _, layer := range vuln.CreatedBy {
		resourceTable = append(resourceTable, []string{layer})
	}
	return
}

func buildVulnerabilityDetailsReportTable(details vulnerabilityDetailsReport) string {
	report := &strings.Builder{}

	if vulCmdState.Details || vulCmdState.Packages || vulFiltersEnabled() || vulCmdState.Cve != "" {
		if vulCmdState.Packages {
			vulnPackagesTable := vulContainerImagePackagesToTable(details.Packages)

			report.WriteString(
				renderSimpleTable(
					[]string{"CVE Count", "Severity", "Package", "Current Version", "Fix Version"},
					vulnPackagesTable,
				),
			)

			if vulFiltersEnabled() {
				filteredOutput := fmt.Sprintf("%v of %v packages showing \n",
					details.Packages.totalPackages, details.Packages.totalUnfiltered)
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
	var (
		vulnKey  string
		vulnKeys = []string{}
	)

	for _, i := range image {
		vulnKey = fmt.Sprintf("%s-%s-%s", i.VulnID, i.FeatureKey.Name, i.FeatureKey.Version)
		if array.ContainsStr(vulnKeys, vulnKey) {
			continue
		}
		vulnKeys = append(vulnKeys, vulnKey)
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

		// filter severity
		if vulCmdState.Severity != "" {
			if filterSeverity(i.Severity, vulCmdState.Severity) {
				filteredPackages = aggregatePackages(filteredPackages, pack)
				continue
			}
		}
		aggregatedPackages = aggregatePackages(aggregatedPackages, pack)
	}

	return filteredPackageTable{packages: aggregatedPackages, totalPackages: len(aggregatedPackages)}
}

func countVulnContainerImagePackages(image []api.VulnerabilityContainer) int {
	var aggregatedPackages []packageTable

	for _, i := range image {
		pack := packageTable{
			cveCount:       1,
			severity:       cases.Title(language.English).String(i.Severity),
			packageName:    i.FeatureKey.Name,
			currentVersion: i.FeatureKey.Version,
			fixVersion:     i.FixInfo.FixedVersion,
		}
		aggregatedPackages = aggregatePackages(aggregatedPackages, pack)
	}

	return len(aggregatedPackages)
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

	// order by severity and package name
	sort.Slice(out, func(i, j int) bool {
		if api.SeverityOrder(out[i][1]) == api.SeverityOrder(out[j][1]) {
			return out[i][2] < out[j][2]
		}
		return api.SeverityOrder(out[i][1]) < api.SeverityOrder(out[j][1])
	})

	return out
}

func filterVulnerabilityContainer(image []api.VulnerabilityContainer) filteredImageTable {
	var (
		vulns           = make(map[string]vulnTable)
		introducedInMap = make(map[string][]string)
		vulnIDs         []string
		vulnsCount      int
		vulnList        []vulnTable
		filtered        []api.VulnerabilityContainer
	)

	for _, i := range image {
		vulnKey := fmt.Sprintf("%s-%s-%s", i.VulnID, i.FeatureKey.Name, i.FeatureKey.Version)
		vulnIDs = append(vulnIDs, vulnKey)
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

		introducedInMap[vulnKey] = append(introducedInMap[vulnKey], introducedIn)

		if _, ok := vulns[vulnKey]; !ok {
			vulns[vulnKey] = vulnTable{
				Name:           i.VulnID,
				Severity:       i.Severity,
				PackageName:    i.FeatureKey.Name,
				CurrentVersion: i.FeatureKey.Version,
				FixVersion:     i.FixInfo.FixedVersion,
				Namespace:      i.FeatureKey.Namespace,
				Feed:           i.FeatureProps.Feed,
				StartTime:      i.StartTime.Format(time.RFC3339),
				VersionFormat:  i.FeatureProps.VersionFormat,
				Src:            i.FeatureProps.Src,
				// Todo(v2): CVSSv3Score is missing from V2
				CVSSv3Score: 0,
				// Todo(v2): CVSSv2Score is missing from V2
				CVSSv2Score: 0,
				Status:      i.Status,
			}
		}

		filtered = append(filtered, i)
	}

	// Set the aggregated introduced by layers for each vuln
	for k, v := range introducedInMap {
		vulnTable := vulns[k]
		vulnTable.CreatedBy = v
		vulns[k] = vulnTable
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

func vulContainerImageLayersToCSV(assessment []api.VulnerabilityContainer) [][]string {
	var out [][]string
	for _, vuln := range assessment {
		out = append(out, []string{
			vuln.VulnID,
			vuln.Severity,
			strconv.FormatFloat(0.0, 'f', 1, 64),
			strconv.FormatFloat(0.0, 'f', 1, 64),
			vuln.FeatureKey.Name,
			vuln.FeatureKey.Version,
			vuln.FixInfo.FixedVersion,
			vuln.FeatureProps.VersionFormat,
			vuln.FeatureProps.Feed,
			vuln.FeatureProps.Src,
			vuln.StartTime.Format(time.RFC3339),
			vuln.FeatureKey.Namespace,
			vuln.EvalCtx.ImageInfo.ID,
			vuln.EvalCtx.ImageInfo.Digest,
			vuln.EvalCtx.ImageInfo.Repo,
			vuln.EvalCtx.ImageInfo.Registry,
			strconv.Itoa(vuln.EvalCtx.ImageInfo.Size),
			vuln.FeatureProps.IntroducedIn,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if api.SeverityOrder(out[i][1]) == api.SeverityOrder(out[j][1]) {
			return out[i][4] < out[j][4]
		}
		return api.SeverityOrder(out[i][1]) < api.SeverityOrder(out[j][1])
	})

	return out
}

func vulContainerImageLayersToTable(imageTable filteredImageTable) [][]string {
	var out [][]string
	var createdByKeys = make(map[string]bool)

	for _, vuln := range imageTable.Vulnerabilities {
		introducedBy := strings.Join(vuln.CreatedBy, ",")
		// if the same vuln is introduced in more than 1 layer, only display the number of layers
		if len(vuln.CreatedBy) > 1 {
			introducedBy = fmt.Sprintf("introduced in %d layers...", len(vuln.CreatedBy))
		}

		if !createdByKeys[fmt.Sprintf("%s-%s", vuln.Name, vuln.CurrentVersion)] {
			out = append(out, []string{
				vuln.Name,
				vuln.Severity,
				vuln.PackageName,
				vuln.CurrentVersion,
				vuln.FixVersion,
				introducedBy,
				vuln.Status,
			})
		}

		createdByKeys[fmt.Sprintf("%s-%s", vuln.Name, vuln.CurrentVersion)] = true
	}

	sort.Slice(out, func(i, j int) bool {
		if api.SeverityOrder(out[i][1]) == api.SeverityOrder(out[j][1]) {
			return out[i][2] < out[j][2]
		}
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
	CreatedBy      []string
	CVSSv2Score    float64
	CVSSv3Score    float64
	Status         string
	Namespace      string
	Feed           string
	Src            string
	StartTime      string
	VersionFormat  string
}
