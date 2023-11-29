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
	// queryShowCmd represents the lql show command
	queryShowCmd = &cobra.Command{
		Use:   "show <query_id>",
		Short: "Show a query",
		Long:  `Show a query in your Lacework account.`,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			b, err := cmd.Flags().GetBool("yaml")
			if err != nil {
				return errors.Wrap(err, "unable to parse --yaml flag")
			}
			if b {
				cli.EnableYAMLOutput()
			}
			return nil
		},
		RunE: showQuery,
	}
)

func init() {
	queryCmd.AddCommand(queryShowCmd)

	queryShowCmd.Flags().Bool(
		"yaml", false, "output query in YAML format",
	)
}

func showQuery(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("retrieving query", "id", args[0])

	cli.StartProgress("Retrieving query...")
	queryResponse, err := cli.LwApi.V2.Query.Get(args[0])
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, "unable to show query")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(queryResponse.Data)
	}

	if cli.YAMLOutput() {
		return cli.OutputYAML(&api.NewQuery{
			QueryID:       queryResponse.Data.QueryID,
			QueryLanguage: queryResponse.Data.QueryLanguage,
			QueryText:     queryResponse.Data.QueryText,
		})
	}

	cli.OutputHuman(queryResponse.Data.QueryText)
	return nil
}
