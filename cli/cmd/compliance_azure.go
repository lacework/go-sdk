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
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// complianceAzureListSubsCmd represents the list-subscriptions sub-command inside the azure command
	complianceAzureListSubsCmd = &cobra.Command{
		Use:     "list-subscriptions <tenant_id>",
		Aliases: []string{"list-subs", "list"},
		Short:   "list subscriptions from tenant",
		Long: `List all Azure subscriptions from the provided tenant ID.

Use the following command to list all Azure integrations in your account:

  $ lacework integrations list --type AZURE_CFG

Then, select one GUID from an integration and visualize its details using the command:

  $ lacework integration show <int_guid>
`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := cli.LwApi.Compliance.ListAzureSubscriptions(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to list azure subscriptions")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response.Data[0])
			}

			cli.OutputHuman(buildAzureSubscriptionsTable(response.Data))
			return nil
		},
	}

	// complianceAzureGetReportCmd represents the get-report sub-command inside the azure command
	complianceAzureGetReportCmd = &cobra.Command{
		Use:     "get-report <tenant_id> <subscriptions_id>",
		Aliases: []string{"get"},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			switch compCmdState.Type {
			case "CIS", "SOC", "PCI":
				compCmdState.Type = fmt.Sprintf("AZURE_%s", compCmdState.Type)
				return nil
			case "AZURE_CIS", "AZURE_SOC", "AZURE_PCI":
				return nil
			default:
				return errors.New("supported report types are: CIS, SOC, or PCI")
			}
		},
		Short: "get the latest Azure compliance report",
		Long: `Get the latest Azure compliance assessment report, these reports run on a regular schedule,
typically once a day. The available report formats are human-readable (default), json and pdf.

To run an ad-hoc compliance assessment use the command:

  $ lacework compliance azure run-assessment <tenant_id>
`,
		Args: cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			config := api.ComplianceAzureReportConfig{
				TenantID:       args[0],
				SubscriptionID: args[1],
				Type:           compCmdState.Type,
			}

			if compCmdState.Pdf || compCmdState.PdfName != "" {
				pdfName := fmt.Sprintf(
					"%s_Report_%s_%s_%s_%s.pdf",
					config.Type,
					config.TenantID,
					config.SubscriptionID,
					cli.Account, time.Now().Format("20060102150405"),
				)
				if compCmdState.PdfName != "" {
					cli.OutputHuman("(DEPRECATED) This flag has been replaced by '--pdf'\n\n")
					pdfName = compCmdState.PdfName
				}

				cli.StartProgress(" Downloading compliance report...")
				err := cli.LwApi.Compliance.DownloadAzureReportPDF(pdfName, config)
				cli.StopProgress()
				if err != nil {
					return errors.Wrap(err, "unable to get azure pdf compliance report")
				}

				cli.OutputHuman("The Azure compliance report was downloaded at '%s'\n", pdfName)
				return nil
			}

			cli.StartProgress(" Getting compliance report...")
			response, err := cli.LwApi.Compliance.GetAzureReport(config)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get azure compliance report")
			}

			if len(response.Data) == 0 {
				return errors.New("there is no data found in the report")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response.Data[0])
			}

			report := response.Data[0]
			cli.OutputHuman("\n")
			cli.OutputHuman(
				buildComplianceReportTable(
					complianceAzureReportDetailsTable(&report),
					complianceReportSummaryTable(report.Summary),
					complianceReportRecommendationsTable(report.Recommendations),
				),
			)
			return nil
		},
	}

	// complianceAzureRunAssessmentCmd represents the run-assessment sub-command inside the azure command
	complianceAzureRunAssessmentCmd = &cobra.Command{
		Use:     "run-assessment <tenant_id>",
		Aliases: []string{"run"},
		Short:   "run a new Azure compliance assessment",
		Long:    `Run a compliance assessment of the provided Azure tenant.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := cli.LwApi.Compliance.RunAzureReport(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to run azure compliance assessment")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			cli.OutputHuman("A new Azure compliance assessment has been initiated.\n")
			cli.OutputHuman("\n")
			cli.OutputHuman(buildAzureRunAssessmentTable(response.IntgGuid, args[0]))
			return nil
		},
	}
)

func init() {
	// add sub-commands to the azure command
	complianceAzureCmd.AddCommand(complianceAzureListSubsCmd)
	complianceAzureCmd.AddCommand(complianceAzureGetReportCmd)
	complianceAzureCmd.AddCommand(complianceAzureRunAssessmentCmd)

	complianceAzureGetReportCmd.Flags().BoolVar(&compCmdState.Details, "details", false,
		"increase details about the compliance report",
	)
	complianceAzureGetReportCmd.Flags().StringVar(&compCmdState.PdfName, "pdf-file", "",
		"(DEPRECATED) use --pdf",
	)
	complianceAzureGetReportCmd.Flags().BoolVar(&compCmdState.Pdf, "pdf", false,
		"download report in PDF format",
	)

	// Azure report types: AZURE_CIS, AZURE_SOC, or AZURE_PCI
	complianceAzureGetReportCmd.Flags().StringVar(&compCmdState.Type, "type", "CIS",
		"report type to display, supported types: CIS, SOC, or PCI",
	)
}

func buildAzureRunAssessmentTable(intGuid, id string) string {
	var (
		tBuilder = &strings.Builder{}
		t        = tablewriter.NewWriter(tBuilder)
	)

	t.SetHeader([]string{"INTEGRATION GUID", "TENANT ID"})
	t.SetBorder(false)
	t.SetAutoWrapText(false)
	t.Append([]string{intGuid, id})
	t.Render()

	return tBuilder.String()
}

func buildAzureSubscriptionsTable(azureSubs []api.CompAzureSubscriptions) string {
	var (
		tBuilder = &strings.Builder{}
		t        = tablewriter.NewWriter(tBuilder)
	)

	t.SetHeader([]string{"Subscriptions"})
	t.SetBorder(false)
	t.SetAutoWrapText(false)
	for _, azure := range azureSubs {
		for _, subs := range azure.Subscriptions {
			t.Append([]string{subs})
		}
	}
	t.Render()

	return tBuilder.String()
}

func complianceAzureReportDetailsTable(report *api.ComplianceAzureReport) [][]string {
	return [][]string{
		[]string{"Report Type", report.ReportType},
		[]string{"Report Title", report.ReportTitle},
		[]string{"Tenant ID", report.TenantID},
		[]string{"Tenant Name", report.TenantName},
		[]string{"Subscription ID", report.SubscriptionID},
		[]string{"Subscription Name", report.SubscriptionName},
		[]string{"Report Time", report.ReportTime.UTC().Format(time.RFC3339)},
	}
}
