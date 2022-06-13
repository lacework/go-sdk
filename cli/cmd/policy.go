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
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
)

var (
	policyCmdState = struct {
		AlertEnabled  bool
		Enabled       bool
		File          string
		Severity      string
		Tag           string
		URL           string
		CUFromLibrary string
		CascadeDelete bool
	}{}

	policyTableHeaders = []string{
		"Policy ID",
		"Severity",
		"Title",
		"State",
		"Alert State",
		"Frequency",
		"Query ID",
		"Tags",
	}

	// policyCmd represents the policy parent command
	policyCmd = &cobra.Command{
		Use:     "policy",
		Aliases: []string{"policies"},
		Short:   "Manage policies",
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

**NOTE: LQL syntax may change.**
`,
	}

	// policyListCmd represents the policy list command
	policyListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List policies",
		Long:    `List all the registered policies in your Lacework account.`,
		Args:    cobra.NoArgs,
		RunE:    listPolicies,
	}
	// policyListTagsCmd represents the policy list command
	policyListTagsCmd = &cobra.Command{
		Use:     "list-tags",
		Aliases: []string{"ls"},
		Short:   "List policy tags",
		Long:    `List all tags associated with policies in your Lacework account.`,
		Args:    cobra.NoArgs,
		RunE:    listPolicyTags,
	}
	// policyShowCmd represents the policy show command
	policyShowCmd = &cobra.Command{
		Use:     "show <policy_id>",
		Aliases: []string{"ls"},
		Short:   "Show policy",
		Long:    `Show details about a single policy.`,
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			b, err := cmd.Flags().GetBool("yaml")
			if err != nil {
				return errors.Wrap(err, "unable to parse --yaml flag")
			}
			if b {
				cli.EnableYAMLOutput()
			}
			return nil
		},
		RunE: showPolicy,
	}

	// This is an experimental command.
	// policyDisableTagCmd represents the policy disable command
	policyDisableTagCmd = &cobra.Command{
		Use:   "disable [policy_id]",
		Short: "Disable Policies",
		Long: `Disable Policies by ID or all policies matching a tag.

To disable a single policy by it's ID:

	lacework policy disable lacework-policy-id

To disable all policies for Aws CIS 1.4.0:

	lacework policy disable --tag framework:cis-aws-1-4-0

To disable all policies for Gcp CIS 1.3.0:

	lacework policy disable --tag framework:cis-gcp-1-3-0

`,
		Args: cobra.RangeArgs(0, 1),
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 && policyCmdState.Tag != "" {
				return errors.New("'--tag' flag may not be use in conjunction with 'policy_id' arg")
			}
			return nil
		},
		RunE: disablePolicy,
	}

	// This is an experimental command.
	// policyEnableTagCmd represents the policy enable command
	policyEnableTagCmd = &cobra.Command{
		Use:   "enable [policy_id]",
		Short: "Enable Policies",
		Long: `Enable Policies by ID or all policies matching a tag.

To enable a single policy by it's ID:

	lacework policy enable lacework-policy-id

To enable all policies for Aws CIS 1.4.0:

	lacework policy enable --tag framework:cis-aws-1-4-0

To enable all policies for Gcp CIS 1.3.0:

	lacework policy enable --tag framework:cis-gcp-1-3-0

`,
		Args: cobra.RangeArgs(0, 1),
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 && policyCmdState.Tag != "" {
				return errors.New("'--tag' flag may not be use in conjunction with 'policy_id' arg")
			}
			return nil
		},
		RunE: enablePolicy,
	}

	policyIDIntRE = regexp.MustCompile(`^(.*-)(\d+)$`)
)

func init() {
	// add the policy command
	rootCmd.AddCommand(policyCmd)

	// add sub-commands to the policy command
	policyCmd.AddCommand(policyListCmd)
	policyCmd.AddCommand(policyListTagsCmd)
	policyCmd.AddCommand(policyShowCmd)
	// experimental commands
	policyCmd.AddCommand(policyDisableTagCmd)
	policyCmd.AddCommand(policyEnableTagCmd)

	// Lacework Content Library
	if cli.IsLCLInstalled() {
		policyCreateCmd.Flags().StringVarP(
			&policyCmdState.CUFromLibrary,
			"library", "l", "",
			"create policy from Lacework Content Library",
		)
		policyUpdateCmd.Flags().StringVarP(
			&policyCmdState.CUFromLibrary,
			"library", "l", "",
			"update policy from Lacework Content Library",
		)
	}

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
	policyListCmd.Flags().StringVar(
		&policyCmdState.Tag,
		"tag", "", "only show policies with the specified tag",
	)
	// policy show specific flags
	policyShowCmd.Flags().Bool(
		"yaml", false, "output query in YAML format",
	)
	// policy disable specific flags
	policyDisableTagCmd.Flags().StringVar(
		&policyCmdState.Tag,
		"tag", "", "disable all policies with the specified tag",
	)
	// policy enable specific flags
	policyEnableTagCmd.Flags().StringVar(
		&policyCmdState.Tag,
		"tag", "", "enable all policies with the specified tag",
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
	// if running via library (CU)
	if policyCmdState.CUFromLibrary != "" {
		return inputPolicyFromLibrary(policyCmdState.CUFromLibrary)
	}
	// if running via file
	if policyCmdState.File != "" {
		return inputPolicyFromFile(policyCmdState.File)
	}
	// if running via URL
	if policyCmdState.URL != "" {
		return inputPolicyFromURL(policyCmdState.URL)
	}
	stat, err := os.Stdin.Stat()
	if err != nil {
		cli.Log.Debugw("error retrieving stdin mode", "error", err.Error())
	} else if (stat.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		return string(bytes), err
	}
	// if running via editor
	action := strings.Split(cmd.Use, " ")[0]
	return inputPolicyFromEditor(action)
}

func inputPolicyFromLibrary(id string) (string, error) {
	var (
		lcl *LaceworkContentLibrary
		err error
	)

	if lcl, err = cli.LoadLCL(); err != nil {
		return "", err
	}
	return lcl.GetPolicy(id)
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

func sortPolicyTable(out [][]string, policyIDIndex int) {
	// order by ID (special handling for policy ID numbers)
	sort.Slice(out, func(i, j int) bool {
		iMatch := policyIDIntRE.FindStringSubmatch(out[i][policyIDIndex])
		jMatch := policyIDIntRE.FindStringSubmatch(out[j][policyIDIndex])
		// both regexes must match
		// both regexes must have proper lengths since we'll be using...
		// ...direct access from here on out
		if iMatch == nil || jMatch == nil || len(iMatch) != 3 || len(jMatch) != 3 {
			return out[i][policyIDIndex] < out[j][policyIDIndex]
		}
		// if string portions aren't the same
		if iMatch[1] != jMatch[1] {
			return out[i][policyIDIndex] < out[j][policyIDIndex]
		}
		// if string portions are the same; compare based on ints
		// no error checking needed for Atoi since use regexp \d+
		iNum, _ := strconv.Atoi(iMatch[2])
		jNum, _ := strconv.Atoi(jMatch[2])
		return iNum < jNum
	})
}

func policyTable(policies []api.Policy) (out [][]string) {
	for _, policy := range policies {
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
			policy.Severity,
			policy.Title,
			state,
			alertState,
			policy.EvalFrequency,
			policy.QueryID,
			strings.Join(policy.Tags, "\n"),
		})
	}
	sortPolicyTable(out, 0)

	return
}

func filterPolicies(policies []api.Policy) []api.Policy {
	newPolicies := []api.Policy{}
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
		// filter tag
		if policyCmdState.Tag != "" && !policy.HasTag(policyCmdState.Tag) {
			continue
		}
		newPolicies = append(newPolicies, policy)
	}
	return newPolicies
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
	policiesResponse, err := cli.LwApi.V2.Policy.List()
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, "unable to list policies")
	}

	policies := filterPolicies(policiesResponse.Data)
	if cli.JSONOutput() {
		return cli.OutputJSON(policies)
	}
	if len(policies) == 0 {
		cli.OutputHuman("There were no policies found.")
		return nil
	}
	cli.OutputHuman(renderSimpleTable(policyTableHeaders, policyTable(policies)))
	return nil
}

func showPolicy(cmd *cobra.Command, args []string) error {
	var (
		msg            string = "unable to show policy"
		policyResponse api.PolicyResponse
		err            error
	)

	cli.Log.Debugw("retrieving policy", "policyID", args[0])
	cli.StartProgress(" Retrieving policy...")
	policyResponse, err = cli.LwApi.V2.Policy.Get(args[0])
	cli.StopProgress()

	// output policy
	if err != nil {
		return errors.Wrap(err, msg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(policyResponse.Data)
	}

	if cli.YAMLOutput() {
		return cli.OutputYAML(&api.NewPolicy{
			PolicyID:      policyResponse.Data.PolicyID,
			PolicyType:    policyResponse.Data.PolicyType,
			QueryID:       policyResponse.Data.QueryID,
			Title:         policyResponse.Data.Title,
			Enabled:       policyResponse.Data.Enabled,
			Description:   policyResponse.Data.Description,
			Remediation:   policyResponse.Data.Remediation,
			Severity:      policyResponse.Data.Severity,
			Limit:         policyResponse.Data.Limit,
			EvalFrequency: policyResponse.Data.EvalFrequency,
			AlertEnabled:  policyResponse.Data.AlertEnabled,
			AlertProfile:  policyResponse.Data.AlertProfile,
			Tags:          policyResponse.Data.Tags,
		})
	}

	cli.OutputHuman(
		renderSimpleTable(policyTableHeaders, policyTable([]api.Policy{policyResponse.Data})))
	cli.OutputHuman("\n")
	cli.OutputHuman(buildPolicyDetailsTable(policyResponse.Data))
	return nil
}

func buildPolicyDetailsTable(policy api.Policy) string {
	details := [][]string{
		{"DESCRIPTION", policy.Description},
		{"REMEDIATION", policy.Remediation},
		{"POLICY TYPE", policy.PolicyType},
		{"LIMIT", fmt.Sprintf("%d", policy.Limit)},
		{"ALERT PROFILE", policy.AlertProfile},
		{"TAGS", strings.Join(policy.Tags, "\n")},
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

func policyTagsTable(pt []string) (out [][]string) {
	for _, tag := range pt {
		out = append(out, []string{tag})
	}

	// order by Tag
	sort.Slice(out, func(i, j int) bool {
		return out[i][0] < out[j][0]
	})

	return
}

func listPolicyTags(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("listing policy tags")

	cli.StartProgress(" Retrieving policy tags...")
	policyTagsResponse, err := cli.LwApi.V2.Policy.ListTags()
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to list policy tags")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(policyTagsResponse.Data)
	}
	if len(policyTagsResponse.Data) == 0 {
		cli.OutputHuman("There were no policy tags found.\n")
		return nil
	}
	cli.OutputHuman(renderSimpleTable([]string{"Tag"}, policyTagsTable(policyTagsResponse.Data)))
	return nil
}

func disablePolicy(_ *cobra.Command, args []string) error {
	state := false
	if len(args) > 0 {
		policy := api.UpdatePolicy{PolicyID: args[0], Enabled: &state}
		_, err := cli.LwApi.V2.Policy.Update(policy)
		if err != nil {
			return err
		}
	}

	// if tag is provided disable all policies matching
	if policyCmdState.Tag != "" {
		return setPolicyStateByTag(policyCmdState.Tag, state)
	}
	return nil
}

func enablePolicy(_ *cobra.Command, args []string) error {
	state := true
	if len(args) > 0 {
		policy := api.UpdatePolicy{PolicyID: args[0], Enabled: &state}
		_, err := cli.LwApi.V2.Policy.Update(policy)
		if err != nil {
			return err
		}
	}

	// if tag is provided enable all policies matching
	if policyCmdState.Tag != "" {
		return setPolicyStateByTag(policyCmdState.Tag, state)
	}
	return nil
}

func setPolicyStateByTag(tag string, policyState bool) error {
	msg := "disable"
	if policyState {
		msg = "enable"
	}

	cli.StartProgress(" Retrieving policies...")
	policyTagsResponse, err := cli.LwApi.V2.Policy.List()
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, "unable to list policies")
	}

	var (
		policiesUpdated  []string
		matchingPolicies []api.Policy
	)

	for _, p := range policyTagsResponse.Data {
		if p.HasTag(tag) && p.Enabled != policyState {
			matchingPolicies = append(matchingPolicies, p)
		}
	}

	if len(matchingPolicies) == 0 {
		cli.OutputHuman("No policies found with tag '%s'\n", tag)
		return nil
	}

	for i, p := range matchingPolicies {
		cli.StartProgress(fmt.Sprintf(" %sing policies %d/%d (%s)...", strings.TrimSuffix(msg, "e"), i+1, len(matchingPolicies), p.PolicyID))
		policy := api.UpdatePolicy{PolicyID: p.PolicyID, Enabled: &policyState}
		resp, err := cli.LwApi.V2.Policy.Update(policy)
		if err != nil {
			if len(policiesUpdated) > 0 {
				return errors.Wrapf(err, "failed to complete bulk %s. %d policies have been %sd: %s",
					msg, len(policiesUpdated), msg, strings.Join(policiesUpdated, ","))
			}
			return errors.Wrapf(err, "failed to %s any policies", msg)

		}
		policiesUpdated = append(policiesUpdated, resp.Data.PolicyID)
		cli.StopProgress()

	}
	cli.OutputHuman("%d policies tagged with %q have been %sd\n", len(policiesUpdated), tag, msg)
	return nil
}
