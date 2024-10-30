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

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	queryListLibraryCmd = &cobra.Command{
		Use:   "list-library",
		Short: "List queries from library",
		Long:  `List all LQL queries in your Lacework Content Library.`,
		Args:  cobra.NoArgs,
		RunE:  listQueryLibrary,
	}
	queryShowLibraryCmd = &cobra.Command{
		Use:   "show-library <query_id>",
		Short: "Show a query from library",
		Long:  `Show a query in your Lacework Content Library.`,
		Args:  cobra.ExactArgs(1),
		RunE:  showQueryLibrary,
	}
)

func init() {
	if cli.isLCLInstalled() {
		queryCmd.AddCommand(queryListLibraryCmd)
		queryCmd.AddCommand(queryShowLibraryCmd)
	}
}

func getListQueryLibraryTable(queries map[string]LCLQuery) (out [][]string) {
	for id := range queries {
		out = append(out, []string{id})
	}
	// order by ID
	sort.Slice(out, func(i, j int) bool {
		return out[i][0] < out[j][0]
	})
	return
}

func listQueryLibrary(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("listing queries from library")

	cli.StartProgress(" Retrieving queries...")
	lcl, err := cli.LoadLCL()
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, "unable to list queries")
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(lcl.Queries)
	}
	if len(lcl.Queries) == 0 {
		cli.OutputHuman("There were no queries found.")
		return nil
	}
	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Query ID"},
			getListQueryLibraryTable(lcl.Queries),
		),
	)
	return nil
}

func showQueryLibrary(_ *cobra.Command, args []string) error {
	var (
		msg         string = "unable to show query"
		queryString string
		newQuery    api.NewQuery
		err         error
	)
	cli.Log.Debugw("retrieving query", "id", args[0])

	cli.StartProgress(" Retrieving query...")
	// input query
	if queryString, err = inputQueryFromLibrary(args[0]); err != nil {
		cli.StopProgress()
		return errors.Wrap(err, msg)
	}
	// parse query
	newQuery, err = api.ParseNewQuery(queryString)
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, msg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(newQuery)
	}
	cli.OutputHuman(newQuery.QueryText)
	return nil
}
