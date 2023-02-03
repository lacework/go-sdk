//
// Author:: David McTavish(<david.mctavish@lacework.net>)
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

// GetGcpAlPubSub gets a single GcpAlPubSub integration matching the provided integration guid
func (svc *CloudAccountsService) GetGcpAlPubSub(guid string) (
	response GcpAlPubSubIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateGcpAtSes updates a single GcpAtSes integration on the Lacework Server
func (svc *CloudAccountsService) UpdateGcpAlPubSub(data CloudAccount) (
	response GcpAlPubSubIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type GcpAlPubSubIntegrationResponse struct {
	Data V2GcpAlPubSubIntegration `json:"data"`
}

type V2GcpAlPubSubIntegration struct {
	v2CommonIntegrationData
	Data GcpAlPubSubSesData `json:"data"`
}

type GcpAlPubSubSesData struct {
	Credentials     GcpAlPubSubCredentials `json:"credentials"`
	IntegrationType string                 `json:"integrationType"`
	// OrganizationId is optional for a project level integration, therefore we omit if empty
	OrganizationId   string `json:"organizationId,omitempty"`
	ProjectId        string `json:"projectId"`
	SubscriptionName string `json:"subscriptionName"`
	TopicId          string `json:"topicId"`
}

type GcpAlPubSubCredentials struct {
	ClientId     string `json:"clientId"`
	ClientEmail  string `json:"clientEmail"`
	PrivateKeyId string `json:"privateKeyId"`
	PrivateKey   string `json:"privateKey,omitempty"`
}
