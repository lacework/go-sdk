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

// NewGcpConfigIntegration returns an instance of gcpConfigIntegration
//
// Basic usage: Initialize a new gcpConfigIntegration struct, then
//              use the new instance to do CRUD operations
//
//   gcp, err := api.NewGcpConfigIntegration("abc",
//     api.GcpIntegrationData{
//       ID: "1234",
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
//   integrationResponse, err := api.CreateGcpConfigIntegration(gcp)
//   if err != nil {
//     return err
//   }
//
func NewGcpConfigIntegration(name string, data GcpIntegrationData) gcpConfigIntegration {
	return gcpConfigIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    GcpCfgIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateGcpConfigIntegration creates a single GCP_CFG integration on the Lacework Server
func (c *Client) CreateGcpConfigIntegration(data gcpConfigIntegration) (response gcpIntegrationsResponse, err error) {
	err = c.createIntegration(data, &response)
	return
}

// GetGcpConfigIntegration gets a single integration matching the integration guid available on the server
func (c *Client) GetGcpConfigIntegration(intgGuid string) (response gcpIntegrationsResponse, err error) {
	err = c.getIntegration(intgGuid, &response)
	return
}

// UpdateGcpConfigIntegration updates a single integration on the server
func (c *Client) UpdateGcpConfigIntegration(data gcpConfigIntegration) (response gcpIntegrationsResponse, err error) {
	err = c.updateIntegration(data.IntgGuid, data, &response)
	return
}

// DeleteGcpConfigIntegration gets a single integration matching the integration guid available on the server
func (c *Client) DeleteGcpConfigIntegration(intgGuid string) (response gcpIntegrationsResponse, err error) {
	err = c.deleteIntegration(intgGuid, &response)
	return
}

func (c *Client) GetGcpIntegrations() (response gcpIntegrationsResponse, err error) {
	return
}

type gcpIntegrationsResponse struct {
	Data    []gcpConfigIntegration `json:"data"`
	Ok      bool                   `json:"ok"`
	Message string                 `json:"message"`
}

type gcpConfigIntegration struct {
	commonIntegrationData
	Data GcpIntegrationData `json:"DATA"`
}

type GcpIntegrationData struct {
	ID               string         `json:"ID"`
	IdType           string         `json:"ID_TYPE"`
	Credentials      GcpCredentials `json:"CREDENTIALS"`
	SubscriptionName string         `json:"SUBSCRIPTION_NAME,omitempty"`
}

type GcpCredentials struct {
	ClientId     string `json:"CLIENT_ID"`
	ClientEmail  string `json:"CLIENT_EMAIL"`
	PrivateKeyId string `json:"PRIVATE_KEY_ID"`
	PrivateKey   string `json:"PRIVATE_KEY"`
}
