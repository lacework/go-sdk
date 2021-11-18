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
	// queryUpdateCmd represents the lql update command
	queryUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a query",
		Long: `
There are multiple ways you can update a query:

  * Typing the query into your default editor (via $EDITOR)
  * From a local file on disk using the flag '--file'
  * From a URL using the flag '--url'

There are also multiple formats you can use to define a query:

  * Javascript Object Notation (JSON)
  * YAML Ain't Markup Language (YAML)

To launch your default editor and update a query.

    lacework query update
`,
		Args: cobra.NoArgs,
		RunE: updateQuery,
	}
)

func init() {
	// add sub-commands to the lql command
	queryCmd.AddCommand(queryUpdateCmd)

	setQuerySourceFlags(queryUpdateCmd)
}

func updateQuery(cmd *cobra.Command, args []string) error {
	msg := "unable to update query"

	// input query
	queryString, err := inputQuery(cmd)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	// parse query
	newQuery, err := parseQuery(queryString)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	updateQuery := api.UpdateQuery{
		QueryText: newQuery.QueryText,
	}

	cli.Log.Debugw("updating query", "query", queryString)
	cli.StartProgress(" Updating query...")
	update, err := cli.LwApi.V2.Query.Update(newQuery.QueryID, updateQuery)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, msg)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(update.Data)
	}

	cli.OutputHuman("The query %s was updated.\n", update.Data.QueryID)
	return nil
}
