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

func createAzureResourceGroup() error {
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
			Name:     "tenant",
			Prompt:   &survey.Input{Message: "Tenant: "},
			Validate: survey.Required,
		},
		{
			Name:     "subscriptions",
			Prompt:   &survey.Multiline{Message: "List of Subscriptions: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name          string
		Description   string `survey:"description"`
		Tenant        string `survey:"tenant"`
		Subscriptions string `survey:"subscriptions"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	azure := api.NewResourceGroup(
		answers.Name,
		api.AzureResourceGroup,
		api.AzureResourceGroupProps{
			Description:   answers.Description,
			Tenant:        answers.Tenant,
			Subscriptions: strings.Split(answers.Subscriptions, "\n"),
		})

	cli.StartProgress(" Creating resource group...")
	_, err = cli.LwApi.V2.ResourceGroups.Create(azure)
	cli.StopProgress()
	return err
}

func setAzureProps(group string) [][]string {
	var (
		azProps api.AzureResourceGroupProps
		details [][]string
	)
	err := json.Unmarshal([]byte(group), &azProps)
	if err != nil {
		return [][]string{}
	}

	details = append(details, []string{"TENANT", azProps.Tenant})
	details = append(details, []string{"SUBSCRIPTIONS", strings.Join(azProps.Subscriptions, ",")})
	return details
}
