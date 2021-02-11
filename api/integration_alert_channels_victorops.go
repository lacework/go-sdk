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

// NewVictorOpsAlertChannel returns an instance of VictorOpsAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new VictorOpsAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   datadog := api.NewVictorOpsAlertChannel("foo",
//   api.VictorOpsChannelData{
// 		WebhookURL: "https://alert.victorops.com/integrations/generic/20131114/alert/31e945ee-5cad-44e7-afb0-97c20ea80dd8/database,
//   },
//   )
//
//   client.Integrations.CreateVictorOpsAlertChannel(datadogChannel)
//
func NewVictorOpsAlertChannel(name string, data VictorOpsChannelData) VictorOpsAlertChannel {
	return VictorOpsAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    VictorOpsChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateVictorOpsAlertChannel creates a datadog alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateVictorOpsAlertChannel(integration VictorOpsAlertChannel) (
	response VictorOpsAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetVictorOpsAlertChannel gets a datadog alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetVictorOpsAlertChannel(guid string) (response VictorOpsAlertChannelResponse,
	err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateVictorOpsAlertChannel updates a single datadog alert channel integration
func (svc *IntegrationsService) UpdateVictorOpsAlertChannel(data VictorOpsAlertChannel) (
	response VictorOpsAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListVictorOpsAlertChannel lists the datadog alert channel integrations available on the Lacework Server
func (svc *IntegrationsService) ListVictorOpsAlertChannel() (response VictorOpsAlertChannelResponse, err error) {
	err = svc.listByType(VictorOpsChannelIntegration, &response)
	return
}

type VictorOpsAlertChannelResponse struct {
	Data    []VictorOpsAlertChannel `json:"data"`
	Ok      bool                    `json:"ok"`
	Message string                  `json:"message"`
}

// VictorOpsSite returns the datadogSite type for the corresponding string input
func VictorOpsSite(site string) (datadogSite, error) {
	if val, ok := datadogSites[site]; ok {
		return val, nil
	}
	return "", errors.Errorf("%v is not a valid VictorOps Site", site)
}

// VictorOpsService returns the datadogService type for the corresponding string input
func VictorOpsService(service string) (datadogService, error) {
	if val, ok := datadogServices[service]; ok {
		return val, nil
	}
	return "", errors.Errorf("%v is not a valid VictorOps Service", service)
}

type VictorOpsAlertChannel struct {
	commonIntegrationData
	Data VictorOpsChannelData `json:"DATA"`
}

type VictorOpsChannelData struct {
	WebhookURL string `json:"INTG_URL" mapstructure:"INTG_URL"`
}
