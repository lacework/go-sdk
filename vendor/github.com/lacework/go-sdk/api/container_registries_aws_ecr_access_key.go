//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

// GetAwsEcrAccessKey gets a single AwsEcrAccessKey integration with access key credentials matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetAwsEcrAccessKey(guid string) (
	response AwsEcrAccessKeyIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsEcrAccessKey updates a single AwsEcrAccessKey integration with access key credential on the Lacework Server
func (svc *ContainerRegistriesService) UpdateAwsEcrAccessKey(data ContainerRegistry) (
	response AwsEcrAccessKeyIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsEcrAccessKeyIntegrationResponse struct {
	Data AwsEcrIntegration `json:"data"`
}

type AwsEcrIntegration struct {
	v2CommonIntegrationData
	Data AwsEcrAccessKeyData `json:"data"`
}

type AwsEcrAccessKeyData struct {
	AccessKeyCredentials AwsEcrAccessKeyCredentials `json:"accessKeyCredentials,omitempty"`
	RegistryDomain       string                     `json:"registryDomain"`
	LimitByTag           []string                   `json:"limitByTag,omitempty"`
	LimitByLabel         []map[string]string        `json:"limitByLabel,omitempty"`
	LimitByRep           []string                   `json:"limitByRep,omitempty"`
	LimitNumImg          int                        `json:"limitNumImg"`
	NonOSPackageEval     bool                       `json:"nonOsPackageEval"`
	AwsAuthType          string                     `json:"awsAuthType"`
	RegistryType         string                     `json:"registryType"`
}

func verifyAwsEcrContainerRegistry(data interface{}) interface{} {
	if ecr, ok := data.(AwsEcrAccessKeyData); ok {
		ecr.RegistryType = AwsEcrContainerRegistry.String()
		ecr.AwsAuthType = AwsEcrAccessKey.String()
		return ecr
	}

	if ecr, ok := data.(AwsEcrIamRoleData); ok {
		ecr.RegistryType = AwsEcrContainerRegistry.String()
		ecr.AwsAuthType = AwsEcrIAM.String()
		return ecr
	}

	return data
}

type AwsEcrAccessKeyCredentials struct {
	AccessKeyID     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}
