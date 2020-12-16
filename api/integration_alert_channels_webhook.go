//
// Author:: Darren Murray(<darren.murray@lacework.net>)
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

// NewWebhookAlertChannel returns an instance of WebhookAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new WebhookAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   webhookChannel := api.NewWebhookAlertChannel("foo",
//     api.WebhookChannelData{
//       WebhookUrl: "https://mywebhook.com/?api-token=123",
//     },
//   )
//
//   client.Integrations.CreateWebhookAlertChannel(webhookChannel)
//
func NewWebhookAlertChannel(name string, data WebhookChannelData) WebhookAlertChannel {
	return WebhookAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    WebhookIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateWebhookAlertChannel creates a webhook alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateWebhookAlertChannel(integration WebhookAlertChannel) (
	response WebhookAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetWebhookAlertChannel gets a webhook alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetWebhookAlertChannel(guid string) (response WebhookAlertChannelResponse,
	err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateWebhookAlertChannel updates a single webhook alert channel integration
func (svc *IntegrationsService) UpdateWebhookAlertChannel(data WebhookAlertChannel) (
	response WebhookAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListWebhookAlertChannel lists the WEBHOOK external integrationS available on the Lacework Server
func (svc *IntegrationsService) ListWebhookAlertChannel() (response WebhookAlertChannelResponse, err error) {
	err = svc.listByType(WebhookIntegration, &response)
	return
}

type WebhookAlertChannelResponse struct {
	Data    []WebhookAlertChannel `json:"data"`
	Ok      bool                  `json:"ok"`
	Message string                `json:"message"`
}

type WebhookAlertChannel struct {
	commonIntegrationData
	Data WebhookChannelData `json:"DATA"`
}

type WebhookChannelData struct {
	WebhookUrl string `json:"WEBHOOK_URL" mapstructure:"WEBHOOK_URL"`
}
