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
	"sort"

	"github.com/AlecAivazis/survey/v2"

	"github.com/lacework/go-sdk/api"
)

type jiraAlertChannelIntegrationSurvey struct {
	Name          string
	Url           string
	Issue         string
	Project       string
	Username      string
	Token         string
	Password      string
	Grouping      string
	Bidirectional bool
}

func getJiraGroupingOptions() []string {
	options := make([]string, 0, len(api.JiraIssueGroupingsSurvey))

	for option := range api.JiraIssueGroupingsSurvey {
		options = append(options, option)
	}

	sort.SliceStable(options, func(i, j int) bool {
		return api.JiraIssueGroupingsSurvey[options[i]] < api.JiraIssueGroupingsSurvey[options[j]]
	})

	return options
}

func createJiraAlertChannelIntegration(jiraType string) error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name: "bidirectional",
			Prompt: &survey.Confirm{
				Message: "Would you like a bidirectional integration?",
				Default: false,
				Help:    "See https://docs.lacework.com/onboarding/jira#bidirectional-integration for more detail.",
			},
			Validate: survey.Required,
		},
		{
			Name: "grouping",
			Prompt: &survey.Select{
				Message: "Group Issues by:",
				Options: getJiraGroupingOptions(),
			},
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
	}

	switch jiraType {
	case api.JiraCloudAlertType, "":
		jiraType = api.JiraCloudAlertType
		questions = append(questions, &survey.Question{
			Name:     "token",
			Prompt:   &survey.Password{Message: "API Token: "},
			Validate: survey.Required,
		})
	case api.JiraServerAlertType:
		questions = append(questions, &survey.Question{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Password: "},
			Validate: survey.Required,
		})
	}

	var answers jiraAlertChannelIntegrationSurvey
	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	grouping := api.JiraIssueGroupingsSurvey[answers.Grouping]
	jira := api.JiraDataV2{
		ApiToken:      answers.Token,
		IssueGrouping: grouping.String(),
		IssueType:     answers.Issue,
		JiraType:      jiraType,
		JiraUrl:       answers.Url,
		ProjectID:     answers.Project,
		Username:      answers.Username,
		Password:      answers.Password,
	}
	if answers.Bidirectional {
		jira.Configuration = api.BidirectionalJiraConfiguration
	}

	// ask the user if they would like to configure a Custom Template
	custom := false
	err = survey.AskOne(&survey.Confirm{
		Message: "Configure a Custom Template File?",
	}, &custom)

	if err != nil {
		return err
	}

	if custom {
		var content string

		err = survey.AskOne(&survey.Editor{
			Message:  "Provide the Custom Template File in JSON format",
			FileName: "*.json",
		}, &content)

		if err != nil {
			return err
		}

		jira.EncodeCustomTemplateFile(content)
	}

	jiraCloudAlertChan := api.NewAlertChannel(answers.Name, api.JiraAlertChannelType, jira)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.AlertChannels.Create(jiraCloudAlertChan)
	cli.StopProgress()
	return err
}
