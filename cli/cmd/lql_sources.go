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
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	queryListSourcesCmd = &cobra.Command{
		Aliases: []string{"sources"},
		Use:     "list-sources",
		Short:   "List Lacework query data sources",
		Long:    `List Lacework query data sources.`,
		Args:    cobra.NoArgs,
		RunE:    listQuerySources,
	}

	queryShowSourceCmd = &cobra.Command{
		Aliases: []string{"describe"},
		Use:     "show-source <datasource_id>",
		Short:   "Show Lacework query data source",
		Long:    `Show Lacework query data source.`,
		Args:    cobra.ExactArgs(1),
		RunE:    showQuerySource,
	}
)

func init() {
	queryCmd.AddCommand(queryListSourcesCmd)
	queryCmd.AddCommand(queryShowSourceCmd)
}

func querySourcesTable(datasources []api.Datasource) (out [][]string) {
	for _, source := range datasources {
		out = append(out, []string{
			source.Name,
			source.Description,
		})
	}

	// order by Name
	sort.Slice(out, func(i, j int) bool {
		return out[i][0] < out[j][0]
	})

	return
}

func listQuerySources(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("retrieving LQL data sources")
	lqlSourcesUnableMsg := "unable to retrieve LQL data sources"
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
		renderCustomTable(
			[]string{"Datasource", "Description"},
			querySourcesTable(datasourcesResponse.Data),
			tableFunc(func(t *tablewriter.Table) {
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
	cli.OutputHuman("\nUse 'lacework query show-source <datasource_id>' to show details about the data source.\n")
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

func showQuerySource(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("retrieving datasource", "id", args[0])

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
		renderSimpleTable(
			[]string{"Datasource", "Description"},
			querySourcesTable([]api.Datasource{datasourceResponse.Data}),
		),
	)
	cli.OutputHuman("\n")
	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Field Name", "Data Type", "Description"},
			getShowQuerySourceTable(datasourceResponse.Data.ResultSchema),
		),
	)
	cli.OutputHuman("\nUse 'lacework query preview-source <datasource_id>' to see an actual result from the data source.\n")
	return nil
}
