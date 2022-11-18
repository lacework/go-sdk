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

// GetInlineScanner gets a single InlineScanner integration matching the
// provided integration guid
func (svc *ContainerRegistriesService) GetInlineScanner(guid string) (
	response InlineScannerIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateInlineScanner updates a single InlineScanner integration on the Lacework Server
func (svc *ContainerRegistriesService) UpdateInlineScanner(data ContainerRegistry) (
	response InlineScannerIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type InlineScannerIntegrationResponse struct {
	Data InlineScannerIntegration `json:"data"`
}

type InlineScannerIntegration struct {
	v2CommonIntegrationData
	Data        InlineScannerData `json:"data"`
	ServerToken V2ServerToken     `json:"serverToken"`
	Props       interface{}       `json:"props,omitempty"`
}

func (reg InlineScannerIntegration) ContainerRegistryType() containerRegistryType {
	t, _ := FindContainerRegistryType(reg.Data.RegistryType)
	return t
}

type InlineScannerData struct {
	RegistryType  string              `json:"registryType"` // always "INLINE_SCANNER"
	IdentifierTag []map[string]string `json:"identifierTag"`
	LimitNumScan  string              `json:"limitNumScan,omitempty"`
}

func verifyInlineScannerContainerRegistry(data interface{}) interface{} {
	if inlineScanner, ok := data.(InlineScannerData); ok {
		inlineScanner.RegistryType = InlineScannerContainerRegistry.String()
		return inlineScanner
	}
	return data
}
