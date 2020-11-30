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

	"go.uber.org/zap"
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
//       Credentials: api.AwsIntegrationCreds {
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
	if err != nil {
		return
	}

	// WORKAROUND (@afiune) The backend is currently not triggering an initial
	// report automatically after creation of Cloud Account (CFG) Integrations,
	// we are implementing this trigger here until we implement it in the backend
	// with RAIN-13422
	if len(response.Data) == 0 {
		return
	}
	if integration.Type != AwsCfgIntegration.String() {
		return
	}

	intgGuid := response.Data[0].IntgGuid
	svc.client.log.Info("triggering compliance report",
		zap.String("cloud_integration", integration.Type),
		zap.String("int_guid", intgGuid),
	)
	_, errComplianceReport := svc.client.Compliance.RunIntegrationReport(intgGuid)
	if errComplianceReport != nil {
		svc.client.log.Warn("unable to trigger compliance report",
			zap.String("cloud_integration", integration.Type),
			zap.String("int_guid", intgGuid),
			zap.String("error", errComplianceReport.Error()),
		)
	}

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
	Credentials AwsIntegrationCreds `json:"CROSS_ACCOUNT_CREDENTIALS" mapstructure:"CROSS_ACCOUNT_CREDENTIALS"`

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
}

func (aws *AwsIntegrationData) EncodeAccountMappingFile(mapping string) {
	encodedMappings := base64.StdEncoding.EncodeToString([]byte(mapping))
	aws.AccountMappingFile = fmt.Sprintf("data:application/json;name=i.json;base64,%s", encodedMappings)
}

func (aws *AwsIntegrationData) DecodeAccountMappingFile() (string, error) {
	if len(aws.AccountMappingFile) == 0 {
		return "", nil
	}

	var (
		b64      = strings.Split(aws.AccountMappingFile, ",")
		raw, err = base64.StdEncoding.DecodeString(b64[1])
	)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

type AwsIntegrationCreds struct {
	RoleArn    string `json:"ROLE_ARN" mapstructure:"ROLE_ARN"`
	ExternalID string `json:"EXTERNAL_ID" mapstructure:"EXTERNAL_ID"`
}
