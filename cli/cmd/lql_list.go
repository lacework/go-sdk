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

	"github.com/lacework/go-sdk/api"
)

var (
	// queryListCmd represents the lql list command
	queryListCmd = &cobra.Command{
		Use:   "list",
		Short: "list queries",
		Long:  `List queries.`,
		Args:  cobra.NoArgs,
		RunE:  listQueries,
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
	return
}

func listQueries(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("listing queries")

	queryResponse, err := cli.LwApi.V2.Query.List()

	if err != nil {
		return errors.Wrap(err, "unable to list queries")
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(queryResponse.Data)
	}
	if len(queryResponse.Data) == 0 {
		cli.OutputHuman("There were no queries found.")
		return nil
	}
	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Query ID", "Owner", "Last Update Time", "Last Update User"},
			queryTable(queryResponse.Data),
		),
	)
	return nil
}
