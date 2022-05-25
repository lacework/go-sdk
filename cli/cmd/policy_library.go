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

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	policyLibraryTableHeaders []string = []string{"Policy ID", "Title", "Query ID", "TAGS"}

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
	policySyncLibraryCmd = &cobra.Command{
		Use:   "sync-library",
		Short: "Synchronize library policies",
		Long: `Synchronize library policies with your Lacework account.
		
Requirements:
Specify the --tag flag to select policies and queries.

Behavior:
1. Policies and queries that exist in the library but not in your account will be created.
2. Policies and queries that exist in both the library and your account will be updated.
3. Nothing will be deleted.

To view all policies in the library and their associated tags.

    lacework policy list-library`,
		Args: cobra.NoArgs,
		RunE: syncPolicyLibrary,
	}
)

func init() {
	if !cli.IsLCLInstalled() {
		return
	}

	policyCmd.AddCommand(policyListLibraryCmd)
	policyListLibraryCmd.Flags().StringVar(
		&policyCmdState.Tag,
		"tag", "", "only show policies with the specified tag",
	)

	policyCmd.AddCommand(policyShowLibraryCmd)

	policyCmd.AddCommand(policySyncLibraryCmd)
	policySyncLibraryCmd.Flags().StringVar(
		&policyCmdState.Tag,
		"tag", "", "sync policies and queries with the specified tag",
	)
}

func policyLibraryTable(policies map[string]LCLPolicy) (out [][]string) {
	for _, policy := range policies {
		out = append(out, []string{
			policy.PolicyID,
			policy.Title,
			policy.QueryID,
			strings.Join(policy.Tags, "\n"),
		})
	}
	sortPolicyTable(out, 0)
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
	lcl, err := cli.LoadLCL()
	cli.StopProgress()

	var policies map[string]LCLPolicy
	if policyCmdState.Tag == "" {
		policies = lcl.Policies
	} else {
		policies = lcl.GetPoliciesByTag(policyCmdState.Tag)
	}

	if err != nil {
		return errors.Wrap(err, "unable to list policies")
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(policies)
	}
	if len(policies) == 0 {
		cli.OutputHuman("There were no policies found.")
		return nil
	}
	cli.OutputHuman(
		renderCustomTable(
			policyLibraryTableHeaders,
			policyLibraryTable(policies),
			tableFunc(func(t *tablewriter.Table) {
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
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

type PolicySyncOperation struct {
	ID          string
	ContentType string
	Operation   string
}

func getPolicySyncOperations(policies map[string]LCLPolicy) ([]PolicySyncOperation, error) {
	policyOps := []PolicySyncOperation{}

	cli.StartProgress(" Retrieving platform policies...")
	policiesResponse, err := cli.LwApi.V2.Policy.List()
	cli.StopProgress()

	if err != nil {
		return policyOps, err
	}
	var platformPolicyIDs = make([]string, len(policiesResponse.Data))
	for i, policy := range policiesResponse.Data {
		platformPolicyIDs[i] = policy.PolicyID
	}

	cli.StartProgress(" Retrieving platform queires...")
	queryResponse, err := cli.LwApi.V2.Query.List()
	cli.StopProgress()

	if err != nil {
		return policyOps, err
	}
	var platformQueryIDs = make([]string, len(queryResponse.Data))
	for i, query := range queryResponse.Data {
		platformQueryIDs[i] = query.QueryID
	}

	for _, lclPolicy := range policies {
		qso := PolicySyncOperation{
			ID:          lclPolicy.QueryID,
			ContentType: "query",
			Operation:   "create",
		}
		if array.ContainsStr(platformQueryIDs, lclPolicy.QueryID) {
			qso.Operation = "update"
		}
		policyOps = append(policyOps, qso)

		pso := PolicySyncOperation{
			ID:          lclPolicy.PolicyID,
			ContentType: "policy",
			Operation:   "create",
		}
		suf := fmt.Sprintf("-%s", lclPolicy.PolicyID)
		for _, platformPolicyID := range platformPolicyIDs {
			// TODO: proper handling for $account based ids
			if platformPolicyID == lclPolicy.PolicyID || strings.HasSuffix(platformPolicyID, suf) {
				pso.Operation = "update"
				break
			}
		}
		policyOps = append(policyOps, pso)
	}

	return policyOps, nil
}

func policySyncOpsDetails(psos []PolicySyncOperation) string {
	var (
		ops             = []string{"Operation details:"}
		detailTemplate  = "  %s %s will be %sd."
		validOperations = []string{"policy-create", "policy-update", "query-create", "query-update"}
		caser           = cases.Title(language.Und)
	)

	for _, pso := range psos {
		key := fmt.Sprintf("%s-%s", pso.ContentType, pso.Operation)

		if !array.ContainsStr(validOperations, key) {
			continue
		}

		ops = append(ops, fmt.Sprintf(
			detailTemplate,
			caser.String(pso.ContentType),
			pso.ID,
			pso.Operation,
		))
	}

	return strings.Join(ops, "\n") + "\n\n"
}

func policySyncOpsSummary(psos []PolicySyncOperation) string {
	var (
		msg = "Policy sync-library will create %d policies, update %d policies, create %d queries, update %d queries."
		obs = map[string]int{
			"policy-create": 0,
			"policy-update": 0,
			"query-create":  0,
			"query-update":  0,
		}
	)

	for _, pso := range psos {
		key := fmt.Sprintf("%s-%s", pso.ContentType, pso.Operation)
		if v, ok := obs[key]; ok {
			obs[key] = v + 1
		}
	}

	return fmt.Sprintf(
		msg,
		obs["policy-create"],
		obs["policy-update"],
		obs["query-create"],
		obs["query-update"],
	)
}

// Simple helper to prompt for next steps after TF plan
func policySyncPrompt(psos []PolicySyncOperation, previewShown *bool) (int, error) {
	options := []string{
		"Apply sync",
	}

	// Omit option to show details if we already have
	if !*previewShown {
		options = append(options, "Show details")
	}
	options = append(options, "Quit")

	var answer int
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Select{
			Message: policySyncOpsSummary(psos),
			Options: options,
		},
		Response: &answer,
	})

	return answer, err
}

func policySyncDisplayChanges(psos []PolicySyncOperation) (bool, error) {
	// Prompt for next steps
	prompt := true
	previewShown := false
	var answer int

	// Displaying resources
	for prompt {
		id, err := policySyncPrompt(psos, &previewShown)
		if err != nil {
			return false, err
		}

		switch {
		case id == 1 && !previewShown:
			cli.OutputHuman(policySyncOpsDetails(psos))
		default:
			answer = id
			prompt = false
		}

		if id == 1 && !previewShown {
			previewShown = true
		}
	}

	// Run apply
	if answer == 0 {
		return true, nil
	}

	// Quit
	return false, nil
}

func policySyncExecuteChanges(lcl *LaceworkContentLibrary, psos []PolicySyncOperation) error {
	for _, pso := range psos {
		msg := fmt.Sprintf("unable to %s %s", pso.Operation, pso.ContentType)

		if pso.ContentType == "query" {
			// input query
			queryString, err := lcl.GetQuery(pso.ID)
			if err != nil {
				return errors.Wrap(err, msg)
			}

			if pso.Operation == "create" {
				cli.Log.Debugw("creating query", "query", queryString)

				// parse query
				newQuery, err := api.ParseNewQuery(queryString)
				if err != nil {
					return errors.Wrap(queryErrorCrumbs(queryString), msg)
				}
				cli.StartProgress(" Creating query...")
				create, err := cli.LwApi.V2.Query.Create(newQuery)
				cli.StopProgress()

				// output
				if err != nil {
					return errors.Wrap(err, msg)
				}
				cli.OutputHuman("The query %s was created.\n", create.Data.QueryID)
			}

			if pso.Operation == "update" {
				cli.Log.Debugw("updating query", "query", queryString)

				// parse query
				newQuery, err := api.ParseNewQuery(queryString)
				if err != nil {
					return errors.Wrap(queryErrorCrumbs(queryString), msg)
				}
				updateQuery := api.UpdateQuery{
					QueryText: newQuery.QueryText,
				}

				// update query
				cli.StartProgress(" Updating query...")
				update, err := cli.LwApi.V2.Query.Update(newQuery.QueryID, updateQuery)
				cli.StopProgress()

				// output
				if err != nil {
					return errors.Wrap(err, msg)
				}
				cli.OutputHuman("The query %s was updated.\n", update.Data.QueryID)
			}
		}

		if pso.ContentType == "policy" {
			// input policy
			policyString, err := lcl.GetPolicy(pso.ID)
			if err != nil {
				return errors.Wrap(err, msg)
			}

			if pso.Operation == "create" {
				cli.Log.Debugw("creating policy", "policy", policyString)

				// parse policy
				newPolicy, err := api.ParseNewPolicy(policyString)
				if err != nil {
					return errors.Wrap(err, msg)
				}

				// create policy
				cli.StartProgress(" Creating policy...")
				createResponse, err := cli.LwApi.V2.Policy.Create(newPolicy)
				cli.StopProgress()

				// output policy
				if err != nil {
					return errors.Wrap(err, msg)
				}
				cli.OutputHuman(fmt.Sprintf("The policy %s was created.\n", createResponse.Data.PolicyID))
			}

			if pso.Operation == "update" {
				cli.Log.Debugw("updating policy", "policy", policyString)

				// parse policy
				updatePolicy, err := api.ParseUpdatePolicy(policyString)
				if err != nil {
					return errors.Wrap(err, msg)
				}
				// remove state from policies we're updating
				updatePolicy.Enabled = nil
				updatePolicy.AlertEnabled = nil

				cli.StartProgress(" Updating policy...")
				updateResponse, err := cli.LwApi.V2.Policy.Update(updatePolicy)
				cli.StopProgress()

				if err != nil {
					return errors.Wrap(err, msg)
				}
				cli.OutputHuman("The policy %s was updated.\n", updateResponse.Data.PolicyID)
			}
		}
	}
	return nil
}

func syncPolicyLibrary(_ *cobra.Command, args []string) error {
	msg := "unable to sync policies"

	// check tag
	if policyCmdState.Tag == "" {
		return errors.New("must specify the --tag flag when performing library sync")
	}
	// check json output
	if cli.JSONOutput() {
		return errors.New("json output format not supported for sync-library")
	}

	// load content library
	cli.StartProgress(" Retrieving library policies...")
	lcl, err := cli.LoadLCL()
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, msg)
	}

	// get policies for tag
	policies := lcl.GetPoliciesByTag(policyCmdState.Tag)
	if len(policies) == 0 {
		cli.OutputHuman("There were no policies found.")
		return nil
	}

	// build smart list of changes
	psos, err := getPolicySyncOperations(policies)
	if err != nil {
		return errors.Wrap(err, msg)
	}

	// prompt for changes
	proceed, err := policySyncDisplayChanges(psos)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	if !proceed {
		return nil
	}

	// execute changes
	err = policySyncExecuteChanges(lcl, psos)
	if err != nil {
		return errors.Wrap(err, msg)
	}

	return nil
}
