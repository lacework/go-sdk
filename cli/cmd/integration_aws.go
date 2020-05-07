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

func createAwsConfigIntegration(lacework *api.Client) error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "role_arn",
			Prompt:   &survey.Input{Message: "Role ARN: "},
			Validate: survey.Required,
		},
		{
			Name:     "external_id",
			Prompt:   &survey.Input{Message: "External ID: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name       string
		RoleArn    string `survey:"role_arn"`
		ExternalID string `survey:"external_id"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	aws := api.NewAwsCfgIntegration(answers.Name,
		api.AwsIntegrationData{
			Credentials: api.AwsIntegrationCreds{
				RoleArn:    answers.RoleArn,
				ExternalID: answers.ExternalID,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = lacework.Integrations.CreateAws(aws)
	cli.StopProgress()
	return err
}

func createAwsCloudTrailIntegration(lacework *api.Client) error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "role_arn",
			Prompt:   &survey.Input{Message: "Role ARN:"},
			Validate: survey.Required,
		},
		{
			Name:     "external_id",
			Prompt:   &survey.Input{Message: "External ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "queue_url",
			Prompt:   &survey.Input{Message: "SQS Queue URL:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name       string
		RoleArn    string `survey:"role_arn"`
		ExternalID string `survey:"external_id"`
		QueueUrl   string `survey:"queue_url"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	aws := api.NewAwsCloudTrailIntegration(answers.Name,
		api.AwsIntegrationData{
			QueueUrl: answers.QueueUrl,
			Credentials: api.AwsIntegrationCreds{
				RoleArn:    answers.RoleArn,
				ExternalID: answers.ExternalID,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = lacework.Integrations.CreateAws(aws)
	cli.StopProgress()
	return err
}
