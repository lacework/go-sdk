//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

// GetIbmQRadar gets a single IbmQRadar alert channel matching the
// provided integration guid
func (svc *AlertChannelsService) GetIbmQRadar(guid string) (
	response IbmQRadarAlertChannelResponseV2,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateIbmQRadar updates a single IbmQRadar integration on the Lacework Server
func (svc *AlertChannelsService) UpdateIbmQRadar(data AlertChannel) (
	response IbmQRadarAlertChannelResponseV2,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type IbmQRadarAlertChannelResponseV2 struct {
	Data IbmQRadarAlertChannelV2 `json:"data"`
}

type IbmQRadarAlertChannelV2 struct {
	v2CommonIntegrationData
	Data IbmQRadarDataV2 `json:"data"`
}

type IbmQRadarDataV2 struct {
	QRadarCommType qradarComm `json:"qradarCommType"`
	HostURL        string     `json:"qradarHostUrl"`
	HostPort       int        `json:"qradarHostPort,omitempty"`
}
