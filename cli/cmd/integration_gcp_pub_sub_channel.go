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

func createGcpPubSubChannelIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "topic_id",
			Prompt:   &survey.Input{Message: "Topic ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "project_id",
			Prompt:   &survey.Input{Message: "Project ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "client_id",
			Prompt:   &survey.Input{Message: "Client ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "client_email",
			Prompt:   &survey.Input{Message: "Client Email:"},
			Validate: survey.Required,
		},
		{
			Name:     "private_key_id",
			Prompt:   &survey.Input{Message: "Private Key ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "private_key",
			Prompt:   &survey.Editor{Message: "Enter properly formatted Private Key:"},
			Validate: survey.Required,
		},
		{
			Name: "issue_grouping",
			Prompt: &survey.Select{Message: "Issue Grouping:",
				Options: []string{"Events", "Resources"},
			},
		},
	}

	answers := struct {
		Name          string
		ClientID      string `survey:"client_id"`
		PrivateKeyID  string `survey:"private_key_id"`
		ClientEmail   string `survey:"client_email"`
		PrivateKey    string `survey:"private_key"`
		ProjectID     string `survey:"project_id"`
		TopicID       string `survey:"topic_id"`
		IssueGrouping string `survey:"issue_grouping"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	gcp := api.NewAlertChannel(answers.Name,
		api.GcpPubSubAlertChannelType,
		api.GcpPubSubDataV2{
			ProjectId:     answers.ProjectID,
			TopicId:       answers.TopicID,
			IssueGrouping: answers.IssueGrouping,
			Credentials: api.GcpPubSubCredentials{
				ClientId:     answers.ClientID,
				ClientEmail:  answers.ClientEmail,
				PrivateKeyId: answers.PrivateKeyID,
				PrivateKey:   answers.PrivateKey,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.AlertChannels.Create(gcp)
	cli.StopProgress()
	return err
}
