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

type jiraAlertChannelIntegrationSurvey struct {
	Name     string
	Url      string
	Issue    string
	Project  string
	Username string
	Token    string
	Password string
}

func createJiraCloudAlertChannelIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "url",
			Prompt:   &survey.Input{Message: "Jira URL: "},
			Validate: survey.Required,
		},
		{
			Name:     "issue",
			Prompt:   &survey.Input{Message: "Issue Type: "},
			Validate: survey.Required,
		},
		{
			Name:     "project",
			Prompt:   &survey.Input{Message: "Project Key: "},
			Validate: survey.Required,
		},
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "Username: "},
			Validate: survey.Required,
		},
		{
			Name:     "token",
			Prompt:   &survey.Password{Message: "API Token: "},
			Validate: survey.Required,
		},
	}

	var answers jiraAlertChannelIntegrationSurvey
	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	jira := api.NewJiraCloudAlertChannel(answers.Name,
		api.JiraAlertChannelData{
			JiraUrl:   answers.Url,
			IssueType: answers.Issue,
			ProjectID: answers.Project,
			Username:  answers.Username,
			ApiToken:  answers.Token,
		},
	)

	return createJiraAlertChannelIntegration(jira)
}

func createJiraServerAlertChannelIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "url",
			Prompt:   &survey.Input{Message: "Jira URL: "},
			Validate: survey.Required,
		},
		{
			Name:     "issue",
			Prompt:   &survey.Input{Message: "Issue Type: "},
			Validate: survey.Required,
		},
		{
			Name:     "project",
			Prompt:   &survey.Input{Message: "Project Key: "},
			Validate: survey.Required,
		},
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "Username: "},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Password: "},
			Validate: survey.Required,
		},
	}

	var answers jiraAlertChannelIntegrationSurvey
	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	jira := api.NewJiraServerAlertChannel(answers.Name,
		api.JiraAlertChannelData{
			JiraUrl:   answers.Url,
			IssueType: answers.Issue,
			ProjectID: answers.Project,
			Username:  answers.Username,
			Password:  answers.Password,
		},
	)
	return createJiraAlertChannelIntegration(jira)
}

func createJiraAlertChannelIntegration(jira api.JiraAlertChannel) error {
	cli.StartProgress(" Creating integration...")
	_, err := cli.LwApi.Integrations.CreateJiraAlertChannel(jira)
	cli.StopProgress()
	return err
}
