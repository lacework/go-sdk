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
	"io"
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
	"gopkg.in/yaml.v3"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/pointer"
	"github.com/lacework/go-sdk/lwseverity"
)

var (
	policyCmdState = struct {
		AlertEnabled  bool
		Enabled       bool
		File          string
		Severity      string
		Tag           string
		State         *bool
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
		"Tags",
	}

	// policyCmd represents the policy parent command
	policyCmd = &cobra.Command{
		Use:     "policy",
		Aliases: []string{"policies"},
		Short:   "Manage policies",
		Long: `Manage policies in your Lacework account.

Policies add annotated metadata to queries for improving the context of alerts,
reports, and information displayed in the Lacework Console.

Policies also facilitate the scheduled execution of Lacework queries.

Queries let you interactively request information from specified
curated datasources. Queries have a defined structure for authoring detections.

Lacework ships a set of default LQL policies that are available in your account.

Limitations:
  * The maximum number of records that each policy will return is 1000
  * The maximum number of API calls is 120 per hour for on-demand LQL query executions

To view all the policies in your Lacework account.

    lacework policy ls

To view more details about a single policy.

    lacework policy show <policy_id>

To view the LQL query associated with the policy, use the query ID.

    lacework query show <query_id>

**Note: LQL syntax may change.**
`,
	}

	// policyListCmd represents the policy list command
	policyListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all policies",
		Long:    `List all registered policies in your Lacework account.`,
		Args:    cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if policyCmdState.Severity != "" &&
				!lwseverity.IsValid(policyCmdState.Severity) {
				return errors.Errorf("the severity %s is not valid, use one of %s",
					policyCmdState.Severity, lwseverity.ValidSeverities.String(),
				)
			}
			return nil
		},
		RunE: listPolicies,
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
		Short:   "Show details about a policy",
		Long:    `Show details about the provided policy ID.`,
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
	policyIDIntRE = regexp.MustCompile(`^(.*-)(\d+)$`)
)

func init() {
	// add the policy command
	rootCmd.AddCommand(policyCmd)

	// add sub-commands to the policy command
	policyCmd.AddCommand(policyListCmd)
	policyCmd.AddCommand(policyListTagsCmd)
	policyCmd.AddCommand(policyShowCmd)

	// Lacework Content Library
	if cli.isLCLInstalled() {
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
func inputPolicy(cmd *cobra.Command, args ...string) (string, error) {
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
		bytes, err := io.ReadAll(os.Stdin)
		return string(bytes), err
	}
	// if running via editor
	action := strings.Split(cmd.Use, " ")[0]

	if action == "create" {
		return inputPolicyFromEditor(action, "")
	}

	policyYaml, err := fetchExistingPolicy(args)
	if err != nil {
		return "", err
	}
	return inputPolicyFromEditor(action, policyYaml)

}

func fetchExistingPolicy(args []string) (string, error) {
	var policyID string

	if len(args) > 0 && len(args[0]) > 0 {
		policyID = args[0]
	} else {
		if err := promptSetPolicyID(&policyID); err != nil {
			return "", err
		}
	}

	cli.StartProgress(fmt.Sprintf("Retrieving policy '%s'...", policyID))
	policy, err := cli.LwApi.V2.Policy.Get(policyID)
	cli.StopProgress()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to retrieve %s", policyID))
	}

	policyYaml, err := yaml.Marshal(policy.Data)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to yaml marshall %s", policyID))
	}
	return string(policyYaml), nil
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
	fileData, err := os.ReadFile(filePath)

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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		err = errors.Wrap(err, msg)
		return
	}
	policy = string(body)
	return
}

func inputPolicyFromEditor(action string, policyYaml string) (policy string, err error) {
	prompt := &survey.Editor{
		Message:       fmt.Sprintf("Use the editor to %s your policy", action),
		FileName:      "policy*.yaml",
		HideDefault:   true,
		AppendDefault: true,
		Default:       policyYaml,
	}

	err = survey.AskOne(prompt, &policy)
	return
}

func promptSetPolicyID(policyID *string) error {
	cli.StartProgress("Retrieving policies...")
	policiesResponse, err := cli.LwApi.V2.Policy.List()
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve policies to select")
	}

	var policyIds []string
	for _, policy := range policiesResponse.Data {
		policyIds = append(policyIds, policy.PolicyID)
	}

	return survey.AskOne(&survey.Select{
		Message: "Select policy to update:",
		Options: policyIds,
	}, policyID)
}

func promptSetPolicyIDs() ([]string, error) {
	cli.StartProgress("Retrieving policies...")
	policiesResponse, err := cli.LwApi.V2.Policy.List()
	cli.StopProgress()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve policies to select")
	}

	var policyIDs []string
	for _, policy := range policiesResponse.Data {
		// exclude manual PolicyTypes and exclude policies whose state matches the new state
		// eg. do not show already enabled policyIDs when running 'lacework policy enable'
		if policy.PolicyType != api.PolicyTypeManual.String() {
			if pointer.CompareBoolPtr(policyCmdState.State, policy.Enabled) {
				continue
			}
			policyIDs = append(policyIDs, policy.PolicyID)
		}
	}

	var response []string
	err = survey.AskOne(&survey.MultiSelect{
		Message: "Select which policy ids to update:",
		Options: policyIDs,
	}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
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
		state := "Disabled"
		if policy.Enabled {
			state = "Enabled"
		}
		alertState := "Disabled"
		if policy.AlertEnabled {
			alertState = "Enabled"
		}
		out = append(out, []string{
			policy.PolicyID,
			policy.Severity,
			policy.Title,
			state,
			alertState,
			policy.EvalFrequency,
			strings.Join(policy.Tags, ", "),
		})
	}
	sortPolicyTable(out, 0)

	return
}

func filterPolicies(policies []api.Policy) []api.Policy {
	newPolicies := []api.Policy{}

	for _, policy := range policies {
		// filter severity if desired
		if lwseverity.ShouldFilter(policy.Severity, policyCmdState.Severity) {
			continue
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

func listPolicies(_ *cobra.Command, _ []string) error {
	cli.Log.Info("listing policies")
	cli.StartProgress("Retrieving policies...")
	policiesResponse, err := cli.LwApi.V2.Policy.List()
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to list policies")
	}

	cli.Log.Infow("total policies", "count", len(policiesResponse.Data))
	policies := filterPolicies(policiesResponse.Data)
	if cli.JSONOutput() {
		return cli.OutputJSON(policies)
	}

	if len(policies) == 0 {
		cli.OutputHuman("There were no policies found.")
	} else {
		cli.OutputHuman(renderCustomTable(policyTableHeaders, policyTable(policies),
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetRowLine(true)
				t.SetColumnSeparator(" ")
				t.SetAutoWrapText(true)
			}),
		))

		if policyCmdState.Tag == "" {
			cli.OutputHuman(
				"\nTry using '--tag <string>' to only show policies with the specified tag.\n",
			)
		} else if policyCmdState.Severity == "" {
			cli.OutputHuman(
				"\nTry using '--severity <string>' to filter policies by severity threshold.\n",
			)
		}
	}
	return nil
}

func showPolicy(_ *cobra.Command, args []string) error {
	var (
		msg            string = "unable to show policy"
		policyResponse api.PolicyResponse
		err            error
	)

	cli.Log.Infow("retrieving policy", "id", args[0])
	cli.StartProgress("Retrieving policy...")
	policyResponse, err = cli.LwApi.V2.Policy.Get(args[0])
	cli.StopProgress()
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
			QueryLanguage: policyResponse.Data.QueryLanguage,
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

	cli.OutputHuman(renderSimpleTable(
		policyTableHeaders, policyTable([]api.Policy{policyResponse.Data}),
	))
	cli.OutputHuman("\n")
	cli.OutputHuman(buildPolicyDetailsTable(policyResponse.Data))
	cli.OutputHuman("\n")
	cli.OutputHuman(renderOneLineCustomTable("DESCRIPTION",
		policyResponse.Data.Description,
		tableFunc(func(t *tablewriter.Table) {
			t.SetAlignment(tablewriter.ALIGN_LEFT)
			t.SetColWidth(120)
			t.SetBorder(false)
			t.SetAutoWrapText(true)
		}),
	))
	cli.OutputHuman("\n")
	cli.OutputHuman(renderOneLineCustomTable("REMEDIATION",
		policyResponse.Data.Remediation,
		tableFunc(func(t *tablewriter.Table) {
			t.SetAlignment(tablewriter.ALIGN_LEFT)
			t.SetColWidth(120)
			t.SetBorder(false)
			t.SetAutoWrapText(true)
			t.SetReflowDuringAutoWrap(false)
		}),
	))
	cli.OutputHuman("\n")

	if policyResponse.Data.QueryID != "" {
		cli.StartProgress("Retrieving query...")
		queryResponse, err := cli.LwApi.V2.Query.Get(policyResponse.Data.QueryID)
		cli.StopProgress()
		if err != nil {
			// something went wrong trying to fetch the LQL query, since this is not
			// the main purpose of this command, we don't error out but instead, log
			// the error and show breadcrumbs to manually fetch the query
			cli.Log.Warnw("unable to get query", "error", err)
			cli.OutputHuman(
				fmt.Sprintf(
					"\nUse 'lacework query show %s' to see the query used by this policy.\n",
					policyResponse.Data.QueryID,
				),
			)
		}
		// we know we are in human-readable format
		if queryResponse.Data.QueryLanguage != nil {
			cli.OutputHuman(renderOneLineCustomTable("QUERY LANGUAGE",
				*queryResponse.Data.QueryLanguage,
				tableFunc(func(t *tablewriter.Table) {
					t.SetAlignment(tablewriter.ALIGN_LEFT)
					t.SetColWidth(120)
					t.SetBorder(false)
					t.SetAutoWrapText(false)
				}),
			))
		}
		cli.OutputHuman(renderOneLineCustomTable("QUERY TEXT",
			queryResponse.Data.QueryText,
			tableFunc(func(t *tablewriter.Table) {
				t.SetAlignment(tablewriter.ALIGN_LEFT)
				t.SetColWidth(120)
				t.SetBorder(false)
				t.SetAutoWrapText(false)
			}),
		))
		cli.OutputHuman("\n")
	}
	return nil
}

func buildPolicyDetailsTable(policy api.Policy) string {
	details := [][]string{
		{"QUERY ID", policy.QueryID},
		{"POLICY TYPE", policy.PolicyType},
		{"LIMIT", fmt.Sprintf("%d", policy.Limit)},
		{"ALERT PROFILE", policy.AlertProfile},
		{"TAGS", strings.Join(policy.Tags, "\n")},
		{"OWNER", policy.Owner},
		{"UPDATED AT", policy.LastUpdateTime},
		{"UPDATED BY", policy.LastUpdateUser},
		{"EVALUATION FREQUENCY", policy.EvalFrequency},
	}
	// Append VALID EXCEPTION CONSTRAINTS to the table
	// Add "None" when ExceptionConfiguration is empty
	exceptionConstraints := strings.Join(
		getPolicyExceptionConstraintsSlice(policy.ExceptionConfiguration), ", ")
	if exceptionConstraints == "" {
		exceptionConstraints = "None"
	}
	entry := []string{"VALID EXCEPTION CONSTRAINTS", exceptionConstraints}
	details = append(details, entry)

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

func getPolicyExceptionConstraintsSlice(exceptionConfiguration map[string][]api.
	PolicyExceptionConfigurationConstraints) []string {
	var exceptionConstraints []string
	constraintFields := exceptionConfiguration["constraintFields"]
	for _, constraint := range constraintFields {
		exceptionConstraints = append(exceptionConstraints, constraint.FieldKey)
	}
	return exceptionConstraints
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

	cli.StartProgress("Retrieving policy tags...")
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

func setPoliciesState(_ *cobra.Command, args []string) error {
	var (
		state = policyCmdState.State
		msg   = "enable"
	)

	if !*state {
		msg = "disable"
	}

	// if tag is provided enable/disable all policies matching the tag
	if policyCmdState.Tag != "" {
		return setPolicyStateByTag(policyCmdState.Tag, *state)
	}

	var (
		bulkPolicies api.BulkUpdatePolicies
		err          error
	)

	if len(args) > 0 {
		for _, policyID := range args {
			bulkPolicies = append(bulkPolicies, api.BulkUpdatePolicy{PolicyID: policyID, Enabled: state})
		}
	}

	// if no arguments are provided enter prompt
	if len(args) == 0 {
		bulkPolicies, err = promptSetPoliciesState()
		if err != nil {
			return err
		}
	}

	if len(bulkPolicies) == 0 {
		cli.OutputHuman("unable to find policies to update\n")
		return nil
	}

	resp, err := cli.LwApi.V2.Policy.UpdateMany(bulkPolicies)
	if err != nil {
		return err
	}
	cli.Log.Debugw("bulk policy updated response:", resp)
	cli.OutputHuman("%d policies have been %sd \n", len(bulkPolicies), msg)

	return nil
}

func promptSetPoliciesState() (api.BulkUpdatePolicies, error) {
	var (
		policyIDs []string
		err       error
	)

	if policyIDs, err = promptSetPolicyIDs(); err != nil {
		return nil, err
	}

	var bulkPolicies api.BulkUpdatePolicies
	for _, policyID := range policyIDs {
		state := true
		bulkPolicies = append(bulkPolicies, api.BulkUpdatePolicy{PolicyID: policyID, Enabled: &state})
	}

	return bulkPolicies, nil
}

func setPolicyStateByTag(tag string, policyState bool) error {
	msg := "disable"
	if policyState {
		msg = "enable"
	}

	cli.StartProgress("Retrieving policies...")
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

	// perform bulk update
	var bulkPolicies api.BulkUpdatePolicies
	for _, p := range matchingPolicies {
		bulkPolicies = append(bulkPolicies, api.BulkUpdatePolicy{
			PolicyID: p.PolicyID,
			Enabled:  &policyState,
		})
	}

	cli.StartProgress(fmt.Sprintf("%sing %d policies...", strings.TrimSuffix(msg, "e"), len(bulkPolicies)))
	resp, err := cli.LwApi.V2.Policy.UpdateMany(bulkPolicies)
	cli.StopProgress()

	if err != nil {
		return errors.Wrapf(err, "failed to complete bulk %s.", msg)
	}

	cli.Log.Debugw("bulk policy updated response:", resp)
	cli.OutputHuman("%d policies tagged with %q have been %sd\n", len(policiesUpdated), tag, msg)
	return nil
}
