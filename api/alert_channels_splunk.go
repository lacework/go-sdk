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

// GetSplunk gets a single Splunk alert channel matching the
// provided integration guid
func (svc *AlertChannelsService) GetSplunk(guid string) (
	response SplunkAlertChannelResponseV2,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateSplunk updates a single Splunk integration on the Lacework Server
func (svc *AlertChannelsService) UpdateSplunk(data AlertChannel) (
	response SplunkAlertChannelResponseV2,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type SplunkAlertChannelResponseV2 struct {
	Data SplunkAlertChannelV2 `json:"data"`
}

type SplunkAlertChannelV2 struct {
	v2CommonIntegrationData
	Data SplunkDataV2 `json:"data"`
}

type SplunkDataV2 struct {
	HecToken  string            `json:"hecToken"`
	Channel   string            `json:"channel,omitempty"`
	Host      string            `json:"host"`
	Port      int               `json:"port"`
	Ssl       bool              `json:"ssl,omitempty"`
	EventData SplunkEventDataV2 `json:"eventData"`
}

type SplunkEventDataV2 struct {
	Index  string `json:"index"`
	Source string `json:"source"`
}
