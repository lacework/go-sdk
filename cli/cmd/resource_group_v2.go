//
// Author:: Zeki Sherif(<zeki.sherif@lacework.net>)
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
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/v2/api"
)

func createResourceGroup(resourceType string) error {
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
			Name:     "query",
			Prompt:   inputRGQueryFromEditor(resourceType),
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name        string `survey:"name"`
		Description string `survey:"description"`
		Query       string `survey:"query"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)

	if err != nil {
		return err
	}

	var rgQuery api.RGQuery
	err = json.Unmarshal([]byte(answers.Query), &rgQuery)
	if err != nil {
		return err
	}

	groupType, isValid := api.FindResourceGroupType(resourceType)
	if !isValid {
		// This should never reach this. The type is controlled by us in cmd/resource_groups
		return errors.New("internal error")
	}
	resourceGroup := api.NewResourceGroup(answers.Name, groupType, answers.Description, &rgQuery)
	cli.StartProgress(" Creating resource group...")
	_, err = cli.LwApi.V2.ResourceGroups.Create(resourceGroup)
	cli.StopProgress()
	return err
}

func inputRGQueryFromEditor(resourceType string) *survey.Editor {
	iType, _ := api.FindResourceGroupType(resourceType)

	prompt := &survey.Editor{
		Message:  fmt.Sprintf("Type a query for the new %s Resource Group", resourceType),
		FileName: "resourceGroupQuery*.json",
		Help: `For examples of a query, refer to: https://docs.fortinet.com/document/lacework-forticnapp/` +
			`latest/api-reference/690087/using-the-resource-groups-api`,
	}
	prompt.Default = iType.QueryTemplate()
	prompt.HideDefault = true
	prompt.AppendDefault = true

	return prompt
}
