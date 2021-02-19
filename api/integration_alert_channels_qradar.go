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

// NewQRadarAlertChannel returns an instance of QRadarAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new QRadarAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   qradarChannel := api.NewQRadarAlertChannel("foo",
//     api.QRadarChannelData{
//       CommunicationType: "HTTPS",
//       HostURL:           "https://qradar-lacework.com",
//       HostPort:          8080,
//     },
//   )
//
//   client.Integrations.CreateQRadarAlertChannel(qradarChannel)
//
func NewQRadarAlertChannel(name string, data QRadarChannelData) QRadarAlertChannel {
	return QRadarAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    QRadarChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateQRadarAlertChannel creates a qradar alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateQRadarAlertChannel(integration QRadarAlertChannel) (
	response QRadarAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetQRadarAlertChannel gets a qradar alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetQRadarAlertChannel(guid string) (response QRadarAlertChannelResponse,
	err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateQRadarAlertChannel updates a single qradar alert channel integration
func (svc *IntegrationsService) UpdateQRadarAlertChannel(data QRadarAlertChannel) (
	response QRadarAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListQRadarAlertChannel lists the qradar alert channel integrations available on the Lacework Server
func (svc *IntegrationsService) ListQRadarAlertChannel() (response QRadarAlertChannelResponse, err error) {
	err = svc.listByType(QRadarChannelIntegration, &response)
	return
}

type QRadarAlertChannelResponse struct {
	Data    []QRadarAlertChannel `json:"data"`
	Ok      bool                 `json:"ok"`
	Message string               `json:"message"`
}

type QRadarAlertChannel struct {
	commonIntegrationData
	Data QRadarChannelData `json:"DATA"`
}

type QRadarChannelData struct {
	CommunicationType qradarComm `json:"QRADAR_COMM_TYPE,omitempty" mapstructure:"QRADAR_COMM_TYPE"`
	HostURL           string     `json:"QRADAR_HOST_URL" mapstructure:"QRADAR_HOST_URL"`
	HostPort          int        `json:"QRADAR_HOST_PORT,omitempty" mapstructure:"QRADAR_HOST_PORT"`
}

type qradarComm string

const (
	// The list of valid inputs for QRadar Communication Type field
	QRadarCommHttps           qradarComm = "HTTPS"
	QRadarCommHttpsSelfSigned qradarComm = "HTTPS Self Signed Cert"
)

var qradarCommTypes = map[string]qradarComm{
	string(QRadarCommHttps):           QRadarCommHttps,
	string(QRadarCommHttpsSelfSigned): QRadarCommHttpsSelfSigned,
}

// QRadarComm returns the qradarComm type for the corresponding string input
func QRadarComm(site string) (qradarComm, error) {
	if val, ok := qradarCommTypes[site]; ok {
		return val, nil
	}
	return "", errors.Errorf("%v is not a valid QRadar communication type", site)
}
