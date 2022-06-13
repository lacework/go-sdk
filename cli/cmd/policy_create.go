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
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// policyCreateCmd represents the policy create command
	policyCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a policy",
		Long: `Create a policy.

A policy is represented in either JSON or YAML format.

The following attributes are minimally required:

    ---
    policyType: Violation
    queryId: MyQuery
    title: My Policy
    enabled: false
    description: My Policy Description
    remediation: My Policy Remediation
    severity: high
    evalFrequency: Daily
    alertEnabled: false
    alertProfile: LW_CloudTrail_Alerts
`,
		Args: cobra.NoArgs,
		RunE: createPolicy,
	}
)

func init() {
	// add sub-commands to the policy command
	policyCmd.AddCommand(policyCreateCmd)

	// policy source specific flags
	setPolicySourceFlags(policyCreateCmd)
}

func createQueryFromLibrary(id string) error {
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

	// create query
	_, err = cli.LwApi.V2.Query.Create(newQuery)
	return err
}

func createPolicy(cmd *cobra.Command, _ []string) error {
	var (
		msg         string = "unable to create policy"
		err         error
		queryExists bool
	)

	// input policy
	policyString, err := inputPolicy(cmd)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	cli.Log.Debugw("creating policy", "policy", policyString)

	// parse policy
	newPolicy, err := api.ParseNewPolicy(policyString)
	if err != nil {
		return errors.Wrap(err, msg)
	}

	// if creating policy from library also create query
	if policyCmdState.CUFromLibrary != "" {
		cli.StartProgress(" Creating query (then policy)...")
		err = createQueryFromLibrary(newPolicy.QueryID)
		cli.StopProgress()

		if err != nil {
			if queryExists = strings.Contains(err.Error(), "already exists"); !queryExists {
				return errors.Wrap(err, "unable to create query")
			}
		}
	}

	// create policy
	cli.StartProgress(" Creating policy...")
	createResponse, err := cli.LwApi.V2.Policy.Create(newPolicy)
	cli.StopProgress()

	// output policy
	if err != nil {
		return errors.Wrap(err, msg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(createResponse.Data)
	}
	// if human output mode, creating from library, and query exists
	if queryExists {
		cli.OutputHuman(fmt.Sprintf("The query %s already exists.\n", newPolicy.QueryID))
	}
	if policyCmdState.CUFromLibrary != "" {
		cli.OutputHuman(fmt.Sprintf("The query %s was created.\n", newPolicy.QueryID))
	}
	cli.OutputHuman(fmt.Sprintf("The policy %s was created.\n", createResponse.Data.PolicyID))
	return nil
}
