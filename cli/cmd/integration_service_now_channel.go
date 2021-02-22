//
// Author:: Darren Murray(<darren.murray@lacework.net>)
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
	"github.com/AlecAivazis/survey/v2"

	"github.com/lacework/go-sdk/api"
)

func createServiceNowAlertChannelIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "instance_url",
			Prompt:   &survey.Input{Message: "InstanceURL:"},
			Validate: survey.Required,
		},
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "Username:"},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Input{Message: "Password:"},
			Validate: survey.Required,
		},
		{
			Name: "issue_grouping",
			Prompt: &survey.Select{Message: "Issue Grouping:",
				Options: []string{"Events", "Resources"},
			},
		},
	}

	answers := struct {
		Name          string
		InstanceURL   string `survey:"instance_url"`
		Username      string `survey:"username"`
		Password      string `survey:"password"`
		IssueGrouping string `survey:"issue_grouping"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	snow := api.ServiceNowChannelData{
		InstanceURL:   answers.InstanceURL,
		Username:      answers.Username,
		Password:      answers.Password,
		IssueGrouping: answers.IssueGrouping,
	}

	// ask the user if they would like to configure a Custom Template
	custom := false
	err = survey.AskOne(&survey.Confirm{
		Message: "Configure a Custom Template File?",
	}, &custom)

	if err != nil {
		return err
	}

	if custom {
		var content string

		err = survey.AskOne(&survey.Editor{
			Message:  "Provide the Custom Template File in JSON format",
			FileName: "*.json",
		}, &content)

		if err != nil {
			return err
		}

		snow.EncodeCustomTemplateFile(content)
	}

	snowAlert := api.NewServiceNowAlertChannel(answers.Name, snow)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.Integrations.CreateServiceNowAlertChannel(snowAlert)
	cli.StopProgress()
	return err
}
