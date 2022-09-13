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

func createAwsGovCloudConfigIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "account_id",
			Prompt:   &survey.Input{Message: "AWS Account ID: "},
			Validate: survey.Required,
		},
		{
			Name:     "access_key_id",
			Prompt:   &survey.Input{Message: "Access Key ID: "},
			Validate: survey.Required,
		},
		{
			Name:     "secret_access_key",
			Prompt:   &survey.Password{Message: "Secret Access Key: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name            string
		AccountID       string `survey:"account_id"`
		AccessKeyID     string `survey:"access_key_id"`
		SecretAccessKey string `survey:"secret_access_key"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	awsCfg := api.NewCloudAccount(answers.Name,
		api.AwsUsGovCfgCloudAccount,
		api.AwsUsGovCfgData{
			Credentials: api.AwsUsGovCfgCredentials{
				AwsAccountID:    answers.AccountID,
				AccessKeyID:     answers.AccessKeyID,
				SecretAccessKey: answers.SecretAccessKey,
			},
		},
	)
	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.CloudAccounts.Create(awsCfg)
	cli.StopProgress()
	return err
}

func createAwsGovCloudCTIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "account_id",
			Prompt:   &survey.Input{Message: "AWS Account ID: "},
			Validate: survey.Required,
		},
		{
			Name:     "access_key_id",
			Prompt:   &survey.Input{Message: "Access Key ID: "},
			Validate: survey.Required,
		},
		{
			Name:     "secret_access_key",
			Prompt:   &survey.Password{Message: "Secret Access Key: "},
			Validate: survey.Required,
		},
		{
			Name:     "queue_url",
			Prompt:   &survey.Input{Message: "SQS Queue URL:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name            string
		AccountID       string `survey:"account_id"`
		AccessKeyID     string `survey:"access_key_id"`
		SecretAccessKey string `survey:"secret_access_key"`
		QueueUrl        string `survey:"queue_url"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	awsCfg := api.NewCloudAccount(answers.Name,
		api.AwsUsGovCtSqsCloudAccount,
		api.AwsUsGovCtSqsData{
			QueueUrl: answers.QueueUrl,
			Credentials: api.AwsUsGovCtSqsCredentials{
				AwsAccountID:    answers.AccountID,
				AccessKeyID:     answers.AccessKeyID,
				SecretAccessKey: answers.SecretAccessKey,
			},
		},
	)
	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.CloudAccounts.Create(awsCfg)
	cli.StopProgress()
	return err
}
