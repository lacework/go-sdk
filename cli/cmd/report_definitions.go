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
	"fmt"
	"strconv"
	"strings"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	CreateReportDefinitionQuestion              = "Create from an existing report definition template?"
	CreateReportDefinitionReportNameQuestion    = "Report Name: "
	CreateReportDefinitionDisplayNameQuestion   = "Display Name: "
	CreateReportDefinitionReportSubTypeQuestion = "Report SubType: "
	CreateReportDefinitionAddSectionQuestion    = "Add another policy section?"
	CreateReportDefinitionSectionTitleQuestion  = "Section Title: "
	CreateReportDefinitionPoliciesQuestion      = "Select Policies in this Section: "
	SelectReportDefinitionQuestion              = "Select an existing report definition as a template?"

	UpdateReportDefinitionQuestion                   = "Update report definition in editor?"
	UpdateReportDefinitionReportNameQuestion         = "Report Name: "
	UpdateReportDefinitionDisplayNameQuestion        = "Display Name: "
	UpdateReportDefinitionEditSectionQuestion        = "Update an existing policy section?"
	UpdateReportDefinitionEditAnotherSectionQuestion = "Update another existing policy section?"
	UpdateReportDefinitionAddSectionQuestion         = "Add a new policy section?"
	UpdateReportDefinitionSelectSectionQuestion      = "Select a section to edit"

	reportDefinitionsCmdState = struct {
		// filter report definitions by subtype. 'AWS', 'GCP' or 'Azure'
		SubType string
		// create report definitions from a file input
		File string
		// retrieve report definition by version
		Version string
	}{}

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
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if reportDefinitionsCmdState.SubType != "" && !array.ContainsStr(api.ReportDefinitionSubtypes, reportDefinitionsCmdState.SubType) {
				return errors.Errorf("'%s' is not valid. Report definitions subtype can be %s", reportDefinitionsCmdState.SubType, api.ReportDefinitionSubtypes)
			}
			return nil
		},
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

			// filter definitions by subtype
			if reportDefinitionsCmdState.SubType != "" {
				filterReportDefinitions(&reportDefinitions)
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
		Long: `Show a single report definition by it's ID.
To show specific report definition version:

    lacework report-definition show <report_definition_id> --version <version>

To show all versions of a report definition:

    lacework report-definition show <report_definition_id> --version all

`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if reportDefinitionsCmdState.Version != "" {
				return fetchReportDefinitionVersion(args[0])
			}

			cli.StartProgress(" Fetching report definition...")
			response, err := cli.LwApi.V2.ReportDefinitions.Get(args[0])
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "unable to get report definition")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}
			outputReportDefinitionTable(response.Data)

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
			cli.StartProgress("Deleting report definition...")
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

func outputReportDefinitionTable(reportDefinition api.ReportDefinition) {
	headers := [][]string{
		{reportDefinition.ReportDefinitionGuid, reportDefinition.ReportName, reportDefinition.ReportType,
			reportDefinition.SubReportType},
	}

	cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "TYPE", "SUB-TYPE"}, headers))
	cli.OutputHuman("\n")
	cli.OutputHuman(buildReportDefinitionDetailsTable(reportDefinition))
}

func outputReportVersionsList(guid string, versions []string) {
	details := [][]string{{"VERSIONS", strings.Join(versions, ", ")}}

	detailsTable := &strings.Builder{}
	detailsTable.WriteString(renderOneLineCustomTable(guid,
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

	cli.OutputHuman(detailsTable.String())
}

func fetchReportDefinitionVersion(id string) error {
	var (
		err              error
		version          int
		reportDefinition api.ReportDefinition
		versions         []string
	)

	// if no version is supplied return all previous versions
	if reportDefinitionsCmdState.Version == "all" {
		cli.StartProgress("Fetching all report definition versions...")
		response, err := cli.LwApi.V2.ReportDefinitions.GetVersions(id)
		cli.StopProgress()

		if err != nil {
			return err
		}

		if cli.JSONOutput() {
			cli.OutputJSON(response)
		}

		for _, reportVersion := range response.Data {
			versions = append(versions, strconv.Itoa(reportVersion.Version))
		}

		outputReportVersionsList(response.Data[0].ReportDefinitionGuid, versions)
		return nil
	}

	if version, err = strconv.Atoi(reportDefinitionsCmdState.Version); err != nil {
		return errors.Wrap(err, "unable to parse version")
	}

	cli.StartProgress(fmt.Sprintf("Fetching report definition version %d...", version))
	response, err := cli.LwApi.V2.ReportDefinitions.GetVersions(id)
	cli.StopProgress()

	if err != nil {
		return err
	}

	for _, r := range response.Data {
		if r.Version == version {
			reportDefinition = r
		}
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(reportDefinition)
	}

	outputReportDefinitionTable(reportDefinition)

	return nil
}

func filterReportDefinitions(reportDefinitions *api.ReportDefinitionsResponse) {
	var filteredDefinitions []api.ReportDefinition
	for _, rd := range reportDefinitions.Data {
		if rd.SubReportType == reportDefinitionsCmdState.SubType {
			filteredDefinitions = append(filteredDefinitions, rd)
		}
	}
	reportDefinitions.Data = filteredDefinitions
}

func init() {
	// add the report-definition command
	rootCmd.AddCommand(reportDefinitionsCommand)

	// add sub-commands to the report-definition command
	reportDefinitionsCommand.AddCommand(reportDefinitionsCreateCommand)
	reportDefinitionsCommand.AddCommand(reportDefinitionsListCommand)
	reportDefinitionsCommand.AddCommand(reportDefinitionsShowCommand)
	reportDefinitionsCommand.AddCommand(reportDefinitionsDeleteCommand)
	reportDefinitionsCommand.AddCommand(reportDefinitionsUpdateCommand)
	reportDefinitionsCommand.AddCommand(reportDefinitionsRevertCommand)
	reportDefinitionsCommand.AddCommand(reportDefinitionsDiffCommand)

	// add flags to report-definition commands
	reportDefinitionsShowCommand.Flags().StringVar(&reportDefinitionsCmdState.Version,
		"version", "", "show a version of a report definition",
	)
	reportDefinitionsListCommand.Flags().StringVar(&reportDefinitionsCmdState.SubType,
		"subtype", "", "filter report definitions by subtype. 'AWS', 'GCP' or 'Azure'",
	)
	reportDefinitionsCreateCommand.Flags().StringVar(&reportDefinitionsCmdState.File,
		"file", "", "create a report definition from an existing definition file",
	)
	reportDefinitionsUpdateCommand.Flags().StringVar(&reportDefinitionsCmdState.File,
		"file", "", "update a report definition from an existing definition file",
	)
}

func buildReportDefinitionDetailsTable(definition api.ReportDefinition) string {
	var (
		details      [][]string
		engine       = ""
		releaseLabel = ""
	)

	details = append(details, []string{"FREQUENCY", definition.Frequency})
	if definition.Props != nil {
		engine = definition.Props.Engine
		releaseLabel = definition.Props.ReleaseLabel
	}

	details = append(details, []string{"ENGINE", engine})
	details = append(details, []string{"RELEASE LABEL", releaseLabel})
	details = append(details, []string{"UPDATED BY", definition.CreatedBy})
	details = append(details, []string{"LAST UPDATED", definition.CreatedTime.String()})
	details = append(details, []string{"VERSION", strconv.Itoa(definition.Version)})

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
