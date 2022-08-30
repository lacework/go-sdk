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

// GetAwsEcr gets a single AwsEcr integration matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetAwsEcr(guid string) (
	response AwsEcrIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsEcr updates a single AwsEcr integration on the Lacework Server
func (svc *ContainerRegistriesService) UpdateAwsEcr(data ContainerRegistry) (
	response AwsEcrIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsEcrIntegrationResponse struct {
	Data AwsEcrIntegration `json:"data"`
}

type AwsEcrIntegration struct {
	v2CommonIntegrationData
	Data AwsEcrData `json:"data"`
}

func (reg AwsEcrIntegration) ContainerRegistryType() containerRegistryType {
	t, _ := FindContainerRegistryType(reg.Data.RegistryType)
	return t
}

type AwsEcrData struct {
	Credentials      EcrCredentials      `json:"credentials"`
	RegistryDomain   string              `json:"registryDomain"`
	RegistryType     string              `json:"registryType"`
	LimitByTag       []string            `json:"limitByTag"`
	LimitByLabel     []map[string]string `json:"limitByLabel"`
	LimitByRep       []string            `json:"limitByRep"`
	LimitNumImg      int                 `json:"limitNumImg"`
	NonOSPackageEval bool                `json:"nonOsPackageEval"`
}

func verifyAwsEcrContainerRegistry(data interface{}) interface{} {
	if ecr, ok := data.(AwsEcrData); ok {
		ecr.RegistryType = AwsEcrContainerRegistry.String()
		return ecr
	}
	return data
}

type EcrCredentials struct {
	RoleArn         string `json:"roleArn,omitempty"`
	ExternalID      string `json:"externalId,omitempty"`
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}
