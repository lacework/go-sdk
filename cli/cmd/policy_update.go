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
	"encoding/json"
	"reflect"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

A policy identifier specifed via command argument will always take precedence over
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

type updatePoliciesYAML struct {
	Policies []api.UpdatePolicy `yaml:"policies"`
}

func parseUpdatePolicy(s string) (api.UpdatePolicy, error) {
	var policy api.UpdatePolicy
	var err error

	// valid json
	if err = json.Unmarshal([]byte(s), &policy); err == nil {
		return policy, err
	}
	// nested yaml
	var policies updatePoliciesYAML

	if err = yaml.Unmarshal([]byte(s), &policies); err == nil {
		if len(policies.Policies) > 0 {
			return policies.Policies[0], err
		}
	}
	// straight yaml
	policy = api.UpdatePolicy{}
	err = yaml.Unmarshal([]byte(s), &policy)
	if err == nil && !reflect.DeepEqual(policy, api.UpdatePolicy{}) { // empty string unmarshals w/o error
		return policy, nil
	}
	// invalid policy
	return policy, errors.New("policy must be valid JSON or YAML")
}

func updatePolicy(cmd *cobra.Command, args []string) error {
	msg := "unable to update policy"

	// input policy
	policyStr, err := inputPolicy(cmd)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	// parse policy
	updatePolicy, err := parseUpdatePolicy(policyStr)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	// set policy id
	if len(args) != 0 {
		updatePolicy.PolicyID = args[0]
	}

	cli.Log.Debugw("updating policy", "policy", policyStr)
	cli.StartProgress(" Updating policy...")
	updateResponse, err := cli.LwApi.V2.Policy.Update(updatePolicy)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, msg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(updateResponse.Data)
	}
	cli.OutputHuman("The policy %s was updated.\n", updateResponse.Data.PolicyID)
	return nil
}
