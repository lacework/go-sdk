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

// NewDatadogAlertChannel returns an instance of DatadogAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new DatadogAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   datadogChannel := api.NewDatadogAlertChannel("foo",
//     api.DatadogChannelData{
//       DatadogSite: "com",
//       DatadogType: "Events Summary",
//       ApiKey: 	  "key",
//     },
//   )
//
//   client.Integrations.CreateDatadogAlertChannel(datadogChannel)
//
func NewDatadogAlertChannel(name string, data DatadogChannelData) DatadogAlertChannel {
	return DatadogAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    DatadogChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateDatadogAlertChannel creates a datadog alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateDatadogAlertChannel(integration DatadogAlertChannel) (
	response DatadogAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetDatadogAlertChannel gets a datadog alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetDatadogAlertChannel(guid string) (response DatadogAlertChannelResponse,
	err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateDatadogAlertChannel updates a single datadog alert channel integration
func (svc *IntegrationsService) UpdateDatadogAlertChannel(data DatadogAlertChannel) (
	response DatadogAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListDatadogAlertChannel lists the datadog alert channel integrations available on the Lacework Server
func (svc *IntegrationsService) ListDatadogAlertChannel() (response DatadogAlertChannelResponse, err error) {
	err = svc.listByType(DatadogChannelIntegration, &response)
	return
}

type DatadogAlertChannelResponse struct {
	Data    []DatadogAlertChannel `json:"data"`
	Ok      bool                  `json:"ok"`
	Message string                `json:"message"`
}

type DatadogAlertChannel struct {
	commonIntegrationData
	Data DatadogChannelData `json:"DATA"`
}

type DatadogChannelData struct {
	DatadogSite string `json:"DATADOG_SITE" mapstructure:"DATADOG_SITE"`
	DatadogType string `json:"DATADOG_TYPE" mapstructure:"DATADOG_TYPE"`
	ApiKey      string `json:"API_KEY" mapstructure:"API_KEY"`
}
