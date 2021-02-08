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

// NewMicrosoftTeamsAlertChannel returns an instance of MicrosoftTeamsAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new MicrosoftTeamsAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   microsoftTeamsChannel := api.NewMicrosoftTeamsAlertChannel("foo",
//     api.MicrosoftTeamsChannelData{
//       TeamsURL: "https://outlook.office.com/webhook/api-token",
//     },
//   )
//
//   client.Integrations.CreateMicrosoftTeamsAlertChannel(microsoftTeamsChannel)
//
func NewMicrosoftTeamsAlertChannel(name string, data MicrosoftTeamsChannelData) MicrosoftTeamsAlertChannel {
	return MicrosoftTeamsAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    MicrosoftTeamsChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateMicrosoftTeamsAlertChannel creates a msTeams alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateMicrosoftTeamsAlertChannel(integration MicrosoftTeamsAlertChannel) (
	response MicrosoftTeamsAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetMicrosoftTeamsAlertChannel gets a msTeams alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetMicrosoftTeamsAlertChannel(guid string) (response MicrosoftTeamsAlertChannelResponse,
	err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateMicrosoftTeamsAlertChannel updates a single msTeams alert channel integration
func (svc *IntegrationsService) UpdateMicrosoftTeamsAlertChannel(data MicrosoftTeamsAlertChannel) (
	response MicrosoftTeamsAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListMicrosoftTeamsAlertChannel lists the Microsoft Teams external integrations available on the Lacework Server
func (svc *IntegrationsService) ListMicrosoftTeamsAlertChannel() (response MicrosoftTeamsAlertChannelResponse, err error) {
	err = svc.listByType(MicrosoftTeamsChannelIntegration, &response)
	return
}

type MicrosoftTeamsAlertChannelResponse struct {
	Data    []MicrosoftTeamsAlertChannel `json:"data"`
	Ok      bool                         `json:"ok"`
	Message string                       `json:"message"`
}

type MicrosoftTeamsAlertChannel struct {
	commonIntegrationData
	Data MicrosoftTeamsChannelData `json:"DATA"`
}

type MicrosoftTeamsChannelData struct {
	TeamsURL string `json:"TEAMS_URL" mapstructure:"TEAMS_URL"`
}
