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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// queryListCmd represents the lql list command
	queryListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List queries",
		Long:    `List all LQL queries in your Lacework account.`,
		Args:    cobra.NoArgs,
		RunE:    listQueries,
	}
)

func init() {
	queryCmd.AddCommand(queryListCmd)
}

func queryTable(queryData []api.Query) (out [][]string) {
	for _, query := range queryData {
		out = append(out, []string{
			query.QueryID,
			query.Owner,
			query.LastUpdateTime,
			query.LastUpdateUser,
		})
	}

	// order by ID
	sort.Slice(out, func(i, j int) bool {
		return out[i][0] < out[j][0]
	})

	return
}

func listQueries(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("listing queries")

	cli.StartProgress(" Retrieving queries...")
	queryResponse, err := cli.LwApi.V2.Query.List()
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to list queries")
	}

	if len(queryResponse.Data) == 0 {
		cli.OutputHuman("There were no queries found.")
		return nil
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(queryResponse.Data)
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Query ID", "Owner", "Last Update Time", "Last Update User"},
			queryTable(queryResponse.Data),
		),
	)
	return nil
}
