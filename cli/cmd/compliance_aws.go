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
	// complianceAwsGetReportCmd represents the get-report sub-command inside the aws command
	complianceAwsGetReportCmd = &cobra.Command{
		Use:     "get-report <account_id>",
		Aliases: []string{"get"},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			switch compCmdState.Type {
			case "CIS":
				compCmdState.Type = fmt.Sprintf("AWS_%s_S3", compCmdState.Type)
				return nil
			case "AWS_CIS_S3", "NIST_800-53_Rev4", "ISO_2700", "HIPAA", "SOC", "PCI":
				return nil
			default:
				return errors.New("supported report types are: CIS, NIST_800-53_Rev4, ISO_2700, HIPAA, SOC, or PCI")
			}
		},
		Short: "get the latest AWS compliance report",
		Long: `Get the latest AWS compliance assessment report, these reports run on a regular schedule,
typically once a day. The available report formats are human-readable (default), json and pdf.

To find out which AWS accounts are connected to you Lacework account, use the following command:

  $ lacework integrations list --type AWS_CFG

Then, choose one integration, copy the GUID and visialize its details using the command:

  $ lacework integration show <int_guid>

To run an ad-hoc compliance assessment use the command:

  $ lacework compliance aws run-assessment <account_id>
`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := api.ComplianceAwsReportConfig{
				AccountID: args[0],
				Type:      compCmdState.Type,
			}

			if compCmdState.Pdf || compCmdState.PdfName != "" {
				pdfName := fmt.Sprintf(
					"%s_Report_%s_%s_%s.pdf",
					config.Type,
					config.AccountID,
					cli.Account, time.Now().Format("20060102150405"),
				)
				if compCmdState.PdfName != "" {
					cli.OutputHuman("(DEPRECATED) This flag has been replaced by '--pdf'\n\n")
					pdfName = compCmdState.PdfName
				}

				cli.StartProgress(" Downloading compliance report...")
				err := cli.LwApi.Compliance.DownloadAwsReportPDF(pdfName, config)
				cli.StopProgress()
				if err != nil {
					return errors.Wrap(err, "unable to get aws pdf compliance report")
				}

				cli.OutputHuman("The AWS compliance report was downloaded at '%s'\n", pdfName)
				return nil
			}

			cli.StartProgress(" Getting compliance report...")
			response, err := cli.LwApi.Compliance.GetAwsReport(config)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get aws compliance report")
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
					complianceAwsReportDetailsTable(&report),
					complianceReportSummaryTable(report.Summary),
					complianceReportRecommendationsTable(report.Recommendations),
				),
			)
			return nil
		},
	}

	// complianceAwsRunAssessmentCmd represents the run-assessment sub-command inside the aws command
	complianceAwsRunAssessmentCmd = &cobra.Command{
		Use:     "run-assessment <account_id>",
		Aliases: []string{"run"},
		Short:   "run a new AWS compliance report",
		Long:    `Run a compliance assessment for the provided AWS account.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := cli.LwApi.Compliance.RunAwsReport(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to run aws compliance report")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			cli.OutputHuman("A new AWS compliance report has been initiated.\n")
			// @afiune not consistent with the other cloud providers
			for key := range response {
				cli.OutputHuman("\n")
				cli.OutputHuman(buildAwsRunAssessmentTable(key, args[0]))
			}
			return nil
		},
	}
)

func init() {
	// add sub-commands to the aws command
	complianceAwsCmd.AddCommand(complianceAwsGetReportCmd)
	complianceAwsCmd.AddCommand(complianceAwsRunAssessmentCmd)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Details, "details", false,
		"increase details about the compliance report",
	)
	complianceAwsGetReportCmd.Flags().StringVar(&compCmdState.PdfName, "pdf-file", "",
		"(DEPRECATED) use --pdf",
	)
	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Pdf, "pdf", false,
		"download report in PDF format",
	)

	// AWS report types: AWS_CIS_S3, NIST_800-53_Rev4, ISO_2700, HIPAA, SOC, or PCI
	complianceAwsGetReportCmd.Flags().StringVar(&compCmdState.Type, "type", "CIS",
		"report type to display, supported types: CIS, NIST_800-53_Rev4, ISO_2700, HIPAA, SOC, or PCI",
	)
}

func buildAwsRunAssessmentTable(intGuid, id string) string {
	var (
		tBuilder = &strings.Builder{}
		t        = tablewriter.NewWriter(tBuilder)
	)

	t.SetHeader([]string{"INTEGRATION GUID", "ACCOUNT ID"})
	t.SetBorder(false)
	t.SetAutoWrapText(false)
	t.Append([]string{intGuid, id})
	t.Render()

	return tBuilder.String()
}

func complianceAwsReportDetailsTable(report *api.ComplianceAwsReport) [][]string {
	return [][]string{
		[]string{"Report Type", report.ReportType},
		[]string{"Report Title", report.ReportTitle},
		[]string{"Account ID", report.AccountID},
		[]string{"Account Alias", report.AccountAlias},
		[]string{"Report Time", report.ReportTime.UTC().Format(time.RFC3339)},
	}
}
