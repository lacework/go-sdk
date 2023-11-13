//
// Author:: Ao Zhang (<ao.zhang@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.n
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

const (
	AzureSubscriptionIntegration string = "SUBSCRIPTION"
	AzureTenantIntegration       string = "TENANT"
)

// GetAzureSidekick gets a single AzureSidekick integration matching the provided integration guid
func (svc *CloudAccountsService) GetAzureSidekick(guid string) (
	response AzureSidekickIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// CreateAzureSidekick creates an AzureSidekick Cloud Account integration
func (svc *CloudAccountsService) CreateAzureSidekick(data CloudAccount) (
	response AzureSidekickIntegrationResponse,
	err error,
) {
	err = svc.create(data, &response)
	return
}

// UpdateAzureSidekick updates a single AzureSidekick integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAzureSidekick(data CloudAccount) (
	response AzureSidekickIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AzureSidekickIntegrationResponse struct {
	Data V2AzureSidekickIntegration `json:"data"`
}

type AzureSidekickToken struct {
	ServerToken string `json:"serverToken"`
	Uri         string `json:"uri"`
}

type V2AzureSidekickIntegration struct {
	v2CommonIntegrationData
	AzureSidekickToken `json:"serverToken"`
	Data               AzureSidekickData `json:"data"`
}

type AzureSidekickData struct {
	Credentials               AzureSidekickCredentials `json:"credentials"`
	IntegrationLevel          string                   `json:"integrationLevel"` // SUBSCRIPTION or TENANT
	ScanningSubscriptionId    string                   `json:"scanningSubscriptionId"`
	TenantId                  string                   `json:"tenantId"`
	BlobContainerName         string                   `json:"blobContainerName"`
	ScanningResourceGroupName string                   `json:"scanningResourceGroupName"`
	StorageAccountUrl         string                   `json:"storageAccountUrl"`
	SubscriptionsList         string                   `json:"subscriptionsList,omitempty"`
	QueryText                 string                   `json:"queryText,omitempty"`
	ScanFrequency             int                      `json:"scanFrequency"` // in hours
	ScanContainers            bool                     `json:"scanContainers"`
	ScanHostVulnerabilities   bool                     `json:"scanHostVulnerabilities"`
	ScanMultiVolume           bool                     `json:"scanMultiVolume"`
	ScanStoppedInstances      bool                     `json:"scanStoppedInstances"`
}

type AzureSidekickCredentials struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret,omitempty"`
}
