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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/lacework/go-sdk/lwseverity"
)

var (
	compCmdState = struct {
		// download report in PDF format
		Pdf bool

		// output report in CSV format
		Csv bool

		// display extended details about a compliance report
		Details bool

		// Filter the recommendations table by category
		Category []string

		// Filter the recommendations table by service
		Service []string

		// Filter the recommendations table by severity
		Severity string

		// Filter the recommendations table by status
		Status string

		// output resources affected by recommendationID
		RecommendationID string
	}{}

	// complianceCmd represents the compliance command
	complianceCmd = &cobra.Command{
		Use:     "compliance",
		Aliases: []string{"comp"},
		Short:   "Manage compliance reports",
		Long: `Manage compliance reports for Google, Azure, or AWS cloud providers.

Lacework cloud security platform provides continuous Compliance monitoring against
cloud security best practices and compliance standards as CIS, PCI DSS, SoC II and
HIPAA benchmark standards.

Get started by integrating one or more cloud accounts using the command:

    lacework cloud-account create

If you prefer to configure the integration via the WebUI, log in to your account at:

    https://<ACCOUNT>.lacework.net

Then navigate to Settings > Integrations > Cloud Accounts.

Use the following command to list all available integrations in your account:

    lacework cloud-account list
`,
	}

	// complianceAzureCmd represents the azure sub-command inside the compliance command
	complianceAzureCmd = &cobra.Command{
		Use:     "azure",
		Aliases: []string{"az"},
		Short:   "Compliance for Azure Cloud",
		Long: `Manage compliance reports for Azure Cloud.

To list all Azure tenants configured in your account:

    lacework compliance azure list-tenants

To list all Azure subscriptions from a tenant, use the command:

    lacework compliance azure list-subscriptions <tenant_id>

To get the latest Azure compliance assessment report, use the command:

    lacework compliance azure get-report <tenant_id> <subscription_id>

These reports run on a regular schedule, typically once a day.
`,
	}

	// complianceGcpCmd represents the gcp sub-command inside the compliance command
	complianceGcpCmd = &cobra.Command{
		Use:     "google",
		Aliases: []string{"gcp"},
		Short:   "Compliance for Google Cloud",
		Long: `Manage compliance reports for Google Cloud.

To list all GCP organizations and projects configured in your account:

    lacework compliance gcp list

To list all GCP projects from an organization, use the command:

    lacework compliance gcp list-projects <organization_id>

To get the latest GCP compliance assessment report, use the command:

    lacework compliance gcp get-report <organization_id> <project_id>

These reports run on a regular schedule, typically once a day.
`,
	}

	// complianceAwsCmd represents the aws sub-command inside the compliance command
	complianceAwsCmd = &cobra.Command{
		Use:   "aws",
		Short: "Compliance for AWS",
		Long: `Manage compliance reports for Amazon Web Services (AWS).

To list all AWS accounts configured in your account:

    lacework compliance aws list-accounts

To get the latest AWS compliance assessment report:

    lacework compliance aws get-report <account_id>

These reports run on a regular schedule, typically once a day.
`,
	}
)

func init() {
	// add the compliance command
	rootCmd.AddCommand(complianceCmd)

	// add sub-commands to the compliance command
	complianceCmd.AddCommand(complianceAzureCmd)
	complianceCmd.AddCommand(complianceAwsCmd)
	complianceCmd.AddCommand(complianceGcpCmd)
}

func complianceReportSummaryTable(summaries []api.ReportSummary) [][]string {
	if len(summaries) == 0 {
		return [][]string{}
	}
	summary := summaries[0]
	return [][]string{
		{"Critical", fmt.Sprint(summary.NumSeverity1NonCompliance)},
		{"High", fmt.Sprint(summary.NumSeverity2NonCompliance)},
		{"Medium", fmt.Sprint(summary.NumSeverity3NonCompliance)},
		{"Low", fmt.Sprint(summary.NumSeverity4NonCompliance)},
		{"Info", fmt.Sprint(summary.NumSeverity5NonCompliance)},
	}
}

func complianceReportRecommendationsTable(recommendations []api.RecommendationV2) [][]string {
	out := [][]string{}
	for _, recommend := range recommendations {
		out = append(out, []string{
			recommend.RecID,
			recommend.Title,
			recommend.Status,
			recommend.SeverityString(),
			recommend.Service,
			fmt.Sprint(len(recommend.Violations)),
			fmt.Sprint(recommend.AssessedResourceCount),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return api.SeverityOrder(out[i][3]) < api.SeverityOrder(out[j][3])
	})

	return out
}

type complianceCSVReportDetails struct {
	// For clouds with tenant models, supply tenant ID
	TenantID string

	// For clouds with tenant models, supply tenant name/alias
	TenantName string

	// Supply the account id for the cloud enviornment
	AccountID string

	// Supply the account name/alias for the cloud enviornment, if available
	AccountName string

	// The type of report being rendered
	ReportType string

	// The time of the report execution
	ReportTime time.Time

	// Recommendations
	Recommendations []api.RecommendationV2
}

func (c complianceCSVReportDetails) GetAccountDetails() []string {
	accountAlias := c.AccountID
	if c.AccountName != "" {
		accountAlias = fmt.Sprintf("%s(%s)", c.AccountName, c.AccountID)
	}

	tenantAlias := c.TenantID
	if c.TenantName != "" {
		tenantAlias = fmt.Sprintf("%s(%s)", c.TenantName, c.TenantID)
	}
	out := []string{}
	if tenantAlias != "" {
		out = append(out, tenantAlias)
	}

	if accountAlias != "" {
		out = append(out, accountAlias)
	}
	return out
}

func (c complianceCSVReportDetails) GetReportMetaData() []string {
	return append([]string{c.ReportType, c.ReportTime.Format(time.RFC3339)}, c.GetAccountDetails()...)
}

func (c complianceCSVReportDetails) SortRecommendations() {
	sort.Slice(c.Recommendations, func(i, j int) bool {
		return c.Recommendations[i].Category < c.Recommendations[j].Category
	})

}

func complianceCSVReportRecommendationsTable(details *complianceCSVReportDetails) [][]string {
	details.SortRecommendations()
	out := [][]string{}

	for _, recommendation := range details.Recommendations {
		// GROW-1266: Do not add if status flag filters suppressed
		if compCmdState.Status == "" || compCmdState.Status == "suppressed" {
			for _, suppression := range recommendation.Suppressions {
				out = append(out,
					append(details.GetReportMetaData(),
						recommendation.Category,
						recommendation.RecID,
						recommendation.Title,
						"Suppressed",
						recommendation.SeverityString(),
						suppression,
						"",
						""))
			}
		}
		for _, violation := range recommendation.Violations {
			out = append(out,
				append(details.GetReportMetaData(),
					recommendation.Category,
					recommendation.RecID,
					recommendation.Title,
					recommendation.Status,
					recommendation.SeverityString(),
					violation.Resource,
					violation.Region,
					strings.Join(violation.Reasons, ",")))
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return api.SeverityOrder(out[i][3]) < api.SeverityOrder(out[j][3])
	})

	return out
}

func buildComplianceReportTable(detailsTable, summaryTable, recommendationsTable [][]string, filteredOutput string) string {
	mainReport := &strings.Builder{}
	mainReport.WriteString(
		renderCustomTable(
			[]string{
				"Compliance Report Details",
				"Non-Compliant Recommendations",
			},
			[][]string{[]string{
				renderCustomTable([]string{}, detailsTable,
					tableFunc(func(t *tablewriter.Table) {
						t.SetBorder(false)
						t.SetColumnSeparator("")
						t.SetAlignment(tablewriter.ALIGN_LEFT)
					}),
				),
				renderCustomTable([]string{"Severity", "Count"}, summaryTable,
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

	if compCmdState.Details || complianceFiltersEnabled() {
		mainReport.WriteString(
			renderCustomTable(
				[]string{"ID", "Recommendation", "Status", "Severity",
					"Service", "Affected", "Assessed"},
				recommendationsTable,
				tableFunc(func(t *tablewriter.Table) {
					t.SetBorder(false)
					t.SetRowLine(true)
					t.SetColumnSeparator(" ")
				}),
			),
		)
		if filteredOutput != "" {
			mainReport.WriteString(filteredOutput)
		}
		mainReport.WriteString("\n")

		if compCmdState.Status == "" {
			mainReport.WriteString(
				"Try adding '--status non-compliant' to show only non-compliant recommendations.",
			)
		} else if compCmdState.Severity == "" {
			mainReport.WriteString(
				"Try adding '--severity high' to show only high and critical recommendations.",
			)
		} else {
			mainReport.WriteString(
				"Try adding [recommendation_id] to show affected resources.",
			)
		}
		mainReport.WriteString("\n")
	} else {
		mainReport.WriteString(
			"Try adding '--details' to increase details shown about the compliance report.\n",
		)
	}
	return mainReport.String()
}

func filterRecommendations(recommendations []api.RecommendationV2) ([]api.RecommendationV2, string) {
	var filtered []api.RecommendationV2
	for _, r := range recommendations {
		if matchRecommendationsFilters(r) {
			filtered = append(filtered, r)
		}
	}
	if len(filtered) == 0 {
		return filtered, "There are no recommendations with the specified filter(s).\n"
	}

	cli.Log.Debugw("filtered recommendations", "recommendations", filtered)
	return filtered, fmt.Sprintf("%v of %v recommendations showing \n", len(filtered), len(recommendations))
}

func matchRecommendationsFilters(r api.RecommendationV2) bool {
	var results []bool

	// severity returns specified threshold and above
	if compCmdState.Severity != "" {
		sevThreshold, _ := lwseverity.Normalize(compCmdState.Severity)
		results = append(results, r.Severity <= sevThreshold)
	}

	if len(compCmdState.Category) > 0 {
		var categories []string
		for _, c := range compCmdState.Category {
			categories = append(categories, strings.ReplaceAll(c, "-", " "))
		}
		results = append(results, array.ContainsStrCaseInsensitive(categories, r.Category))
	}

	if len(compCmdState.Service) > 0 {
		results = append(results, array.ContainsStrCaseInsensitive(compCmdState.Service, r.Service))
	}

	if compCmdState.Status != "" {
		results = append(results, r.Status == statusToProperTypes(compCmdState.Status))
	}

	return !array.ContainsBool(results, false)
}

func complianceFiltersEnabled() bool {
	return len(compCmdState.Category) > 0 || compCmdState.Status != "" || compCmdState.Severity != "" || len(compCmdState.Service) > 0
}

func statusToProperTypes(status string) string {
	switch strings.ToLower(status) {
	case "non-compliant", "noncompliant":
		return "NonCompliant"
	case "compliant":
		return "Compliant"
	case "could-not-assess", "couldnotassess":
		return "CouldNotAssess"
	case "suppressed":
		return "Suppressed"
	case "requires-manual-assessment", "requiresmanualassessment":
		return "RequiresManualAssessment"
	default:
		return "Unknown"
	}
}

func outputResourcesByRecommendationID(report api.CloudComplianceReportV2) error {
	recommendation, found := report.GetComplianceRecommendation(compCmdState.RecommendationID)
	if !found || recommendation == nil {
		return errors.Errorf("recommendation id '%s' not found.", compCmdState.RecommendationID)
	}
	violations := recommendation.Violations
	affectedResources := len(recommendation.Violations)

	if cli.JSONOutput() {
		return cli.OutputJSON(recommendation)
	}

	cli.OutputHuman(
		renderOneLineCustomTable("RECOMMENDATION DETAILS",
			renderCustomTable([]string{},
				[][]string{
					{"ID", compCmdState.RecommendationID},
					{"SEVERITY", recommendation.SeverityString()},
					{"SERVICE", recommendation.Service},
					{"CATEGORY", recommendation.Category},
					{"STATUS", recommendation.Status},
					{"ASSESSED RESOURCES", strconv.Itoa(recommendation.AssessedResourceCount)},
					{"AFFECTED RESOURCES", strconv.Itoa(affectedResources)},
				},
				tableFunc(func(t *tablewriter.Table) {
					t.SetBorder(false)
					t.SetColumnSeparator(" ")
					t.SetAutoWrapText(false)
					t.SetAlignment(tablewriter.ALIGN_LEFT)
				}),
			), tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetAutoWrapText(false)
			}),
		))

	if affectedResources == 0 {
		cli.OutputHuman("\nNo resources found affected by '%s'\n", compCmdState.RecommendationID)
		return nil
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"AFFECTED RESOURCE", "REGION", "REASON"},
			violationsToTable(violations),
		),
	)
	return nil
}

func violationsToTable(violations []api.ComplianceViolationV2) (resourceTable [][]string) {
	for _, v := range violations {
		resourceTable = append(resourceTable, []string{v.Resource, v.Region, strings.Join(v.Reasons, ",")})
	}
	return
}

// nolint
func getReportTypes(reportSubType string) (validTypes []string, err error) {
	cacheKey := fmt.Sprintf("reports/definitions/%s", reportSubType)
	expired := cli.ReadCachedAsset(cacheKey, &validTypes)

	if expired {
		cli.StartProgress("fetching valid report types...")
		reportDefinitions, err := cli.LwApi.V2.ReportDefinitions.List()
		cli.StopProgress()

		if err != nil {
			return nil, err
		}

		for _, report := range reportDefinitions.Data {
			if report.SubReportType == reportSubType {
				validTypes = append(validTypes, report.ReportNotificationType)
			}
		}
		cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), validTypes)
	}

	return validTypes, err
}

func prettyPrintReportTypes(reportTypes []string) string {
	var sb strings.Builder
	for i, r := range reportTypes {
		if i%5 == 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("'%s',", r))
	}
	return sb.String()
}

func validReportName(cloud string, name string) error {
	var validReportNames []string
	definitions, err := cli.LwApi.V2.ReportDefinitions.List()
	if err != nil {
		return errors.Wrap(err, "unable to list report definitions")
	}

	for _, d := range definitions.Data {
		if d.SubReportType == cloud {
			validReportNames = append(validReportNames, d.ReportName)
		}
	}

	if array.ContainsStr(validReportNames, name) {
		return nil
	} else {
		return errors.Errorf("'%s' is not a valid report name.\nRun 'lacework report-definition list --subtype %s' for a list of valid report names", name, cloud)
	}
}
