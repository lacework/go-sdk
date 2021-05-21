//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

package api

// NewAwsGovCloudIntegration returns an instance of AwsGovCloudIntegration with the provided
// integration type, name and data.
//
// Basic usage: Initialize a new AwsGovCloudIntegration struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   aws := api.NewAwsGovCloudIntegration("foo",
//     api.AwsGovCloudIntegration,
//     api.AwsGovCloudIntegrationData{
//       Credentials: api.AwsGovCloudCreds {
//         AccountID: "553453453",
//         AccessKeyID: "AWS123abcAccessKeyID",
//         SecretAccessKey: "AWS123abc123abcSecretAccessKey0000000000",
//       },
//     },
//   )
//
//   client.Integrations.CreateAwsGovCloudCfg(aws)
//
func NewAwsGovCloudCfgIntegration(name string, iType integrationType, data AwsGovCloudIntegrationData) AwsGovCloudIntegration {
	return AwsGovCloudIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    iType.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateAwsGovCloud creates a single AWS Gov Cloud integration on the Lacework Server
func (svc *IntegrationsService) CreateAwsGovCloudCfg(integration AwsGovCloudIntegration) (
	response AwsGovCloudIntegrationsResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAwsGovCloudCfg gets a single AWS Gov Cloud integration matching the integration guid on
// the Lacework Server
func (svc *IntegrationsService) GetAwsGovCloudCfg(guid string) (
	response AwsGovCloudIntegrationsResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsGovCloudCfg updates a single AWS Gov Cloud integration on the Lacework Server
func (svc *IntegrationsService) UpdateAwsGovCloudCfg(data AwsGovCloudIntegration) (
	response AwsGovCloudIntegrationsResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// DeleteAwsGovCloudCfg deletes a single AWS Gov Cloud Cfg integration matching the integration guid on
// the Lacework Server
func (svc *IntegrationsService) DeleteAwsGovCloudCfg(guid string) (
	response AwsGovCloudIntegrationsResponse,
	err error,
) {
	err = svc.delete(guid, &response)
	return
}

// ListAwsGovCloudCfg lists the AWS_US_GOV_CFG external integrations available on the Lacework Server
func (svc *IntegrationsService) ListAwsGovCloudCfg() (response AwsIntegrationsResponse, err error) {
	err = svc.listByType(AwsGovCloudCfgIntegration, &response)
	return
}

type AwsGovCloudIntegrationsResponse struct {
	Data    []AwsGovCloudIntegration `json:"data"`
	Ok      bool                     `json:"ok"`
	Message string                   `json:"message"`
}

type AwsGovCloudIntegration struct {
	commonIntegrationData
	Data AwsGovCloudIntegrationData `json:"DATA"`
}

type AwsGovCloudIntegrationData struct {
	Credentials AwsGovCloudCreds `json:"ACCESS_KEY_CREDENTIALS" mapstructure:"ACCESS_KEY_CREDENTIALS"`
}

type AwsGovCloudCreds struct {
	AccountID       string `json:"ACCOUNT_ID" mapstructure:"ACCOUNT_ID"`
	AccessKeyID     string `json:"ACCESS_KEY_ID" mapstructure:"ACCESS_KEY_ID"`
	SecretAccessKey string `json:"SECRET_ACCESS_KEY" mapstructure:"SECRET_ACCESS_KEY"`
}
