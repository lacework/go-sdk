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

// GetGcpGcr gets a single GcpGcr integration matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetGcpGcr(guid string) (
	response GcpGcrIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateGcpGcr updates a single GcpGcr integration on the Lacework Server
func (svc *ContainerRegistriesService) UpdateGcpGcr(data ContainerRegistry) (
	response GcpGcrIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type GcpGcrIntegrationResponse struct {
	Data GcpGcrIntegration `json:"data"`
}

type GcpGcrIntegration struct {
	v2CommonIntegrationData
	Data GcpGcrData `json:"data"`
}

func (reg GcpGcrIntegration) ContainerRegistryType() containerRegistryType {
	t, _ := FindContainerRegistryType(reg.Data.RegistryType)
	return t
}

type GcpGcrData struct {
	Credentials      GcpCredentialsV2    `json:"credentials"`
	RegistryDomain   string              `json:"registryDomain"`
	RegistryType     string              `json:"registryType"`
	LimitByTag       []string            `json:"limitByTag,omitempty"`
	LimitByLabel     []map[string]string `json:"limitByLabel,omitempty"`
	LimitByRep       []string            `json:"limitByRep,omitempty"`
	LimitNumImg      int                 `json:"limitNumImg"`
	NonOSPackageEval bool                `json:"nonOsPackageEval"`
}

func verifyGcpGcrContainerRegistry(data interface{}) interface{} {
	if gar, ok := data.(GcpGcrData); ok {
		gar.RegistryType = GcpGcrContainerRegistry.String()
		return gar
	}
	return data
}
