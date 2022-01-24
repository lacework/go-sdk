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
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// queryShowCmd represents the lql show command
	queryShowCmd = &cobra.Command{
		Use:   "show <query_id>",
		Short: "Show a query",
		Long:  `Show a query.`,
		Args:  cobra.ExactArgs(1),
		RunE:  showQuery,
	}
)

func init() {
	queryCmd.AddCommand(queryShowCmd)

	if IsLCLInstalled(*cli.LwComponents) {
		queryShowCmd.Flags().BoolVarP(
			&queryCmdState.ShowFromLibrary,
			"library", "l", false,
			"show query from Lacework Content Library",
		)
	}
}

func showQuery(_ *cobra.Command, args []string) error {
	var (
		lcl           *LaceworkContentLibrary
		newQuery      api.NewQuery
		queryResponse api.QueryResponse
		err           error
	)
	cli.Log.Debugw("retrieving query", "id", args[0])

	cli.StartProgress(" Retrieving query...")
	if queryCmdState.ShowFromLibrary {
		if lcl, err = LoadLCL(*cli.LwComponents); err == nil {
			newQuery, err = lcl.GetNewQuery(args[0])
			queryResponse.Data = api.Query{
				QueryID:     newQuery.QueryID,
				QueryText:   newQuery.QueryText,
				EvaluatorID: newQuery.EvaluatorID,
			}
		}
	} else {
		queryResponse, err = cli.LwApi.V2.Query.Get(args[0])
	}
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, "unable to show query")
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(queryResponse.Data)
	}
	cli.OutputHuman(queryResponse.Data.QueryText)
	return nil
}
