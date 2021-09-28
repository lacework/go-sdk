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

func createGcrIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name: "domain",
			Prompt: &survey.Select{
				Message: "Registry Domain:",
				Options: []string{
					"gcr.io",
					"us.gcr.io",
					"eu.gcr.io",
					"asia.gcr.io",
				},
			},
			Validate: survey.Required,
		},
		{
			Name:     "client_id",
			Prompt:   &survey.Input{Message: "Client ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "private_key_id",
			Prompt:   &survey.Input{Message: "Private Key ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "client_email",
			Prompt:   &survey.Input{Message: "Client Email:"},
			Validate: survey.Required,
		},
		{
			Name:     "private_key",
			Prompt:   &survey.Editor{Message: "Enter properly formatted Private Key:"},
			Validate: survey.Required,
		},
		{
			Name: "non_os_package_support",
			Prompt: &survey.Confirm{
				Message: "Enable Scanning for non-os packages: "},
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
		Name                string
		Domain              string
		ClientID            string `survey:"client_id"`
		PrivateKeyID        string `survey:"private_key_id"`
		ClientEmail         string `survey:"client_email"`
		PrivateKey          string `survey:"private_key"`
		NonOSPackageSupport bool   `survey:"non_os_package_support"`
		LimitTag            string `survey:"limit_tag"`
		LimitLabel          string `survey:"limit_label"`
		LimitRepos          string `survey:"limit_repos"`
		LimitMaxImages      string `survey:"limit_max_images"`
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

	gcr := api.NewGcrRegistryIntegration(answers.Name,
		api.ContainerRegData{
			Credentials: api.ContainerRegCreds{
				ClientEmail:  answers.ClientEmail,
				ClientID:     answers.ClientID,
				PrivateKey:   answers.PrivateKey,
				PrivateKeyID: answers.PrivateKeyID,
			},
			RegistryDomain:   answers.Domain,
			NonOSPackageEval: answers.NonOSPackageSupport,
			LimitByTag:       answers.LimitTag,
			LimitByLabel:     answers.LimitLabel,
			LimitByRep:       answers.LimitRepos,
			LimitNumImg:      limitMaxImages,
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.Integrations.CreateContainerRegistry(gcr)
	cli.StopProgress()
	return err
}
