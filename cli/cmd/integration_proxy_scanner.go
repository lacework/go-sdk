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
	}

	answers := struct {
		Name string
	}{}

	if err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	); err != nil {
		return err
	}

	// default values
	repositoriesLimit := make([]string, 0)
	tagsLimit := make([]string, 0)
	labelLimit := make([]map[string]string, 0)
	limitNumScan := 5

	proxy := api.NewContainerRegistry(
		answers.Name,
		api.ProxyScannerContainerRegistry,
		api.ProxyScannerData{
			RegistryType: api.ProxyScannerContainerRegistry.String(),
			LimitNumImg:  limitNumScan,
			LimitByRep:   repositoriesLimit,
			LimitByTag:   tagsLimit,
			LimitByLabel: labelLimit,
		},
	)

	cli.StartProgress("Creating integration...")
	_, err := cli.LwApi.V2.ContainerRegistries.Create(proxy)
	cli.StopProgress()
	return err
}
