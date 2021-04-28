//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
)

var (
	// vulContainerScanCmd represents the scan sub-command inside the container vulnerability command
	vulContainerScanCmd = &cobra.Command{
		Use:   "scan <registry> <repository> <tag|digest>",
		Short: "request an on-demand container vulnerability assessment",
		Long: `Request on-demand container vulnerability assessments and view the generated results.

NOTE: Scans can take up to 15 minutes to return results.

Arguments:
  <registry>    container registry where the container image has been published
  <repository>  repository name that contains the container image
  <tag|digest>  either a tag or an image digest to scan (digest format: sha256:1ee...1d3b)

To list all container registries configured in your account:

    $ lacework vulnerability container list-registries`,
		Args: cobra.ExactArgs(3),
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}

			err := requestOnDemandContainerVulnerabilityScan(args)
			var e *vulnerabilityPolicyError
			if errors.As(err, &e) {
				c.SilenceUsage = true
			}

			return err
		},
	}

	// vulContainerScanStatusCmd represents the scan-status sub-command inside the container
	// vulnerability command
	vulContainerScanStatusCmd = &cobra.Command{
		Use:     "scan-status <request_id>",
		Aliases: []string{"status"},
		Short:   "check the status of an on-demand container vulnerability assessment",
		Long:    "Check the status of an on-demand container vulnerability assessment.",
		Args:    cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}
			err := checkOnDemandContainerVulnerabilityStatus(args[0])
			var e *vulnerabilityPolicyError
			if errors.As(err, &e) {
				c.SilenceUsage = true
			}

			return err
		},
	}

	// vulContainerListRegistriesCmd represents the list-registries sub-command inside the container
	// vulnerability command
	vulContainerListRegistriesCmd = &cobra.Command{
		Use:     "list-registries",
		Aliases: []string{"list-reg", "registries"},
		Short:   "list all container registries configured",
		Long:    `List all container registries configured in your account.`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			registries, err := getContainerRegistries()
			if err != nil {
				return err
			}
			if len(registries) == 0 {
				msg := `There are no container registries configured in your account.

Get started by integrating your container registry using the command:

    $ lacework integration create

If you prefer to configure the integration via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Integrations > Container Registry.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(registries)
			}

			var rows [][]string
			for _, acc := range registries {
				rows = append(rows, []string{acc})
			}

			cli.OutputHuman(renderSimpleTable([]string{"Container Registries"}, rows))
			return nil
		},
	}

	// vulContainerListAssessmentsCmd represents the list-assessments sub-command inside the container
	// vulnerability command
	vulContainerListAssessmentsCmd = &cobra.Command{
		Use:     "list-assessments",
		Aliases: []string{"list", "ls"},
		Short:   "list container vulnerability assessments (default last 7 days)",
		Long: `List all container vulnerability assessments for the last 7 days by default, or
pass --start and --end to specify a custom time range. You can also pass --active
to filter on active containers in your environment, as well as pass --fixable to
filter on containers with vulnerabilities that have fixes available.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			var (
				response api.VulnContainerAssessmentsResponse
				err      error
			)
			if vulCmdState.Start != "" || vulCmdState.End != "" {
				start, end, errT := parseStartAndEndTime(vulCmdState.Start, vulCmdState.End)
				if errT != nil {
					return errors.Wrap(errT, "unable to parse time range")
				}

				cli.Log.Infow("requesting list of assessments from custom time range",
					"start_time", start, "end_time", end,
				)
				response, err = cli.LwApi.Vulnerabilities.Container.ListAssessmentsDateRange(start, end)
			} else {
				cli.Log.Info("requesting list of assessments from the last 7 days")
				response, err = cli.LwApi.Vulnerabilities.Container.ListAssessments()
			}

			if err != nil {
				return errors.Wrap(err, "unable to get assessments")
			}

			cli.Log.Debugw("assessments", "raw", response)

			if len(response.Assessments) == 0 {
				cli.OutputHuman("There are no container assessments for this environment.\n")
				return nil
			}

			// filter assessments by repositories, if the user doesn't provide a filter
			// the function returns all the assessments
			assessments := filterAssessmentsByRepositories(response.Assessments)

			// if the user wants to show only assessments of running containers
			// order them by that field, number of running containers
			if vulCmdState.Active {
				// Sort the assessments by running containers
				sort.Slice(assessments, func(i, j int) bool {
					return stringToInt(assessments[i].NdvContainers) > stringToInt(assessments[j].NdvContainers)
				})
			} else {
				// Sort the assessments by date
				sort.Slice(assessments, func(i, j int) bool {
					return assessments[i].StartTime.ToTime().After(assessments[j].StartTime.ToTime())
				})
			}

			filteredAssessments := filterAssessments(assessments)

			if cli.JSONOutput() {
				return cli.OutputJSON(filteredAssessments)
			}

			rows := vulAssessmentsToTable(filteredAssessments)

			// if the user wants to show only assessments of containers running
			// and we don't have any, show a friendly message
			if len(rows) == 0 {
				cli.OutputHuman(buildContainerAssessmentsError())
			} else {
				cli.OutputHuman(
					renderSimpleTable(
						[]string{"Registry", "Repository", "Last Scan", "Status",
							"Containers", "Vulnerabilities", "Image Digest"},
						rows,
					),
				)
				if !vulCmdState.Active {
					cli.OutputHuman(
						"\nTry adding '--active' to only show assessments of containers actively running with vulnerabilities.\n",
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

	// vulContainerShowAssessmentCmd represents the show-assessment sub-command inside the container
	// vulnerability command
	vulContainerShowAssessmentCmd = &cobra.Command{
		Use:     "show-assessment <sha256:hash>",
		Aliases: []string{"show"},
		Short:   "show results of a container vulnerability assessment",
		Long: `Show the results from a vulnerability assessment of a specified container.

Arguments:
  <sha256:hash> a sha256 hash of a container image (format: sha256:1ee...1d3b)

By default, this command expects a sha256 image digest or tag. To lookup an
assessment by its image id, use the flag '--image_id' followed by the sha256
image id.

To request an on-demand vulnerability scan:

    $ lacework vulnerability container scan <registry> <repository> <tag|digest>`,
		Args: cobra.ExactArgs(1),
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

func init() {
	// add sub-commands to the 'vulnerability container' command
	vulContainerCmd.AddCommand(vulContainerScanCmd)
	vulContainerCmd.AddCommand(vulContainerScanStatusCmd)
	vulContainerCmd.AddCommand(vulContainerListAssessmentsCmd)
	vulContainerCmd.AddCommand(vulContainerListRegistriesCmd)
	vulContainerCmd.AddCommand(vulContainerShowAssessmentCmd)

	// add start flag to list-assessments command
	vulContainerListAssessmentsCmd.Flags().StringVar(&vulCmdState.Start,
		"start", "", "start of the time range in UTC (format: yyyy-MM-ddTHH:mm:ssZ)",
	)
	// add end flag to list-assessments command
	vulContainerListAssessmentsCmd.Flags().StringVar(&vulCmdState.End,
		"end", "", "end of the time range in UTC (format: yyyy-MM-ddTHH:mm:ssZ)",
	)
	// add active flag to list-assessments command
	vulContainerListAssessmentsCmd.Flags().BoolVar(&vulCmdState.Active,
		"active", false, "only show assessments of containers actively running with vulnerabilities in your environment",
	)
	// add repository flag to list-assessments command
	vulContainerListAssessmentsCmd.Flags().StringSliceVarP(&vulCmdState.Repositories,
		"repository", "r", []string{}, "filter assessments for specific repositories",
	)

	setPollFlag(
		vulContainerScanCmd.Flags(),
		vulContainerScanStatusCmd.Flags(),
	)

	setHtmlFlag(
		vulContainerScanCmd.Flags(),
		vulContainerScanStatusCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setDetailsFlag(
		vulContainerScanCmd.Flags(),
		vulContainerScanStatusCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setSeverityFlag(
		vulContainerScanCmd.Flags(),
		vulContainerScanStatusCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setFailOnSeverityFlag(
		vulContainerScanCmd.Flags(),
		vulContainerScanStatusCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setFailOnFixableFlag(
		vulContainerScanCmd.Flags(),
		vulContainerScanStatusCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setFixableFlag(
		vulContainerScanCmd.Flags(),
		vulContainerScanStatusCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
		vulContainerListAssessmentsCmd.Flags(),
	)

	setPackagesFlag(
		vulContainerScanCmd.Flags(),
		vulContainerScanStatusCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	vulContainerShowAssessmentCmd.Flags().BoolVar(
		&vulCmdState.ImageID, "image_id", false,
		"tread the provided sha256 hash as image id",
	)
}

func filterAssessmentsByRepositories(assessments []api.VulnContainerAssessmentSummary) []api.VulnContainerAssessmentSummary {
	if len(vulCmdState.Repositories) == 0 {
		return assessments
	}

	var filtered []api.VulnContainerAssessmentSummary
	for _, assessment := range assessments {
		// for every repository that the user is filtering for
		for _, repo := range vulCmdState.Repositories {
			if strings.Contains(assessment.ImageRepo, repo) {
				filtered = append(filtered, assessment)
			}
		}

	}

	return filtered
}

func requestOnDemandContainerVulnerabilityScan(args []string) error {
	cli.Log.Debugw("requesting vulnerability scan",
		"registry", args[0],
		"repository", args[1],
		"tag_or_digest", args[2],
	)
	scan, err := cli.LwApi.Vulnerabilities.Container.Scan(args[0], args[1], args[2])
	if err != nil {
		return userFriendlyErrorForOnDemandCtrVulnScan(err, args[0], args[1], args[2])
	}

	cli.Log.Debugw("vulnerability scan", "details", scan)
	if !scan.Ok {
		return errors.Errorf(
			"there is a problem with the vulnerability scan: %s",
			scan.Message,
		)
	}

	cli.OutputHuman(
		"A new vulnerability scan has been requested. (request_id: %s)\n\n",
		scan.Data.RequestID,
	)

	if vulCmdState.Poll {
		cli.Log.Infow("tracking scan progress",
			"param", "--poll",
			"request_id", scan.Data.RequestID,
		)
		return pollScanStatus(scan.Data.RequestID)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(scan.Data)
	}

	cli.OutputHuman("To track the progress of the scan, use the command:\n")
	cli.OutputHuman("  $ lacework vulnerability container scan-status %s\n", scan.Data.RequestID)
	return nil
}

func checkOnDemandContainerVulnerabilityStatus(reqID string) error {
	if vulCmdState.Poll {
		cli.Log.Infow("tracking scan progress",
			"param", "--poll",
			"request_id", reqID,
		)
		return pollScanStatus(reqID)
	}

	results, err, scanning := checkScanStatus(reqID)
	if err != nil {
		return err
	}

	// if the scan is still running, display a nice message
	if scanning {
		cli.OutputHuman(
			"The vulnerability scan is still running. (request_id: %s)\n\n",
			reqID,
		)
		cli.OutputHuman("Use '--poll' to poll until the vulnerability scan completes.\n")
		return nil
	}

	if err := buildVulnContainerAssessmentReports(results); err != nil {
		return err
	}

	if vulFailureFlagsEnabled() {
		cli.Log.Infow("failure flags enabled",
			"fail_on_severity", vulCmdState.FailOnSeverity,
			"fail_on_fixable", vulCmdState.FailOnFixable,
		)
		vulnPolicy := NewVulnerabilityPolicyError(
			results,
			vulCmdState.FailOnSeverity,
			vulCmdState.FailOnFixable,
		)
		if vulnPolicy.NonCompliant() {
			return vulnPolicy
		}
	}

	return nil
}

func showContainerAssessmentsWithSha256(sha string) error {
	var (
		assessment  api.VulnContainerAssessmentResponse
		searchField string
		err         error
	)

	if vulCmdState.ImageID {
		searchField = "image_id"
		cli.Log.Debugw("retrieve image assessment", searchField, sha)
		assessment, err = cli.LwApi.Vulnerabilities.Container.AssessmentFromImageID(sha)
	} else {
		searchField = "digest"
		cli.Log.Debugw("retrieve image assessment", searchField, sha)
		assessment, err = cli.LwApi.Vulnerabilities.Container.AssessmentFromImageDigest(sha)
	}
	if err != nil {
		return errors.Wrap(err, "unable to show vulnerability assessment")
	}

	cli.Log.Debugw("image assessment", "details", assessment)
	status := assessment.CheckStatus()
	switch status {
	case "Success":
		if err := buildVulnContainerAssessmentReports(&assessment.Data); err != nil {
			return err
		}
	case "Unsupported":
		return errors.Errorf(
			`unable to retrieve assessment for the provided container image. (unsupported distribution)

For more information about supported distributions, visit:
    https://support.lacework.com/hc/en-us/articles/360035472393-Container-Vulnerability-Assessment-Overview
`,
		)
	case "NotFound":
		msg := fmt.Sprintf(
			"unable to find any assessment from a container image with %s '%s'",
			searchField, sha,
		)

		// add a suggestion to the user in regards of the image_id vs digest
		if !vulCmdState.ImageID {
			msg = fmt.Sprintf("%s\n\n(?) Are you trying to lookup a container vulnerability assessment using an image id?", msg)
			msg = fmt.Sprintf("%s\n(?) Try using the flag '--image_id'", msg)
		}

		return errors.New(msg)
	case "Failed":
		return errors.New(
			"the assessment failed to execute. Use '--debug' to troubleshoot.",
		)
	default:
		return errors.New(
			"unable to get assessment status from the container image. Use '--debug' to troubleshoot.",
		)
	}

	if vulFailureFlagsEnabled() {
		cli.Log.Infow("failure flags enabled",
			"fail_on_severity", vulCmdState.FailOnSeverity,
			"fail_on_fixable", vulCmdState.FailOnFixable,
		)
		vulnPolicy := NewVulnerabilityPolicyError(
			&assessment.Data,
			vulCmdState.FailOnSeverity,
			vulCmdState.FailOnFixable,
		)
		if vulnPolicy.NonCompliant() {
			return vulnPolicy
		}
	}

	return nil
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
						"Fix Version", "Introduced in Layer"},
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

func buildVulnerabilitySummaryReportTable(assessment *api.VulnContainerAssessment) string {
	if assessment.TotalVulnerabilities == 0 {
		return fmt.Sprintf("Great news! This container image has no vulnerabilities... (time for %s)\n", randomEmoji())
	}

	mainReport := &strings.Builder{}
	mainReport.WriteString(
		renderCustomTable(
			[]string{
				"Container Image Details",
				"Vulnerabilities",
			},
			[][]string{{
				renderCustomTable([]string{},
					vulContainerImageToTable(assessment.Image),
					tableFunc(func(t *tablewriter.Table) {
						t.SetBorder(false)
						t.SetColumnSeparator("")
						t.SetAlignment(tablewriter.ALIGN_LEFT)
					}),
				),
				renderCustomTable([]string{"Severity", "Count", "Fixable"},
					vulContainerAssessmentToCountsTable(assessment),
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

type packageTable struct {
	cveCount       int
	severity       string
	packageName    string
	currentVersion string
	fixVersion     string
	packageStatus  string
	hostCount      int
}

type vulnTable struct {
	Name           string
	Severity       string
	PackageName    string
	CurrentVersion string
	FixVersion     string
	CreatedBy      string
}

type filteredPackageTable struct {
	packages        []packageTable
	totalPackages   int
	totalUnfiltered int
}

type filteredImageTable struct {
	Vulnerabilities             []vulnTable
	TotalVulnerabilitiesShowing int
	TotalVulnerabilities        int
	ImageLayers                 []api.VulnContainerImageLayer
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

func filterVulnContainerImagePackages(image *api.VulnContainerImage) filteredPackageTable {
	if image == nil {
		return filteredPackageTable{}
	}
	var filteredPackages []packageTable
	var aggregatedPackages []packageTable

	for _, layer := range image.ImageLayers {
		for _, pkg := range layer.Packages {
			for _, vul := range pkg.Vulnerabilities {
				pack := packageTable{
					cveCount:       1,
					severity:       strings.Title(vul.Severity),
					packageName:    pkg.Name,
					currentVersion: pkg.Version,
					fixVersion:     vul.FixVersion,
				}

				if vulCmdState.Fixable && vul.FixVersion == "" {
					filteredPackages = aggregatePackages(filteredPackages, pack)
					continue
				}

				if vulCmdState.Severity != "" {
					if filterSeverity(vul.Severity, vulCmdState.Severity) {
						filteredPackages = aggregatePackages(filteredPackages, pack)
						continue
					}
				}
				aggregatedPackages = aggregatePackages(aggregatedPackages, pack)
			}
		}
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
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	return out
}

func filterVulContainerImageLayers(image *api.VulnContainerImage) filteredImageTable {
	if image == nil {
		return filteredImageTable{}
	}
	var (
		vulns               []vulnTable
		vulnsCount          int
		filteredImageLayers []api.VulnContainerImageLayer
		filteredPkg         api.VulnContainerPackage
		filteredImageLayer  api.VulnContainerImageLayer
	)

	for _, layer := range image.ImageLayers {
		filteredImageLayer = layer
		filteredImageLayer.Packages = []api.VulnContainerPackage{}
		for _, pkg := range layer.Packages {
			filteredPkg = pkg
			filteredPkg.Vulnerabilities = []api.ContainerVulnerability{}
			for _, vul := range pkg.Vulnerabilities {
				vulnsCount++
				if vulCmdState.Fixable && vul.FixVersion == "" {
					continue
				}

				if vulCmdState.Severity != "" {
					if filterSeverity(vul.Severity, vulCmdState.Severity) {
						continue
					}
				}

				space := regexp.MustCompile(`\s+`)
				createdBy := space.ReplaceAllString(layer.CreatedBy, " ")

				vulns = append(vulns, vulnTable{
					Name:           vul.Name,
					Severity:       strings.Title(vul.Severity),
					PackageName:    pkg.Name,
					CurrentVersion: pkg.Version,
					FixVersion:     vul.FixVersion,
					CreatedBy:      createdBy,
				})

				filteredPkg.Vulnerabilities = append(filteredPkg.Vulnerabilities, vul)
			}
			if len(filteredPkg.Vulnerabilities) > 0 {
				filteredImageLayer.Packages = append(filteredImageLayer.Packages, filteredPkg)
			}
		}
		if len(filteredImageLayer.Packages) > 0 {
			filteredImageLayers = append(filteredImageLayers, filteredImageLayer)
		}
	}

	return filteredImageTable{
		Vulnerabilities:             vulns,
		TotalVulnerabilitiesShowing: len(vulns),
		TotalVulnerabilities:        vulnsCount,
		ImageLayers:                 filteredImageLayers}
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
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	return out
}

func vulContainerAssessmentToCountsTable(assessment *api.VulnContainerAssessment) [][]string {
	return [][]string{
		{"Critical", fmt.Sprint(assessment.CriticalVulnerabilities),
			fmt.Sprint(assessment.VulnFixableCount("critical"))},
		{"High", fmt.Sprint(assessment.HighVulnerabilities),
			fmt.Sprint(assessment.VulnFixableCount("high"))},
		{"Medium", fmt.Sprint(assessment.MediumVulnerabilities),
			fmt.Sprint(assessment.VulnFixableCount("medium"))},
		{"Low", fmt.Sprint(assessment.LowVulnerabilities),
			fmt.Sprint(assessment.VulnFixableCount("low"))},
		{"Info", fmt.Sprint(assessment.InfoVulnerabilities),
			fmt.Sprint(assessment.VulnFixableCount("info"))},
	}
}

func vulContainerImageToTable(image *api.VulnContainerImage) [][]string {
	if image == nil || image.ImageInfo == nil {
		return [][]string{}
	}

	info := image.ImageInfo
	return [][]string{
		{"ID", info.ImageID},
		{"Digest", info.ImageDigest},
		{"Registry", info.Registry},
		{"Repository", info.Repository},
		{"Size", byteCountBinary(info.Size)},
		{"Created At", info.CreatedTime},
		{"Tags", strings.Join(info.Tags, ",")},
	}
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

	if vulCmdState.Fixable {
		msg = fmt.Sprintf("%s with fixable vulnerabilities", msg)
	}

	return fmt.Sprintf("%s in your environment.\n", msg)
}

func filterAssessments(assessments []api.VulnContainerAssessmentSummary) []assessmentOutput {
	var out []assessmentOutput
	for _, assessment := range assessments {
		// do not add assessments that doesn't have running containers
		// if the user wants to show only assessments of containers running
		if vulCmdState.Active && assessment.NdvContainers == "0" {
			continue
		}
		if vulCmdState.Fixable && assessment.NumFixes == "0" {
			continue
		}

		// if an assessment is unsupported, the summary should not be generated
		var (
			assessmentSummary  = "-"
			hasVulnerabilities bool
		)
		if assessment.ImageScanStatus != "Unsupported" {
			assessmentSummary, hasVulnerabilities = vulSummaryFromAssessment(&assessment)
			if vulCmdState.Active && !hasVulnerabilities {
				continue
			}
		}

		if vulCmdState.Active && assessment.ImageScanStatus == "Unsupported" {
			continue
		}

		out = append(out, assessmentOutput{
			imageRegistry:     assessment.ImageRegistry,
			imageRepo:         assessment.ImageRepo,
			startTime:         assessment.StartTime.UTC().Format(time.RFC3339),
			imageScanStatus:   assessment.ImageScanStatus,
			ndvContainers:     assessment.NdvContainers,
			assessmentSummary: assessmentSummary,
			imageDigest:       assessment.ImageDigest,
		})
	}
	return out
}

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

func vulSummaryFromAssessment(assessment *api.VulnContainerAssessmentSummary) (string, bool) {
	var summary []string

	summary = addToAssessmentSummary(summary, assessment.NumVulnerabilitiesSeverity1, "Critical")
	summary = addToAssessmentSummary(summary, assessment.NumVulnerabilitiesSeverity2, "High")
	summary = addToAssessmentSummary(summary, assessment.NumVulnerabilitiesSeverity3, "Medium")
	summary = addToAssessmentSummary(summary, assessment.NumVulnerabilitiesSeverity4, "Low")
	summary = addToAssessmentSummary(summary, assessment.NumVulnerabilitiesSeverity5, "Info")

	if len(summary) == 0 {
		return fmt.Sprintf("None! Time for %s", randomEmoji()), false
	}

	if assessment.NumFixes != "" {
		summary = append(summary, fmt.Sprintf("%s Fixable", assessment.NumFixes))
	}

	return strings.Join(summary, " "), true
}

func addToAssessmentSummary(text []string, num, severity string) []string {
	if len(text) == 0 {
		if num != "" && num != "0" {
			return append(text, fmt.Sprintf("%s %s", num, severity))
		}
	}
	return text
}

func filterSeverity(severity string, threshold string) bool {
	thresholdValue, _ := severityToProperTypes(threshold)
	severityValue, _ := severityToProperTypes(severity)
	return severityValue > thresholdValue
}

func getContainerRegistries() ([]string, error) {
	var (
		registries            = make([]string, 0)
		regsIntegrations, err = cli.LwApi.Integrations.ListContainerRegistryIntegrations()
	)
	if err != nil {
		return registries, errors.Wrap(err, "unable to get container registry integrations")
	}

	for _, i := range regsIntegrations.Data {
		// avoid adding empty registries coming from the new local_scanner and avoid adding duplicate registries
		if i.Data.RegistryDomain == "" || array.ContainsStr(registries, i.Data.RegistryDomain) {
			continue
		}

		registries = append(registries, i.Data.RegistryDomain)
	}

	return registries, nil
}

// Creates a user-friendly error message
func userFriendlyErrorForOnDemandCtrVulnScan(err error, registry, repo, tag string) error {
	if strings.Contains(err.Error(),
		"Could not find integration matching the registry provided",
	) || strings.Contains(err.Error(),
		"Could not find vulnerability integrations",
	) {

		registries, errReg := getContainerRegistries()
		if errReg != nil {
			cli.Log.Debugw("error trying to retrieve configured registries", "error", errReg)
			return errors.Errorf("container registry '%s' not found", registry)
		}

		if len(registries) == 0 {
			msg := `there are no container registries configured in your account.

Get started by integrating your container registry using the command:

    $ lacework integration create

If you prefer to configure the integration via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Integrations > Container Registry.
`
			return errors.New(fmt.Sprintf(msg, cli.Account))
		}

		msg := `container registry '%s' not found

Your account has the following container registries configured:

    > %s

To integrate a new container registry use the command:

    $ lacework integration create
`
		return errors.New(fmt.Sprintf(msg, registry, strings.Join(registries, "\n    > ")))
	}

	if strings.Contains(
		err.Error(),
		"Could not successfully send scan request to available integrations for given repo and label",
	) {

		msg := `container image '%s:%s' not found in registry '%s'.

This error is likely due to a problem with the container registry integration 
configured in your account. Verify that the integration was configured with 
Lacework using the correct permissions, and that the repository belongs
to the provided registry.

To view all container registries configured in your account use the command:

    $ lacework vulnerability container list-registries
`
		return errors.Errorf(msg, repo, tag, registry)
	}

	return errors.Wrap(err, "unable to request on-demand vulnerability scan")
}
