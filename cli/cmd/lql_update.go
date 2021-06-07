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
	// lqlUpdateCmd represents the lql update command
	lqlUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "update an LQL query",
		Long:  `Update an LQL query.`,
		Args:  cobra.NoArgs,
		RunE:  updateQuery,
	}
)

func init() {
	// add sub-commands to the lql command
	lqlCmd.AddCommand(lqlUpdateCmd)

	setQuerySourceFlags(lqlUpdateCmd)
}

func updateQuery(cmd *cobra.Command, args []string) error {
	lqlUpdateUnableMsg := "unable to update LQL query"

	query, err := inputQuery(cmd, args)
	if err != nil {
		return errors.Wrap(err, lqlUpdateUnableMsg)
	}

	cli.Log.Debugw("updating LQL query", "query", query)
	update, err := cli.LwApi.LQL.Update(query)

	if err != nil {
		err = queryErrorCrumbs(query, err)
		return errors.Wrap(err, lqlUpdateUnableMsg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(update.Message)
	}
	cli.OutputHuman(
		fmt.Sprintf("LQL query (%s) updated successfully.\n", update.Message.ID))
	return nil
}
