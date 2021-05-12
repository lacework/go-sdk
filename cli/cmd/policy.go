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
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
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

	// policyCmd represents the policy parent command
	policyCmd = &cobra.Command{
		Hidden: true,
		Use:    "policy",
		Short:  "manage policies",
		Long: `Manage policies.

NOTE: This feature is not yet available!`,
	}

	// policyCreateCmd represents the policy create command
	policyCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create a policy",
		Long: `Create a policy.

A policy is represented in JSON format.
The following attributes are minimally required:
{
    "policy_id": "lacework-example-1",
    "title": "My Policy",
    "enabled": false,
    "lql_id": "MyQuery",
    "severity": "high",
    "description": "My Policy Description",
    "remediation": "My Policy Remediation"
}`,
		Args: cobra.NoArgs,
		RunE: createPolicy,
	}

	// policyListCmd represents the policy list command
	policyListCmd = &cobra.Command{
		Use:   "list",
		Short: "list policies",
		Long:  `List policies.`,
		Args:  cobra.NoArgs,
		RunE:  listPolicies,
	}

	// policyListCmd represents the policy list command
	policyShowCmd = &cobra.Command{
		Use:   "show <policy_id>",
		Short: "show policy",
		Long:  `Show policy.`,
		Args:  cobra.ExactArgs(1),
		RunE:  showPolicy,
	}

	// policyDeleteCmd represents the policy delete command
	policyDeleteCmd = &cobra.Command{
		Use:   "delete <policy_id>",
		Short: "delete policy",
		Long:  `Delete policy.`,
		Args:  cobra.ExactArgs(1),
		RunE:  deletePolicy,
	}
)

func init() {
	// add the policy command
	rootCmd.AddCommand(policyCmd)

	// add sub-commands to the policy command
	policyCmd.AddCommand(policyCreateCmd)
	policyCmd.AddCommand(policyListCmd)
	policyCmd.AddCommand(policyShowCmd)
	policyCmd.AddCommand(policyDeleteCmd)

	// policy source specific flags
	setPolicySourceFlags(policyCreateCmd)

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

func createPolicy(cmd *cobra.Command, _ []string) error {
	policy, err := inputPolicy(cmd)
	if err != nil {
		return errors.Wrap(err, "unable to create policy")
	}

	cli.Log.Debugw("creating policy", "policy", policy)

	var create api.PolicyResponse
	if create, err = cli.LwApi.Policy.Create(policy); err != nil {
		return errors.Wrap(err, "unable to create policy")

	}

	if cli.JSONOutput() {
		return cli.OutputJSON(create.Data)
	}
	policyID := ""
	if len(create.Data) > 0 {
		policyID = create.Data[0].ID
	}
	cli.OutputHuman(fmt.Sprintf("Policy (%s) created successfully.\n", policyID))
	return nil
}

var policyTableHeaders = []string{"Policy ID", "Severity", "Title", "Enabled", "Alert Enabled", "Frequency", "Query ID"}

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
		out = append(out, []string{
			policy.ID,
			policy.Severity,
			policy.Title,
			strconv.FormatBool(policy.Enabled),
			strconv.FormatBool(policy.AlertEnabled),
			policy.Frequency,
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

	policyResponse, err := cli.LwApi.Policy.GetAll()
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

	policyResponse, err := cli.LwApi.Policy.GetByID(args[0])
	if err != nil {
		return errors.Wrap(err, "unable to show policy")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(policyResponse.Data)
	}
	if len(policyResponse.Data) == 0 {
		return yikes("unable to show policy")
	}
	cli.OutputHuman(renderSimpleTable(policyTableHeaders, policyTable(policyResponse.Data)))
	return nil
}

func deletePolicy(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("deleting policy", "policyID", args[0])

	delete, err := cli.LwApi.Policy.Delete(args[0])
	if err != nil {
		return errors.Wrap(err, "unable to delete policy")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(delete)
	}
	cli.OutputHuman(
		fmt.Sprintf("Policy (%s) deleted successfully.\n", args[0]))
	return nil
}
