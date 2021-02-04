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

import "github.com/pkg/errors"

type datadogSite string

type datadogService string

const (
	//The list of valid inputs for DatadogService field
	DatadogSiteEu  datadogSite = "eu"
	DatadogSiteCom datadogSite = "com"

	//The list of valid inputs for DatadogService field
	DatadogServiceLogsDetails   datadogService = "Logs Detail"
	DatadogServiceEventsSummary datadogService = "Events Summary"
	DatadogServiceLogsSummary   datadogService = "Logs Summary"
)

var datadogSites = map[string]datadogSite{
	string(DatadogSiteEu):  DatadogSiteEu,
	string(DatadogSiteCom): DatadogSiteCom,
}

var datadogServices = map[string]datadogService{
	string(DatadogServiceLogsDetails):   DatadogServiceLogsDetails,
	string(DatadogServiceEventsSummary): DatadogServiceEventsSummary,
	string(DatadogServiceLogsSummary):   DatadogServiceLogsSummary,
}

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
//   datadog := api.NewDatadogAlertChannel("foo",
//   api.DatadogChannelData{
// 		DatadogSite:    api.DatadogSiteEu.String(),
//  	DatadogService: api.DatadogServiceEventsSummary.String(),
// 	  	ApiKey:      	"datadog-key",
//   },
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

// DatadogSite returns the datadogSite type for the corresponding string input
func DatadogSite(site string) (datadogSite, error) {
	if val, ok := datadogSites[site]; ok {
		return val, nil
	}
	return "", errors.Errorf("%v is not a valid Datadog Site", site)
}

// DatadogService returns the datadogService type for the corresponding string input
func DatadogService(service string) (datadogService, error) {
	if val, ok := datadogServices[service]; ok {
		return val, nil
	}
	return "", errors.Errorf("%v is not a valid Datadog Site", service)
}

type DatadogAlertChannel struct {
	commonIntegrationData
	Data DatadogChannelData `json:"DATA"`
}

type DatadogChannelData struct {
	DatadogSite    datadogSite    `json:"DATADOG_SITE,omitempty" mapstructure:"DATADOG_SITE"`
	DatadogService datadogService `json:"DATADOG_TYPE,omitempty" mapstructure:"DATADOG_TYPE"`
	ApiKey         string         `json:"API_KEY" mapstructure:"API_KEY"`
}
