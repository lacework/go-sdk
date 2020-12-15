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

func (svc *IntegrationsService) CreateWebhookAlertChannel(integration WebhookAlertChannel) (
	response WebhookAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

func (svc *IntegrationsService) GetWebhookAlertChannel(guid string) (response WebhookAlertChannelResponse,
	err error) {
	err = svc.get(guid, &response)
	return
}

func (svc *IntegrationsService) UpdateWebhookAlertChannel(data WebhookAlertChannel) (
	response WebhookAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

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
	IssueGrouping string `json:"ISSUE_GROUPING,omitempty" mapstructure:"ISSUE_GROUPING"`
	WebhookUrl    string `json:"WEBHOOK_URL" mapstructure:"WEBHOOK_URL"`
}
