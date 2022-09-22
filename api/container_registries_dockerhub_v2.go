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

// GetDockerhubV2 gets a single DockerhubV2 integration matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetDockerhubV2(guid string) (
	response DockerhubV2IntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateDockerhubV2 updates a single DockerhubV2 integration on the Lacework Server
func (svc *ContainerRegistriesService) UpdateDockerhubV2(data ContainerRegistry) (
	response DockerhubV2IntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type DockerhubV2IntegrationResponse struct {
	Data DockerhubV2Integration `json:"data"`
}

type DockerhubV2Integration struct {
	v2CommonIntegrationData
	Data DockerhubV2Data `json:"data"`
}

func (reg DockerhubV2Integration) ContainerRegistryType() containerRegistryType {
	t, _ := FindContainerRegistryType(reg.Data.RegistryType)
	return t
}

type DockerhubV2Data struct {
	Credentials           DockerhubV2Credentials `json:"credentials"`
	RegistryDomain        string                 `json:"registryDomain"`
	RegistryType          string                 `json:"registryType"`
	RegistryNotifications *bool                  `json:"registryNotifications,omitempty"`
	LimitByTag            []string               `json:"limitByTag,omitempty"`
	LimitByLabel          []map[string]string    `json:"limitByLabel,omitempty"`
	NonOSPackageEval      bool                   `json:"nonOsPackageEval"`
}

func verifyDockerhubV2ContainerRegistry(data interface{}) interface{} {
	if ecr, ok := data.(DockerhubV2Data); ok {
		ecr.RegistryType = DockerhubV2ContainerRegistry.String()
		return ecr
	}
	return data
}

type DockerhubV2Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}
