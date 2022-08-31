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
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/lacework/go-sdk/api"
)

func createDockerV2Integration() error {
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
			Name:   "ssl",
			Prompt: &survey.Confirm{Message: "Enable SSL?"},
		},
		{
			Name:   "non_os_package_support",
			Prompt: &survey.Confirm{Message: "Enable scanning for Non-OS packages: "},
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
	}

	answers := struct {
		Name                string
		Domain              string
		Username            string
		Password            string
		SSL                 bool
		NonOSPackageSupport bool   `survey:"non_os_package_support"`
		LimitTag            string `survey:"limit_tag"`
		LimitLabel          string `survey:"limit_label"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	docker := api.NewContainerRegistry(answers.Name,
		api.DockerhubV2ContainerRegistry,
		api.DockerhubV2Data{
			Credentials: api.DockerhubV2Credentials{
				Username: answers.Username,
				Password: answers.Password,
				Ssl:      answers.SSL,
			},
			RegistryDomain:   answers.Domain,
			NonOSPackageEval: answers.NonOSPackageSupport,
			LimitByTag:       strings.Split(answers.LimitTag, "\n"),
			LimitByLabel:     castStringToLimitByLabel(answers.LimitLabel),
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.ContainerRegistries.Create(docker)
	cli.StopProgress()
	return err
}
