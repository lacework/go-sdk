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

// GetDockerhub gets a single Dockerhub integration matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetDockerhub(guid string) (
	response DockerhubIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateDockerhub updates a single Dockerhub integration on the Lacework Server
func (svc *ContainerRegistriesService) UpdateDockerhub(data ContainerRegistry) (
	response DockerhubIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type DockerhubIntegrationResponse struct {
	Data DockerhubIntegration `json:"data"`
}

type DockerhubIntegration struct {
	v2CommonIntegrationData
	Data DockerhubData `json:"data"`
}

func (reg DockerhubIntegration) ContainerRegistryType() containerRegistryType {
	t, _ := FindContainerRegistryType(reg.Data.RegistryType)
	return t
}

type DockerhubData struct {
	Credentials           DockerhubCredentials `json:"credentials"`
	RegistryDomain        string               `json:"registryDomain"`
	RegistryType          string               `json:"registryType"`
	RegistryNotifications bool                 `json:"registryNotifications"`
	LimitByTag            []string             `json:"limitByTag,omitempty"`
	LimitByLabel          []map[string]string  `json:"limitByLabel,omitempty"`
	LimitByRep            []string             `json:"limitByRep,omitempty"`
	LimitNumImg           int                  `json:"limitNumImg"`
	NonOSPackageEval      bool                 `json:"nonOsPackageEval"`
}

func verifyDockerhubContainerRegistry(data interface{}) interface{} {
	if ecr, ok := data.(DockerhubData); ok {
		ecr.RegistryType = DockerhubContainerRegistry.String()
		return ecr
	}
	return data
}

type DockerhubCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
