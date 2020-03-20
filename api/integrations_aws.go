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

// NewAwsConfigIntegration returns an instance of awsConfigIntegration
//
// Basic usage: Initialize a new awsConfigIntegration struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   aws, err := api.NewAwsConfigIntegration("foo",
//     api.AwsIntegrationData{
//       Credentials: api.AwsIntegrationCreds {
//         RoleArn: "arn:aws:XYZ",
//         ExternalId: "1",
//       },
//     },
//   )
//   if err != nil {
//     return err
//   }
//
//   client.Integrations.CreateAwsConfig(aws)
//
func NewAwsConfigIntegration(name string, data AwsIntegrationData) awsConfigIntegration {
	return awsConfigIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    AwsCfgIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateAwsConfig creates a single AWS_CFG integration on the Lacework Server
func (svc *IntegrationsService) CreateAwsConfig(integration awsConfigIntegration) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAwsConfig gets a single AWS_CFG integration matching the integration guid
// on the Lacwork Server
func (svc *IntegrationsService) GetAwsConfig(guid string) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsConfig updates a single AWS_CFG integration on the Lacework Server
func (svc *IntegrationsService) UpdateAwsConfig(data awsConfigIntegration) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// DeleteAwsConfig deletes a single AWS_CFG integration matching the integration guid
// on the Lacework Server
func (svc *IntegrationsService) DeleteAwsConfig(guid string) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.delete(guid, &response)
	return
}

type awsIntegrationsResponse struct {
	Data    []awsConfigIntegration `json:"data"`
	Ok      bool                   `json:"ok"`
	Message string                 `json:"message"`
}

type awsConfigIntegration struct {
	commonIntegrationData
	Data AwsIntegrationData `json:"DATA"`
}

type AwsIntegrationData struct {
	Credentials AwsIntegrationCreds `json:"CROSS_ACCOUNT_CREDENTIALS"`
}

type AwsIntegrationCreds struct {
	RoleArn    string `json:"ROLE_ARN"`
	ExternalId string `json:"EXTERNAL_ID"`
}
