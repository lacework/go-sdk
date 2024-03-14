//
// Author:: Rubinder Singh (<rubinder.singh@lacework.net>)
// Copyright:: Copyright 2024, Lacework Inc.
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

// GetAzureAdAl gets a single AzureAdAl integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAzureAdAl(guid string) (
	response AzureAdAlIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAzureAdAl updates a single AzureAdAl integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAzureAdAl(data CloudAccount) (
	response AzureAdAlIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AzureAdAlIntegrationResponse struct {
	Data AzureAdAl `json:"data"`
}

type AzureAdAl struct {
	v2CommonIntegrationData
	Data AzureAdAlData `json:"data"`
}

type AzureAdAlData struct {
	Credentials       AzureAdAlCredentials `json:"credentials"`
	TenantID          string               `json:"tenantId"`
	EventHubNamespace string               `json:"eventHubNamespace"`
	EventHubName      string               `json:"eventHubName"`
}

type AzureAdAlCredentials struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
