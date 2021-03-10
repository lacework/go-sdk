//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

// This file has the abstraction of the Lacework APIs to manage integrations
// of type AWS_ECR using CROSS_ACCOUNT_CREDENTIALS as Authentication Method

type ecrAuthType int

const (
	AwsEcrIAM ecrAuthType = iota
	AwsEcrAccessKey
)

// AwsEcrAuthTypes is the list of available ECR auth types
var AwsEcrAuthTypes = map[ecrAuthType]string{
	AwsEcrIAM:       "AWS_IAM",
	AwsEcrAccessKey: "AWS_ACCESS_KEY",
}

// String returns the string representation of an ECR auth type
func (i ecrAuthType) String() string {
	return AwsEcrAuthTypes[i]
}

func NewAwsEcrWithCrossAccountIntegration(name string, data AwsEcrDataWithCrossAccountCreds) AwsEcrWithCrossAccountIntegration {
	data.RegistryType = EcrRegistry.String()
	data.AwsAuthType = AwsEcrIAM.String()
	return AwsEcrWithCrossAccountIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    ContainerRegistryIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateAwsEcrWithCrossAccount creates an AWS_ECR integration using an IAM Role as
// authenticatin method to access the registry
func (svc *IntegrationsService) CreateAwsEcrWithCrossAccount(integration AwsEcrWithCrossAccountIntegration) (
	response AwsEcrWithCrossAccountIntegrationResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAwsEcrWithCrossAccount gets an AWS_ECR integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetAwsEcrWithCrossAccount(guid string) (
	response AwsEcrWithCrossAccountIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsEcrWithCrossAccount updates a single AWS_ECR integration
func (svc *IntegrationsService) UpdateAwsEcrWithCrossAccount(integration AwsEcrWithCrossAccountIntegration) (
	response AwsEcrWithCrossAccountIntegrationResponse,
	err error,
) {
	err = svc.update(integration.IntgGuid, integration, &response)
	return
}

type AwsEcrWithCrossAccountIntegrationResponse struct {
	Data    []AwsEcrWithCrossAccountIntegration `json:"data"`
	Ok      bool                                `json:"ok"`
	Message string                              `json:"message"`
}

type AwsEcrWithCrossAccountIntegration struct {
	commonIntegrationData
	Data AwsEcrDataWithCrossAccountCreds `json:"DATA"`
}

type AwsEcrDataWithCrossAccountCreds struct {
	Credentials AwsCrossAccountCreds `json:"CROSS_ACCOUNT_CREDENTIALS" mapstructure:"CROSS_ACCOUNT_CREDENTIALS"`
	AwsEcrCommonData
}
