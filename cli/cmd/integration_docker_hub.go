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
	"github.com/lacework/go-sdk/api"
)

func createDockerHubIntegration(lacework *api.Client) error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
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
	}

	answers := struct {
		Name           string
		Username       string
		Password       string
		LimitTag       string `survey:"limit_tag"`
		LimitLabel     string `survey:"limit_label"`
		LimitRepos     string `survey:"limit_repos"`
		LimitMaxImages string `survey:"limit_max_images"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = "▸"
		}),
	)
	if err != nil {
		return err
	}

	// @afiune we should not error here since the survey.Select is letting
	// the user choose 5, 10, or 15. Though, we always check errors!
	// TODO can we do this via `Transformer`?
	limitMaxImages, err := strconv.Atoi(answers.LimitMaxImages)
	if err != nil {
		cli.Log.Warnw("unable to convert limit_max_images, using default",
			"error", err,
			"input", answers.LimitMaxImages,
			"default", "5",
		)
		limitMaxImages = 5
	}

	docker := api.NewContainerRegIntegration(answers.Name,
		api.ContainerRegData{
			Credentials: api.ContainerRegCreds{
				Username: answers.Username,
				Password: answers.Password,
			},
			RegistryType:   api.DockerHubRegistry.String(),
			RegistryDomain: "index.docker.io",
			LimitByTag:     answers.LimitTag,
			LimitByLabel:   answers.LimitLabel,
			LimitByRep:     answers.LimitRepos,
			LimitNumImg:    limitMaxImages,
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = lacework.Integrations.CreateContainerRegistry(docker)
	cli.StopProgress()
	return err
}
