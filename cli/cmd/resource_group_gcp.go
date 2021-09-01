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

func createGcpResourceGroup() error {
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
			Name:     "organization",
			Prompt:   &survey.Input{Message: "organization: "},
			Validate: survey.Required,
		},
		{
			Name:     "projects",
			Prompt:   &survey.Multiline{Message: "List of Projects: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name         string
		Description  string `survey:"description"`
		Organization string `survey:"organization"`
		Projects     string `survey:"projects"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	gcp := api.NewResourceGroup(
		answers.Name,
		api.GcpResourceGroup,
		api.GcpResourceGroupProps{
			Description:  answers.Description,
			Organization: answers.Organization,
			Projects:     strings.Split(answers.Projects, "\n"),
		})

	cli.StartProgress(" Creating resource group...")
	_, err = cli.LwApi.V2.ResourceGroups.Create(gcp)
	cli.StopProgress()
	return err
}

func setGcpProps(group interface{}) [][]string {
	var (
		gcpProps api.GcpResourceGroupProps
		details  [][]string
	)
	err := json.Unmarshal([]byte(group.(string)), &gcpProps)
	if err != nil {
		return [][]string{}
	}

	details = append(details, []string{"ORGANIZATION", gcpProps.Organization})
	details = append(details, []string{"PROJECTS", strings.Join(gcpProps.Projects, ",")})
	return details
}
