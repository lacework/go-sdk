//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/lacework/go-sdk/api"
)

func createEmailAlertChannelIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "recipients",
			Prompt:   &survey.Multiline{Message: "List of Recipients: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name       string
		Recipients string
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	emailAlertChan := api.NewAlertChannel(answers.Name,
		api.EmailUserAlertChannelType,
		api.EmailUserData{
			ChannelProps: api.EmailUserChannelProps{
				Recipients: strings.Split(answers.Recipients, "\n"),
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.AlertChannels.Create(emailAlertChan)
	cli.StopProgress()
	return err
}
