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
	// queryDeleteCmd represents the lql delete command
	queryDeleteCmd = &cobra.Command{
		Use:   "delete <query_id>",
		Short: "delete a query",
		Long:  `Delete a query.`,
		Args:  cobra.ExactArgs(1),
		RunE:  deleteQuery,
	}
)

func init() {
	// add sub-commands to the lql command
	queryCmd.AddCommand(queryDeleteCmd)
}

func deleteQuery(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("deleting query", "id", args[0])

	_, err := cli.LwApi.V2.Query.Delete(args[0])

	if err != nil {
		return errors.Wrap(err, "unable to delete query")
	}
	cli.OutputHuman(
		fmt.Sprintf("query (%s) deleted successfully.\n", args[0]))
	return nil
}
