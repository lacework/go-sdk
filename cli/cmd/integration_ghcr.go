//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/lacework/go-sdk/api"
)

func createGhcrIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "Username:"},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Password:"},
			Validate: survey.Required,
		},
		{
			Name:   "ssl",
			Prompt: &survey.Confirm{Message: "Enable SSL?"},
		},
		{
			Name:   "notifications",
			Prompt: &survey.Confirm{Message: "Subscribe to Registry Notifications?"},
		},
		{
			Name: "non_os_package_support",
			Prompt: &survey.Confirm{
				Message: "Enable Scanning for non-os packages: "},
		},
		{
			Name: "limit_max_images",
			Prompt: &survey.Select{
				Message: "Limit number of images per repository: ",
				Options: []string{"5", "10", "15"},
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name                string
		Username            string
		Password            string
		SSL                 bool
		Notifications       bool
		NonOSPackageSupport bool   `survey:"non_os_package_support"`
		LimitTags           string `survey:"limit_tags"`
		LimitLabels         string `survey:"limit_labels"`
		LimitRepos          string `survey:"limit_repos"`
		LimitMaxImages      string `survey:"limit_max_images"`
	}{}

	if err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	); err != nil {
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

	// @afiune these are the new API v2 limits
	if err := askV2Limits(&answers); err != nil {
		return err
	}

	ghcr := api.NewContainerRegistry(answers.Name,
		api.GhcrContainerRegistry,
		api.GhcrData{
			Credentials: api.GhcrCredentials{
				Username: answers.Username,
				Password: answers.Password,
				Ssl:      answers.SSL,
			},
			NonOSPackageEval: answers.NonOSPackageSupport,
			LimitByTag:       strings.Split(answers.LimitTags, "\n"),
			LimitByLabel:     castStringToLimitByLabel(answers.LimitLabels),
			LimitByRep:       strings.Split(answers.LimitRepos, "\n"),
			LimitNumImg:      limitMaxImages,
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.ContainerRegistries.Create(ghcr)
	cli.StopProgress()
	return err
}
