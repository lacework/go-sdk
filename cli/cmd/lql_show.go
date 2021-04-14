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
	// lqlShowCmd represents the lql show command
	lqlShowCmd = &cobra.Command{
		Use:   "show <query_id>",
		Short: "show an LQL query",
		Long:  `Show an LQL query.`,
		Args:  cobra.ExactArgs(1),
		RunE:  showQuery,
	}
)

func init() {
	lqlCmd.AddCommand(lqlShowCmd)
}

func showQuery(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("retrieving LQL query", "queryID", args[0])

	queryResponse, err := cli.LwApi.LQL.GetQueryByID(args[0])

	if err != nil {
		return errors.Wrap(err, "unable to show LQL query")
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(queryResponse.Data)
	}
	if len(queryResponse.Data) == 0 {
		return yikes("unable to show LQL query")
	}
	cli.OutputHuman(queryResponse.Data[0].QueryText)
	return nil
}
