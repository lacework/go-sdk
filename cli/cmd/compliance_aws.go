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
	// complianceAwsListAccountsCmd represents the list-accounts inside the aws command
	complianceAwsListAccountsCmd = &cobra.Command{
		Use:     "list-accounts",
		Aliases: []string{"list", "ls"},
		Short:   "List all AWS accounts configured",
		Long:    `List all AWS accounts configured in your account.`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cli.StartProgress("Fetching list of configured AWS accounts...")
			awsIntegrations, err := cli.LwApi.Integrations.ListAwsCfg()
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get aws compliance integrations")
			}

			return cliListAwsAccounts(&awsIntegrations)
		},
	}

	// complianceAwsGetReportCmd represents the get-report sub-command inside the aws command
	complianceAwsGetReportCmd = &cobra.Command{
		Use:     "get-report <account_id> [recommendation_id]",
		Aliases: []string{"get"},
		PreRunE: func(_ *cobra.Command, args []string) error {
			if compCmdState.Csv {
				cli.EnableCSVOutput()
			}

			if len(args) > 1 {
				compCmdState.RecommendationID = args[1]
				if !validRecommendationID(compCmdState.RecommendationID) {
					return errors.Errorf("\n'%s' is not a valid recommendation id\n", compCmdState.RecommendationID)
				}
			}

			switch compCmdState.Type {
			case "CIS":
				compCmdState.Type = fmt.Sprintf("AWS_%s_S3", compCmdState.Type)
				return nil
			case "SOC_Rev2":
				compCmdState.Type = fmt.Sprintf("AWS_%s", compCmdState.Type)
				return nil
			case "AWS_CIS_S3", "NIST_800-53_Rev4", "NIST_800-171_Rev2", "ISO_2700", "HIPAA", "SOC", "AWS_SOC_Rev2", "PCI":
				return nil
			default:
				return errors.New("supported report types are: CIS, NIST_800-53_Rev4, NIST_800-171_Rev2, ISO_2700, HIPAA, SOC, SOC_Rev2, or PCI")
			}
		},
		Short: "Get the latest AWS compliance report",
		Long: `Get the latest compliance assessment report from the provided AWS account, these
reports run on a regular schedule, typically once a day. The available report formats
are human-readable (default), json and pdf.

To list all AWS accounts configured in your account:

    lacework compliance aws list-accounts

To run an ad-hoc compliance assessment of an AWS account:

    lacework compliance aws run-assessment <account_id>

To show recommendation details and affected resources for a recommendation id:

    lacework compliance aws get-report <account_id> [recommendation_id]
`,
		Args: cobra.RangeArgs(1, 2),
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

				cli.StartProgress("Downloading compliance report...")
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

			var (
				report   api.ComplianceAwsReport
				cacheKey = fmt.Sprintf("compliance/aws/%s/%s", config.AccountID, config.Type)
			)
			expired := cli.ReadCachedAsset(cacheKey, &report)
			if expired {
				cli.StartProgress("Getting compliance report...")
				response, err := cli.LwApi.Compliance.GetAwsReport(config)
				cli.StopProgress()
				if err != nil {
					return errors.Wrap(err, "unable to get aws compliance report")
				}

				if len(response.Data) == 0 {
					return errors.New("no data found in the report")
				}

				report = response.Data[0]

				cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), report)
			}

			filteredOutput := ""

			if complianceFiltersEnabled() {
				report.Recommendations, filteredOutput = filterRecommendations(report.Recommendations)
			}

			if cli.JSONOutput() && compCmdState.RecommendationID == "" {
				return cli.OutputJSON(report)
			}

			if cli.CSVOutput() {
				recommendations := complianceCSVReportRecommendationsTable(
					&complianceCSVReportDetails{
						AccountName:     report.AccountID,
						AccountID:       report.AccountID,
						ReportType:      report.ReportType,
						ReportTime:      report.ReportTime,
						Recommendations: report.Recommendations,
					},
				)

				return cli.OutputCSV(
					[]string{"Report_Type", "Report_Time", "Account",
						"Section", "ID", "Recommendation", "Status",
						"Severity", "Resource", "Region", "Reason"},
					recommendations,
				)
			}

			// If RecommendationID is provided, output resources matching that id
			if compCmdState.RecommendationID != "" {
				return outputResourcesByRecommendationID(report)
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
		Short:   "Run a new AWS compliance assessment",
		Long:    `Run a compliance assessment for the provided AWS account.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := cli.LwApi.Compliance.RunAwsReport(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to run aws compliance assessment")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			cli.OutputHuman("A new AWS compliance assessment has been initiated.\n")
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

	// complianceAwsDisableReportCmd represents the disable-report sub-command inside the aws command
	// experimental feature
	complianceAwsDisableReportCmd = &cobra.Command{
		Use:     "disable-report <report_type>",
		Aliases: []string{"disable"},
		Hidden:  true,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "CIS":
				args[0] = fmt.Sprintf("AWS_%s_S3", args[0])
				return nil
			case "AWS_CIS_S3":
				return nil
			default:
				return errors.New("CIS is the only supported report type")
			}
		},
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {

			schema, err := fetchCachedAwsComplianceReportSchema(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to get aws compliance report schema")
			}

			// set state of all recommendations in this report to disabled
			patchReq := api.NewRecommendationV1State(schema, false)
			response, err := cli.LwApi.Recommendations.Aws.Patch(patchReq)
			if err != nil {
				return errors.Wrap(err, "unable to patch aws recommendations")
			}

			var cacheKey = fmt.Sprintf("compliance/aws/schema/%s", args[0])
			cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), response.RecommendationList())
			cli.OutputHuman(fmt.Sprintf("All recommendations for report %s have been disabled\n", args[0]))
			return nil
		},
	}

	// complianceAwsEnableReportCmd represents the enable-report sub-command inside the aws command
	// experimental feature
	complianceAwsEnableReportCmd = &cobra.Command{
		Use:     "enable-report <report_type>",
		Aliases: []string{"enable"},
		Hidden:  true,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "CIS":
				args[0] = fmt.Sprintf("AWS_%s_S3", args[0])
				return nil
			case "AWS_CIS_S3":
				return nil
			default:
				return errors.New("CIS is the only supported report type")
			}
		},
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {

			schema, err := fetchCachedAwsComplianceReportSchema(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to get aws compliance report schema")
			}

			// set state of all recommendations in this report to enabled
			patchReq := api.NewRecommendationV1State(schema, true)
			response, err := cli.LwApi.Recommendations.Aws.Patch(patchReq)
			if err != nil {
				return errors.Wrap(err, "unable to patch aws recommendations")
			}

			var cacheKey = fmt.Sprintf("compliance/aws/schema/%s", args[0])
			cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), response.RecommendationList())
			cli.OutputHuman(fmt.Sprintf("All recommendations for report %s have been enabled\n", args[0]))
			return nil
		},
	}

	// complianceAwsReportStatusCmd represents the report-status sub-command inside the aws command
	// experimental feature
	complianceAwsReportStatusCmd = &cobra.Command{
		Use:     "report-status <report_type>",
		Aliases: []string{"status"},
		Hidden:  true,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "CIS":
				args[0] = fmt.Sprintf("AWS_%s_S3", args[0])
				return nil
			case "AWS_CIS_S3":
				return nil
			default:
				return errors.New("CIS is the only supported report type")
			}
		},
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var rows [][]string
			report, err := fetchCachedAwsComplianceReportSchema(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to get Aws compliance report schema")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(api.NewRecommendationV1(report))
			}

			for _, r := range report {
				rows = append(rows, []string{r.ID, strconv.FormatBool(r.State)})
			}

			cli.OutputHuman(renderOneLineCustomTable(args[0],
				renderCustomTable([]string{}, rows,
					tableFunc(func(t *tablewriter.Table) {
						t.SetBorder(false)
						t.SetColumnSeparator(" ")
						t.SetAutoWrapText(false)
						t.SetAlignment(tablewriter.ALIGN_LEFT)
					}),
				),
				tableFunc(func(t *tablewriter.Table) {
					t.SetBorder(false)
					t.SetAutoWrapText(false)
				}),
			))
			return nil
		},
	}
)

func init() {
	// add sub-commands to the aws command
	complianceAwsCmd.AddCommand(complianceAwsGetReportCmd)
	complianceAwsCmd.AddCommand(complianceAwsListAccountsCmd)
	complianceAwsCmd.AddCommand(complianceAwsRunAssessmentCmd)

	// Experimental Commands
	complianceAwsCmd.AddCommand(complianceAwsReportStatusCmd)
	complianceAwsCmd.AddCommand(complianceAwsDisableReportCmd)
	complianceAwsCmd.AddCommand(complianceAwsEnableReportCmd)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Details, "details", false,
		"increase details about the compliance report",
	)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Pdf, "pdf", false,
		"download report in PDF format",
	)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Csv, "csv", false,
		"output report in CSV format",
	)

	// AWS report types: AWS_CIS_S3, NIST_800-53_Rev4, NIST_800-171_Rev2, ISO_2700, HIPAA, SOC, AWS_SOC_Rev2, or PCI
	complianceAwsGetReportCmd.Flags().StringVar(&compCmdState.Type, "type", "CIS",
		"report type to display, supported types: CIS, NIST_800-53_Rev4, NIST_800-171_Rev2, ISO_2700, HIPAA, SOC, SOC_Rev2, or PCI",
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

type awsAccount struct {
	AccountID string `json:"account_id"`
	Status    string `json:"status"`
}

func cliListAwsAccounts(awsIntegrations *api.AwsIntegrationsResponse) error {
	awsAccounts := make([]awsAccount, 0)
	jsonOut := struct {
		Accounts []awsAccount `json:"aws_accounts"`
	}{Accounts: awsAccounts}

	if awsIntegrations == nil || len(awsIntegrations.Data) == 0 {
		if cli.JSONOutput() {
			return cli.OutputJSON(jsonOut)
		}

		msg := `There are no AWS accounts configured in your account.

Get started by integrating your AWS accounts to analyze configuration compliance using the command:

    lacework integration create

If you prefer to configure the integration via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Integrations > Cloud Accounts.
`
		cli.OutputHuman(msg, cli.Account)
		return nil
	}

	for _, i := range awsIntegrations.Data {
		if containsDuplicateAccountID(awsAccounts, i.Data.AwsAccountID) {
			cli.Log.Warnw("duplicate aws account", "integration_guid", i.IntgGuid, "account", i.Data.AwsAccountID)
			continue
		}
		awsAccounts = append(awsAccounts, awsAccount{
			AccountID: i.Data.AwsAccountID,
			Status:    i.Status(),
		})
	}

	if cli.JSONOutput() {
		jsonOut.Accounts = awsAccounts
		return cli.OutputJSON(jsonOut)
	}

	var rows [][]string
	for _, acc := range awsAccounts {
		rows = append(rows, []string{acc.AccountID, acc.Status})
	}

	cli.OutputHuman(renderSimpleTable([]string{"AWS Account", "Status"}, rows))
	return nil
}

func containsDuplicateAccountID(awsAccount []awsAccount, accountID string) bool {
	for _, value := range awsAccount {
		if accountID == value.AccountID {
			return true
		}
	}
	return false
}

func fetchCachedAwsComplianceReportSchema(reportType string) (response []api.RecommendationV1, err error) {
	var cacheKey = fmt.Sprintf("compliance/aws/schema/%s", reportType)
	expired := cli.ReadCachedAsset(cacheKey, &response)
	if expired {
		cli.StartProgress("Fetching compliance report schema...")
		response, err = cli.LwApi.Recommendations.Aws.GetReport(reportType)
		cli.StopProgress()
		if err != nil {
			return nil, errors.Wrap(err, "unable to get aws compliance report schema")
		}
		if len(response) == 0 {
			return nil, errors.New("no data found in the report")
		}

		// write previous state to cache, allowing for revert to previous state
		cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), response)
	}
	return
}
