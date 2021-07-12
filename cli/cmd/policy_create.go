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
	"fmt"
	"reflect"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	// policyCreateCmd represents the policy create command
	policyCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create a policy",
		Long: `Create a policy.

A policy is represented in either JSON or YAML format.
The following attributes are minimally required:
---
evaluatorId: Cloudtrail
policyId: lacework-example-1
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
policyUi:
  domain: AWS
  subdomain: Cloudtrail
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

type newPoliciesYAML struct {
	Policies []api.NewPolicy `yaml:"policies"`
}

func parseNewPolicy(s string) (api.NewPolicy, error) {
	var policy api.NewPolicy
	var err error

	// valid json
	if err = json.Unmarshal([]byte(s), &policy); err == nil {
		return policy, err
	}
	// nested yaml
	var policies newPoliciesYAML

	if err = yaml.Unmarshal([]byte(s), &policies); err == nil {
		if len(policies.Policies) > 0 {
			return policies.Policies[0], err
		}
	}
	// straight yaml
	policy = api.NewPolicy{}
	err = yaml.Unmarshal([]byte(s), &policy)
	if err == nil && !reflect.DeepEqual(policy, api.NewPolicy{}) { // empty string unmarshals w/o error
		return policy, nil
	}
	// invalid policy
	return policy, errors.New("policy must be valid JSON or YAML")
}

func createPolicy(cmd *cobra.Command, _ []string) error {
	msg := "unable to create policy"

	// input policy
	policyStr, err := inputPolicy(cmd)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	// parse policy
	newPolicy, err := parseNewPolicy(policyStr)
	if err != nil {
		return errors.Wrap(err, msg)
	}

	cli.Log.Debugw("creating policy", "policy", policyStr)

	var createResponse api.PolicyResponse
	if createResponse, err = cli.LwApi.V2.Policy.Create(newPolicy); err != nil {
		return errors.Wrap(err, msg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(createResponse.Data)
	}
	cli.OutputHuman(fmt.Sprintf("Policy (%s) created successfully.\n", createResponse.Data.PolicyID))
	return nil
}
