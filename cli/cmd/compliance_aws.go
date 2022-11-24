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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/lacework/go-sdk/lwseverity"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	compAwsCmdState = struct {
		Type string
	}{Type: "AWS_CIS_14"}

	// complianceAwsListAccountsCmd represents the list-accounts inside the aws command
	complianceAwsListAccountsCmd = &cobra.Command{
		Use:     "list-accounts",
		Aliases: []string{"list", "ls"},
		Short:   "List all AWS accounts configured",
		Long:    `List all AWS accounts configured in your account.`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cli.StartProgress("Fetching list of configured AWS accounts...")
			awsAccounts, err := cli.LwApi.V2.CloudAccounts.ListByType(api.AwsCfgCloudAccount)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get aws compliance integrations")
			}

			return cliListAwsAccounts(awsAccounts)
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

			// Todo: Enable dynamic report type validation. Disabled until reportDefinitions api is out of beta
			// validTypes, err := getReportTypes(api.ReportDefinitionNotificationTypeAws)
			// if err != nil {
			//	 return errors.Wrap(err, "unable to retrieve valid report types")
			// }

			validTypes := []string{"AWS_NIST_CSF", "AWS_NIST_800-53_rev5", "AWS_HIPAA", "NIST_800-53_Rev4", "LW_AWS_SEC_ADD_1_0",
				"AWS_SOC_Rev2", "AWS_PCI_DSS_3.2.1", "AWS_CIS_S3", "ISO_2700", "SOC", "AWS_CSA_CCM_4_0_5", "PCI", "AWS_Cyber_Essentials_2_2",
				"AWS_ISO_27001:2013", "AWS_CIS_14", "AWS_CMMC_1.02", "HIPAA", "AWS_SOC_2", "AWS_CIS_1_4_ISO_IEC_27002_2022", "NIST_800-171_Rev2",
				"AWS_NIST_800-171_rev"}

			if array.ContainsStr(validTypes, compAwsCmdState.Type) {
				return nil
			} else {
				return errors.Errorf(`supported report types are: %s'`, strings.Join(validTypes, ", "))
			}
		},
		Short: "Get the latest AWS compliance report",
		Long: `Get the latest compliance assessment report from the provided AWS account, these
reports run on a regular schedule, typically once a day. The available report formats
are human-readable (default), json and pdf.

To list all AWS accounts configured in your account:

    lacework compliance aws list-accounts

To show recommendation details and affected resources for a recommendation id:

    lacework compliance aws get-report <account_id> [recommendation_id]
`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			reportType, err := api.NewAwsReportType(compAwsCmdState.Type)
			if err != nil {
				return errors.Errorf("invalid report type %q", compAwsCmdState.Type)
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
				if !lwseverity.IsValid(compCmdState.Severity) {
					return errors.Errorf("the severity %s is not valid, use one of %s",
						compCmdState.Severity, lwseverity.ValidSeverities.String(),
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

	// complianceAwsListAccountsCmd represents the search inside the aws command
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
			err := api.WindowedSearchFirst(cli.LwApi.V2.Inventory.Search, api.V2ApiMaxSearchWindowDays, api.V2ApiMaxSearchHistoryDays, &awsInventorySearchResponse, &filter)
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

			err = api.WindowedSearchFirst(cli.LwApi.V2.ComplianceEvaluations.Search, api.V2ApiMaxSearchWindowDays, api.V2ApiMaxSearchHistoryDays, &awsComplianceEvaluationSearchResponse, &complianceFilter)
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
	complianceAwsCmd.AddCommand(complianceAwsSearchCmd)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Details, "details", false,
		"increase details about the compliance report",
	)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Pdf, "pdf", false,
		"download report in PDF format",
	)

	complianceAwsGetReportCmd.Flags().BoolVar(&compCmdState.Csv, "csv", false,
		"output report in CSV format",
	)

	// AWS report types: AWS_NIST_CSF, AWS_NIST_800-53_rev5, AWS_HIPAA, NIST_800-53_Rev4, LW_AWS_SEC_ADD_1_0,
	//AWS_SOC_Rev2, AWS_PCI_DSS_3.2.1, AWS_CIS_S3, ISO_2700, SOC, AWS_CSA_CCM_4_0_5, PCI, AWS_Cyber_Essentials_2_2,
	//AWS_ISO_27001:2013, AWS_CIS_14, AWS_CMMC_1.02, HIPAA, AWS_SOC_2, AWS_CIS_1_4_ISO_IEC_27002_2022, NIST_800-171_Rev2,
	//AWS_NIST_800-171_rev2'
	complianceAwsGetReportCmd.Flags().StringVar(&compAwsCmdState.Type, "type", "AWS_CIS_14",
		"report type to display, run 'lacework report-definitions list' for valid types",
	)

	complianceAwsGetReportCmd.Flags().StringSliceVar(&compCmdState.Category, "category", []string{},
		"filter report details by category (identity-and-access-management, s3, logging...)",
	)

	complianceAwsGetReportCmd.Flags().StringSliceVar(&compCmdState.Service, "service", []string{},
		"filter report details by service (aws:s3, aws:iam, aws:cloudtrail, ...)",
	)

	complianceAwsGetReportCmd.Flags().StringVar(&compCmdState.Severity, "severity", "",
		fmt.Sprintf("filter report details by severity threshold (%s)",
			lwseverity.ValidSeverities.String()),
	)

	complianceAwsGetReportCmd.Flags().StringVar(&compCmdState.Status, "status", "",
		fmt.Sprintf("filter report details by status (%s)",
			strings.Join(api.ValidComplianceStatus, ", ")),
	)
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

func cliListAwsAccounts(awsIntegrations api.CloudAccountsResponse) error {
	awsAccounts := make([]awsAccount, 0)
	errCount := 0
	jsonOut := struct {
		Accounts []awsAccount `json:"aws_accounts"`
	}{Accounts: awsAccounts}

	if len(awsIntegrations.Data) == 0 {
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
		var (
			account     string
			accountData api.AwsCfgData
		)

		awsJson, err := json.Marshal(i.Data)
		if err != nil {
			continue
		}
		err = json.Unmarshal(awsJson, &accountData)
		if err != nil {
			continue
		}

		if accountData.AwsAccountID == "" {
			errCount++
			cli.Log.Debugf(fmt.Sprintf("unable to find account id for cloud account %s\n", i.IntgGuid))
			continue
		}
		account = accountData.AwsAccountID

		if containsDuplicateAccountID(awsAccounts, account) {
			cli.Log.Warnw("duplicate aws account", "integration_guid", i.IntgGuid, "account", account)
			continue
		}
		awsAccounts = append(awsAccounts, awsAccount{
			AccountID: account,
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
	if errCount > 0 {
		cli.OutputHuman(fmt.Sprintf("\n unable to find Aws Account ID's for %d integration(s)\n", errCount))
	}
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
