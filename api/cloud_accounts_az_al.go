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

// GetAzAl gets a single AzAl integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAzAl(guid string) (
	response AzAlIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAzAl updates a single AzAl integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAzAl(data CloudAccount) (
	response AzAlIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AzAlIntegrationResponse struct {
	Data AzAl `json:"data"`
}

type AzAl struct {
	v2CommonIntegrationData
	Data AzAlData `json:"data"`
}

type AzAlData struct {
	Credentials AzAlCredentials `json:"crossAccountCredentials"`
	TenantID    string          `json:"tenantId"`
	QueueUrl    string          `json:"queueUrl"`
}

type AzAlCredentials struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
