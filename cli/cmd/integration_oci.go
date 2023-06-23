//
// Author:: Kolbeinn Karlsson (<kolbeinn.karlsson@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

func createOciConfigIntegration() error {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "tenant_id",
			Prompt:   &survey.Input{Message: "Tenant ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "tenant_name",
			Prompt:   &survey.Input{Message: "Tenant Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "home_region",
			Prompt:   &survey.Input{Message: "Home Region:"},
			Validate: survey.Required,
		},
		{
			Name:     "user_ocid",
			Prompt:   &survey.Input{Message: "User OCID:"},
			Validate: survey.Required,
		},
		{
			Name:     "fingerprint",
			Prompt:   &survey.Input{Message: "Public Key Fingerprint:"},
			Validate: survey.Required,
		},
		{
			Name:     "private_key",
			Prompt:   &survey.Input{Message: "Private Key:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name        string
		TenantID    string `survey:"tenant_id"`
		TenantName  string `survey:"tenant_name"`
		HomeRegion  string `survey:"home_region"`
		UserOCID    string `survey:"user_ocid"`
		Fingerprint string `survey:"fingerprint"`
		PrivateKey  string `survey:"private_key"`
	}{}

	err := survey.Ask(questions, &answers, survey.WithIcons(promptIconsFunc))
	if err != nil {
		return err
	}

	oci := api.NewCloudAccount(
		answers.Name,
		api.OciCfgCloudAccount,
		api.OciCfgData{
			TenantID:   answers.TenantID,
			TenantName: answers.TenantName,
			HomeRegion: answers.HomeRegion,
			UserOCID:   answers.UserOCID,
			Credentials: api.OciCfgCredentials{
				Fingerprint: answers.Fingerprint,
				PrivateKey:  answers.PrivateKey,
			},
		},
	)

	cli.StartProgress(" Creating integration...")
	_, err = cli.LwApi.V2.CloudAccounts.Create(oci)
	cli.StopProgress()
	return err
}
