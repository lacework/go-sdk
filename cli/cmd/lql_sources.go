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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// querySourcesCmd represents the lql data sources command
	querySourcesCmd = &cobra.Command{
		Hidden:  true, // no api exists for this
		Aliases: []string{"sources"},
		Use:     "list-sources",
		Short:   "list LQL data sources",
		Long:    `List LQL data sources.`,
		Args:    cobra.NoArgs,
		RunE:    getQuerySources,
	}
)

func init() {
	queryCmd.AddCommand(querySourcesCmd)
}

func dataSourcesToTable(dataSources []string) (out [][]string) {
	for _, source := range dataSources {
		out = append(out, []string{
			source,
		})
	}
	return
}

func getQuerySources(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("retrieving LQL data sources")

	lqlSourcesUnableMsg := "unable to retrieve LQL data sources"
	dataSourcesResponse, err := cli.LwApi.V2.Query.DataSources()

	if err != nil {
		return errors.Wrap(err, lqlSourcesUnableMsg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(dataSourcesResponse.Data)
	}
	if len(dataSourcesResponse.Data) == 0 {
		return yikes(lqlSourcesUnableMsg)
	}
	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Data Source"},
			dataSourcesToTable(dataSourcesResponse.Data[0].DataSources),
		),
	)
	return nil
}
