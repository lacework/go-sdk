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

const (
	lqlValidateUnableMsg string = "unable to validate query"
)

var (
	// queryValidateCmd represents the lql validate command
	queryValidateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate a query",
		Long: `Use this command to validate a single LQL query before creating it.

There are multiple ways you can validate a query:

  * Typing the query into your default editor (via $EDITOR)
  * From a local file on disk using the flag '--file'
  * From a URL using the flag '--url'

There are also multiple formats you can use to define a query:

  * Javascript Object Notation (JSON)
  * YAML Ain't Markup Language (YAML)

To launch your default editor and validate a query.

    lacework query validate
`,
		Args: cobra.NoArgs,
		RunE: validateQuery,
	}
)

func init() {
	queryCmd.AddCommand(queryValidateCmd)

	setQuerySourceFlags(queryValidateCmd)

	if cli.IsLCLInstalled() {
		queryValidateCmd.Flags().StringVarP(
			&queryCmdState.CURVFromLibrary,
			"library", "l", "",
			"validate query from Lacework Content Library",
		)
	}
}

func validateQuery(cmd *cobra.Command, args []string) error {
	// input query
	queryString, err := inputQuery(cmd)
	if err != nil {
		return errors.Wrap(err, lqlValidateUnableMsg)
	}
	// parse query
	newQuery, err := api.ParseNewQuery(queryString)
	if err != nil {
		return errors.Wrap(queryErrorCrumbs(queryString), lqlValidateUnableMsg)
	}

	cli.Log.Debugw("validating query", "query", queryString)

	return validateQueryAndOutput(newQuery)
}

func validateQueryAndOutput(nq api.NewQuery) error {
	vq := api.ValidateQuery{
		QueryText:   nq.QueryText,
		EvaluatorID: nq.EvaluatorID,
	}

	cli.StartProgress(" Validating query...")
	validate, err := cli.LwApi.V2.Query.Validate(vq)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, lqlValidateUnableMsg)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(validate.Data)
	}

	cli.OutputHuman("Query validated successfully. Time for %s\n", randomEmoji())
	return nil
}
