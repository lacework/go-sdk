//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

// gcpResourceLevel determines Project or Organization level integration
type gcpResourceLevel int

const (
	// Project level integration with GCP
	GcpProjectIntegration gcpResourceLevel = iota

	// Organization level integration with GCP
	GcpOrganizationIntegration
)

var gcpResourceLevels = map[gcpResourceLevel]string{
	GcpProjectIntegration:      "PROJECT",
	GcpOrganizationIntegration: "ORGANIZATION",
}

func (g gcpResourceLevel) String() string {
	return gcpResourceLevels[g]
}

// NewGcpIntegration returns an instance of gcpIntegration with the provided
// integration type, name and data. The type can only be GcpCfgIntegration or
// GcpAuditLogIntegration
//
//
// Basic usage: Initialize a new gcpIntegration struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   gcp, err := api.NewGcpIntegration("abc",
//     api.GcpCfgIntegration,
//     api.GcpIntegrationData{
//       ID: "1234",
//       IdType: "id_type",
//       Credentials: api.GcpCredentials{
//         ClientId: "id",
//         ClientEmail: "email",
//         PrivateKeyId: "key_id",
//         PrivateKey: "key",
//       },
//     },
//   )
//   if err != nil {
//     return err
//   }
//
//   client.Integrations.CreateGcp(gcp)
//
func NewGcpIntegration(name string, iType integrationType, data GcpIntegrationData) gcpIntegration {
	return gcpIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    iType.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// NewGcpCfgIntegration returns an instance of gcpIntegration of type GCP_CFG
func NewGcpCfgIntegration(name string, data GcpIntegrationData) gcpIntegration {
	return NewGcpIntegration(name, GcpCfgIntegration, data)
}

// NewGcpAuditLogIntegration returns an instance of gcpIntegration of type GCP_AT_SES
func NewGcpAuditLogIntegration(name string, data GcpIntegrationData) gcpIntegration {
	return NewGcpIntegration(name, GcpAuditLogIntegration, data)
}

// CreateGcp creates a single Gcp integration on the Lacework Server
func (svc *IntegrationsService) CreateGcp(data gcpIntegration) (
	response gcpIntegrationsResponse,
	err error,
) {
	err = svc.create(data, &response)
	return
}

// GetGcp gets a single Gcp integration matching the integration guid
// on the Lacework Server
func (svc *IntegrationsService) GetGcp(guid string) (
	response gcpIntegrationsResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateGcp updates a single Gcp integration on the Lacework Server
func (svc *IntegrationsService) UpdateGcp(data gcpIntegration) (
	response gcpIntegrationsResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// DeleteGcp deletes a single Gcp integration matching the integration guid
// on the Lacework Server
func (svc *IntegrationsService) DeleteGcp(guid string) (
	response gcpIntegrationsResponse,
	err error,
) {
	err = svc.delete(guid, &response)
	return
}

// ListGcpCfg lists the GCP_CFG external integrations available on the Lacework Server
func (svc *IntegrationsService) ListGcpCfg() (response gcpIntegrationsResponse, err error) {
	err = svc.listByType(GcpCfgIntegration, &response)
	return
}

// ListGcpAuditLog lists the GCP_AT_SES external integrations available on the Lacework Server
func (svc *IntegrationsService) ListGcpAuditLog() (response gcpIntegrationsResponse, err error) {
	err = svc.listByType(GcpAuditLogIntegration, &response)
	return
}

type gcpIntegrationsResponse struct {
	Data    []gcpIntegration `json:"data"`
	Ok      bool             `json:"ok"`
	Message string           `json:"message"`
}

type gcpIntegration struct {
	commonIntegrationData
	Data GcpIntegrationData `json:"DATA"`
}

type GcpIntegrationData struct {
	ID          string         `json:"ID"`
	IdType      string         `json:"ID_TYPE"`
	Credentials GcpCredentials `json:"CREDENTIALS"`

	// SubscriptionName is a field that exists and is required for the GCP_AT_SES
	// integration, though, it doesn't exist for GCP_CFG integrations, that's why
	// we omit it if empty
	SubscriptionName string `json:"SUBSCRIPTION_NAME,omitempty"`
}

type GcpCredentials struct {
	ClientId     string `json:"CLIENT_ID"`
	ClientEmail  string `json:"CLIENT_EMAIL"`
	PrivateKeyId string `json:"PRIVATE_KEY_ID"`
	PrivateKey   string `json:"PRIVATE_KEY"`
}
