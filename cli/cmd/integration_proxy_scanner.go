//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

func createProxyScannerIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
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
			Prompt: &survey.Multiline{
				Message: "List of 'key:value' labels:",
				Default: "*",
			},
		},
		{
			Name: "limit_repos",
			Prompt: &survey.Input{
				Message: "Limit by Repository: ",
				Default: "*",
			},
		},
		{
			Name: "limit_max_images",
			Prompt: &survey.Select{
				Message: "Limit Number of Images per Repo: ",
				Options: []string{
					"5", "10", "15",
				},
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name           string
		LimitTag       string `survey:"limit_tag"`
		LimitLabel     string `survey:"limit_label"`
		LimitRepos     string `survey:"limit_repos"`
		LimitMaxImages string `survey:"limit_max_images"`
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

	proxy := api.NewContainerRegistry(
		answers.Name,
		api.ProxyScannerContainerRegistry,
		api.ProxyScannerData{
			RegistryType: api.ProxyScannerContainerRegistry.String(),
			LimitByTag:   strings.Split(answers.LimitTag, "\n"),
			LimitByLabel: castStringToLimitByLabel(answers.LimitLabel),
			LimitByRep:   strings.Split(answers.LimitRepos, "\n"),
			LimitNumImg:  limitMaxImages,
		},
	)

	cli.StartProgress("Creating integration...")
	_, err = cli.LwApi.V2.ContainerRegistries.Create(proxy)
	cli.StopProgress()
	return err
}
