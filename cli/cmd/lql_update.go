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

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/lacework/go-sdk/api"
)

var (
	// queryUpdateCmd represents the lql update command
	queryUpdateCmd = &cobra.Command{
		Use:   "update [query_id]",
		Short: "Update a query",
		Args:  cobra.RangeArgs(0, 1),
		Long: `
There are multiple ways you can update a query:

  * Typing the query into your default editor (via $EDITOR)
  * Passing a query id to load it into your default editor
  * From a local file on disk using the flag '--file'
  * From a URL using the flag '--url'

There are also multiple formats you can use to define a query:

  * Javascript Object Notation (JSON)
  * YAML Ain't Markup Language (YAML)

To launch your default editor and update a query.

    lacework query update
`,
		RunE: updateQuery,
	}
)

func init() {
	// add sub-commands to the lql command
	queryCmd.AddCommand(queryUpdateCmd)

	setQuerySourceFlags(queryUpdateCmd)

	if cli.IsLCLInstalled() {
		queryUpdateCmd.Flags().StringVarP(
			&queryCmdState.CUVFromLibrary,
			"library", "l", "",
			"update query from Lacework Content Library",
		)
	}
}

func updateQuery(cmd *cobra.Command, args []string) error {
	msg := "unable to update query"

	var (
		queryString string
		err         error
	)

	if len(args) != 0 {
		// query id via argument
		cli.StartProgress("Retrieving query...")
		queryRes, err := cli.LwApi.V2.Query.Get(args[0])
		cli.StopProgress()
		if err != nil {
			return errors.Wrap(err, "unable to load query from your account")
		}

		queryYaml, err := yaml.Marshal(&api.NewQuery{
			QueryID:     queryRes.Data.QueryID,
			QueryText:   queryRes.Data.QueryText,
			EvaluatorID: queryRes.Data.EvaluatorID,
		})
		if err != nil {
			return errors.Wrap(err, msg)
		}

		prompt := &survey.Editor{
			Message:       fmt.Sprintf("Update query %s", args[0]),
			Default:       string(queryYaml),
			HideDefault:   true,
			AppendDefault: true,
			FileName:      "query*.yaml",
		}
		var queryStr string
		err = survey.AskOne(prompt, &queryStr)
		if err != nil {
			return errors.Wrap(err, msg)
		}

		queryString = queryStr
	} else {
		// input query
		queryString, err = inputQuery(cmd)
		if err != nil {
			return errors.Wrap(err, msg)
		}
	}

	// parse query
	newQuery, err := api.ParseNewQuery(queryString)
	if err != nil {
		return errors.Wrap(queryErrorCrumbs(queryString), msg)
	}

	// avoid letting the user change the query id
	if len(args) != 0 && newQuery.QueryID != args[0] {
		return errors.New("changes to query id not supported")
	}

	// update query
	cli.Log.Debugw("updating query", "query", queryString)
	cli.StartProgress(" Updating query...")
	update, err := cli.LwApi.V2.Query.Update(newQuery.QueryID, api.UpdateQuery{
		QueryText: newQuery.QueryText,
	})
	cli.StopProgress()

	// output
	if err != nil {
		return errors.Wrap(err, msg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(update.Data)
	}
	cli.OutputHuman("The query %s was updated.\n", update.Data.QueryID)
	return nil
}
