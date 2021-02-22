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

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// NewServiceNowAlertChannel returns an instance of ServiceNowAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new ServiceNowAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   serviceNowChannel := api.NewServiceNowAlertChannel("foo",
//     api.ServiceNowChannelData{
//       InstanceURL:   "snow-lacework.com",
//       Username:      "snow-user",
//       Password:      "snow-password",
//       IssueGrouping: "Events",
//     },
//   )
//
//   client.Integrations.CreateServiceNowAlertChannel(serviceNowChannel)
//
func NewServiceNowAlertChannel(name string, data ServiceNowChannelData) ServiceNowAlertChannel {
	return ServiceNowAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    ServiceNowChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateServiceNowAlertChannel creates a serviceNow alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateServiceNowAlertChannel(integration ServiceNowAlertChannel) (
	response ServiceNowAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetServiceNowAlertChannel gets a serviceNow alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetServiceNowAlertChannel(guid string) (response ServiceNowAlertChannelResponse,
	err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateServiceNowAlertChannel updates a single serviceNow alert channel integration
func (svc *IntegrationsService) UpdateServiceNowAlertChannel(data ServiceNowAlertChannel) (
	response ServiceNowAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListServiceNowAlertChannel lists the serviceNow alert channel integrations available on the Lacework Server
func (svc *IntegrationsService) ListServiceNowAlertChannel() (response ServiceNowAlertChannelResponse, err error) {
	err = svc.listByType(ServiceNowChannelIntegration, &response)
	return
}

type ServiceNowAlertChannelResponse struct {
	Data    []ServiceNowAlertChannel `json:"data"`
	Ok      bool                     `json:"ok"`
	Message string                   `json:"message"`
}

type ServiceNowAlertChannel struct {
	commonIntegrationData
	Data ServiceNowChannelData `json:"DATA"`
}

type ServiceNowChannelData struct {
	InstanceURL        string `json:"INSTANCE_URL" mapstructure:"INSTANCE_URL"`
	Username           string `json:"USER_NAME" mapstructure:"USER_NAME"`
	Password           string `json:"PASSWORD" mapstructure:"PASSWORD"`
	CustomTemplateFile string `json:"CUSTOM_TEMPLATE_FILE,omitempty" mapstructure:"CUSTOM_TEMPLATE_FILE"`
	IssueGrouping      string `json:"ISSUE_GROUPING,omitempty" mapstructure:"ISSUE_GROUPING"`
}

func (snow *ServiceNowChannelData) EncodeCustomTemplateFile(template string) {
	encodedTemplate := base64.StdEncoding.EncodeToString([]byte(template))
	snow.CustomTemplateFile = fmt.Sprintf("data:application/json;name=i.json;base64,%s", encodedTemplate)
}

func (snow *ServiceNowChannelData) DecodeCustomTemplateFile() (string, error) {
	if len(snow.CustomTemplateFile) == 0 {
		return "", nil
	}

	var (
		b64      = strings.Split(snow.CustomTemplateFile, ",")
		raw, err = base64.StdEncoding.DecodeString(b64[1])
	)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}
