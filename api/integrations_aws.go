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

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// NewAwsIntegration returns an instance of AwsIntegration with the provided
// integration type, name and data. The type can only be AwsCfgIntegration or
// AwsCloudTrailIntegration
//
// Basic usage: Initialize a new AwsIntegration struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   aws := api.NewAwsIntegration("foo",
//     api.AwsCfgIntegration,
//     api.AwsIntegrationData{
//       Credentials: api.AwsCrossAccountCreds {
//         RoleArn: "arn:aws:XYZ",
//         ExternalID: "1",
//       },
//     },
//   )
//
//   client.Integrations.CreateAws(aws)
//
func NewAwsIntegration(name string, iType integrationType, data AwsIntegrationData) AwsIntegration {
	return AwsIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    iType.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// NewAwsCfgIntegration returns an instance of AwsIntegration of type AWS_CFG
func NewAwsCfgIntegration(name string, data AwsIntegrationData) AwsIntegration {
	return NewAwsIntegration(name, AwsCfgIntegration, data)
}

// NewAwsCloudTrailIntegration returns an instance of AwsIntegration of type AWS_CT_SQS
func NewAwsCloudTrailIntegration(name string, data AwsIntegrationData) AwsIntegration {
	return NewAwsIntegration(name, AwsCloudTrailIntegration, data)
}

// CreateAws creates a single AWS integration on the Lacework Server
func (svc *IntegrationsService) CreateAws(integration AwsIntegration) (
	response AwsIntegrationsResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAws gets a single AWS integration matching the integration guid on
// the Lacework Server
func (svc *IntegrationsService) GetAws(guid string) (
	response AwsIntegrationsResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAws updates a single AWS integration on the Lacework Server
func (svc *IntegrationsService) UpdateAws(data AwsIntegration) (
	response AwsIntegrationsResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// DeleteAws deletes a single AWS integration matching the integration guid on
// the Lacework Server
func (svc *IntegrationsService) DeleteAws(guid string) (
	response AwsIntegrationsResponse,
	err error,
) {
	err = svc.delete(guid, &response)
	return
}

// ListAwsCfg lists the AWS_CFG external integrations available on the Lacework Server
func (svc *IntegrationsService) ListAwsCfg() (response AwsIntegrationsResponse, err error) {
	err = svc.listByType(AwsCfgIntegration, &response)
	return
}

// ListAwsCloudTrail lists the AWS_CT_SQS external integrations available on the Lacework Server
func (svc *IntegrationsService) ListAwsCloudTrail() (response AwsIntegrationsResponse, err error) {
	err = svc.listByType(AwsCloudTrailIntegration, &response)
	return
}

type AwsIntegrationsResponse struct {
	Data    []AwsIntegration `json:"data"`
	Ok      bool             `json:"ok"`
	Message string           `json:"message"`
}

type AwsIntegration struct {
	commonIntegrationData
	Data AwsIntegrationData `json:"DATA"`
}

type AwsIntegrationData struct {
	Credentials AwsCrossAccountCreds `json:"CROSS_ACCOUNT_CREDENTIALS,omitempty" mapstructure:"CROSS_ACCOUNT_CREDENTIALS"`

	// QueueUrl is a field that exists and is required for the AWS_CT_SQS integration,
	// though, it doesn't exist for AWS_CFG integrations, that's why we omit it if empty
	QueueUrl string `json:"QUEUE_URL,omitempty" mapstructure:"QUEUE_URL"`

	// This field must be a base64 encode with the following format:
	//
	// "data:application/json;name=i.json;base64,[ENCODING]"
	//
	// [ENCODING] is the the base64 encode, use EncodeAccountMappingFile() to encode a JSON mapping file
	AccountMappingFile string `json:"ACCOUNT_MAPPING_FILE,omitempty" mapstructure:"ACCOUNT_MAPPING_FILE"`

	// AwsAccountID is the AWS account that owns the IAM role credentials
	AwsAccountID string `json:"AWS_ACCOUNT_ID,omitempty" mapstructure:"AWS_ACCOUNT_ID"`

	// GovCloudCredentials represents the credential structure for AWS_US_GOV_CFG and AWS_US_GOV_CT_SQS integrations
	GovCloudCredentials AwsGovCloudCreds `json:"ACCESS_KEY_CREDENTIALS,omitempty" mapstructure:"ACCESS_KEY_CREDENTIALS"`
}

func (aws *AwsIntegrationData) EncodeAccountMappingFile(mapping []byte) {
	encodedMappings := base64.StdEncoding.EncodeToString(mapping)
	aws.AccountMappingFile = fmt.Sprintf("data:application/json;name=i.json;base64,%s", encodedMappings)
}

func (aws *AwsIntegrationData) DecodeAccountMappingFile() ([]byte, error) {
	if len(aws.AccountMappingFile) == 0 {
		return []byte{}, nil
	}

	var (
		b64      = strings.Split(aws.AccountMappingFile, ",")
		raw, err = base64.StdEncoding.DecodeString(b64[1])
	)
	if err != nil {
		return []byte{}, err
	}

	return raw, nil
}

type AwsCrossAccountCreds struct {
	RoleArn    string `json:"ROLE_ARN" mapstructure:"ROLE_ARN"`
	ExternalID string `json:"EXTERNAL_ID" mapstructure:"EXTERNAL_ID"`
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
