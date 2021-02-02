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

func createDatadogIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name: "datadog_site",
			Prompt: &survey.Select{Message: "Datadog Site: ",
				Options: []string{api.DatadogSiteEu.String(), api.DatadogSiteCom.String()},
				Default: api.DatadogSiteCom.String(),
			},
		},
		{
			Name: "datadog_type",
			Prompt: &survey.Select{Message: "Datadog Type: ",
				Options: []string{api.DatadogServiceLogsDetails.String(), api.DatadogServiceEventsSummary.String(), api.DatadogServiceLogsSummary.String()},
				Default: api.DatadogServiceLogsDetails.String(),
			},
		},
		{
			Name:     "api_key",
			Prompt:   &survey.Input{Message: "Api Key: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name           string
		DatadogSite    string `survey:"datadog_site"`
		DatadogService string `survey:"datadog_type"`
		ApiKey         string `survey:"api_key"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	datadog := api.NewDatadogAlertChannel(answers.Name,
		api.DatadogChannelData{
			DatadogSite:    answers.DatadogSite,
			DatadogService: answers.DatadogService,
			ApiKey:         answers.ApiKey,
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.Integrations.CreateDatadogAlertChannel(datadog)
	cli.StopProgress()
	return err
}
