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

func createQRadarAlertChannelIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "host_url",
			Prompt:   &survey.Input{Message: "Url:"},
			Validate: survey.Required,
		},
		{
			Name:     "host_port",
			Prompt:   &survey.Input{Message: "Port:"},
			Validate: survey.Required,
		},
		{
			Name: "communication_type",
			Prompt: &survey.Select{Message: "Communication Type:",
				Options: []string{string(api.QRadarCommHttps), string(api.QRadarCommSelfSigned)},
				Default: string(api.QRadarCommHttps),
			},
		},
	}

	answers := struct {
		Name              string
		HostURL           string `survey:"host_url"`
		HostPort          int    `survey:"host_port"`
		CommunicationType string `survey:"communication_type"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	commType, err := api.DatadogService(answers.CommunicationType)
	if err != nil {
		return err
	}

	qradar := api.NewQRadarAlertChannel(answers.Name,
		api.QRadarChannelData{
			HostURL:           answers.HostURL,
			HostPort:          answers.HostPort,
			CommunicationType: commType,
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.Integrations.CreateQRadarAlertChannel(qradar)
	cli.StopProgress()
	return err
}
