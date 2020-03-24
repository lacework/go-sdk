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

import "fmt"

// NewAwsIntegration returns an instance of awsIntegration with the provided
// integration type, name and data. The type can only be AwsCfgIntegration or
// AwsCloudTrailIntegration
//
// Basic usage: Initialize a new awsIntegration struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   aws, err := api.NewAwsIntegration("foo",
//     api.AwsCfgIntegration,
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
//   client.Integrations.CreateAws(aws)
//
func NewAwsIntegration(name string, iType integrationType, data AwsIntegrationData) awsIntegration {
	return awsIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    iType.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// NewAwsCfgIntegration returns an instance of awsIntegration of type AWS_CFG
func NewAwsCfgIntegration(name string, data AwsIntegrationData) awsIntegration {
	return NewAwsIntegration(name, AwsCfgIntegration, data)
}

// NewAwsCloudTrailIntegration returns an instance of awsIntegration of type AWS_CT_SQS
func NewAwsCloudTrailIntegration(name string, data AwsIntegrationData) awsIntegration {
	return NewAwsIntegration(name, AwsCloudTrailIntegration, data)
}

// CreateAws creates a single AWS integration on the Lacework Server
func (svc *IntegrationsService) CreateAws(integration awsIntegration) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAws gets a single AWS integration matching the integration guid on
// the Lacework Server
func (svc *IntegrationsService) GetAws(guid string) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAws updates a single AWS integration on the Lacework Server
func (svc *IntegrationsService) UpdateAws(data awsIntegration) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// DeleteAws deletes a single AWS integration matching the integration guid on
// the Lacework Server
func (svc *IntegrationsService) DeleteAws(guid string) (
	response awsIntegrationsResponse,
	err error,
) {
	err = svc.delete(guid, &response)
	return
}

// ListAwsCfg lists the AWS_CFG external integrations available on the Lacework Server
func (svc *IntegrationsService) ListAwsCfg() (response awsIntegrationsResponse, err error) {
	apiPath := fmt.Sprintf(apiIntegrationsByType, AwsCfgIntegration.String())
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

// ListAwsCloudTrail lists the AWS_CT_SQS external integrations available on the Lacework Server
func (svc *IntegrationsService) ListAwsCloudTrail() (response awsIntegrationsResponse, err error) {
	err = svc.listByType(AwsCloudTrailIntegration, &response)
	return
}

type awsIntegrationsResponse struct {
	Data    []awsIntegration `json:"data"`
	Ok      bool             `json:"ok"`
	Message string           `json:"message"`
}

type awsIntegration struct {
	commonIntegrationData
	Data AwsIntegrationData `json:"DATA"`
}

type AwsIntegrationData struct {
	Credentials AwsIntegrationCreds `json:"CROSS_ACCOUNT_CREDENTIALS"`

	// QueueUrl is a field that exists and is required for the AWS_CT_SQS integration,
	// though, it doesn't exist for AWS_CFG integrations, that's why we omit it if empty
	QueueUrl string `json:"QUEUE_URL,omitempty"`
}

type AwsIntegrationCreds struct {
	RoleArn    string `json:"ROLE_ARN"`
	ExternalId string `json:"EXTERNAL_ID"`
}
