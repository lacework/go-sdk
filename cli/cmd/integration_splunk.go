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

func createSplunkIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:   "channel",
			Prompt: &survey.Input{Message: "Channel: "},
		},
		{
			Name:     "hec_token",
			Prompt:   &survey.Input{Message: "Hec Token: "},
			Validate: survey.Required,
		},
		{
			Name:     "host",
			Prompt:   &survey.Input{Message: "Host: "},
			Validate: survey.Required,
		},
		{
			Name:     "port",
			Prompt:   &survey.Input{Message: "Port: "},
			Validate: survey.Required,
		},
		{
			Name:     "source",
			Prompt:   &survey.Input{Message: "Source: "},
			Validate: survey.Required,
		},
		{
			Name:     "index",
			Prompt:   &survey.Input{Message: "Index: "},
			Validate: survey.Required,
		},
		{
			Name:   "ssl",
			Prompt: &survey.Confirm{Message: "Enable SSL?"},
		},
	}

	answers := struct {
		Name     string
		Channel  string
		HecToken string `survey:"hec_token"`
		Host     string
		Port     int
		Source   string
		Index    string
		Ssl      bool
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	splunk := api.NewAlertChannel(answers.Name,
		api.SplunkHecAlertChannelType,
		api.SplunkHecDataV2{
			Channel:  answers.Channel,
			HecToken: answers.HecToken,
			Host:     answers.Host,
			Port:     answers.Port,
			Ssl:      answers.Ssl,
			EventData: api.SplunkHecEventDataV2{
				Index:  answers.Index,
				Source: answers.Source,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.AlertChannels.Create(splunk)
	cli.StopProgress()
	return err
}
