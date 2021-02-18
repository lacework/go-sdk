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

// NewNewRelicAlertChannel returns an instance of NewRelicAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new NewRelicAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//	newRelicChannel := api.NewNewRelicAlertChannel("foo",
//		api.NewRelicChannelData{
//			AccountID: 2338053,
//			InsertKey: "x-xx-xxxxxxxxxxxxxxxxxx",
//		},
//	)
//
//   client.Integrations.CreateNewRelicAlertChannel(newRelicChannel)
//
func NewNewRelicAlertChannel(name string, data NewRelicChannelData) NewRelicAlertChannel {
	return NewRelicAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    NewRelicChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateNewRelicAlertChannel creates an NEW_RELIC_INSIGHTS alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateNewRelicAlertChannel(integration NewRelicAlertChannel) (
	response NewRelicAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetNewRelicAlertChannel gets an NEW_RELIC_INSIGHTS alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetNewRelicAlertChannel(guid string) (
	response NewRelicAlertChannelResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateNewRelicAlertChannel updates a single NEW_RELIC_INSIGHTS alert channel integration
func (svc *IntegrationsService) UpdateNewRelicAlertChannel(data NewRelicAlertChannel) (
	response NewRelicAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListNewRelicAlertChannel lists the NEW_RELIC_INSIGHTS external integrations available on the Lacework Server
func (svc *IntegrationsService) ListNewRelicAlertChannel() (response NewRelicAlertChannelResponse, err error) {
	err = svc.listByType(NewRelicChannelIntegration, &response)
	return
}

type NewRelicAlertChannelResponse struct {
	Data    []NewRelicAlertChannel `json:"data"`
	Ok      bool                   `json:"ok"`
	Message string                 `json:"message"`
}

type NewRelicAlertChannel struct {
	commonIntegrationData
	Data NewRelicChannelData `json:"DATA"`
}

type NewRelicChannelData struct {
	AccountID int    `json:"ACCOUNT_ID" mapstructure:"ACCOUNT_ID"`
	InsertKey string `json:"INSERT_KEY" mapstructure:"INSERT_KEY"`
}
