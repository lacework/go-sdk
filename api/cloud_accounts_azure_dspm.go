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

// GetAzureDspm gets a single AzureDspm integration matching the provided integration guid
func (svc *CloudAccountsService) GetAzureDspm(guid string) (
	response AzureDspmResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// CreateAzureDspm creates an AzureDspm Cloud Account integration
func (svc *CloudAccountsService) CreateAzureDspm(data CloudAccount) (
	response AzureDspmResponse,
	err error,
) {
	err = svc.create(data, &response)
	return
}

// UpdateAzureDspm updates a single AzureDspm integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAzureDspm(data CloudAccount) (
	response AzureDspmResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AzureDspmResponse struct {
	Data AzureDspm `json:"data"`
}

type AzureDspm struct {
	v2CommonIntegrationData
	azureDspmToken `json:"serverToken"`
	Data           AzureDspmData `json:"data"`
}

type azureDspmToken struct {
	ServerToken string `json:"serverToken"`
	Uri         string `json:"uri"`
}

// AzureDspmData contains the data needed by Lacework platform services.
type AzureDspmData struct {
	TenantID          string               `json:"tenantId,omitempty"`
	StorageAccountUrl string               `json:"storageAccountUrl,omitempty"`
	BlobContainerName string               `json:"blobContainerName,omitempty"`
	Credentials       AzureDspmCredentials `json:"credentials"`
}

type AzureDspmCredentials struct {
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}
