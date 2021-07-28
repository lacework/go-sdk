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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// queryCreateCmd represents the lql create command
	queryCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create a query",
		Long:  `Create a query.`,
		Args:  cobra.NoArgs,
		RunE:  createQuery,
	}
)

func init() {
	// add sub-commands to the lql command
	queryCmd.AddCommand(queryCreateCmd)

	setQuerySourceFlags(queryCreateCmd)
}

func createQuery(cmd *cobra.Command, args []string) error {
	msg := "unable to create query"

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

	cli.Log.Debugw("creating query", "query", queryString)
	create, err := cli.LwApi.V2.Query.Create(newQuery)

	if err != nil {
		return errors.Wrap(err, msg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(create.Data)
	}
	cli.OutputHuman(
		fmt.Sprintf("Query (%s) created successfully.\n", create.Data.QueryID))
	return nil
}
