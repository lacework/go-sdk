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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/lacework/go-sdk/lwseverity"
)

var (
	compAzCmdState = struct {
		Type string
	}{Type: "AZURE_CIS_131"}

	// complianceAzureListSubsCmd represents the list-subscriptions sub-command inside the azure command
	complianceAzureListSubsCmd = &cobra.Command{
		Use:     "list-subscriptions",
		Aliases: []string{"list-subs"},
		Short:   "List subscriptions `<tenant-id>`",
		Long: `List all Azure subscriptions for Tenant.

Use the following command to list all Azure Tenants configured in your account:

    lacework compliance az list`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var (
				response, err             = cli.LwApi.V2.Configs.Azure.ListSubscriptions(args[0])
				cliCompAzureSubscriptions []cliComplianceAzureInfo
			)
			if err != nil {
				return errors.Wrap(err, "unable to list azure subscriptions")
			}

			if len(response.Data) == 0 {
				cli.OutputHuman("There are no azure subscriptions found for tenant %s\n", args[0])
				return nil
			}

			for _, az := range response.Data {
				cliCompAzureSubscriptions = append(cliCompAzureSubscriptions, splitAzureSubscriptionsApiResponse(az))
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(cliCompAzureSubscriptions)
			}

			rows := [][]string{}
			for _, subscriptionList := range cliCompAzureSubscriptions {
				for _, subscription := range subscriptionList.Subscriptions {
					rows = append(rows, []string{subscription.ID, subscription.Alias})
				}
			}

			cli.OutputHuman(renderSimpleTable(
				[]string{"Subscription ID", "Subscription Alias"}, rows),
			)
			return nil
		},
	}

	// complianceAzureListTenantsCmd represents the list-tenants sub-command inside the azure command
	complianceAzureListTenantsCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"list-tenants", "ls"},
		Short:   "List Azure tenants and subscriptions",
		Long:    `List all Azure tenants and subscriptions configured in your account.`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cli.StartProgress("Fetching list of configured Azure tenants...")
			response, err := cli.LwApi.V2.CloudAccounts.ListByType(api.AzureCfgCloudAccount)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get azure integrations")
			}

			if len(response.Data) == 0 {
				cli.OutputHuman("There are no azure cloud accounts configured in your account\n")
				return nil
			}

			return cliListTenantsAndSubscriptions(response)
		},
	}

	// complianceAzureGetReportCmd represents the get-report sub-command inside the azure command
	complianceAzureGetReportCmd = &cobra.Command{
		Use:     "get-report <tenant_id> <subscriptions_id>",
		Aliases: []string{"get", "show"},
		PreRunE: func(_ *cobra.Command, args []string) error {
			if compCmdState.Csv {
				cli.EnableCSVOutput()
			}

			if len(args) > 2 {
				compCmdState.RecommendationID = args[2]
				if !validRecommendationID(compCmdState.RecommendationID) {
					return errors.Errorf("\n'%s' is not a valid recommendation id\n", compCmdState.RecommendationID)
				}
			}
			// Todo: Enable dynamic report type validation. Disabled until reportDefinitions api is out of beta
			//validTypes, err := getReportTypes(api.ReportDefinitionNotificationTypeAzure)
			//if err != nil {
			//	return errors.Wrap(err, "unable to retrieve valid report types")
			//}

			validTypes := []string{"AZURE_CIS_131", "AZURE_NIST_800_171_REV2", "AZURE_NIST_800_53_REV5", "AZURE_NIST_CSF",
				"AZURE_PCI", "AZURE_SOC_Rev2", "AZURE_ISO_27001", "AZURE_SOC", "AZURE_HIPAA", "AZURE_CIS", "AZURE_PCI_Rev2"}
			if array.ContainsStr(validTypes, compAzCmdState.Type) {
				return nil
			} else {
				return errors.Errorf("supported report types are: %s", strings.Join(validTypes, ", "))
			}
		},
		Short: "Get the latest Azure compliance report",
		Long: `Get the latest Azure compliance assessment report, these reports run on a regular schedule,
typically once a day. The available report formats are human-readable (default), json and pdf.

To list all Azure tenants and subscriptions configured in your account:

    lacework compliance azure list

To show recommendation details and affected resources for a recommendation id:

    lacework compliance azure get-report <tenant_id> <subscriptions_id> [recommendation_id]
`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(_ *cobra.Command, args []string) error {
			reportType, err := api.NewAzureReportType(compAzCmdState.Type)
			if err != nil {
				return errors.Errorf("invalid report type %q", compAzCmdState.Type)
			}

			var (
				// clean tenantID and subscriptionID if they were provided
				// with an Alias in between parentheses
				tenantID, _       = splitIDAndAlias(args[0])
				subscriptionID, _ = splitIDAndAlias(args[1])
				config            = api.AzureReportConfig{
					TenantID:       tenantID,
					SubscriptionID: subscriptionID,
					Type:           reportType,
				}
			)

			if compCmdState.Pdf {
				pdfName := fmt.Sprintf(
					"%s_Report_%s_%s_%s_%s.pdf",
					config.Type,
					config.TenantID,
					config.SubscriptionID,
					cli.Account, time.Now().Format("20060102150405"),
				)

				cli.StartProgress("Downloading compliance report...")
				err := cli.LwApi.V2.Reports.Azure.DownloadPDF(pdfName, config)
				cli.StopProgress()
				if err != nil {
					return errors.Wrap(err, "unable to get azure pdf compliance report")
				}

				cli.OutputHuman("The Azure compliance report was downloaded at '%s'\n", pdfName)
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
				report   api.AzureReport
				cacheKey = fmt.Sprintf("compliance/azure/v2/%s/%s/%s",
					config.TenantID, config.SubscriptionID, config.Type)
			)
			expired := cli.ReadCachedAsset(cacheKey, &report)
			if expired {
				cli.StartProgress("Getting compliance report...")
				response, err := cli.LwApi.V2.Reports.Azure.Get(config)
				cli.StopProgress()
				if err != nil {
					return errors.Wrap(err, "unable to get azure compliance report")
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
						AccountName:     report.SubscriptionName,
						AccountID:       report.SubscriptionID,
						TenantName:      report.TenantName,
						TenantID:        report.TenantID,
						ReportType:      reportType.String(),
						ReportTime:      report.ReportTime,
						Recommendations: report.Recommendations,
					},
				)

				return cli.OutputCSV(
					[]string{"Report_Type", "Report_Time", "Tenant",
						"Subscription", "Section", "ID", "Recommendation",
						"Status", "Severity", "Resource", "Region", "Reason"},
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
					complianceAzureReportDetailsTable(&report),
					complianceReportSummaryTable(report.Summary),
					recommendations,
					filteredOutput,
				),
			)
			return nil
		},
	}

	// complianceAzureDisableReportCmd represents the disable-report sub-command inside the azure command
	// experimental feature
	complianceAzureDisableReportCmd = &cobra.Command{
		Use:     "disable-report <report_type>",
		Hidden:  true,
		Aliases: []string{"disable"},
		Short:   "Disable all recommendations for a given report type",
		Long: `Disable all recommendations for a given report type.
Supported report types are: CIS_1_0, CIS_1_3_1

To show the current status of recommendations in a report run:
	lacework compliance azure status CIS_1_3_1

To disable all recommendations for CIS_1_3_1 report run:
	lacework compliance azure disable CIS_1_3_1
`,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "CIS", "CIS_1_0", "AZURE_CIS":
				args[0] = "CIS_1_0"
				return nil
			case "CIS_1_3_1", "AZURE_CIS_131":
				args[0] = "CIS_1_3_1"
				return nil
			default:
				return errors.New("supported report types are: CIS_1_0, CIS_1_3_1")
			}
		},
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// prompt for changes
			proceed, err := complianceAzureDisableReportDisplayChanges(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to confirm disable")
			}
			if !proceed {
				return nil
			}

			schema, err := fetchCachedAzureComplianceReportSchema(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to fetch azure compliance report schema")
			}

			// set state of all recommendations in this report to disabled
			patchReq := api.NewRecommendationV2State(schema, false)
			cli.StartProgress("disabling recommendations...")
			response, err := cli.LwApi.V2.Recommendations.Azure.Patch(patchReq)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to patch azure recommendations")
			}

			var cacheKey = fmt.Sprintf("compliance/azure/schema/%s", args[0])
			cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), response.RecommendationList())
			cli.OutputHuman("All recommendations for report %s have been disabled\n", args[0])
			return nil
		},
	}

	// complianceAzureEnableReportCmd represents the enable-report sub-command inside the azure command
	// experimental feature
	complianceAzureEnableReportCmd = &cobra.Command{
		Use:     "enable-report <report_type>",
		Hidden:  true,
		Aliases: []string{"enable"},
		Short:   "Enable all recommendations for a given report type",
		Long: `Enable all recommendations for a given report type.
Supported report types are: CIS_1_0, CIS_1_3_1

To show the current status of recommendations in a report run:
	lacework compliance azure status CIS_1_3_1

To enable all recommendations for CIS_1_3_1 report run:
	lacework compliance azure enable CIS_1_3_1
`,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "CIS", "CIS_1_0", "AZURE_CIS":
				args[0] = "CIS_1_0"
				return nil
			case "CIS_1_3_1", "AZURE_CIS_131":
				args[0] = "CIS_1_3_1"
				return nil
			default:
				return errors.New("supported report types are: CIS_1_0, CIS_1_3_1")
			}
		},
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {

			schema, err := fetchCachedAzureComplianceReportSchema(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to fetch azure compliance report schema")
			}

			// set state of all recommendations in this report to enabled
			patchReq := api.NewRecommendationV2State(schema, true)
			cli.StartProgress("enabling recommendations...")
			response, err := cli.LwApi.V2.Recommendations.Azure.Patch(patchReq)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to patch azure recommendations")
			}

			var cacheKey = fmt.Sprintf("compliance/azure/schema/%s", args[0])
			cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), response.RecommendationList())
			cli.OutputHuman("All recommendations for report %s have been enabled\n", args[0])
			return nil
		},
	}

	// complianceAzureReportStatusCmd represents the report-status sub-command inside the azure command
	// experimental feature
	complianceAzureReportStatusCmd = &cobra.Command{
		Use:     "report-status <report_type>",
		Hidden:  true,
		Aliases: []string{"status"},
		Short:   "Show the status of recommendations for a given report type",
		Long: `Show the status of recommendations for a given report type.
Supported report types are: CIS_1_0, CIS_1_3_1

To show the current status of recommendations in a report run:
	lacework compliance azure status CIS_1_3_1

The output from status with the --json flag can be used in the body of PATCH api/v1/external/recommendations/azure
	lacework compliance azure status CIS_1_3_1 --json
`,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case "CIS", "CIS_1_0", "AZURE_CIS":
				args[0] = "CIS_1_0"
				return nil
			case "CIS_1_3_1", "AZURE_CIS_131":
				args[0] = "CIS_1_3_1"
				return nil
			default:
				return errors.New("supported report types are: CIS_1_0, CIS_1_3_1")
			}
		},
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var rows [][]string
			report, err := fetchCachedAzureComplianceReportSchema(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to fetch azure compliance report schema")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(api.NewRecommendationV2(report))
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
	// add sub-commands to the azure command
	complianceAzureCmd.AddCommand(complianceAzureListSubsCmd)
	complianceAzureCmd.AddCommand(complianceAzureListTenantsCmd)
	complianceAzureCmd.AddCommand(complianceAzureGetReportCmd)

	// Experimental Commands
	complianceAzureCmd.AddCommand(complianceAzureReportStatusCmd)
	complianceAzureCmd.AddCommand(complianceAzureDisableReportCmd)
	complianceAzureCmd.AddCommand(complianceAzureEnableReportCmd)

	complianceAzureGetReportCmd.Flags().BoolVar(&compCmdState.Details, "details", false,
		"increase details about the compliance report",
	)
	complianceAzureGetReportCmd.Flags().BoolVar(&compCmdState.Pdf, "pdf", false,
		"download report in PDF format",
	)

	// Output the report in CSV format
	complianceAzureGetReportCmd.Flags().BoolVar(&compCmdState.Csv, "csv", false,
		"output report in CSV format",
	)

	// Azure report types: AZURE_CIS_131, AZURE_NIST_800_171_REV2, AZURE_NIST_800_53_REV5, AZURE_NIST_CSF,
	//AZURE_PCI, AZURE_SOC_Rev2, AZURE_ISO_27001, AZURE_SOC, AZURE_HIPAA, AZURE_CIS, AZURE_PCI_Rev2
	complianceAzureGetReportCmd.Flags().StringVar(&compAzCmdState.Type, "type", "AZURE_CIS_131",
		"report type to display, run 'lacework report-definitions list' for valid types",
	)

	complianceAzureGetReportCmd.Flags().StringSliceVar(&compCmdState.Category, "category", []string{},
		"filter report details by category (networking, storage, ...)",
	)

	complianceAzureGetReportCmd.Flags().StringSliceVar(&compCmdState.Service, "service", []string{},
		"filter report details by service (azure:ms:storage, azure:ms:sql, azure:ms:network, ...)",
	)

	complianceAzureGetReportCmd.Flags().StringVar(&compCmdState.Severity, "severity", "",
		fmt.Sprintf("filter report details by severity threshold (%s)",
			lwseverity.ValidSeverities.String()),
	)

	complianceAzureGetReportCmd.Flags().StringVar(&compCmdState.Status, "status", "",
		fmt.Sprintf("filter report details by status (%s)",
			strings.Join(api.ValidComplianceStatus, ", ")),
	)
}

// Simple helper to prompt for approval after disable request
func complianceAzureDisableReportCmdPrompt(arg string) (int, error) {
	var message string
	switch arg {
	case "CIS", "CIS_1_0", "AZURE_CIS":
		message = `WARNING! Disabling all recommendations for CIS_1_0 will disable the following reports and its corresponding compliance alerts:
 AZURE CIS Benchmark
 PCI Benchmark
 SOC 2 Report

 Would you like to proceed?
 `
	case "CIS_1_3_1", "AZURE_CIS_131":
		message = `WARNING! Disabling all recommendations for CIS_1_3_1 will disable the following reports and its corresponding compliance alerts:
 AZURE CIS Benchmark 1.3.1
 PCI Benchmark Rev2
 SOC 2 Report Rev2
 HIPAA Report
 ISO27001 Report (+ couple CIS 1.0 controls)
 NIST 800-171 rev2 Report
 NIST 800-53 rev5 Report
 NIST CSF rev2 Report

 Would you like to proceed?
 `
	}

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

func complianceAzureDisableReportDisplayChanges(arg string) (bool, error) {
	answer, err := complianceAzureDisableReportCmdPrompt(arg)
	if err != nil {
		return false, err
	}
	return answer == 0, nil
}

func complianceAzureReportDetailsTable(report *api.AzureReport) [][]string {
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

func splitAzureSubscriptionsApiResponse(azInfo api.AzureConfigData) cliComplianceAzureInfo {
	var (
		tenantID, tenantAlias = splitIDAndAlias(azInfo.Tenant)
		cliAzureInfo          = cliComplianceAzureInfo{
			Tenant:        cliComplianceIDAlias{tenantID, tenantAlias},
			Subscriptions: make([]cliComplianceIDAlias, 0),
		}
	)

	for _, subscription := range azInfo.Subscriptions {
		id, alias := splitIDAndAlias(subscription)
		cliAzureInfo.Subscriptions = append(cliAzureInfo.Subscriptions, cliComplianceIDAlias{id, alias})
	}

	return cliAzureInfo
}

type cliComplianceAzureInfo struct {
	Tenant        cliComplianceIDAlias   `json:"tenant"`
	Subscriptions []cliComplianceIDAlias `json:"subscriptions"`
}

func cliListTenantsAndSubscriptions(azureIntegrations api.CloudAccountsResponse) error {
	jsonOut := struct {
		Subscriptions []azureSubscription `json:"azure_subscriptions"`
	}{Subscriptions: make([]azureSubscription, 0)}

	if len(azureIntegrations.Data) == 0 {
		if cli.JSONOutput() {
			return cli.OutputJSON(jsonOut)
		}

		msg := `There are no Azure Tenants configured in your account.

Get started by integrating your Azure Tenants to analyze configuration compliance using the command:

    lacework integration create

If you prefer to configure the integration via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Integrations > Cloud Accounts.
`
		cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
		return nil
	}

	if cli.JSONOutput() {
		jsonOut.Subscriptions = extractAzureSubscriptions(azureIntegrations)
		return cli.OutputJSON(jsonOut)
	}

	var rows [][]string
	for _, az := range extractAzureSubscriptions(azureIntegrations) {
		rows = append(rows, []string{az.TenantID, az.SubscriptionID, az.Status})
	}

	cli.OutputHuman(renderSimpleTable([]string{"Azure Tenant", "Azure Subscription", "Status"}, rows))
	return nil
}

type azureSubscription struct {
	TenantID       string `json:"tenant_id"`
	SubscriptionID string `json:"subscription_id"`
	Status         string `json:"status"`
}

func extractAzureSubscriptions(response api.CloudAccountsResponse) []azureSubscription {
	var azureSubscriptions []azureSubscription
	var azureData api.AzureCfgData
	if len(response.Data) == 0 {
		return azureSubscriptions
	}

	for _, az := range response.Data {
		azJson, err := json.Marshal(az.Data)
		if err != nil {
			continue
		}

		err = json.Unmarshal(azJson, &azureData)
		if err != nil {
			continue
		}
		// fetch the subscription ids from tenant id
		azureSubscriptions = append(azureSubscriptions, getAzureSubscriptions(azureData.TenantID, az.Status())...)
	}

	sort.Slice(azureSubscriptions, func(i, j int) bool {
		switch strings.Compare(azureSubscriptions[i].TenantID, azureSubscriptions[j].TenantID) {
		case -1:
			return true
		case 1:
			return false
		}
		return azureSubscriptions[i].SubscriptionID < azureSubscriptions[j].SubscriptionID
	})

	return azureSubscriptions
}

func getAzureSubscriptions(tenantID, status string) []azureSubscription {
	var subs []azureSubscription
	cli.StartProgress(fmt.Sprintf("Fetching subscriptions from tenant (%s)...", tenantID))
	subsResponse, err := cli.LwApi.V2.Configs.Azure.ListSubscriptions(tenantID)
	cli.StopProgress()
	if err != nil {
		cli.Log.Warnw("unable to list azure subscriptions", "tenant_id", tenantID, "error", err.Error())
		return subs
	}
	for _, subsRes := range subsResponse.Data {
		for _, subRes := range subsRes.Subscriptions {
			subscriptionID, _ := splitIDAndAlias(subRes)
			subs = append(subs, azureSubscription{
				TenantID:       tenantID,
				SubscriptionID: subscriptionID,
				Status:         status,
			})
		}
	}
	return subs
}

func fetchCachedAzureComplianceReportSchema(reportType string) (response []api.RecV2, err error) {
	var cacheKey = fmt.Sprintf("compliance/azure/schema/%s", reportType)

	expired := cli.ReadCachedAsset(cacheKey, &response)
	if expired {
		cli.StartProgress("Fetching compliance report schema...")
		response, err = cli.LwApi.V2.Recommendations.Azure.GetReport(reportType)
		cli.StopProgress()
		if err != nil {
			return nil, errors.Wrap(err, "unable to get Azure compliance report schema")
		}

		if len(response) == 0 {
			return nil, errors.New("no data found in the report")
		}

		cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Minute*30), response)
	}
	return
}
