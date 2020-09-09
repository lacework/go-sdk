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

func NewAwsEcrRegistryIntegration(name string, data AwsEcrData) AwsEcrIntegration {
	data.RegistryType = EcrRegistry.String()
	return AwsEcrIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    ContainerRegistryIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateAwsEcrRegistry creates an AWS_ECR integration on the Lacework Server
func (svc *IntegrationsService) CreateAwsEcrRegistry(integration AwsEcrIntegration) (
	response AwsEcrResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAwsEcrRegistry gets an AWS_ECR integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetAwsEcrRegistry(guid string) (
	response AwsEcrResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsEcrRegistry updates a single AWS_ECR integration
func (svc *IntegrationsService) UpdateAwsEcrRegistry(integration AwsEcrIntegration) (
	response AwsEcrResponse,
	err error,
) {
	err = svc.update(integration.IntgGuid, integration, &response)
	return
}

type AwsEcrResponse struct {
	Data    []AwsEcrIntegration `json:"data"`
	Ok      bool                `json:"ok"`
	Message string              `json:"message"`
}

// For AWS_ECR registry
type AwsEcrIntegration struct {
	commonIntegrationData
	Data AwsEcrData `json:"DATA"`
}

type AwsEcrData struct {
	Credentials    AwsEcrCreds `json:"ACCESS_KEY_CREDENTIALS" mapstructure:"ACCESS_KEY_CREDENTIALS"`
	RegistryType   string      `json:"REGISTRY_TYPE" mapstructure:"REGISTRY_TYPE"`
	RegistryDomain string      `json:"REGISTRY_DOMAIN" mapstructure:"REGISTRY_DOMAIN"`
	LimitByTag     string      `json:"LIMIT_BY_TAG" mapstructure:"LIMIT_BY_TAG"`
	LimitByLabel   string      `json:"LIMIT_BY_LABEL" mapstructure:"LIMIT_BY_LABEL"`
	LimitByRep     string      `json:"LIMIT_BY_REP,omitempty" mapstructure:"LIMIT_BY_REP"`
	LimitNumImg    int         `json:"LIMIT_NUM_IMG,omitempty" mapstructure:"LIMIT_NUM_IMG"`
}

type AwsEcrCreds struct {
	AccessKeyID     string `json:"ACCESS_KEY_ID" mapstructure:"ACCESS_KEY_ID"`
	SecretAccessKey string `json:"SECRET_ACCESS_KEY" mapstructure:"SECRET_ACCESS_KEY"`
}
