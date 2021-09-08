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

func createMachineResourceGroup() error {
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
			Prompt:   &survey.Multiline{Message: "List of 'key:value' Machine Tags:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name        string
		Description string `survey:"description"`
		MachineTags string `survey:"tags"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	machine := api.NewResourceGroup(
		answers.Name,
		api.MachineResourceGroup,
		api.MachineResourceGroupProps{
			Description: answers.Description,
			MachineTags: castStringToLimitByLabel(answers.MachineTags),
		})

	cli.StartProgress(" Creating resource group...")
	_, err = cli.LwApi.V2.ResourceGroups.Create(machine)
	cli.StopProgress()
	return err
}

func setMachineProps(group []byte) []string {
	machineProps, err := unmarshallMachineProps(group)
	if err != nil {
		return []string{}
	}

	var tags []string
	for _, tagMap := range machineProps.MachineTags {
		for key, val := range tagMap {
			tags = append(tags, fmt.Sprintf("%s: %v", key, val))

		}
	}
	return []string{"MACHINE TAGS", strings.Join(tags, ",")}
}

func unmarshallMachinePropString(group []byte) (props api.MachineResourceGroupProps, err error) {
	var rawProps api.MachineResourceGroupJsonStringProps
	err = json.Unmarshal(group, &rawProps)
	props = api.MachineResourceGroupProps(rawProps)
	return
}

func unmarshallMachineProps(group []byte) (props api.MachineResourceGroupProps, err error) {
	err = json.Unmarshal(group, &props)
	return
}
