//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	"strings"

	"github.com/lacework/go-sdk/lwseverity"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// policyUpdateCmd represents the policy update command
	policyUpdateCmd = &cobra.Command{
		Use:   "update [policy_id...]",
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
a policy identifer specified via payload.

The severity of many policies can be updated at once by passing a list of policy identifiers:

	lacework policy update my-policy-1 my-policy-2 --severity critical

`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// return error if multiple policy-ids are supplied without severity flag
			if len(args) > 1 && policyCmdState.Severity == "" {
				return errors.Errorf(`policy bulk update is only supported with the '--severity' flag 

For example: 

     lacework policy update %s --severity critical
					`, strings.Join(args, " "))
			}

			if policyCmdState.Severity != "" && !lwseverity.IsValid(policyCmdState.Severity) {
				return errors.Errorf("invalid severity %q valid severities are: %s",
					policyCmdState.Severity, lwseverity.ValidSeverities.String())
			}

			return nil
		},
		RunE: updatePolicy,
	}
)

func init() {
	// add sub-commands to the policy command
	policyCmd.AddCommand(policyUpdateCmd)

	// policy source specific flags
	setPolicySourceFlags(policyUpdateCmd)

	// add severity flag
	policyUpdateCmd.Flags().StringVar(&policyCmdState.Severity, "severity", "",
		"update the policy severity",
	)
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
		msg          = "unable to update policy"
		err          error
		queryUpdated bool
		policyID     string
	)

	// if severity flag is provided, attempt bulk update
	if policyCmdState.Severity != "" {
		err = policyBulkUpdate(args)
		if err != nil {
			return err
		}
		return nil
	}

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

func policyBulkUpdate(args []string) error {
	var (
		policyIds []string
		err       error
	)
	// if no policy ids are provided; prompt a list of policy ids
	if len(args) == 0 {
		policyIds, err = promptSetPolicyIDs()
		if err != nil {
			return err
		}
	} else {
		policyIds = args
	}

	var bulkPolicies api.BulkUpdatePolicies
	for _, p := range policyIds {
		bulkPolicies = append(bulkPolicies, api.BulkUpdatePolicy{
			PolicyID: p,
			Severity: policyCmdState.Severity,
		})
	}

	response, err := cli.LwApi.V2.Policy.UpdateMany(bulkPolicies)
	if err != nil {
		return errors.Wrap(err, "unable to update policies")
	}

	cli.Log.Debugw("bulk policy updated", "response", response)
	cli.OutputHuman("%d policies updated with new severity %q\n", len(policyIds), policyCmdState.Severity)

	return nil
}
