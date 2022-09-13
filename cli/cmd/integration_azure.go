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

func createAzureConfigIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "client_id",
			Prompt:   &survey.Input{Message: "Client ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "client_secret",
			Prompt:   &survey.Input{Message: "Client Secret:"},
			Validate: survey.Required,
		},
		{
			Name:     "tenant_id",
			Prompt:   &survey.Input{Message: "Tenant ID:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name         string
		ClientID     string `survey:"client_id"`
		ClientSecret string `survey:"client_secret"`
		TenantID     string `survey:"tenant_id"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	azure := api.NewCloudAccount(answers.Name,
		api.AzureCfgCloudAccount,
		api.AzureCfgData{
			TenantID: answers.TenantID,
			Credentials: api.AzureCfgCredentials{
				ClientID:     answers.ClientID,
				ClientSecret: answers.ClientSecret,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.CloudAccounts.Create(azure)
	cli.StopProgress()
	return err
}

func createAzureActivityLogIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "client_id",
			Prompt:   &survey.Input{Message: "Client ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "client_secret",
			Prompt:   &survey.Input{Message: "Client Secret:"},
			Validate: survey.Required,
		},
		{
			Name:     "tenant_id",
			Prompt:   &survey.Input{Message: "Tenant ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "queue_url",
			Prompt:   &survey.Input{Message: "Queue URL:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name         string
		ClientID     string `survey:"client_id"`
		ClientSecret string `survey:"client_secret"`
		TenantID     string `survey:"tenant_id"`
		QueueUrl     string `survey:"queue_url"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	azure := api.NewCloudAccount(answers.Name,
		api.AzureAlSeqCloudAccount,
		api.AzureAlSeqData{
			QueueUrl: answers.QueueUrl,
			TenantID: answers.TenantID,
			Credentials: api.AzureAlSeqCredentials{
				ClientID:     answers.ClientID,
				ClientSecret: answers.ClientSecret,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.CloudAccounts.Create(azure)
	cli.StopProgress()
	return err
}
