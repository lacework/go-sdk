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

// GetAzureCfg gets a single AzureCfg integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAzureCfg(guid string) (
	response AzureCfgIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAzureCfg updates a single AzureCfg integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAzureCfg(data CloudAccount) (
	response AzureCfgIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AzureCfgIntegrationResponse struct {
	Data AzureCfg `json:"data"`
}

type AzureCfg struct {
	v2CommonIntegrationData
	Data AzureCfgData `json:"data"`
}

type AzureCfgData struct {
	Credentials AzureCfgCredentials `json:"credentials"`
	TenantID    string              `json:"tenantId"`
}

type AzureCfgCredentials struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
