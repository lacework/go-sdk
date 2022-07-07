//
// Author:: Darren Murray(<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// policyExceptionCmd represents the policy parent command
	policyExceptionCmd = &cobra.Command{
		Use:     "policy-exception",
		Aliases: []string{"policy-exceptions", "pe"},
		Short:   "Manage policy exceptions",
		Long: `Manage policies in your Lacework account.

A policy exception is used to exclude .... from a policy

To view all the policies in your Lacework account.

    lacework policy ls
`,
	}

	// policyExceptionListCmd represents the policy exception list command
	policyExceptionListCmd = &cobra.Command{
		Use:     "list <policy-id>",
		Aliases: []string{"ls"},
		Short:   "List policy exceptions",
		Long: `List all of the policy exceptions in your Lacework account.

To list the exceptions for a single policy, provide the policy id argument:

	lacework policy-exception list [lacework-policy-id]`,
		Args: cobra.ExactArgs(1),
		RunE: listPolicyExceptions,
	}

	// policyExceptionShowCmd represents the policy exception show command
	policyExceptionShowCmd = &cobra.Command{
		Use:     "show <policy-id> <exception-id>",
		Aliases: []string{"get"},
		Short:   "Show policy exception details",
		Long: `Show the details of a policy exception. 

To show details of a policy exception, run the show command with policy id and exception id arguments:
	
	lacework policy-exception show <policy-id> <exception-id>`,
		Args: cobra.ExactArgs(2),
		RunE: showPolicyException,
	}

	// policyExceptionDeleteCmd represents the policy exception delete command
	policyExceptionDeleteCmd = &cobra.Command{
		Use:     "delete <policy-id> <exception-id>",
		Aliases: []string{"rm"},
		Short:   "Delete a policy exception",
		Long: `Delete a policy exception. 

To remove a policy exception, run the delete command with policy id and exception id arguments:

	lacework policy-exception delete <policy-id> <exception-id>`,
		Args: cobra.ExactArgs(2),
		RunE: deletePolicyException,
	}

	// policyExceptionCreateCmd represents the policy exception create command
	policyExceptionCreateCmd = &cobra.Command{
		Use:     "create [policy-id]",
		Aliases: []string{"rm"},
		Short:   "Create a policy exception",
		Long: `Create a new policy exception. 

To create a new policy exception, run the create command:

	lacework policy-exception create [policy-id]`,
		Args: cobra.MaximumNArgs(1),
		RunE: createPolicyException,
	}
)

func init() {
	// add the policy exception command
	rootCmd.AddCommand(policyExceptionCmd)

	// add sub-commands to the policy exception command
	policyExceptionCmd.AddCommand(policyExceptionListCmd)
	policyExceptionCmd.AddCommand(policyExceptionShowCmd)
	policyExceptionCmd.AddCommand(policyExceptionDeleteCmd)
	policyExceptionCmd.AddCommand(policyExceptionCreateCmd)
}

func listPolicyExceptions(_ *cobra.Command, args []string) error {
	var policyExceptions [][]string
	if len(args) > 0 {
		cli.StartProgress(" Retrieving policy exceptions...")
		policyExceptionResponse, err := cli.LwApi.V2.Policy.Exceptions.List(args[0])
		cli.StopProgress()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to list policy exceptions for id %s", args[0]))
		}
		policyExceptions = append(policyExceptions, policyExceptionTable(policyExceptionResponse.Data, args[0])...)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(policyExceptions)
	}
	if len(policyExceptions) == 0 {
		cli.OutputHuman("There were no policy exceptions found.\n")
		return nil
	}

	cli.OutputHuman(renderSimpleTable(policyExceptionTableHeaders, policyExceptions))
	return nil
}

func showPolicyException(_ *cobra.Command, args []string) error {
	var policyException api.PolicyExceptionResponse
	cli.StartProgress(" Fetching policy exception...")
	err := cli.LwApi.V2.Policy.Exceptions.Get(args[0], args[1], &policyException)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to fetch policy exception for id %s", args[0]))
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(policyException)
	}

	cli.OutputHuman(policyExceptionDetailsTable(policyException.Data, args[0]))
	return nil
}

func deletePolicyException(_ *cobra.Command, args []string) error {
	cli.StartProgress(" Deleting policy exception...")
	err := cli.LwApi.V2.Policy.Exceptions.Delete(args[0], args[1])
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to remove policy exception for id %s", args[0]))
	}

	cli.OutputHuman("Policy Exception %s %s deleted\n", args[0], args[1])
	return nil
}

func createPolicyException(_ *cobra.Command, args []string) error {
	res, policyID, err := promptCreatePolicyException(args)

	if err != nil {
		return errors.Wrap(err, "unable to create policy exception")
	}

	cli.OutputHuman("New Policy Exception %s %s created \n", policyID, res.Data.ExceptionID)
	return nil
}

func promptCreatePolicyException(args []string) (api.PolicyExceptionResponse, string, error) {
	var (
		policy     api.PolicyResponse
		policyList []string
		policyID   string
		err        error
	)

	if len(args) > 0 {
		policy, err = cli.LwApi.V2.Policy.Get(args[0])
		if err != nil {
			return api.PolicyExceptionResponse{}, "", errors.Wrap(err, fmt.Sprintf("invalid policy id %s", args[0]))
		}
		policyID = policy.Data.PolicyID
	} else {
		cli.StartProgress(" Retrieving list of policies...")
		policies, err := cli.LwApi.V2.Policy.List()
		cli.StopProgress()
		if err != nil {
			return api.PolicyExceptionResponse{}, "", errors.Wrap(err, "unable to fetch policies")
		}
		for _, p := range policies.Data {
			policyList = append(policyList, p.PolicyID)
		}
		if err = survey.AskOne(&survey.Select{
			Message: "Policy ID:",
			Options: policyList,
		}, &policyID); err != nil {
			return api.PolicyExceptionResponse{}, "", err
		}
		policy, err = cli.LwApi.V2.Policy.Get(policyID)
		if err != nil {
			return api.PolicyExceptionResponse{}, "", errors.Wrap(err, fmt.Sprintf("invalid policy id %s", policyID))
		}
	}

	questions := []*survey.Question{
		{
			Name:     "description",
			Prompt:   &survey.Input{Message: "Exception Description: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Description string `json:"description"`
	}{}

	err = survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return api.PolicyExceptionResponse{}, "", err
	}

	var constraints []api.PolicyExceptionConstraint
	constraints = append(constraints, promptAddExceptionConstraint(policy.Data))
	addConstraint := false
	for {
		if err := survey.AskOne(&survey.Confirm{
			Message: "Add another constraint?",
		}, &addConstraint); err != nil {
			return api.PolicyExceptionResponse{}, "", err
		}

		if addConstraint {
			constraints = append(constraints, promptAddExceptionConstraint(policy.Data))
		} else {
			break
		}
	}
	exception := api.PolicyException{Description: answers.Description, Constraints: constraints}

	cli.StartProgress(" Creating policy exception ...")
	response, err := cli.LwApi.V2.Policy.Exceptions.Create(policyID, exception)

	cli.StopProgress()
	return response, policyID, err
}

var policyExceptionTableHeaders = []string{"POLICY ID", "EXCEPTION ID", "DESCRIPTION", "UPDATED AT", "UPDATED BY"}

func policyExceptionTable(policyException []api.PolicyException, policyID string) (out [][]string) {
	for _, exception := range policyException {
		out = append(out, []string{
			policyID,
			exception.ExceptionID,
			exception.Description,
			exception.LastUpdateTime,
			exception.LastUpdateUser,
		})
	}
	return
}

func policyExceptionDetailsTable(policyException api.PolicyException, policyID string) string {
	var (
		table   strings.Builder
		out     [][]string
		details [][]string
	)

	out = append(out, []string{
		policyID,
		policyException.ExceptionID,
		policyException.Description,
		policyException.LastUpdateTime,
		policyException.LastUpdateUser,
	})

	table.WriteString(renderSimpleTable(policyExceptionTableHeaders, out))
	table.WriteString("\n")

	for _, constraint := range policyException.Constraints {
		jsonFieldValues, _ := json.Marshal(constraint.FieldValues)
		details = append(details, []string{constraint.FieldKey, string(jsonFieldValues)})
	}

	table.WriteString(renderOneLineCustomTable("CONSTRAINTS",
		renderSimpleTable([]string{"FIELD KEY", "FIELD VALUES"}, details),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		})))

	return table.String()
}

func promptAddExceptionConstraint(policy api.Policy) api.PolicyExceptionConstraint {
	help := "Visit the lacework docs to see exception criteria https://docs.lacework.com/aws-compliance-policy-exceptions-criteria"
	msg := ""
	switch {
	case array.ContainsStr(policy.Tags, "domain:AWS"):
		msg = "Valid constraint keys for aws polices are 'accountIds', 'resourceNames', 'regionNames' and 'resourceTags'`\n"
	}

	questions := []*survey.Question{
		{
			Name: "fieldKey",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("%s Constraint Field Key: ", msg),
				Help:    help,
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		FieldKey string `json:"fieldKey"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return api.PolicyExceptionConstraint{}
	}

	var values []any

	if strings.Contains(strings.ToLower(answers.FieldKey), "tags") {
		constraintMap, err := promptAddExceptionConstraintMap()
		if err != nil {
			return api.PolicyExceptionConstraint{}
		}
		values = append(values, constraintMap)

		addTag := false
		for {
			if err := survey.AskOne(&survey.Confirm{
				Message: "Add another tag?",
			}, &addTag); err != nil {
				return api.PolicyExceptionConstraint{}
			}

			if addTag {
				constraintMap, err = promptAddExceptionConstraintMap()
				if err != nil {
					return api.PolicyExceptionConstraint{}
				}
				values = append(values, constraintMap)

			} else {
				break
			}
		}
	} else if answers.FieldKey == "accountIds" {
		values, err = promptAddExceptionConstraintAwsAccountsList()
		if err != nil {
			return api.PolicyExceptionConstraint{}
		}
	} else {
		values, err = promptAddExceptionConstraintList()
		if err != nil {
			return api.PolicyExceptionConstraint{}
		}
	}

	return api.PolicyExceptionConstraint{
		FieldKey:    answers.FieldKey,
		FieldValues: values,
	}
}

func promptAddExceptionConstraintList() ([]any, error) {
	var (
		values      []any
		fieldValues string
	)
	err := survey.AskOne(&survey.Multiline{Message: "Constraint Field Values:"}, &fieldValues)
	if err != nil {
		return nil, err
	}
	for _, v := range strings.Split(fieldValues, "\n") {
		values = append(values, v)
	}
	return values, nil
}

func promptAddExceptionConstraintAwsAccountsList() ([]any, error) {
	var (
		values      []any
		fieldValues []string
	)

	cli.StartProgress(" Retrieving aws accounts...")
	accountIds, err := cli.LwApi.Integrations.AwsAccountIDs()
	cli.StopProgress()

	err = survey.AskOne(&survey.MultiSelect{Message: "Select Aws Accounts:", Options: accountIds}, &fieldValues)
	if err != nil {
		return nil, err
	}
	for _, v := range fieldValues {
		values = append(values, v)
	}
	return values, nil
}

func promptAddExceptionConstraintMap() (any, error) {
	mapQuestions := []*survey.Question{
		{
			Name:     "key",
			Prompt:   &survey.Input{Message: "Tag Key: "},
			Validate: survey.Required,
		},
		{
			Name:     "value",
			Prompt:   &survey.Input{Message: "Tag Value: "},
			Validate: survey.Required,
		},
	}

	mapAnswers := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}

	err := survey.Ask(mapQuestions, &mapAnswers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return nil, err
	}

	return mapAnswers, nil
}
