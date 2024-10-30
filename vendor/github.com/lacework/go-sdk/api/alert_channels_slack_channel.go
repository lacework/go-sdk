//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

// GetSlackChannel gets a single SlackChannel alert channel matching the
// provided integration guid
func (svc *AlertChannelsService) GetSlackChannel(guid string) (
	response SlackChannelAlertChannelResponseV2,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateSlackChannel updates a single SlackChannel integration on the Lacework Server
func (svc *AlertChannelsService) UpdateSlackChannel(data AlertChannel) (
	response SlackChannelAlertChannelResponseV2,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type SlackChannelAlertChannelResponseV2 struct {
	Data SlackChannelAlertChannelV2 `json:"data"`
}

type SlackChannelAlertChannelV2 struct {
	v2CommonIntegrationData
	Data SlackChannelDataV2 `json:"data"`
}

type SlackChannelDataV2 struct {
	SlackUrl string `json:"slackUrl"`
}
