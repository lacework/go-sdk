//
// Author:: Darren Murray(<darren.murray@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/lacework/go-sdk/api"
)

func createAwsResourceGroup() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "description",
			Prompt:   &survey.Input{Message: "Description: "},
			Validate: survey.Required,
		},
		{
			Name:     "account_ids",
			Prompt:   &survey.Multiline{Message: "List of Account IDs: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name        string
		Description string `survey:"description"`
		AccountIDs  string `survey:"account_ids"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	aws := api.NewResourceGroup(
		answers.Name,
		api.AwsResourceGroup,
		api.AwsResourceGroupProps{
			Description: answers.Description,
			AccountIDs:  strings.Split(answers.AccountIDs, "\n"),
		})

	cli.StartProgress(" Creating resource group...")
	_, err = cli.LwApi.V2.ResourceGroups.Create(aws)
	cli.StopProgress()
	return err
}

func setAwsProps(group string) []string {
	var awsProps api.AwsResourceGroupProps
	err := json.Unmarshal([]byte(group), &awsProps)
	if err != nil {
		return []string{}
	}

	return []string{"ACCOUNT IDS", strings.Join(awsProps.AccountIDs, ",")}
}
