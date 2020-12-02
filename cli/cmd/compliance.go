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
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	compCmdState = struct {
		// the report type to display, supported are: CIS, SOC, or PCI
		// default: CIS
		Type string

		// download the report as PDF format with the provided filename
		// (DEPRECATED)
		PdfName string

		// download report in PDF format
		Pdf bool

		// display extended details about a compliance report
		Details bool
	}{Type: "CIS"}

	// complianceCmd represents the compliance command
	complianceCmd = &cobra.Command{
		Use:     "compliance",
		Aliases: []string{"comp"},
		Short:   "manage compliance reports",
		Long: `Manage compliance reports for Google, Azure, or AWS cloud providers.

Lacework cloud security platform provides continuous Compliance monitoring against
cloud security best practices and compliance standards as CIS, PCI DSS, SoC II and
HIPAA benchmark standards.

Get started by integrating one or more cloud accounts using the command:

    $ lacework integration create

Or, if you prefer to do it via the WebUI, log in to your account at:

    https://<ACCOUNT>.lacework.net

Then navigate to Settings > Integrations > Cloud Accounts.

Use the following command to list all available integrations in your account:

    $ lacework integrations list
`,
	}

	// complianceAzureCmd represents the azure sub-command inside the compliance command
	complianceAzureCmd = &cobra.Command{
		Use:     "azure",
		Aliases: []string{"az"},
		Short:   "compliance for Azure Cloud",
		Long: `Manage compliance reports for Azure Cloud.

To get the latest Azure compliance assessment report, use the command:

    $ lacework compliance azure get-report <tenant_id> <subscriptions_id>

These reports run on a regular schedule, typically once a day.

To find out which Azure tenants/subscriptions are connected to your
Lacework account, use the following command:

    $ lacework integrations list --type AZURE_CFG

Then, choose one integration, copy the GUID and visualize its details
using the command:

    $ lacework integration show <int_guid>

To list all Azure subscriptions from a tenant, use the command:

    $ lacework compliance azure list-subscriptions <tenant_id>

To run an ad-hoc compliance assessment use the command:

    $ lacework compliance azure run-assessment <tenant_id>
`,
	}

	// complianceGcpCmd represents the gcp sub-command inside the compliance command
	complianceGcpCmd = &cobra.Command{
		Use:     "google",
		Aliases: []string{"gcp"},
		Short:   "compliance for Google Cloud",
		Long: `Manage compliance reports for Google Cloud.

To get the latest GCP compliance assessment report, use the command:

    $ lacework compliance gcp get-report <organization_id> <project_id>

These reports run on a regular schedule, typically once a day.

To find out which GCP organizations/projects are connected to your
Lacework account, use the following command:

    $ lacework integrations list --type GCP_CFG

Then, choose one integration, copy the GUID and visualize its details
using the command:

    $ lacework integration show <int_guid>

To list all GCP projects from an organization, use the command:

    $ lacework compliance gcp list-projects <organization_id>

To run an ad-hoc compliance assessment use the command:

    $ lacework compliance gcp run-assessment <org_or_project_id>
`,
	}

	// complianceAwsCmd represents the aws sub-command inside the compliance command
	complianceAwsCmd = &cobra.Command{
		Use:   "aws",
		Short: "compliance for AWS",
		Long: `Manage compliance reports for Amazon Web Services (AWS).

To list all AWS accounts configured in your account:

    $ lacework compliance aws list-accounts

To get the latest AWS compliance assessment report:

    $ lacework compliance aws get-report <account_id>

These reports run on a regular schedule, typically once a day.

To run an ad-hoc compliance assessment:

    $ lacework compliance aws run-assessment <account_id>
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

func complianceReportSummaryTable(summaries []api.ComplianceSummary) [][]string {
	if len(summaries) == 0 {
		return [][]string{}
	}
	summary := summaries[0]
	return [][]string{
		[]string{"Critical", fmt.Sprint(summary.NumSeverity1NonCompliance)},
		[]string{"High", fmt.Sprint(summary.NumSeverity2NonCompliance)},
		[]string{"Medium", fmt.Sprint(summary.NumSeverity3NonCompliance)},
		[]string{"Low", fmt.Sprint(summary.NumSeverity4NonCompliance)},
		[]string{"Info", fmt.Sprint(summary.NumSeverity5NonCompliance)},
	}
}

func complianceReportRecommendationsTable(recommendations []api.ComplianceRecommendation) [][]string {
	out := [][]string{}
	for _, recommend := range recommendations {
		out = append(out, []string{
			recommend.RecID,
			recommend.Title,
			recommend.Status,
			recommend.SeverityString(),
			recommend.Service,
			fmt.Sprint(recommend.ResourceCount),
			fmt.Sprint(recommend.AssessedResourceCount),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][3]) < severityOrder(out[j][3])
	})

	return out
}

func buildComplianceReportRecommandations(recommendationsTable [][]string) string {
	var (
		detailsTable = &strings.Builder{}
		t            = tablewriter.NewWriter(detailsTable)
	)

	t.SetRowLine(true)
	t.SetBorders(tablewriter.Border{
		Left:   false,
		Right:  false,
		Top:    true,
		Bottom: true,
	})
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader([]string{
		"ID",
		"Recommendation",
		"Status",
		"Severity",
		"Service",
		"Affected",
		"Assessed",
	})
	t.AppendBulk(recommendationsTable)
	t.Render()

	return detailsTable.String()
}

func buildComplianceReportTable(detailsTable, summaryTable, recommendationsTable [][]string) string {
	var (
		t             *tablewriter.Table
		mainReport    = &strings.Builder{}
		summaryReport = &strings.Builder{}
		reportDetails = &strings.Builder{}
	)

	t = tablewriter.NewWriter(reportDetails)
	t.SetBorder(false)
	t.SetColumnSeparator("")
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.AppendBulk(detailsTable)
	t.Render()

	t = tablewriter.NewWriter(summaryReport)
	t.SetBorder(false)
	t.SetColumnSeparator(" ")
	t.SetHeader([]string{
		"Severity", "Count",
	})
	t.AppendBulk(summaryTable)
	t.Render()

	t = tablewriter.NewWriter(mainReport)
	t.SetBorder(false)
	t.SetAutoWrapText(false)
	t.SetHeader([]string{
		"Compliance Report Details",
		"Non-Compliant Recommendations",
	})
	t.Append([]string{
		reportDetails.String(),
		summaryReport.String(),
	})
	t.Render()

	if compCmdState.Details {
		mainReport.WriteString(buildComplianceReportRecommandations(recommendationsTable))
		mainReport.WriteString("\n")
		mainReport.WriteString(
			"Try using '--pdf' to download the report in PDF format.",
		)
		mainReport.WriteString("\n")
	} else {
		mainReport.WriteString(
			"Try using '--details' to increase details shown about the compliance report.\n",
		)
	}
	return mainReport.String()
}
