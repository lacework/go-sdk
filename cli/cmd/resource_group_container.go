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
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/lacework/go-sdk/api"
)

func createContainerResourceGroup() error {
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
			Name:     "tags",
			Prompt:   &survey.Multiline{Message: "List of Tags: "},
			Validate: survey.Required,
		}, {
			Name:     "labels",
			Prompt:   &survey.Multiline{Message: "List of 'key:value' Labels:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name        string
		Description string `survey:"description"`
		Tags        string `survey:"tags"`
		Labels      string `survey:"labels"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	container := api.NewResourceGroup(
		answers.Name,
		api.ContainerResourceGroup,
		api.ContainerResourceGroupProps{
			Description:     answers.Description,
			ContainerTags:   strings.Split(answers.Tags, "\n"),
			ContainerLabels: castStringToLimitByLabel(answers.Labels),
		})

	cli.StartProgress(" Creating resource group...")
	_, err = cli.LwApi.V2.ResourceGroups.Create(container)
	cli.StopProgress()
	return err
}

func setContainerProps(group interface{}) [][]string {
	var (
		ctrProps api.ContainerResourceGroupProps
		labels   []string
		details  [][]string
	)
	err := json.Unmarshal([]byte(group.(string)), &ctrProps)
	if err != nil {
		return [][]string{}
	}

	for _, labelMap := range ctrProps.ContainerLabels {
		for key, val := range labelMap {
			labels = append(labels, fmt.Sprintf("%s: %v", key, val))
		}
	}
	details = append(details, []string{"CONTAINER LABELS", strings.Join(labels, ",")})
	details = append(details, []string{"CONTAINER TAGS", strings.Join(ctrProps.ContainerTags, ",")})
	return details
}
