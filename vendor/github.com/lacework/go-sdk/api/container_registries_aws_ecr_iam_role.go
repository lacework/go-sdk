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

// GetAwsEcrIamRole gets a single AwsEcr with Iam Role credentials integration matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetAwsEcrIamRole(guid string) (
	response AwsEcrIamRoleIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsEcrIamRole updates a single AwsEcr with Iam Role credentials integration on the Lacework Server
func (svc *ContainerRegistriesService) UpdateAwsEcrIamRole(data ContainerRegistry) (
	response AwsEcrIamRoleIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsEcrIamRoleIntegrationResponse struct {
	Data AwsEcrIamRoleIntegration `json:"data"`
}

type AwsEcrIamRoleIntegration struct {
	v2CommonIntegrationData
	Data AwsEcrIamRoleData `json:"data"`
}

func (reg AwsEcrIamRoleIntegration) ContainerRegistryType() containerRegistryType {
	t, _ := FindContainerRegistryType(reg.Data.RegistryType)
	return t
}

type AwsEcrIamRoleData struct {
	CrossAccountCredentials AwsEcrCrossAccountCredentials `json:"crossAccountCredentials,omitempty"`
	RegistryDomain          string                        `json:"registryDomain"`
	RegistryType            string                        `json:"registryType"`
	LimitByTag              []string                      `json:"limitByTag,omitempty"`
	LimitByLabel            []map[string]string           `json:"limitByLabel,omitempty"`
	LimitByRep              []string                      `json:"limitByRep,omitempty"`
	LimitNumImg             int                           `json:"limitNumImg"`
	NonOSPackageEval        bool                          `json:"nonOsPackageEval"`
	AwsAuthType             string                        `json:"awsAuthType"`
}

type AwsEcrCrossAccountCredentials struct {
	RoleArn    string `json:"roleArn,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
}
