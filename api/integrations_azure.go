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

package api

// NewAzureIntegration returns an instance of azureIntegration
//
// Basic usage: Initialize a new azureIntegration struct, then
//              use the new instance to do CRUD operations
//
//   azure, err := api.NewAzureIntegration("bar",
//     api.AzureIntegrationData{},
//   )
//   if err != nil {
//     return err
//   }
//
//   integrationResponse, err := api.CreateAzureConfigIntegration(azure)
//   if err != nil {
//     return err
//   }
//
func NewAzureIntegration(name string, data AzureIntegrationData) azureIntegration {
	return azureIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    AzureCfgIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

func (c *Client) GetAzureIntegrations() (response azureIntegrationsResponse, err error) {
	return
}

type azureIntegrationsResponse struct {
	Data    []azureIntegration `json:"data"`
	Ok      bool               `json:"ok"`
	Message string             `json:"message"`
}

type azureIntegration struct {
	commonIntegrationData
	Data AzureIntegrationData `json:"DATA"`
}

type AzureIntegrationData struct{}
