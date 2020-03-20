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

// NewAzureConfigIntegration returns an instance of azureConfigIntegration
//
// Basic usage: Initialize a new azureConfigIntegration struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   azure, err := api.NewAzureConfigIntegration("bar",
//     api.AzureIntegrationData{},
//   )
//   if err != nil {
//     return err
//   }
//
//   client.Integrations.CreateAzureConfig(azure)
//
func NewAzureConfigIntegration(name string, data AzureIntegrationData) azureConfigIntegration {
	return azureConfigIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    AzureCfgIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateAzureConfig creates a single AZURE_CFG integration on the Lacework Server
func (svc *IntegrationsService) CreateAzureConfig(integration awsConfigIntegration) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAzureConfig gets a single AZURE_CFG integration matching the integration guid
// on the Lacework Server
func (svc *IntegrationsService) GetAzureConfig(guid string) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAzureConfig updates a single AZURE_CFG integration on the Lacework Server
func (svc *IntegrationsService) UpdateAzureConfig(data awsConfigIntegration) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// DeleteAzureConfig deletes a single AZURE_CFG integration matching the integration
// guid on the Lacework Server
func (svc *IntegrationsService) DeleteAzureConfig(guid string) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.delete(guid, &response)
	return
}

type azureIntegrationsResponse struct {
	Data    []azureConfigIntegration `json:"data"`
	Ok      bool                     `json:"ok"`
	Message string                   `json:"message"`
}

type azureConfigIntegration struct {
	commonIntegrationData
	Data AzureIntegrationData `json:"DATA"`
}

type AzureIntegrationData struct{}
