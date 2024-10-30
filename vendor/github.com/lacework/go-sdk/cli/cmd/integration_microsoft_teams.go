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

func createMicrosoftTeamsChannelIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "webhook_url",
			Prompt:   &survey.Input{Message: "Webhook URL: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name       string
		WebhookUrl string `survey:"webhook_url"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	teams := api.NewAlertChannel(answers.Name,
		api.MicrosoftTeamsAlertChannelType,
		api.MicrosoftTeamsData{
			TeamsURL: answers.WebhookUrl,
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.AlertChannels.Create(teams)
	cli.StopProgress()
	return err
}
