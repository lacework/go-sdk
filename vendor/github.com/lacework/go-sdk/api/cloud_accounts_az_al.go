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

// GetAzureAlSeq gets a single AzureAlSeq integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAzureAlSeq(guid string) (
	response AzureAlSeqIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAzureAlSeq updates a single AzureAlSeq integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAzureAlSeq(data CloudAccount) (
	response AzureAlSeqIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AzureAlSeqIntegrationResponse struct {
	Data AzureAlSeq `json:"data"`
}

type AzureAlSeq struct {
	v2CommonIntegrationData
	Data AzureAlSeqData `json:"data"`
}

type AzureAlSeqData struct {
	Credentials AzureAlSeqCredentials `json:"credentials"`
	TenantID    string                `json:"tenantId"`
	QueueUrl    string                `json:"queueUrl"`
}

type AzureAlSeqCredentials struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
