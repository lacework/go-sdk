//
// Author:: Darren Murray(<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"strings"

	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// report-definitions command is used to manage lacework report definitions
	reportDefinitionsCommand = &cobra.Command{
		Use:     "report-definition",
		Aliases: []string{"report-definitions", "rd"},
		Short:   "Manage report definitions",
		Long: `Manage report definitions to configure the data retrieval and layout information for a report.
`,
	}

	// list command is used to list all lacework report definitions
	reportDefinitionsListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all report definitions",
		Long:    "List all report definitions configured in your Lacework account.",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cli.StartProgress(" Fetching report definitions...")
			reportDefinitions, err := cli.LwApi.V2.ReportDefinitions.List()
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get report definitions")
			}

			if len(reportDefinitions.Data) == 0 {
				cli.OutputHuman("There are no report definitions configured in your account.\n")
				return nil
			}
			if cli.JSONOutput() {
				return cli.OutputJSON(reportDefinitions)
			}

			var rows [][]string
			for _, definition := range reportDefinitions.Data {
				rows = append(rows, []string{definition.ReportDefinitionGuid, definition.ReportName,
					definition.ReportType, definition.SubReportType})
			}

			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "TYPE", "SUB-TYPE"}, rows))
			return nil
		},
	}
	// show command is used to retrieve a lacework report definition by guid
	reportDefinitionsShowCommand = &cobra.Command{
		Use:   "show <report_definition_id>",
		Short: "Show a report definition by ID",
		Long:  "Show a single report definition by it's ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cli.StartProgress(" Fetching report definition...")
			response, err := cli.LwApi.V2.ReportDefinitions.Get(args[0])
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "unable to get report definition")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			reportDefinition := response.Data
			headers := [][]string{
				{reportDefinition.ReportDefinitionGuid, reportDefinition.ReportName, reportDefinition.ReportType,
					reportDefinition.SubReportType},
			}

			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "TYPE", "SUB-TYPE"}, headers))
			cli.OutputHuman("\n")
			cli.OutputHuman(buildReportDefinitionDetailsTable(reportDefinition))

			return nil
		},
	}

	// delete command is used to remove a lacework report definition by id
	reportDefinitionsDeleteCommand = &cobra.Command{
		Use:   "delete <report_definition_id>",
		Short: "Delete a report definition",
		Long:  "Delete a single report definition by it's ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cli.StartProgress(" Deleting report definition...")
			err := cli.LwApi.V2.ReportDefinitions.Delete(args[0])
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to delete report definition")
			}
			cli.OutputHuman("The report definition with GUID %s was deleted\n", args[0])
			return nil
		},
	}
)

func init() {
	// add the report-definition command
	rootCmd.AddCommand(reportDefinitionsCommand)

	// add sub-commands to the report-definition command
	reportDefinitionsCommand.AddCommand(reportDefinitionsListCommand)
	reportDefinitionsCommand.AddCommand(reportDefinitionsShowCommand)
	reportDefinitionsCommand.AddCommand(reportDefinitionsDeleteCommand)
}

func buildReportDefinitionDetailsTable(definition api.ReportDefinition) string {
	var (
		details [][]string
	)

	details = append(details, []string{"FREQUENCY", definition.Frequency})
	details = append(details, []string{"ENGINE", definition.Props.Engine})
	details = append(details, []string{"RELEASE LABEL", definition.Props.ReleaseLabel})
	details = append(details, []string{"UPDATED BY", definition.CreatedBy})
	details = append(details, []string{"LAST UPDATED", definition.CreatedTime.String()})

	detailsTable := &strings.Builder{}
	detailsTable.WriteString(renderOneLineCustomTable("REPORT DEFINITION DETAILS",
		renderCustomTable([]string{}, details,
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
	),
	)

	policiesTable := &strings.Builder{}
	if len(definition.ReportDefinitionDetails.Sections) > 0 {

		for _, s := range definition.ReportDefinitionDetails.Sections {
			var policies [][]string
			policies = append(policies, []string{s.Title, strings.Join(s.Policies, ", ")})

			policiesTable.WriteString(renderCustomTable([]string{"title", "policy"}, policies,
				tableFunc(func(t *tablewriter.Table) {
					t.SetBorder(false)
					t.SetColumnSeparator(" ")
				}),
			),
			)
		}
		policiesTable.WriteString("\n")

		detailsTable.WriteString(renderOneLineCustomTable("POLICIES", policiesTable.String(), tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetColumnSeparator(" ")
			t.SetAutoWrapText(false)
		})))
	}

	return detailsTable.String()
}