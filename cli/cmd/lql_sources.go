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
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/api"
)

var (
	queryListSourcesCmd = &cobra.Command{
		Aliases: []string{"sources"},
		Use:     "list-sources",
		Short:   "List Lacework query datasources",
		Long:    `List Lacework query datasources.`,
		Args:    cobra.NoArgs,
		RunE:    listQuerySources,
	}

	queryShowSourceCmd = &cobra.Command{
		Aliases: []string{"describe"},
		Use:     "show-source <datasource_id>",
		Short:   "Show Lacework query datasource",
		Long:    `Show Lacework query datasource.`,
		Args:    cobra.ExactArgs(1),
		RunE:    showQuerySource,
	}
)

func init() {
	queryCmd.AddCommand(queryListSourcesCmd)
	queryCmd.AddCommand(queryShowSourceCmd)
}

func querySourcesTable(datasources []api.Datasource) (out [][]string) {
	var preOut [][]string
	for _, source := range datasources {
		preOut = append(preOut, []string{source.Name, source.Description})
	}

	// order by Name
	sort.Slice(preOut, func(i, j int) bool {
		return preOut[i][0] < preOut[j][0]
	})

	// condence output since datasources can be really long,
	// how long? you ask, as of today, we have over 150 characters
	for _, source := range preOut {
		out = append(out, []string{
			fmt.Sprintf("%s\n%s", source[0], source[1]),
		})
	}
	return
}

func listQuerySources(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("retrieving LQL datasources")
	lqlSourcesUnableMsg := "unable to retrieve LQL datasources"
	datasourcesResponse, err := cli.LwApi.V2.Datasources.List()

	if err != nil {
		return errors.Wrap(err, lqlSourcesUnableMsg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(datasourcesResponse.Data)
	}
	if len(datasourcesResponse.Data) == 0 {
		return yikes(lqlSourcesUnableMsg)
	}

	cli.OutputHuman(
		renderCustomTable([]string{"Datasource"},
			querySourcesTable(datasourcesResponse.Data),
			tableFunc(func(t *tablewriter.Table) {
				t.SetAlignment(tablewriter.ALIGN_LEFT)
				t.SetColWidth(120)
				t.SetAutoWrapText(true)
				t.SetRowLine(true)
				t.SetBorder(false)
				t.SetReflowDuringAutoWrap(false)
			}),
		),
	)

	cli.OutputHuman(
		"\nUse 'lacework query show-source <datasource_id>' to show details about the datasource.\n",
	)
	return nil
}

func getShowQuerySourceTable(resultSchema []api.DatasourceSchema) (out [][]string) {
	for _, schemaItem := range resultSchema {
		out = append(out, []string{
			schemaItem.Name,
			schemaItem.DataType,
			schemaItem.Description,
		})
	}
	return
}

func getShowQuerySourceRelationshipsTable(relationships []api.DatasourceRelationship) (out [][]string) {
	for _, relationship := range relationships {
		out = append(out, []string{
			relationship.Name,
			relationship.From,
			relationship.To,
			relationship.ToCardinality,
			relationship.Description,
		})
	}
	return
}

func showQuerySource(_ *cobra.Command, args []string) error {
	cli.Log.Infow("retrieving datasource", "id", args[0])

	cli.StartProgress(" Retrieving datasource...")
	datasourceResponse, err := cli.LwApi.V2.Datasources.Get(args[0])
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show datasource")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(datasourceResponse.Data)
	}
	cli.OutputHuman(
		renderOneLineCustomTable("Datasource",
			datasourceResponse.Data.Name,
			tableFunc(func(t *tablewriter.Table) {
				t.SetAlignment(tablewriter.ALIGN_LEFT)
				t.SetColWidth(120)
				t.SetBorder(false)
			}),
		),
	)
	cli.OutputHuman("\n")
	cli.OutputHuman(renderOneLineCustomTable("DESCRIPTION",
		datasourceResponse.Data.Description,
		tableFunc(func(t *tablewriter.Table) {
			t.SetAlignment(tablewriter.ALIGN_LEFT)
			t.SetColWidth(120)
			t.SetBorder(false)
			t.SetAutoWrapText(true)
		}),
	))
	cli.OutputHuman("\n")
	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Field Name", "Data Type", "Description"},
			getShowQuerySourceTable(datasourceResponse.Data.ResultSchema),
		),
	)
	// if source relationships exist
	if len(datasourceResponse.Data.SourceRelationships) > 0 {
		cli.OutputHuman("\n")
		cli.OutputHuman(
			renderSimpleTable(
				[]string{"Relationship Name", "From", "To", "Cardinality", "Description"},
				getShowQuerySourceRelationshipsTable(datasourceResponse.Data.SourceRelationships),
			),
		)
	}
	// breadcrumb
	cli.OutputHuman(
		fmt.Sprintf(
			"\nUse 'lacework query preview-source %s' to see an actual result from the datasource.\n",
			datasourceResponse.Data.Name,
		),
	)
	return nil
}
