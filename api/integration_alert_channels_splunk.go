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

// NewSplunkAlertChannel returns an instance of SplunkAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new SplunkAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   splunkChannel := api.NewSplunkAlertChannel("foo",
//     api.SplunkChannelData{
//       Channel: "channel-name",
//       HecToken: "AA111111-11AA-1AA1-11AA-11111AA1111A",
//       Host: "localhost",
//       Port: 80,
//       Ssl: false,
//       EventData: api.SplunkEventData{
//         Index: "index",
//         Source: "source",
//        },
//     },
//   )
//
//   client.Integrations.CreateSplunkAlertChannel(splunkChannel)
//
func NewSplunkAlertChannel(name string, data SplunkChannelData) SplunkAlertChannel {
	return SplunkAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    SplunkIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateSplunkAlertChannel creates a splunk alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateSplunkAlertChannel(integration SplunkAlertChannel) (
	response SplunkAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetSplunkAlertChannel gets a splunk alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetSplunkAlertChannel(guid string) (response SplunkAlertChannelResponse,
	err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateSplunkAlertChannel updates a single splunk alert channel integration
func (svc *IntegrationsService) UpdateSplunkAlertChannel(data SplunkAlertChannel) (
	response SplunkAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListSplunkAlertChannel lists the splunk alert channel integrations available on the Lacework Server
func (svc *IntegrationsService) ListSplunkAlertChannel() (response SplunkAlertChannelResponse, err error) {
	err = svc.listByType(SplunkIntegration, &response)
	return
}

type SplunkAlertChannelResponse struct {
	Data    []SplunkAlertChannel `json:"data"`
	Ok      bool                 `json:"ok"`
	Message string               `json:"message"`
}

type SplunkAlertChannel struct {
	commonIntegrationData
	Data SplunkChannelData `json:"DATA"`
}

type SplunkChannelData struct {
	Channel   string          `json:"CHANNEL,omitempty" mapstructure:"CHANNEL"`
	HecToken  string          `json:"HEC_TOKEN" mapstructure:"HEC_TOKEN"`
	Host      string          `json:"HOST" mapstructure:"HOST"`
	Port      int             `json:"PORT" mapstructure:"PORT"`
	Ssl       bool            `json:"SSL" mapstructure:"SSL"`
	EventData SplunkEventData `json:"EVENT_DATA" mapstructure:"EVENT_DATA"`
}

type SplunkEventData struct {
	Index  string `json:"INDEX" mapstructure:"INDEX"`
	Source string `json:"SOURCE" mapstructure:"SOURCE"`
}
