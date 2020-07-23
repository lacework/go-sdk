//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

func createPagerDutyAlertChannelIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "key",
			Prompt:   &survey.Input{Message: "Integration Key: "},
			Validate: survey.Required,
		},
		{
			Name: "alert_severity_level",
			Prompt: &survey.Select{
				Message: "Alert Severity Level: ",
				Options: []string{
					"Critical",
					"High and above",
					"Medium and above",
					"Low and above",
					"All",
				},
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name          string
		Key           string
		AlertSeverity string `survey:"alert_severity_level"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	alert := api.NewPagerDutyAlertChannel(answers.Name,
		api.PagerDutyData{
			IntegrationKey:   answers.Key,
			MinAlertSeverity: alertSeverityToEnum(answers.AlertSeverity),
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.Integrations.CreatePagerDutyAlertChannel(alert)
	cli.StopProgress()
	return err
}
