//
// Author:: Darren Murray(<darren.murray@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

	"github.com/fatih/structs"
	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	CreateReportDistributionReportNameQuestion     = "Report Distribution Name: "
	CreateReportDistributionFrequencyQuestion      = "Select Frequency: "
	CreateReportDistributionDefinitionQuestion     = "Select Report Definition: "
	CreateReportDistributionAlertChannelsQuestion  = "Select Alert Channels: "
	CreateReportDistributionResourceGroupsQuestion = "Select Resource Groups: "
	CreateReportDistributionIntegrationAwsQuestion = "Select Aws Accounts: "
	CreateReportDistributionAddSeveritiesQuestion  = "Add Severities? "
	CreateReportDistributionSeveritiesQuestion     = "Select Severities: "
	CreateReportDistributionAddViolationsQuestion  = "Add Violations? "
	CreateReportDistributionScopeQuestion          = "Select Distribution Scope:"
	CreateReportDistributionViolationsQuestion     = "Select Violations: "
	UpdateReportDistributionReportNameQuestion     = "Update Report Distribution Name? "
	UpdateReportDistributionFrequencyQuestion      = "Update Frequency?"
	UpdateReportDistributionAlertChannelsQuestion  = "Update Alert Channels? "
	UpdateReportDistributionAddSeveritiesQuestion  = "Update Severities? "
	UpdateReportDistributionAddViolationsQuestion  = "Update Violations? "

	// report-distributions command is used to manage lacework report distributions
	reportDistributionsCommand = &cobra.Command{
		Use:     "report-distribution",
		Hidden:  true,
		Aliases: []string{"report-distributions"},
		Short:   "Manage report distributions",
		Long: `Manage report distributions to configure the data retrieval and layout information for a report.
`,
	}

	// list command is used to list all lacework report distributions
	reportDistributionsListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all report distributions",
		Long:    "List all report distributions configured in your Lacework account.",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cli.StartProgress(" Fetching report distributions...")
			reportDistributions, err := cli.LwApi.V2.ReportDistributions.List()
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get report distributions")
			}

			if len(reportDistributions.Data) == 0 {
				cli.OutputHuman("There are no report distributions configured in your account.\n")
				return nil
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(reportDistributions)
			}

			var rows [][]string
			for _, distribution := range reportDistributions.Data {
				rows = append(rows, []string{distribution.ReportDistributionGuid, distribution.DistributionName,
					distribution.Frequency, distribution.ReportDefinitionGuid})
			}

			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "FREQUENCY", "DEFINITION GUID"}, rows))
			return nil
		},
	}
	// show command is used to retrieve a lacework report distribution by guid
	reportDistributionsShowCommand = &cobra.Command{
		Use:   "show <report_distribution_id>",
		Short: "Show a report distribution by ID",
		Long: `Show a single report distribution by it's ID.
To show specific report distribution details:

    lacework report-distribution show <report_distribution_id>

`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.StartProgress(" Fetching report distribution...")
			response, err := cli.LwApi.V2.ReportDistributions.Get(args[0])
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "unable to get report distribution")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}
			outputReportDistributionTable(response.Data)

			return nil
		},
	}

	// delete command is used to remove a lacework report distribution by id
	reportDistributionsDeleteCommand = &cobra.Command{
		Use:   "delete <report_distribution_id>",
		Short: "Delete a report distribution",
		Long:  "Delete a single report distribution by it's ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cli.StartProgress("Deleting report distribution...")
			err := cli.LwApi.V2.ReportDistributions.Delete(args[0])
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to delete report distribution")
			}
			cli.OutputHuman("The report distribution with GUID %s was deleted\n", args[0])
			return nil
		},
	}
)

func outputReportDistributionTable(reportDistribution api.ReportDistribution) {
	headers := [][]string{
		{reportDistribution.ReportDistributionGuid, reportDistribution.DistributionName, reportDistribution.Frequency,
			reportDistribution.ReportDefinitionGuid},
	}

	cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "FREQUENCY", "DEFINITION GUID"}, headers))
	cli.OutputHuman("\n")
	cli.OutputHuman(buildReportDistributionDetailsTable(reportDistribution))
}

func init() {
	// add the report-distribution command
	rootCmd.AddCommand(reportDistributionsCommand)

	// add sub-commands to the report-distribution command
	reportDistributionsCommand.AddCommand(reportDistributionsListCommand)
	reportDistributionsCommand.AddCommand(reportDistributionsShowCommand)
	reportDistributionsCommand.AddCommand(reportDistributionsDeleteCommand)
	reportDistributionsCommand.AddCommand(reportDistributionsCreateCommand)
	reportDistributionsCommand.AddCommand(reportDistributionsUpdateCommand)
}

func buildReportDistributionDetailsTable(distribution api.ReportDistribution) string {
	var (
		details      [][]string
		integrations []string
	)

	if distribution.Data.Integrations != nil {
		for _, integration := range distribution.Data.Integrations {
			var integrationKV strings.Builder
			integrationMap := structs.Map(integration)
			for k, v := range integrationMap {
				if v == "" {
					continue
				}
				integrationKV.WriteString(fmt.Sprintf("%s: %s ", k, v))
			}
			integrations = append(integrations, integrationKV.String())
		}
	}

	details = append(details, []string{"SEVERITY", strings.Join(distribution.Data.Severities, ", ")})
	details = append(details, []string{"VIOLATIONS", strings.Join(distribution.Data.Violations, ", ")})

	detailsTable := &strings.Builder{}
	dataTable := &strings.Builder{}

	dataTable.WriteString(renderCustomTable([]string{}, details,
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetColumnSeparator(" ")
			t.SetAutoWrapText(false)
			t.SetAlignment(tablewriter.ALIGN_LEFT)
		}),
	))
	dataTable.WriteString("\n")

	var data [][]string
	channels := strings.Join(distribution.AlertChannels, "\n")
	groups := strings.Join(distribution.Data.ResourceGroups, "\n")
	integrationList := strings.Join(integrations, "\n")

	data = append(data, []string{channels, groups, integrationList})

	dataTable.WriteString(renderCustomTable([]string{"ALERT CHANNELS", "RESOURCE GROUPS", "INTEGRATIONS"}, data,
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetColumnSeparator(" ")
		}),
	),
	)

	detailsTable.WriteString(renderOneLineCustomTable("REPORT DISTRIBUTION DETAILS",
		dataTable.String(),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		}),
	),
	)

	return detailsTable.String()
}
