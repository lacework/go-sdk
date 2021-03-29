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

	"github.com/lacework/go-sdk/api"
)

const (
	lqlCreateDebugMsg    string = "creating LQL query"
	lqlCreateNotFoundMsg string = "Query created successfully but not returned.\n"
	lqlCreateSuccessMsg  string = "LQL query (%s) created successfully.\n"
	lqlCreateUnableMsg   string = "unable to create LQL query"
)

var (
	// lqlCreateCmd represents the lql create command
	lqlCreateCmd = &cobra.Command{
		Use:   "create [query]",
		Short: "create an LQL query",
		Long:  `Create an LQL query.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  createQuery,
	}
)

func init() {
	// add sub-commands to the lql command
	lqlCmd.AddCommand(lqlCreateCmd)

	setQueryFlags(lqlCreateCmd.Flags())
}

func createQuery(cmd *cobra.Command, args []string) error {
	var create api.LQLQueryResponse

	query, err := inputQuery(cmd, args)
	if err != nil {
		return errors.Wrap(err, lqlCreateUnableMsg)
	}

	cli.Log.Debugw(lqlCreateDebugMsg, "query", query)
	create, err = cli.LwApi.LQL.CreateQuery(query)

	if err != nil {
		return errors.Wrap(err, lqlCreateUnableMsg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(create.Data)
	}
	if len(create.Data) == 0 {
		cli.OutputHuman(lqlCreateNotFoundMsg)
	} else {
		cli.OutputHuman(
			fmt.Sprintf(lqlCreateSuccessMsg, create.Data[0].ID))
	}
	return nil
}
