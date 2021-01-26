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

func createGcpConfigIntegration() error {
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
			Name:     "private_key_id",
			Prompt:   &survey.Input{Message: "Private Key ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "client_email",
			Prompt:   &survey.Input{Message: "Client Email:"},
			Validate: survey.Required,
		},
		{
			Name:     "private_key",
			Prompt:   &survey.Editor{Message: "Private Key:"},
			Validate: survey.Required,
		},
		{
			Name: "integration_level",
			Prompt: &survey.Select{
				Message: "Integration Level:",
				Options: []string{
					api.GcpProjectIntegration.String(),
					api.GcpOrganizationIntegration.String(),
				},
			},
			Validate: survey.Required,
		},
		{
			Name:     "org_project_id",
			Prompt:   &survey.Input{Message: "Organization/Project ID:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name             string
		ClientID         string `survey:"client_id"`
		PrivateKeyID     string `survey:"private_key_id"`
		ClientEmail      string `survey:"client_email"`
		PrivateKey       string `survey:"private_key"`
		IntegrationLevel string `survey:"integration_level"`
		OrgProjectID     string `survey:"org_project_id"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	gcp := api.NewGcpCfgIntegration(answers.Name,
		api.GcpIntegrationData{
			ID:     answers.OrgProjectID,
			IDType: answers.IntegrationLevel,
			Credentials: api.GcpCredentials{
				ClientID:     answers.ClientID,
				ClientEmail:  answers.ClientEmail,
				PrivateKeyID: answers.PrivateKeyID,
				PrivateKey:   answers.PrivateKey,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.Integrations.CreateGcp(gcp)
	cli.StopProgress()
	return err
}

func createGcpAuditLogIntegration() error {
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
			Name:     "private_key_id",
			Prompt:   &survey.Input{Message: "Private Key ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "client_email",
			Prompt:   &survey.Input{Message: "Client Email:"},
			Validate: survey.Required,
		},
		{
			Name:     "private_key",
			Prompt:   &survey.Input{Message: "Private Key:"},
			Validate: survey.Required,
		},
		{
			Name: "integration_level",
			Prompt: &survey.Select{
				Message: "Integration Level:",
				Options: []string{
					api.GcpProjectIntegration.String(),
					api.GcpOrganizationIntegration.String(),
				},
			},
			Validate: survey.Required,
		},
		{
			Name:     "org_project_id",
			Prompt:   &survey.Input{Message: "Organization/Project ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "subscription_name",
			Prompt:   &survey.Input{Message: "Subscription Name:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name             string
		ClientID         string `survey:"client_id"`
		PrivateKeyID     string `survey:"private_key_id"`
		ClientEmail      string `survey:"client_email"`
		PrivateKey       string `survey:"private_key"`
		IntegrationLevel string `survey:"integration_level"`
		OrgProjectID     string `survey:"org_project_id"`
		SubscriptionName string `survey:"subscription_name"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	gcp := api.NewGcpAuditLogIntegration(answers.Name,
		api.GcpIntegrationData{
			ID:               answers.OrgProjectID,
			IDType:           answers.IntegrationLevel,
			SubscriptionName: answers.SubscriptionName,
			Credentials: api.GcpCredentials{
				ClientID:     answers.ClientID,
				ClientEmail:  answers.ClientEmail,
				PrivateKeyID: answers.PrivateKeyID,
				PrivateKey:   answers.PrivateKey,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.Integrations.CreateGcp(gcp)
	cli.StopProgress()
	return err
}
