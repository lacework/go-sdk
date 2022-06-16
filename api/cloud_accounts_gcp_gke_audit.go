//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

// GetGcpGkeAudit gets a single GcpGkeAudit integration matching the provided integration guid
func (svc *CloudAccountsService) GetGcpGkeAudit(guid string) (
	response GcpGkeAuditIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateGcpGkeAudit updates a single GcpGkeAudit integration on the Lacework Server
func (svc *CloudAccountsService) UpdateGcpGkeAudit(data CloudAccount) (
	response GcpGkeAuditIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type GcpGkeAuditIntegrationResponse struct {
	Data GcpGkeAuditIntegration `json:"data"`
}

type GcpGkeAuditIntegration struct {
	v2CommonIntegrationData
	Data GcpGkeAuditData `json:"data"`
}

type GcpGkeAuditData struct {
	Credentials     GcpGkeAuditCredentials `json:"credentials"`
	IntegrationType string                 `json:"integrationType"`
	// OrganizationId is optional for a project level integration, therefore we omit if empty
	OrganizationId   string `json:"organizationId,omitempty"`
	ProjectId        string `json:"projectId"`
	SubscriptionName string `json:"subscriptionName"`
}

type GcpGkeAuditCredentials struct {
	ClientId     string `json:"clientId"`
	ClientEmail  string `json:"clientEmail"`
	PrivateKeyId string `json:"PrivateKeyID"`
	PrivateKey   string `json:"PrivateKey"`
}
