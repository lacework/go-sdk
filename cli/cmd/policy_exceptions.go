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
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
)

var (
	// policyExceptionCmd represents the policy parent command
	policyExceptionCmd = &cobra.Command{
		Use:     "policy-exception",
		Aliases: []string{"policy-exceptions", "pe", "px"},
		Short:   "Manage policy exceptions",
		Long: `Manage policy exceptions in your Lacework account.

To view all the policies in your Lacework account.

    lacework policy list
`,
	}

	// policyExceptionListCmd represents the policy exception list command
	policyExceptionListCmd = &cobra.Command{
		Use:     "list <policy_id>",
		Aliases: []string{"ls"},
		Short:   "List all exceptions from a single policy",
		Long:    `List all of the policy exceptions from the provided policy ID.`,
		Args:    cobra.ExactArgs(1),
		RunE:    listPolicyExceptions,
	}

	// policyExceptionShowCmd represents the policy exception show command
	policyExceptionShowCmd = &cobra.Command{
		Use:     "show <policy_id> <exception_id>",
		Aliases: []string{"get"},
		Short:   "Show details about a policy exception",
		Long:    `Show the details of a policy exception.`,
		Args:    cobra.ExactArgs(2),
		RunE:    showPolicyException,
	}

	// policyExceptionDeleteCmd represents the policy exception delete command
	policyExceptionDeleteCmd = &cobra.Command{
		Use:     "delete <policy_id> <exception_id>",
		Aliases: []string{"rm"},
		Short:   "Delete a policy exception",
		Long: `Delete a policy exception. 

To remove a policy exception, run the delete command with policy ID and exception ID arguments:

    lacework policy-exception delete <policy_id> <exception_id>`,
		Args: cobra.ExactArgs(2),
		RunE: deletePolicyException,
	}

	// policyExceptionCreateCmd represents the policy exception create command
	policyExceptionCreateCmd = &cobra.Command{
		Use:     "create [policy_id]",
		Aliases: []string{"rm"},
		Short:   "Create a policy exception",
		Long: `Create a new policy exception. 

To create a new policy exception, run the command:

    lacework policy-exception create [policy_id]

If you run the command without providing the policy_id, a
list of policies is displayed in an interactive prompt.
`,
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
	if len(args) > 0 {
		cli.StartProgress(fmt.Sprintf(
			"Retrieving policy exceptions from policy ID '%s'...", args[0],
		))
		policyExceptionResponse, err := cli.LwApi.V2.Policy.Exceptions.List(args[0])
		cli.StopProgress()
		if err != nil {
			return errors.Wrapf(err, "unable to list policy exceptions for ID %s", args[0])
		}

		if cli.JSONOutput() {
			return cli.OutputJSON(policyExceptionResponse.Data)
		}

		if len(policyExceptionResponse.Data) == 0 {
			cli.OutputHuman("There were no policy exceptions found.\n")
			return nil
		}

		cli.OutputHuman(renderSimpleTable(policyExceptionTableHeaders, policyExceptionTable(policyExceptionResponse.Data, args[0])))
	}

	return nil
}

func showPolicyException(_ *cobra.Command, args []string) error {
	var policyException api.PolicyExceptionResponse
	cli.StartProgress(fmt.Sprintf(
		"Fetching policy exception '%s' from policy '%s'...", args[0], args[1],
	))
	err := cli.LwApi.V2.Policy.Exceptions.Get(args[0], args[1], &policyException)
	cli.StopProgress()
	if err != nil {
		return errors.Wrapf(err, "unable to fetch policy exception for ID %s", args[0])
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(policyException)
	}

	cli.OutputHuman(policyExceptionDetailsTable(policyException.Data, args[0]))
	return nil
}

func deletePolicyException(_ *cobra.Command, args []string) error {
	cli.StartProgress(fmt.Sprintf(
		"Deleting policy exception '%s' from policy '%s'...", args[0], args[1],
	))
	err := cli.LwApi.V2.Policy.Exceptions.Delete(args[0], args[1])
	cli.StopProgress()
	if err != nil {
		return errors.Wrapf(err, "unable to remove policy exception for ID %s", args[0])
	}

	cli.OutputHuman("Policy exception '%s' deleted from policy '%s'\n", args[0], args[1])
	return nil
}

func createPolicyException(_ *cobra.Command, args []string) error {
	res, policyID, err := promptCreatePolicyException(args)

	if err != nil {
		return errors.Wrap(err, "unable to create policy exception")
	}

	cli.OutputHuman(
		"The policy exception '%s' was created for policy '%s' \n",
		res.Data.ExceptionID, policyID,
	)
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
			return api.PolicyExceptionResponse{}, "", errors.Wrapf(err, "invalid policy ID %s", args[0])
		}
		policyID = policy.Data.PolicyID
	} else {
		cli.StartProgress("Retrieving list of policies...")
		policies, err := cli.LwApi.V2.Policy.List()
		cli.StopProgress()
		if err != nil {
			return api.PolicyExceptionResponse{}, "", errors.Wrap(err, "unable to fetch policies")
		}
		for _, p := range policies.Data {
			if p.PolicyType != "Violation" {
				policyList = append(policyList, p.PolicyID)
			}
		}
		if err = survey.AskOne(&survey.Select{
			Message: "Policy ID:",
			Options: policyList,
		}, &policyID); err != nil {
			return api.PolicyExceptionResponse{}, "", err
		}
		policy, err = cli.LwApi.V2.Policy.Get(policyID)
		if err != nil {
			return api.PolicyExceptionResponse{}, "", errors.Wrapf(err, "invalid policy ID %s", policyID)
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

	constraints, err := promptAddExceptionConstraints(policy.Data)
	if err != nil {
		return api.PolicyExceptionResponse{}, "", err
	}

	exception := api.PolicyException{Description: answers.Description, Constraints: constraints}
	cli.StartProgress("Creating policy exception ...")
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

func promptAddExceptionConstraints(policy api.Policy) ([]api.PolicyExceptionConstraint, error) {
	if len(policy.ExceptionConfiguration) == 0 {
		return []api.PolicyExceptionConstraint{}, nil
	}
	var responses []api.PolicyExceptionConstraint
	questions := buildPromptAddExceptionConstraintListQuestions(policy.ExceptionConfiguration["constraintFields"])

	for _, q := range questions {
		addConstraint, err := promptAddConstraint(q.constraint.FieldKey)
		if err != nil {
			return []api.PolicyExceptionConstraint{}, nil
		}

		if !addConstraint {
			continue
		}

		answer, err := promptAskConstraintsQuestion(q)

		if err != nil {
			return []api.PolicyExceptionConstraint{}, nil
		}

		if answer != nil {
			responses = append(responses, *answer)
		}
	}

	if len(responses) == 0 {
		return []api.PolicyExceptionConstraint{}, errors.New("policy exceptions must have at least 1 constraint")
	}

	return responses, nil
}

func promptAskConstraintsQuestion(constraintQuestion PolicyExceptionSurveyQuestion) (*api.PolicyExceptionConstraint, error) {
	var (
		answer *api.PolicyExceptionConstraint
		err    error
	)

	switch constraintQuestion.constraint.DataType {
	case "String":
		if constraintQuestion.constraint.MultiValue {
			// prompt string list question
			answer, err = promptAddExceptionConstraintList(constraintQuestion.constraint.FieldKey, constraintQuestion.questions)
		} else {
			// prompt string question
			answer, err = promptAddExceptionConstraintString(constraintQuestion.constraint.FieldKey, constraintQuestion.questions)
		}
	case "KVTagPair":
		if constraintQuestion.constraint.MultiValue {
			//prompt kv tag list
			answer, err = promptAddExceptionConstraintMapList(constraintQuestion.constraint.FieldKey, constraintQuestion.questions)
		} else {
			//prompt kv tag
			mapAnswers, err := promptAddExceptionConstraintMap(constraintQuestion.questions)
			if err != nil {
				return nil, err
			}
			return &api.PolicyExceptionConstraint{
				FieldKey:    constraintQuestion.constraint.FieldKey,
				FieldValues: []any{mapAnswers},
			}, nil
		}
	}

	if err != nil {
		return nil, err
	}

	return answer, nil
}

func buildPromptAddExceptionConstraintListQuestions(constraints []api.PolicyExceptionConfigurationConstraints) []PolicyExceptionSurveyQuestion {
	questions := []PolicyExceptionSurveyQuestion{}

	for _, constraint := range constraints {
		switch constraint.DataType {
		case "String":
			if constraint.MultiValue {
				questions = append(questions, PolicyExceptionSurveyQuestion{[]*survey.Question{{
					Name:     constraint.FieldKey,
					Prompt:   &survey.Multiline{Message: fmt.Sprintf("%s values:", constraint.FieldKey)},
					Validate: survey.Required,
				}},
					constraint,
				})
			} else {
				questions = append(questions, PolicyExceptionSurveyQuestion{[]*survey.Question{{
					Name:     constraint.FieldKey,
					Prompt:   &survey.Input{Message: fmt.Sprintf("%s value:", constraint.FieldKey)},
					Validate: survey.Required,
				}},
					constraint,
				})
			}
		case "KVTagPair":
			if constraint.MultiValue {
				questions = append(questions, PolicyExceptionSurveyQuestion{[]*survey.Question{{
					Name:     "key",
					Prompt:   &survey.Input{Message: "key:"},
					Validate: survey.Required,
				}, {
					Name:     "value",
					Prompt:   &survey.Input{Message: "value:"},
					Validate: survey.Required,
				},
				},
					constraint,
				})
			} else {
				questions = append(questions, PolicyExceptionSurveyQuestion{[]*survey.Question{{
					Name:     "key",
					Prompt:   &survey.Input{Message: "key:"},
					Validate: survey.Required,
				}, {
					Name:     "value",
					Prompt:   &survey.Input{Message: "value:"},
					Validate: survey.Required,
				},
				},
					constraint,
				})
			}
		}
	}

	return questions
}

type PolicyExceptionSurveyQuestion struct {
	questions  []*survey.Question
	constraint api.PolicyExceptionConfigurationConstraints
}

func promptAddExceptionConstraintList(key string, questions []*survey.Question) (*api.PolicyExceptionConstraint, error) {
	if key == "accountIds" {
		return promptAddExceptionConstraintAwsAccountsList()
	}

	var (
		values      []any
		fieldValues string
	)

	err := survey.Ask(questions, &fieldValues)
	if err != nil {
		return nil, err
	}

	for _, v := range strings.Split(fieldValues, "\n") {
		values = append(values, v)
	}

	return &api.PolicyExceptionConstraint{
		FieldKey:    key,
		FieldValues: values,
	}, nil
}

func promptAddExceptionConstraintString(key string, questions []*survey.Question) (*api.PolicyExceptionConstraint, error) {
	var fieldValue string

	err := survey.Ask(questions, &fieldValue)
	if err != nil {
		return nil, err
	}

	return &api.PolicyExceptionConstraint{
		FieldKey:    key,
		FieldValues: []any{fieldValue},
	}, nil
}

func promptAddExceptionConstraintAwsAccountsList() (*api.PolicyExceptionConstraint, error) {
	var (
		values      []any
		fieldValues []string
		accountIds  []string
	)

	cli.StartProgress("Retrieving AWS accounts...")
	accounts, err := cli.LwApi.V2.CloudAccounts.ListByType(api.AwsCfgCloudAccount)
	cli.StopProgress()

	if err != nil {
		return nil, err
	}

	if len(accounts.Data) == 0 {
		return nil, errors.New("no aws cloud accounts found")
	}

	for _, ca := range accounts.Data {
		if caMap, ok := ca.GetData().(map[string]interface{}); ok {
			accountIds = append(accountIds, caMap["awsAccountId"].(string))
		}
	}

	question := survey.Question{Name: "awsAccounts",
		Prompt: &survey.MultiSelect{Message: "Select AWS Accounts:",
			Options: array.Unique(accountIds)},
		Validate: survey.MinItems(1)}

	err = survey.Ask([]*survey.Question{&question}, &fieldValues)
	if err != nil {
		return nil, err
	}

	for _, v := range fieldValues {
		values = append(values, v)
	}
	return &api.PolicyExceptionConstraint{
		FieldKey:    "accountIds",
		FieldValues: values,
	}, nil
}

func promptAddExceptionConstraintMap(mapQuestions []*survey.Question) (constraintMapAnswer, error) {
	var mapAnswer constraintMapAnswer
	err := survey.Ask(mapQuestions, &mapAnswer,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return constraintMapAnswer{}, err
	}

	return mapAnswer, nil
}

type constraintMapAnswer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func promptAddExceptionConstraintMapList(key string, mapQuestions []*survey.Question) (*api.PolicyExceptionConstraint, error) {
	var mapAnswers []any

	res, err := promptAddExceptionConstraintMap(mapQuestions)
	if err != nil {
		return nil, err
	}
	mapAnswers = append(mapAnswers, res)

	addTag := false
	for {
		if err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Add another %s constraint?", key),
		}, &addTag); err != nil {
			return nil, err
		}

		if addTag {
			res, err := promptAddExceptionConstraintMap(mapQuestions)
			if err != nil {
				return nil, err
			}
			mapAnswers = append(mapAnswers, res)
		} else {
			break
		}
	}

	return &api.PolicyExceptionConstraint{
		FieldKey:    key,
		FieldValues: mapAnswers,
	}, nil
}

func promptAddConstraint(key string) (bool, error) {
	addConstraint := false

	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Add constraint %q?", key),
	}, &addConstraint); err != nil {
		return false, err
	}

	return addConstraint, nil
}
