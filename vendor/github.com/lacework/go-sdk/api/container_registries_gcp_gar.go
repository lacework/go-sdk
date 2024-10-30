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

// GetGcpGar gets a single GcpGar integration matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetGcpGar(guid string) (
	response GcpGarIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateGcpGar updates a single GcpGar integration on the Lacework Server
func (svc *ContainerRegistriesService) UpdateGcpGar(data ContainerRegistry) (
	response GcpGarIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type GcpGarIntegrationResponse struct {
	Data GcpGarIntegration `json:"data"`
}

type GcpGarIntegration struct {
	v2CommonIntegrationData
	Data GcpGarData `json:"data"`
}

func (reg GcpGarIntegration) ContainerRegistryType() containerRegistryType {
	t, _ := FindContainerRegistryType(reg.Data.RegistryType)
	return t
}

type GcpGarData struct {
	Credentials      GcpCredentialsV2    `json:"credentials"`
	RegistryDomain   string              `json:"registryDomain"`
	RegistryType     string              `json:"registryType"` // always "GCP_GAR"
	LimitByTag       []string            `json:"limitByTag,omitempty"`
	LimitByLabel     []map[string]string `json:"limitByLabel,omitempty"`
	LimitByRep       []string            `json:"limitByRep,omitempty"`
	LimitNumImg      int                 `json:"limitNumImg"`
	NonOSPackageEval bool                `json:"nonOsPackageEval"`
}

func verifyGcpGarContainerRegistry(data interface{}) interface{} {
	if gar, ok := data.(GcpGarData); ok {
		gar.RegistryType = GcpGarContainerRegistry.String()
		return gar
	}
	return data
}

// GcpCredentials is already defined in api/integrations_gcp.go:163
// so we need to add a "V2" at the end to make it clear that this is
// the Google Credentials struct for API v2
type GcpCredentialsV2 struct {
	ClientEmail  string `json:"clientEmail"`
	ClientID     string `json:"clientId"`
	PrivateKeyID string `json:"privateKeyId"`
	PrivateKey   string `json:"privateKey,omitempty"`
}
