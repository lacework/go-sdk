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

// This file has the abstraction of the Lacework APIs to manage integrations
// of type AWS_ECR using ACCESS_KEY_CREDENTIALS as Authentication Method

func NewAwsEcrWithAccessKeyIntegration(name string, data AwsEcrDataWithAccessKeyCreds) AwsEcrWithAccessKeyIntegration {
	data.RegistryType = EcrRegistry.String()
	data.AwsAuthType = AwsEcrAccessKey.String()
	return AwsEcrWithAccessKeyIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    ContainerRegistryIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateAwsEcrWithAccessKey creates an AWS_ECR integration using an AWS Access
// Key as authenticatin method to access the registry
func (svc *IntegrationsService) CreateAwsEcrWithAccessKey(integration AwsEcrWithAccessKeyIntegration) (
	response AwsEcrWithAccessKeyIntegrationResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAwsEcrWithAccessKey gets an AWS_ECR integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetAwsEcrWithAccessKey(guid string) (
	response AwsEcrWithAccessKeyIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsEcrWithAccessKey updates a single AWS_ECR integration
func (svc *IntegrationsService) UpdateAwsEcrWithAccessKey(integration AwsEcrWithAccessKeyIntegration) (
	response AwsEcrWithAccessKeyIntegrationResponse,
	err error,
) {
	err = svc.update(integration.IntgGuid, integration, &response)
	return
}

type AwsEcrCommonData struct {
	AwsAuthType      string `json:"AWS_AUTH_TYPE" mapstructure:"AWS_AUTH_TYPE"`
	RegistryType     string `json:"REGISTRY_TYPE" mapstructure:"REGISTRY_TYPE"`
	RegistryDomain   string `json:"REGISTRY_DOMAIN" mapstructure:"REGISTRY_DOMAIN"`
	LimitByTag       string `json:"LIMIT_BY_TAG" mapstructure:"LIMIT_BY_TAG"`
	LimitByLabel     string `json:"LIMIT_BY_LABEL" mapstructure:"LIMIT_BY_LABEL"`
	LimitByRep       string `json:"LIMIT_BY_REP,omitempty" mapstructure:"LIMIT_BY_REP"`
	LimitNumImg      int    `json:"LIMIT_NUM_IMG,omitempty" mapstructure:"LIMIT_NUM_IMG"`
	NonOSPackageEval bool   `json:"NON_OS_PACKAGE_EVAL" mapstructure:"NON_OS_PACKAGE_EVAL"`
}

type AwsEcrWithAccessKeyIntegrationResponse struct {
	Data    []AwsEcrWithAccessKeyIntegration `json:"data"`
	Ok      bool                             `json:"ok"`
	Message string                           `json:"message"`
}

type AwsEcrWithAccessKeyIntegration struct {
	commonIntegrationData
	Data AwsEcrDataWithAccessKeyCreds `json:"DATA"`
}

type AwsEcrDataWithAccessKeyCreds struct {
	Credentials AwsEcrAccessKeyCreds `json:"ACCESS_KEY_CREDENTIALS" mapstructure:"ACCESS_KEY_CREDENTIALS"`
	AwsEcrCommonData
}

type AwsEcrAccessKeyCreds struct {
	AccessKeyID     string `json:"ACCESS_KEY_ID" mapstructure:"ACCESS_KEY_ID"`
	SecretAccessKey string `json:"SECRET_ACCESS_KEY" mapstructure:"SECRET_ACCESS_KEY"`
}
