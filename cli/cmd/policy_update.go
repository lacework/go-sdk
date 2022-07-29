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
	// policyUpdateCmd represents the policy update command
	policyUpdateCmd = &cobra.Command{
		Use:   "update [policy_id]",
		Short: "Update a policy",
		Long: `Update a policy.

A policy identifier is required to update a policy.

A policy identifier can be specified via:

1.  A policy update command argument

    lacework policy update my-policy-1

2. The policy update payload

    {
        "policy_id": "my-policy-1",
        "severity": "critical"
    }

A policy identifier specified via command argument always takes precedence over
a policy identifer specified via payload.`,
		Args: cobra.MaximumNArgs(1),
		RunE: updatePolicy,
	}
)

func init() {
	// add sub-commands to the policy command
	policyCmd.AddCommand(policyUpdateCmd)

	// policy source specific flags
	setPolicySourceFlags(policyUpdateCmd)
}

func updateQueryFromLibrary(id string) error {
	var (
		queryString string
		err         error
		newQuery    api.NewQuery
	)

	// input query
	queryString, err = inputQueryFromLibrary(id)
	if err != nil {
		return err
	}

	cli.Log.Debugw("creating query", "query", queryString)

	// parse query
	newQuery, err = api.ParseNewQuery(queryString)
	if err != nil {
		return queryErrorCrumbs(queryString)
	}
	updateQuery := api.UpdateQuery{
		QueryText: newQuery.QueryText,
	}

	// update query
	_, err = cli.LwApi.V2.Query.Update(newQuery.QueryID, updateQuery)
	return err
}

func updatePolicy(cmd *cobra.Command, args []string) error {
	var (
		msg          string = "unable to update policy"
		err          error
		queryUpdated bool
		policyID     string
	)

	if len(args) != 0 && len(args[0]) != 0 {
		policyID = args[0]
	}
	// input policy
	policyString, err := inputPolicy(cmd, policyID)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	cli.Log.Debugw("updating policy", "policy", policyString)

	// parse policy
	updatePolicy, err := api.ParseUpdatePolicy(policyString)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	// set policy id if not already set
	if policyID != "" && updatePolicy.PolicyID == "" {
		updatePolicy.PolicyID = policyID
	}

	cli.StartProgress("Updating policy...")
	updateResponse, err := cli.LwApi.V2.Policy.Update(updatePolicy)
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, msg)
	}
	// if updating policy from library also update query
	if policyCmdState.CUFromLibrary != "" {
		cli.StartProgress("Updating query...")
		err = updateQueryFromLibrary(updatePolicy.QueryID)
		cli.StopProgress()

		if err != nil {
			return errors.Wrap(err, msg)
		}
		queryUpdated = true
	}

	// output policy
	if cli.JSONOutput() {
		return cli.OutputJSON(updateResponse.Data)
	}
	if queryUpdated {
		cli.OutputHuman(fmt.Sprintf("The query %s was updated.\n", updatePolicy.QueryID))
	}
	cli.OutputHuman("The policy %s was updated.\n", updateResponse.Data.PolicyID)
	return nil
}
