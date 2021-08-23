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
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	policyCmdState = struct {
		AlertEnabled bool
		Enabled      bool
		File         string
		Repo         bool
		Severity     string
		URL          string
	}{}

	policyTableHeaders = []string{
		"Policy ID", "Evaluator ID", "Severity", "Title", "State", "Alert State", "Frequency", "Query ID"}

	// policyCmd represents the policy parent command
	policyCmd = &cobra.Command{
		Use:     "policy",
		Aliases: []string{"policies"},
		Short:   "manage policies",
		Long: `Manage policies in your Lacework account.

A policy is a mechanism used to add annotated metadata to a Lacework query for improving
the context of alerts, reports, and information displayed in the Lacework Console.

A policy also facilitates the scheduled execution of a Lacework query

A query is a mechanism used to interactively request information from a specific
curated dataset. A query has a defined structure for authoring detections.

Lacework ships a set of default LQL policies that are available in your account.

Limitations:
  * The maximum number of records that each policy will return is 1000
  * The maximum number of API calls is 120 per hour for ad-hoc LQL query executions

To view all the policies in your Lacework account.

    lacework policy ls

To view more details about a single policy.

    lacework policy show <policy_id>

To view the LQL query associated with the policy, use the query id shown.

    lacework query show <query_id>

** NOTE: LQL syntax may change. **
`,
	}

	// policyListCmd represents the policy list command
	policyListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list policies",
		Long:    `List all the registered policies in your Lacework account.`,
		Args:    cobra.NoArgs,
		RunE:    listPolicies,
	}

	// policyListCmd represents the policy list command
	policyShowCmd = &cobra.Command{
		Use:     "show <policy_id>",
		Aliases: []string{"ls"},
		Short:   "show policy",
		Long:    `Show details about a single policy.`,
		Args:    cobra.ExactArgs(1),
		RunE:    showPolicy,
	}

	// policyDeleteCmd represents the policy delete command
	policyDeleteCmd = &cobra.Command{
		Use:   "delete <policy_id>",
		Short: "delete a policy",
		Long: `Delete a policy by providing the policy id.

Use the command 'lacework policy list' to list the registered policies in
your Lacework account.`,
		Args: cobra.ExactArgs(1),
		RunE: deletePolicy,
	}
)

func init() {
	// add the policy command
	rootCmd.AddCommand(policyCmd)

	// add sub-commands to the policy command
	policyCmd.AddCommand(policyListCmd)
	policyCmd.AddCommand(policyShowCmd)
	policyCmd.AddCommand(policyDeleteCmd)

	// policy list specific flags
	policyListCmd.Flags().StringVar(
		&policyCmdState.Severity,
		"severity", "",
		fmt.Sprintf("filter policies by severity threshold (%s)",
			strings.Join(api.ValidPolicySeverities, ", ")),
	)
	policyListCmd.Flags().BoolVar(
		&policyCmdState.Enabled,
		"enabled", false, "only show enabled policies",
	)
	policyListCmd.Flags().BoolVar(
		&policyCmdState.AlertEnabled,
		"alert_enabled", false, "only show alert_enabled policies",
	)
}

func setPolicySourceFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		if cmd == nil {
			return
		}
		action := strings.Split(cmd.Use, " ")[0]

		// file flag to specify a policy from disk
		cmd.Flags().StringVarP(
			&policyCmdState.File,
			"file", "f", "",
			fmt.Sprintf("path to a policy to %s", action),
		)
		/* repo flag to specify a policy from repo
		cmd.Flags().BoolVarP(
			&policyCmdState.Repo,
			"repo", "r", false,
			fmt.Sprintf("id of a policy to %s via active repo", action),
		)*/
		// url flag to specify a policy from url
		cmd.Flags().StringVarP(
			&policyCmdState.URL,
			"url", "u", "",
			fmt.Sprintf("url to a policy to %s", action),
		)
	}
}

// for commands that take a policy as input
func inputPolicy(cmd *cobra.Command) (string, error) {
	// if running via repo
	if policyCmdState.Repo {
		return inputPolicyFromRepo()
	}
	// if running via file
	if policyCmdState.File != "" {
		return inputPolicyFromFile(policyCmdState.File)
	}
	// if running via URL
	if policyCmdState.URL != "" {
		return inputPolicyFromURL(policyCmdState.URL)
	}
	// if running via stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		return string(bytes), err
	}
	// if running via editor
	action := strings.Split(cmd.Use, " ")[0]
	return inputPolicyFromEditor(action)
}

func inputPolicyFromRepo() (policy string, err error) {
	err = errors.New("NotImplementedError")
	return
}

func inputPolicyFromFile(filePath string) (string, error) {
	fileData, err := ioutil.ReadFile(filePath)

	if err != nil {
		return "", errors.Wrap(err, "unable to read file")
	}

	return string(fileData), nil
}

func inputPolicyFromURL(url string) (policy string, err error) {
	msg := "unable to access URL"

	response, err := http.Get(url)
	if err != nil {
		err = errors.Wrap(err, msg)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errors.Wrap(errors.New(response.Status), msg)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.Wrap(err, msg)
		return
	}
	policy = string(body)
	return
}

func inputPolicyFromEditor(action string) (policy string, err error) {
	prompt := &survey.Editor{
		Message:  fmt.Sprintf("Type a policy to %s", action),
		FileName: "policy*.json",
	}
	err = survey.AskOne(prompt, &policy)

	return
}

func policyTable(policies []api.Policy) (out [][]string) {
	sevThreshold, _ := severityToProperTypes(policyCmdState.Severity)

	for _, policy := range policies {
		// filter severity if desired
		if sevThreshold > 0 {
			policySeverity, _ := severityToProperTypes(policy.Severity)

			if policySeverity > sevThreshold {
				continue
			}
		}
		// filter enabled=false if requesting "enabled-only"
		if policyCmdState.Enabled && !policy.Enabled {
			continue
		}
		// filter alert_enabled=false if requesting "alert_enabled-only"
		if policyCmdState.AlertEnabled && !policy.AlertEnabled {
			continue
		}
		state := "disabled"
		if policy.Enabled {
			state = "enabled"
		}
		alertState := "disabled"
		if policy.AlertEnabled {
			alertState = "enabled"
		}
		out = append(out, []string{
			policy.PolicyID,
			policy.EvaluatorID,
			policy.Severity,
			policy.Title,
			state,
			alertState,
			policy.EvalFrequency,
			policy.QueryID,
		})
	}
	return
}

func listPolicies(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("listing policies")

	if policyCmdState.Severity != "" && !array.ContainsStr(
		api.ValidPolicySeverities, policyCmdState.Severity) {
		return errors.Wrap(
			errors.New(fmt.Sprintf("the severity %s is not valid, use one of %s",
				policyCmdState.Severity, strings.Join(api.ValidPolicySeverities, ", "))),
			"unable to list policies",
		)
	}

	cli.StartProgress(" Retrieving policies...")
	policyResponse, err := cli.LwApi.V2.Policy.List()
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to list policies")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(policyResponse.Data)
	}
	if len(policyResponse.Data) == 0 {
		cli.OutputHuman("There were no policies found.")
		return nil
	}
	cli.OutputHuman(renderSimpleTable(policyTableHeaders, policyTable(policyResponse.Data)))
	return nil
}

func showPolicy(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("retrieving policy", "policyID", args[0])
	cli.StartProgress(" Retrieving policy...")
	policyResponse, err := cli.LwApi.V2.Policy.Get(args[0])
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to show policy")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(policyResponse.Data)
	}
	cli.OutputHuman(
		renderSimpleTable(policyTableHeaders, policyTable([]api.Policy{policyResponse.Data})))
	cli.OutputHuman("\n")
	cli.OutputHuman(buildPolicyDetailsTable(policyResponse.Data))
	return nil
}

func deletePolicy(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("deleting policy", "policyID", args[0])
	cli.StartProgress(" Deleting policy...")
	deleted, err := cli.LwApi.V2.Policy.Delete(args[0])
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to delete policy")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(deleted)
	}

	cli.OutputHuman(
		fmt.Sprintf("The policy %s was deleted.\n", args[0]))
	return nil
}

func buildPolicyDetailsTable(policy api.Policy) string {
	details := [][]string{
		{"DESCRIPTION", policy.Description},
		{"REMEDIATION", policy.Remediation},
		{"POLICY TYPE", policy.PolicyType},
		{"LIMIT", fmt.Sprintf("%d", policy.Limit)},
		{"ALERT PROFILE", policy.AlertProfile},
		{"OWNER", policy.Owner},
		{"UPDATED AT", policy.LastUpdateTime},
		{"UPDATED BY", policy.LastUpdateUser},
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
