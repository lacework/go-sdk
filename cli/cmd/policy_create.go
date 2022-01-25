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
	// policyCreateCmd represents the policy create command
	policyCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a policy",
		Long: `Create a policy.

A policy is represented in either JSON or YAML format.

The following attributes are minimally required:

    ---
    evaluatorId: Cloudtrail
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

	if IsLCLInstalled(*cli.LwComponents) {
		policyCreateCmd.Flags().StringVarP(
			&policyCmdState.CreateFromLibrary,
			"library", "l", "",
			"create policy from Lacework Content Library",
		)
	}
}

func createPolicy(cmd *cobra.Command, _ []string) error {
	var (
		msg       string = "unable to create policy"
		lcl       *LaceworkContentLibrary
		err       error
		newPolicy api.NewPolicy
		policyStr string
	)

	if policyCmdState.CreateFromLibrary != "" {
		if lcl, err = LoadLCL(*cli.LwComponents); err == nil {
			newPolicy, err = lcl.GetNewPolicy(policyCmdState.CreateFromLibrary)
		}
	} else {
		// input policy
		policyStr, err = inputPolicy(cmd)
		if err != nil {
			return errors.Wrap(err, msg)
		}
		cli.Log.Debugw("creating policy", "policy", policyStr)
		// parse policy
		newPolicy, err = api.ParseNewPolicy(policyStr)
	}

	if err != nil {
		return errors.Wrap(err, msg)
	}

	cli.StartProgress(" Creating policy...")
	createResponse, err := cli.LwApi.V2.Policy.Create(newPolicy)
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, msg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(createResponse.Data)
	}
	cli.OutputHuman(fmt.Sprintf("The policy %s was created.\n", createResponse.Data.PolicyID))
	return nil
}
