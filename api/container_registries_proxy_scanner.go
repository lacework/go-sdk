//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

// GetProxyScanner gets a single ProxyScanner integration matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetProxyScanner(guid string) (
	response ProxyScannerIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateProxyScanner updates a single ProxyScanner integration on the Lacework Server
func (svc *ContainerRegistriesService) UpdateProxyScanner(data ContainerRegistry) (
	response ProxyScannerIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type ProxyScannerIntegrationResponse struct {
	Data ProxyScannerIntegration `json:"data"`
}

type ProxyScannerIntegration struct {
	v2CommonIntegrationData
	Data        ProxyScannerData        `json:"data"`
	ServerToken ProxyScannerServerToken `json:"serverToken"`
}

type ProxyScannerServerToken struct {
	Token string `json:"serverToken"`
	URI   string `json:"uri"`
}

func (reg ProxyScannerIntegration) ContainerRegistryType() containerRegistryType {
	t, _ := FindContainerRegistryType(reg.Data.RegistryType)
	return t
}

type ProxyScannerData struct {
	RegistryType string              `json:"registryType"` // always "PROXY_SCANNER"
	LimitNumImg  int                 `json:"limitNumImg"`
	LimitByRep   []map[string]string `json:"limitByRep"`
	LimitByTag   []map[string]string `json:"limitByTag"`
	LimitByLabel []map[string]string `json:"limitByLabel"`
}

func verifyProxyScannerContainerRegistry(data interface{}) interface{} {
	if proxyScanner, ok := data.(ProxyScannerData); ok {
		proxyScanner.RegistryType = ProxyScannerContainerRegistry.String()
		return proxyScanner
	}
	return data
}
