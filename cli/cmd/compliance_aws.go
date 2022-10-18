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

	"github.com/AlecAivazis/survey/v2"
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
		Aliases: []string{"get", "show"},
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
			case "AWS_CIS_S3", "NIST_800-53_Rev4", "NIST_800-171_Rev2", "ISO_2700", "HIPAA", "SOC", "AWS_SOC_Rev2",
				"PCI", "AWS_CIS_14", "AWS_CMMC_1.02", "AWS_HIPAA", "AWS_ISO_27001:2013", "AWS_NIST_CSF", "AWS_NIST_800-171_rev2",
				"AWS_NIST_800-53_rev5", "AWS_PCI_DSS_3.2.1", "AWS_SOC_2", "LW_AWS_SEC_ADD_1_0":
				return nil
			default:
				return errors.New(`supported report types are: AWS_CIS_S3', 'NIST_800-53_Rev4', 'NIST_800-171_Rev2', 
'ISO_2700', 'HIPAA', 'SOC', 'AWS_SOC_Rev2', 'PCI', 'AWS_CIS_14', 'AWS_CMMC_1.02', 'AWS_HIPAA', 'AWS_ISO_27001:2013', 
'AWS_NIST_CSF', 'AWS_NIST_800-171_rev2', 'AWS_NIST_800-53_rev5', 'AWS_PCI_DSS_3.2.1', 'AWS_SOC_2', 'LW_AWS_SEC_ADD_1_0'`)
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
			reportType, err := api.NewAwsReportType(compCmdState.Type)
			if err != nil {
				return errors.Errorf("invalid report type %q", compCmdState.Type)
			}

			var (
				// clean the AWS account ID if it was provided
				// with an Alias in between parentheses
				awsAccountID, _ = splitIDAndAlias(args[0])
				config          = api.AwsReportConfig{
					AccountID: awsAccountID,
					Type:      reportType,
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
				err := cli.LwApi.V2.Reports.Aws.DownloadPDF(pdfName, config)
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
				report   api.AwsReport
				cacheKey = fmt.Sprintf("compliance/aws/v2/%s/%s", config.AccountID, config.Type)
			)
			expired := cli.ReadCachedAsset(cacheKey, &report)
			if expired {
				cli.StartProgress("Getting compliance report...")
				response, err := cli.LwApi.V2.Reports.Aws.Get(config)
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

	// Todo(v2): deprecate??
	// complianceAwsRunAssessmentCmd represents the run-assessment sub-command inside the aws command
	complianceAwsRunAssessmentCmd = &cobra.Command{
		Use:     "run-assessment <account_id>",
		Aliases: []string{"run"},
		Short:   "Run a new AWS compliance assessment",
		Long:    `Run a compliance assessment for the provided AWS account.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// Todo(v2): replace with v2
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
		Hidden:  true,
		Aliases: []string{"disable"},
		Short:   "Disable all recommendations for a given report type",
		Long: `Disable all recommendations for a given report type.
Supported report types are CIS_1_1

To show the current status of recommendations in a report run:
	lacework compliance aws status CIS_1_1

To disable all recommendations for CIS_1_1 report run:
	lacework compliance aws disable CIS_1_1
`,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "CIS", "CIS_1_1", "AWS_CIS_S3":
				args[0] = "CIS_1_1"
				return nil
			default:
				return errors.New("CIS_1_1 is the only supported report type")
			}
		},
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// prompt for changes
			proceed, err := complianceAwsDisableReportDisplayChanges()
			if err != nil {
				return errors.Wrap(err, "unable to confirm disable")
			}
			if !proceed {
				return nil
			}

			schema, err := fetchCachedAwsComplianceReportSchema(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to get aws compliance report schema")
			}

			// set state of all recommendations in this report to disabled
			patchReq := api.NewRecommendationV1State(schema, false)
			cli.StartProgress("disabling recommendations...")
			response, err := cli.LwApi.Recommendations.Aws.Patch(patchReq)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to patch aws recommendations")
			}

			var cacheKey = fmt.Sprintf("compliance/aws/schema/%s", "CIS_1_1")
			cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), response.RecommendationList())
			cli.OutputHuman("All recommendations for report %s have been disabled\n", args[0])
			return nil
		},
	}

	// complianceAwsEnableReportCmd represents the enable-report sub-command inside the aws command
	// experimental feature
	complianceAwsEnableReportCmd = &cobra.Command{
		Use:     "enable-report <report_type>",
		Hidden:  true,
		Aliases: []string{"enable"},
		Short:   "Enable all recommendations for a given report type",
		Long: `Enable all recommendations for a given report type.
Supported report types are CIS_1_1

To show the current status of recommendations in a report run:
	lacework compliance aws status CIS_1_1

To enable all recommendations for CIS_1_1 report run:
	lacework compliance aws enable CIS_1_1
`,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "CIS", "CIS_1_1", "AWS_CIS_S3":
				args[0] = "CIS_1_1"
				return nil
			default:
				return errors.New("CIS_1_1 is the only supported report type")
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
			cli.StartProgress("enabling recommendations...")
			response, err := cli.LwApi.Recommendations.Aws.Patch(patchReq)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to patch aws recommendations")
			}

			var cacheKey = fmt.Sprintf("compliance/aws/schema/%s", args[0])
			cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), response.RecommendationList())
			cli.OutputHuman("All recommendations for report %s have been enabled\n", args[0])
			return nil
		},
	}

	// complianceAwsReportStatusCmd represents the report-status sub-command inside the aws command
	// experimental feature
	complianceAwsReportStatusCmd = &cobra.Command{
		Use:     "report-status <report_type>",
		Hidden:  true,
		Aliases: []string{"status"},
		Short:   "Show the status of recommendations for a given report type",
		Long: `Show the status of recommendations for a given report type.
Supported report types are CIS_1_1

To show the current status of recommendations in a report run:
	lacework compliance aws status CIS_1_1

The output from status with the --json flag can be used in the body of PATCH api/v1/external/recommendations/aws
	lacework compliance aws status CIS_1_1 --json
`,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "CIS", "CIS_1_1", "AWS_CIS_S3":
				args[0] = "CIS_1_1"
				return nil
			default:
				return errors.New("CIS_1_1 is the only supported report type")
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

	// complianceAwsListAccountsCmd represents the list-accounts inside the aws command
	complianceAwsSearchCmd = &cobra.Command{
		Use:   "search <resource_arn>",
		Short: "Search for all known violations of a given resource arn",
		Long:  `Search for all known violations of a given resource arn.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cli.StartProgress(fmt.Sprintf("Searching accounts for resource '%s'...", args[0]))
			var (
				now                        = time.Now().UTC()
				before                     = now.AddDate(0, 0, -7) // last 7 days
				awsInventorySearchResponse api.InventoryAwsResponse
				filter                     = api.InventorySearch{
					SearchFilter: api.SearchFilter{
						Filters: []api.Filter{{
							Expression: "eq",
							Field:      "urn",
							Value:      args[0],
						}},
						TimeFilter: &api.TimeFilter{
							StartTime: &before,
							EndTime:   &now,
						},
					},
					Dataset: api.AwsInventoryDataset,
					Csp:     api.AwsInventoryType,
				}
			)
			err := cli.LwApi.V2.Inventory.Search(&awsInventorySearchResponse, filter)
			cli.StopProgress()

			if len(awsInventorySearchResponse.Data) == 0 {
				cli.OutputHuman("Resource '%s' not found.\n\nTo learn how to configure Lacework with AWS visit "+
					"https://docs.lacework.com/onboarding/category/integrate-lacework-with-aws \n", args[0])
				return nil
			}
			cli.StopProgress()
			if err != nil {
				return err
			}

			cli.StartProgress(fmt.Sprintf("Searching for compliance violations for '%s'...", args[0]))
			var (
				awsComplianceEvaluationSearchResponse api.ComplianceEvaluationAwsResponse
				complianceFilter                      = api.ComplianceEvaluationSearch{
					SearchFilter: api.SearchFilter{
						Filters: []api.Filter{{
							Expression: "eq",
							Field:      "resource",
							Value:      args[0],
						}},
						TimeFilter: &api.TimeFilter{
							StartTime: &before,
							EndTime:   &now,
						},
					},
					Dataset: api.AwsComplianceEvaluationDataset,
				}
			)
			err = cli.LwApi.V2.ComplianceEvaluations.Search(&awsComplianceEvaluationSearchResponse, complianceFilter)
			cli.StopProgress()
			if err != nil {
				return err
			}

			var recommendationIDs []string
			var uniqueRecommendations []api.ComplianceEvaluationAws

			for _, recommend := range awsComplianceEvaluationSearchResponse.Data {
				if !array.ContainsStr(recommendationIDs, recommend.Id) {
					recommendationIDs = append(recommendationIDs, recommend.Id)
					uniqueRecommendations = append(uniqueRecommendations, recommend)
				}
			}

			if len(uniqueRecommendations) == 0 {
				cli.OutputHuman("No violations found. Time for %s\n", randomEmoji())
				return nil
			}

			// output table
			var out [][]string
			for _, recommend := range uniqueRecommendations {
				out = append(out, []string{
					recommend.Id,
					recommend.Account.AccountId,
					recommend.Reason,
					recommend.Severity,
					recommend.Status,
				})
			}

			cli.OutputHuman(
				renderSimpleTable(
					[]string{"RECOMMENDATION ID", "ACCOUNT ID", "REASON", "SEVERITY", "STATUS"},
					out,
				),
			)

			return nil
		},
	}
)

func init() {
	// add sub-commands to the aws command
	complianceAwsCmd.AddCommand(complianceAwsGetReportCmd)
	complianceAwsCmd.AddCommand(complianceAwsListAccountsCmd)
	complianceAwsCmd.AddCommand(complianceAwsRunAssessmentCmd)
	complianceAwsCmd.AddCommand(complianceAwsSearchCmd)

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

// Simple helper to prompt for approval after disable request
func complianceAwsDisableReportCmdPrompt() (int, error) {
	message := `WARNING! Disabling all recommendations for CIS_1_1 will disable the following reports and its corresponding compliance alerts:
AWS CIS Benchmark and S3 Report
AWS HIPAA Report
AWS ISO 27001:2013 Report
AWS NIST 800-171 Report
AWS NIST 800-53 Report
AWS PCI DSS Report
AWS SOC 2 Report
AWS SOC 2 Report Rev2

Would you like to proceed?
`
	options := []string{
		"Proceed with disable",
		"Quit",
	}

	var answer int
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Select{
			Message: message,
			Options: options,
		},
		Response: &answer,
	})

	return answer, err
}

func complianceAwsDisableReportDisplayChanges() (bool, error) {
	answer, err := complianceAwsDisableReportCmdPrompt()
	if err != nil {
		return false, err
	}
	return answer == 0, nil
}

func complianceAwsReportDetailsTable(report *api.AwsReport) [][]string {
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
