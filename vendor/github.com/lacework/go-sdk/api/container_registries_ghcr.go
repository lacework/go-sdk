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

// GetGhcr gets a single Ghcr integration matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetGhcr(guid string) (
	response GhcrIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateGhcr updates a single Ghcr integration on the Lacework Server
func (svc *ContainerRegistriesService) UpdateGhcr(data ContainerRegistry) (
	response GhcrIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type GhcrIntegrationResponse struct {
	Data GhcrIntegration `json:"data"`
}

type GhcrIntegration struct {
	v2CommonIntegrationData
	Data GhcrData `json:"data"`
}

func (reg GhcrIntegration) ContainerRegistryType() containerRegistryType {
	t, _ := FindContainerRegistryType(reg.Data.RegistryType)
	return t
}

type GhcrData struct {
	Credentials           GhcrCredentials     `json:"credentials"`
	RegistryNotifications bool                `json:"registryNotifications"`
	RegistryDomain        string              `json:"registryDomain"` // always "ghcr.io"
	RegistryType          string              `json:"registryType"`   // always "GHCR"
	LimitByTag            []string            `json:"limitByTag,omitempty"`
	LimitByLabel          []map[string]string `json:"limitByLabel,omitempty"`
	LimitByRep            []string            `json:"limitByRep,omitempty"`
	LimitNumImg           int                 `json:"limitNumImg"`
	NonOSPackageEval      bool                `json:"nonOsPackageEval"`
}

func verifyGhcrContainerRegistry(data interface{}) interface{} {
	if ghcr, ok := data.(GhcrData); ok {
		ghcr.RegistryType = GhcrContainerRegistry.String()
		ghcr.RegistryDomain = "ghcr.io"
		return ghcr
	}
	return data
}

// GcpCredentials is already defined in api/integrations_gcp.go:163
// so we need to add a "V2" at the end to make it clear that this is
// the Google Credentials struct for API v2
type GhcrCredentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Ssl      bool   `json:"ssl"`
}
