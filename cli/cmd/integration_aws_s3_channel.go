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

func createAwsS3ChannelIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "role_arn",
			Prompt:   &survey.Input{Message: "Role ARN:"},
			Validate: survey.Required,
		},
		{
			Name:     "bucket_arn",
			Prompt:   &survey.Input{Message: "Bucket ARN:"},
			Validate: survey.Required,
		},
		{
			Name:     "external_id",
			Prompt:   &survey.Input{Message: "External ID:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name       string
		RoleArn    string `survey:"role_arn"`
		BucketArn  string `survey:"bucket_arn"`
		ExternalID string `survey:"external_id"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	s3 := api.NewAlertChannel(answers.Name,
		api.AwsS3AlertChannelType,
		api.AwsS3DataV2{
			Credentials: api.AwsS3Credentials{
				RoleArn:    answers.RoleArn,
				BucketArn:  answers.BucketArn,
				ExternalID: answers.ExternalID,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.AlertChannels.Create(s3)
	cli.StopProgress()
	return err
}
