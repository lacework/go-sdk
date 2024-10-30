//
// Author:: Kolbeinn Karlsson (<kolbeinn.karlsson@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

// GetOciCfg gets a single OciCfg integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetOciCfg(guid string) (
	response OciCfgIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateOciCfg updates a single OciCfg integration on the Lacework Server
func (svc *CloudAccountsService) UpdateOciCfg(data CloudAccount) (
	response OciCfgIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type OciCfgIntegrationResponse struct {
	Data OciCfg `json:"data"`
}

type OciCfg struct {
	v2CommonIntegrationData
	Data OciCfgData `json:"data"`
}

type OciCfgData struct {
	Credentials OciCfgCredentials `json:"credentials"`
	HomeRegion  string            `json:"homeRegion"`
	TenantID    string            `json:"tenantId"`
	TenantName  string            `json:"tenantName"`
	UserOCID    string            `json:"userOcid"`
}

type OciCfgCredentials struct {
	Fingerprint string `json:"fingerprint"`
	PrivateKey  string `json:"privateKey,omitempty"`
}
