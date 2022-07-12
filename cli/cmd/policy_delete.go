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

var (
	// policyDeleteCmd represents the policy delete command
	policyDeleteCmd = &cobra.Command{
		Use:   "delete <policy_id>",
		Short: "Delete a policy",
		Long: `Delete a policy by providing the policy ID.

Use the command 'lacework policy list' to list the registered policies in
your Lacework account.`,
		Args: cobra.ExactArgs(1),
		RunE: deletePolicy,
	}
)

func init() {
	// add sub-commands to the policy command
	policyCmd.AddCommand(policyDeleteCmd)

	policyDeleteCmd.Flags().BoolVar(
		&policyCmdState.CascadeDelete,
		"cascade", false, "delete policy and its associated query",
	)
}

func deletePolicy(_ *cobra.Command, args []string) error {
	var (
		getResponse api.PolicyResponse
		err         error
		queryID     string
	)

	if policyCmdState.CascadeDelete {
		cli.Log.Debugw("retrieving policy", "policyID", args[0])
		cli.StartProgress("Retrieving policy...")
		getResponse, err = cli.LwApi.V2.Policy.Get(args[0])
		cli.StopProgress()

		if err != nil {
			return errors.Wrap(err, "unable to retrieve policy")
		}
		queryID = getResponse.Data.QueryID
	}

	cli.Log.Debugw("deleting policy", "policyID", args[0])
	cli.StartProgress(" Deleting policy...")
	_, err = cli.LwApi.V2.Policy.Delete(args[0])
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, "unable to delete policy")
	}
	cli.OutputHuman(
		fmt.Sprintf("The policy %s was deleted.\n", args[0]),
	)
	// delete query
	if policyCmdState.CascadeDelete {
		cli.Log.Debugw("deleting query", "id", queryID)
		cli.StartProgress(" Deleting query...")
		_, err := cli.LwApi.V2.Query.Delete(queryID)
		cli.StopProgress()

		if err != nil {
			return errors.Wrap(err, "unable to delete query")
		}
		cli.OutputHuman("The query %s was deleted.\n", queryID)
	}
	return nil
}
