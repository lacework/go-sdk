//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

// NewGcpPubSubAlertChannel returns an instance of GcpPubSubAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new GcpPubSubAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//	gcpPubSubChannel := api.NewGcpPubSubAlertChannel("foo",
//		api.GcpPubSubChannelData{
//			ProjectID:     "my-sample-project-191923",
//			TopicID:       "mytopic",
//			IssueGrouping: "Resources",
//			Credentials: api.GcpCredentials{
//				ClientID:     "client_id",
//				ClientEmail:  "foo@example.iam.gserviceaccount.com",
//				PrivateKey:   "priv_key",
//				PrivateKeyID: "p_key_id",
//			},
//		},
//		)
//
//   client.Integrations.CreateGcpPubSubAlertChannel(gcpPubSubChannel)
//
func NewGcpPubSubAlertChannel(name string, data GcpPubSubChannelData) GcpPubSubAlertChannel {
	return GcpPubSubAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    GcpPubSubChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateGcpPubSubAlertChannel creates an GCP_PUBSUB alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateGcpPubSubAlertChannel(integration GcpPubSubAlertChannel) (
	response GcpPubSubAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetGcpPubSubAlertChannel gets an GCP_PUBSUB alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetGcpPubSubAlertChannel(guid string) (
	response GcpPubSubAlertChannelResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateGcpPubSubAlertChannel updates a single GCP_PUBSUB alert channel integration
func (svc *IntegrationsService) UpdateGcpPubSubAlertChannel(data GcpPubSubAlertChannel) (
	response GcpPubSubAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListGcpPubSubAlertChannel lists the GCP_PUBSUB external integrations available on the Lacework Server
func (svc *IntegrationsService) ListGcpPubSubAlertChannel() (response GcpPubSubAlertChannelResponse, err error) {
	err = svc.listByType(GcpPubSubChannelIntegration, &response)
	return
}

type GcpPubSubAlertChannelResponse struct {
	Data    []GcpPubSubAlertChannel `json:"data"`
	Ok      bool                    `json:"ok"`
	Message string                  `json:"message"`
}

type GcpPubSubAlertChannel struct {
	commonIntegrationData
	Data GcpPubSubChannelData `json:"DATA"`
}

type GcpPubSubChannelData struct {
	ProjectID     string         `json:"PROJECT_ID" mapstructure:"PROJECT_ID"`
	TopicID       string         `json:"TOPIC_ID" mapstructure:"TOPIC_ID"`
	IssueGrouping string         `json:"ISSUE_GROUPING,omitempty" mapstructure:"ISSUE_GROUPING"`
	Credentials   GcpCredentials `json:"CREDENTIALS" mapstructure:"CREDENTIALS"`
}
