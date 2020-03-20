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
//   integrationResponse, err := api.CreateAwsConfigIntegration(aws)
//   if err != nil {
//     return err
//   }
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

// CreateAwsConfigIntegration creates a single AWS_CFG integration on the Lacework Server
func (c *Client) CreateAwsConfigIntegration(integration awsConfigIntegration) (response awsIntegrationsResponse, err error) {
	err = c.createIntegration(integration, &response)
	return
}

// GetAwsConfigIntegration gets a single AWS_CFG integration matching the integration guid available on the server
func (c *Client) GetAwsConfigIntegration(intgGuid string) (response awsIntegrationsResponse, err error) {
	err = c.getIntegration(intgGuid, &response)
	return
}

// UpdateAwsConfigIntegration updates a single AWS_CFG integration on the Lacework Server
func (c *Client) UpdateAwsConfigIntegration(data awsConfigIntegration) (response awsIntegrationsResponse, err error) {
	err = c.updateIntegration(data.IntgGuid, data, &response)
	return
}

// DeleteAwsConfigIntegration deletes a single integration matching the integration guid on the Lacework Server
func (c *Client) DeleteAwsConfigIntegration(intgGuid string) (response awsIntegrationsResponse, err error) {
	err = c.deleteIntegration(intgGuid, &response)
	return
}

func (c *Client) GetAwsIntegrations() (response awsIntegrationsResponse, err error) {
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
	Credentials  AwsIntegrationCreds `json:"CROSS_ACCOUNT_CREDENTIALS"`
	AwsAccountId string              `json:"AWS_ACCOUNT_ID"`
}

type AwsIntegrationCreds struct {
	RoleArn    string `json:"ROLE_ARN"`
	ExternalId string `json:"EXTERNAL_ID"`
}
