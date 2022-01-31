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

	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	policyLibraryTableHeaders []string = []string{"Policy ID", "Title", "Query ID"}

	policyListLibraryCmd = &cobra.Command{
		Use:   "list-library",
		Short: "List policies from library",
		Long:  `List all LQL policies in your Lacework Content Library.`,
		Args:  cobra.NoArgs,
		RunE:  listPoliciesLibrary,
	}
	policyShowLibraryCmd = &cobra.Command{
		Use:   "show-library <policy_id>",
		Short: "Show a policy from library",
		Long:  `Show a policy in your Lacework Content Library.`,
		Args:  cobra.ExactArgs(1),
		RunE:  showPolicyLibrary,
	}
)

func init() {
	policyCmd.AddCommand(policyListLibraryCmd)
	policyCmd.AddCommand(policyShowLibraryCmd)
}

func policyLibraryTable(policies map[string]LCLPolicy) (out [][]string) {
	for _, policy := range policies {
		out = append(out, []string{
			policy.PolicyID,
			policy.Title,
			policy.QueryID,
		})
	}
	return
}

func buildPolicyLibraryDetailsTable(policy api.Policy) string {
	details := [][]string{
		{"DESCRIPTION", policy.Description},
		{"REMEDIATION", policy.Remediation},
		{"POLICY TYPE", policy.PolicyType},
		{"LIMIT", fmt.Sprintf("%d", policy.Limit)},
		{"ALERT PROFILE", policy.AlertProfile},
		{"EVALUATION FREQUENCY", policy.EvalFrequency},
	}

	return renderOneLineCustomTable("POLICY DETAILS",
		renderCustomTable([]string{}, details,
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
				t.SetAutoWrapText(false)
				t.SetAlignment(tablewriter.ALIGN_LEFT)
			}),
		),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		}),
	)
}

func listPoliciesLibrary(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("listing policies from library")

	cli.StartProgress(" Retrieving policies...")
	lcl, err := LoadLCL(*cli.LwComponents)
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, "unable to list policies")
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(lcl.Policies)
	}
	if len(lcl.Policies) == 0 {
		cli.OutputHuman("There were no policies found.")
		return nil
	}
	cli.OutputHuman(
		renderSimpleTable(
			policyLibraryTableHeaders,
			policyLibraryTable(lcl.Policies),
		),
	)
	return nil
}

func showPolicyLibrary(_ *cobra.Command, args []string) error {
	var (
		msg            string = "unable to show policy"
		policyString   string
		newPolicy      api.NewPolicy
		policyResponse api.PolicyResponse
		err            error
	)
	cli.Log.Debugw("retrieving policy", "id", args[0])

	cli.StartProgress(" Retrieving policy...")
	// input policy
	if policyString, err = inputPolicyFromLibrary(args[0]); err != nil {
		cli.StopProgress()
		return errors.Wrap(err, msg)
	}
	// parse policy
	newPolicy, err = api.ParseNewPolicy(policyString)
	policyResponse.Data = api.Policy{
		EvaluatorID:   newPolicy.EvaluatorID,
		PolicyID:      newPolicy.PolicyID,
		PolicyType:    newPolicy.PolicyType,
		QueryID:       newPolicy.QueryID,
		Title:         newPolicy.Title,
		Enabled:       newPolicy.Enabled,
		Description:   newPolicy.Description,
		Remediation:   newPolicy.Remediation,
		Severity:      newPolicy.Severity,
		Limit:         newPolicy.Limit,
		EvalFrequency: newPolicy.EvalFrequency,
		AlertEnabled:  newPolicy.AlertEnabled,
		AlertProfile:  newPolicy.AlertProfile,
	}
	cli.StopProgress()

	// output policy
	if err != nil {
		return errors.Wrap(err, msg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(newPolicy)
	}
	cli.OutputHuman(
		renderSimpleTable(policyTableHeaders, policyTable([]api.Policy{policyResponse.Data})))
	cli.OutputHuman("\n")
	cli.OutputHuman(buildPolicyLibraryDetailsTable(policyResponse.Data))
	return nil
}
