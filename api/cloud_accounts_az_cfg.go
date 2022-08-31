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

// GetAzCfg gets a single AzCfg integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAzCfg(guid string) (
	response AzCfgIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAzCfg updates a single AzCfg integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAzCfg(data CloudAccount) (
	response AzCfgIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AzCfgIntegrationResponse struct {
	Data AzCfg `json:"data"`
}

type AzCfg struct {
	v2CommonIntegrationData
	Data AzCfgData `json:"data"`
}

type AzCfgData struct {
	Credentials AzCfgCredentials `json:"crossAccountCredentials"`
	TenantID    string           `json:"tenantId"`
}

type AzCfgCredentials struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
