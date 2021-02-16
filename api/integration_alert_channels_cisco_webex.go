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

// NewCiscoWebexAlertChannel returns an instance of CiscoWebexAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new CiscoWebexAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   ciscoWebexChannel := api.NewCiscoWebexAlertChannel("foo",
//     api.CiscoWebexChannelData{
//       WebhookURL: "https://webexapis.com/v1/webhooks/incoming/api-token",
//     },
//   )
//
//   client.Integrations.CreateCiscoWebexAlertChannel(ciscoWebexChannel)
//
func NewCiscoWebexAlertChannel(name string, data CiscoWebexChannelData) CiscoWebexAlertChannel {
	return CiscoWebexAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    CiscoWebexChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateCiscoWebexAlertChannel creates a ciscoWebex alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateCiscoWebexAlertChannel(integration CiscoWebexAlertChannel) (
	response CiscoWebexAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetCiscoWebexAlertChannel gets a ciscoWebex alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetCiscoWebexAlertChannel(guid string) (response CiscoWebexAlertChannelResponse,
	err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateCiscoWebexAlertChannel updates a single ciscoWebex alert channel integration
func (svc *IntegrationsService) UpdateCiscoWebexAlertChannel(data CiscoWebexAlertChannel) (
	response CiscoWebexAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListCiscoWebexAlertChannel lists the WEBHOOK external integrationS available on the Lacework Server
func (svc *IntegrationsService) ListCiscoWebexAlertChannel() (response CiscoWebexAlertChannelResponse, err error) {
	err = svc.listByType(CiscoWebexChannelIntegration, &response)
	return
}

type CiscoWebexAlertChannelResponse struct {
	Data    []CiscoWebexAlertChannel `json:"data"`
	Ok      bool                     `json:"ok"`
	Message string                   `json:"message"`
}

type CiscoWebexAlertChannel struct {
	commonIntegrationData
	Data CiscoWebexChannelData `json:"DATA"`
}

type CiscoWebexChannelData struct {
	WebhookURL string `json:"WEBHOOK" mapstructure:"WEBHOOK"`
}
