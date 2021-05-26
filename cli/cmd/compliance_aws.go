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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
)

var (
	// complianceAwsListAccountsCmd represents the list-accounts inside the aws command
	complianceAwsListAccountsCmd = &cobra.Command{
		Use:     "list-accounts",
		Aliases: []string{"list"},
		Short:   "list all AWS accounts configured",
		Long:    `List all AWS accounts configured in your account.`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			awsIntegrations, err := cli.LwApi.Integrations.ListAwsCfg()
			if err != nil {
				return errors.Wrap(err, "unable to get aws compliance integrations")
			}
			if len(awsIntegrations.Data) == 0 {
				msg := `There are no AWS accounts configured in your account.

Get started by integrating your AWS accounts to analyze configuration compliance using the command:

    $ lacework integration create

If you prefer to configure the integration via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Integrations > Cloud Accounts.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}

			awsAccounts := make([]string, 0)
			for _, i := range awsIntegrations.Data {
				awsAccounts = append(awsAccounts, i.Data.AwsAccountID)
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(awsAccounts)
			}

			rows := [][]string{}
			for _, acc := range awsAccounts {
				rows = append(rows, []string{acc})
			}

			cli.OutputHuman(renderSimpleTable([]string{"AWS Accounts"}, rows))
			return nil
		},
	}

	// complianceAwsGetReportCmd represents the get-report sub-command inside the aws command
	complianceAwsGetReportCmd = &cobra.Command{
		Use:     "get-report <account_id>",
		Aliases: []string{"get"},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if compCmdState.Csv {
				cli.NonInteractive()
			}

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
		Long: `Get the latest compliance assessment report from the provided AWS account, these
reports run on a regular schedule, typically once a day. The available report formats
are human-readable (default), json and pdf.

To list all AWS accounts configured in your account:

    $ lacework compliance aws list-accounts

To run an ad-hoc compliance assessment of an AWS account:

    $ lacework compliance aws run-assessment <account_id>
`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var (
				// clean the AWS account ID if it was provided
				// with an Alias in between parentheses
				awsAccountID, _ = splitIDAndAlias(args[0])
				config          = api.ComplianceAwsReportConfig{
					AccountID: awsAccountID,
					Type:      compCmdState.Type,
				}
			)

			if compCmdState.Pdf {
				pdfName := fmt.Sprintf(
					"%s_Report_%s_%s_%s.pdf",
					config.Type,
					config.AccountID,
					cli.Account, time.Now().Format("20060102150405"),
				)

				cli.StartProgress(" Downloading compliance report...")
				err := cli.LwApi.Compliance.DownloadAwsReportPDF(pdfName, config)
				cli.StopProgress()
				if err != nil {
					return errors.Wrap(err, "unable to get aws pdf compliance report")
				}

				cli.OutputHuman("The AWS compliance report was downloaded at '%s'\n", pdfName)
				return nil
			}

			if compCmdState.Severity != "" {
				if !array.ContainsStr(api.ValidEventSeverities, compCmdState.Severity) {
					return errors.Errorf("the severity %s is not valid, use one of %s",
						compCmdState.Severity, strings.Join(api.ValidEventSeverities, ", "),
					)
				}
			}
			if compCmdState.Status != "" {
				if !array.ContainsStr(api.ValidComplianceStatus, compCmdState.Status) {
					return errors.Errorf("the status %s is not valid, use one of %s",
						compCmdState.Status, strings.Join(api.ValidComplianceStatus, ", "),
					)
				}
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

			report := response.Data[0]
			filteredOutput := ""

			if complianceFiltersEnabled() {
				report.Recommendations, filteredOutput = filterRecommendations(report.Recommendations)
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(report)
			}

			recommendations := complianceReportRecommendationsTable(report.Recommendations)
			cli.OutputHuman("\n")
			cli.OutputHuman(
				buildComplianceReportTable(
					complianceAwsReportDetailsTable(&report),
					complianceReportSummaryTable(report.Summary),
					recommendations,
					filteredOutput,
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
				cli.OutputHuman(
					renderSimpleTable(
						[]string{"INTEGRATION GUID", "ACCOUNT ID"},
						[][]string{[]string{key, args[0]}},
					),
				)
			}
			return nil
		},
	}
)

func init() {
	// add sub-commands to the aws command
	complianceAwsCmd.AddCommand(complianceAwsGetReportCmd)
	complianceAwsCmd.AddCommand(complianceAwsListAccountsCmd)
	complianceAwsCmd.AddCommand(complianceAwsRunAssessmentCmd)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Details, "details", false,
		"increase details about the compliance report",
	)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Pdf, "pdf", false,
		"download report in PDF format",
	)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Csv, "csv", false,
		"download report in CSV format",
	)

	// AWS report types: AWS_CIS_S3, NIST_800-53_Rev4, ISO_2700, HIPAA, SOC, or PCI
	complianceAwsGetReportCmd.Flags().StringVar(&compCmdState.Type, "type", "CIS",
		"report type to display, supported types: CIS, NIST_800-53_Rev4, ISO_2700, HIPAA, SOC, or PCI",
	)

	complianceAwsGetReportCmd.Flags().StringSliceVar(&compCmdState.Category, "category", []string{},
		"filter report details by category (identity-and-access-management, s3, logging...)",
	)

	complianceAwsGetReportCmd.Flags().StringSliceVar(&compCmdState.Service, "service", []string{},
		"filter report details by service (aws:s3, aws:iam, aws:cloudtrail, ...)",
	)

	complianceAwsGetReportCmd.Flags().StringVar(&compCmdState.Severity, "severity", "",
		fmt.Sprintf("filter report details by severity threshold (%s)",
			strings.Join(api.ValidEventSeverities, ", ")),
	)

	complianceAwsGetReportCmd.Flags().StringVar(&compCmdState.Status, "status", "",
		fmt.Sprintf("filter report details by status (%s)",
			strings.Join(api.ValidComplianceStatus, ", ")),
	)
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
