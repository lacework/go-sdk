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
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
)

func createAwsEcrIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "domain",
			Prompt:   &survey.Input{Message: "Registry Domain: "},
			Validate: survey.Required,
		},
		{
			Name: "non_os_package_support",
			Prompt: &survey.Confirm{
				Message: "Enable scanning for Non-OS packages: "},
		},
		{
			Name: "limit_tag",
			Prompt: &survey.Input{
				Message: "Limit by Tag: ",
				Default: "*",
			},
		},
		{
			Name: "limit_label",
			Prompt: &survey.Input{
				Message: "Limit by Label: ",
				Default: "*",
			},
		},
		{
			Name:   "limit_repos",
			Prompt: &survey.Input{Message: "Limit by Repository: "},
		},
		{
			Name: "limit_max_images",
			Prompt: &survey.Select{
				Message: "Limit Number of Images per Repo: ",
				Options: []string{"5", "10", "15"},
			},
			Validate: survey.Required,
		},
		{
			Name: "aws_auth_type",
			Prompt: &survey.Select{
				Message: "Authentication Type: ",
				Options: []string{"AWS IAM Role", "AWS Access Key"},
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name                string
		Domain              string
		AccessKeyID         string `survey:"access_key_id"`
		SecretAccessKey     string `survey:"secret_access_key"`
		LimitTag            string `survey:"limit_tag"`
		LimitLabel          string `survey:"limit_label"`
		LimitRepos          string `survey:"limit_repos"`
		LimitMaxImages      string `survey:"limit_max_images"`
		AwsAuthType         string `survey:"aws_auth_type"`
		NonOSPackageSupport bool   `survey:"non_os_package_support"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	limitMaxImages, err := strconv.Atoi(answers.LimitMaxImages)
	if err != nil {
		cli.Log.Warnw("unable to convert limit_max_images, using default",
			"error", err,
			"input", answers.LimitMaxImages,
			"default", "5",
		)
		limitMaxImages = 5
	}

	switch answers.AwsAuthType {

	case "AWS IAM Role":
		ecrAuthAnswers := struct {
			RoleArn    string `survey:"role_arn"`
			ExternalID string `survey:"external_id"`
		}{}

		questionsAuth := []*survey.Question{
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
		}

		err := survey.Ask(questionsAuth, &ecrAuthAnswers,
			survey.WithIcons(promptIconsFunc),
		)
		if err != nil {
			return err
		}

		ecr := api.NewAwsEcrWithCrossAccountIntegration(answers.Name,
			api.AwsEcrDataWithCrossAccountCreds{
				Credentials: api.AwsCrossAccountCreds{
					RoleArn:    ecrAuthAnswers.RoleArn,
					ExternalID: ecrAuthAnswers.ExternalID,
				},
				AwsEcrCommonData: api.AwsEcrCommonData{
					RegistryDomain:   answers.Domain,
					NonOSPackageEval: answers.NonOSPackageSupport,
					LimitByTag:       answers.LimitTag,
					LimitByLabel:     answers.LimitLabel,
					LimitByRep:       answers.LimitRepos,
					LimitNumImg:      limitMaxImages,
				},
			},
		)

		cli.StartProgress(" Creating integration...")
		_, err = cli.LwApi.Integrations.CreateAwsEcrWithCrossAccount(ecr)
		cli.StopProgress()
		return err

	case "AWS Access Key":
		ecrAuthAnswers := struct {
			AccessKeyID     string `survey:"access_key_id"`
			SecretAccessKey string `survey:"secret_access_key"`
		}{}

		questionsAuth := []*survey.Question{
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

		err := survey.Ask(questionsAuth, &ecrAuthAnswers,
			survey.WithIcons(promptIconsFunc),
		)
		if err != nil {
			return err
		}

		ecr := api.NewAwsEcrWithAccessKeyIntegration(answers.Name,
			api.AwsEcrDataWithAccessKeyCreds{
				Credentials: api.AwsEcrAccessKeyCreds{
					AccessKeyID:     ecrAuthAnswers.AccessKeyID,
					SecretAccessKey: ecrAuthAnswers.SecretAccessKey,
				},
				AwsEcrCommonData: api.AwsEcrCommonData{
					RegistryDomain: answers.Domain,
					LimitByTag:     answers.LimitTag,
					LimitByLabel:   answers.LimitLabel,
					LimitByRep:     answers.LimitRepos,
					LimitNumImg:    limitMaxImages,
				},
			},
		)

		cli.StartProgress(" Creating integration...")
		_, err = cli.LwApi.Integrations.CreateAwsEcrWithAccessKey(ecr)
		cli.StopProgress()
		return err

	default:
		return errors.New("unknown ECR authentication method")
	}
}
