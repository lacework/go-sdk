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

func createGarIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
		{
			Name: "domain",
			Prompt: &survey.Select{
				Message: "Registry Domain:",
				Options: []string{
					"northamerica-northeast1-docker.pkg.dev",
					"us-central1-docker.pkg.dev",
					"us-east1-docker.pkg.dev",
					"us-east4-docker.pkg.dev",
					"us-west1-docker.pkg.dev",
					"us-west2-docker.pkg.dev",
					"us-west3-docker.pkg.dev",
					"us-west4-docker.pkg.dev",
					"southamerica-east1-docker.pkg.dev",
					"europe-north1-docker.pkg.dev",
					"europe-west1-docker.pkg.dev",
					"europe-west2-docker.pkg.dev",
					"europe-west3-docker.pkg.dev",
					"europe-west4-docker.pkg.dev",
					"europe-west6-docker.pkg.dev",
					"asia-east1-docker.pkg.dev",
					"asia-east2-docker.pkg.dev",
					"asia-northeast1-docker.pkg.dev",
					"asia-northeast2-docker.pkg.dev",
					"asia-northeast3-docker.pkg.dev",
					"asia-south1-docker.pkg.dev",
					"asia-southeast1-docker.pkg.dev",
					"asia-southeast2-docker.pkg.dev",
					"australia-southeast1-docker.pkg.dev",
					"asia-docker.pkg.dev",
					"europe-docker.pkg.dev",
					"us-docker.pkg.dev",
				},
				Default: "us-west1-docker.pkg.dev",
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
			Name: "limit_max_images",
			Prompt: &survey.Select{
				Message: "Limit number of images per repository: ",
				Options: []string{"5", "10", "15"},
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name           string
		Domain         string
		ClientID       string `survey:"client_id"`
		PrivateKeyID   string `survey:"private_key_id"`
		ClientEmail    string `survey:"client_email"`
		PrivateKey     string `survey:"private_key"`
		LimitTags      string `survey:"limit_tags"`
		LimitLabels    string `survey:"limit_labels"`
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

	// @afiune these are the new API v2 limits
	if err := askForV2Limits(&answers); err != nil {
		return err
	}

	gar := api.NewContainerRegistry(answers.Name,
		api.GcpGarContainerRegistry,
		api.GcpGarData{
			Credentials: api.GcpCredentialsV2{
				ClientEmail:  answers.ClientEmail,
				ClientID:     answers.ClientID,
				PrivateKey:   answers.PrivateKey,
				PrivateKeyID: answers.PrivateKeyID,
			},
			RegistryDomain: answers.Domain,
			LimitByTag:     strings.Split(answers.LimitTags, "\n"),
			LimitByLabel:   castStringToLimitByLabel(answers.LimitLabels),
			LimitByRep:     strings.Split(answers.LimitRepos, "\n"),
			LimitNumImg:    limitMaxImages,
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.ContainerRegistries.Create(gar)
	cli.StopProgress()
	return err
}
