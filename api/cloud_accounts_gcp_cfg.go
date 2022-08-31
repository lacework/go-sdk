//
// Author:: Darren Murray(<darren.murray@lacework.net>)
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

// GetGcpCfg gets a single GcpCfg integration matching the provided integration guid
func (svc *CloudAccountsService) GetGcpCfg(guid string) (
	response GcpCfgIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateGcpCfg updates a single GcpCfg integration on the Lacework Server
func (svc *CloudAccountsService) UpdateGcpCfg(data CloudAccount) (
	response GcpCfgIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type GcpCfgIntegrationResponse struct {
	Data V2GcpCfgIntegration `json:"data"`
}

type V2GcpCfgIntegration struct {
	v2CommonIntegrationData
	Data GcpCfgData `json:"data"`
}

type GcpCfgData struct {
	Credentials GcpCfgCredentials `json:"credentials"`
	IDType      string            `json:"idType"`
	// Either the org id or project id
	ID string `json:"id"`
}

type GcpCfgCredentials struct {
	ClientID     string `json:"clientId"`
	ClientEmail  string `json:"clientEmail"`
	PrivateKeyID string `json:"privateKeyId"`
	PrivateKey   string `json:"privateKey"`
}
