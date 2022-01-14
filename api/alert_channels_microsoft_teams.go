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

// GetMicrosoftTeams gets a single instance of a MicrosoftTeams alert channel
// with the corresponding integration guid
func (svc *AlertChannelsService) GetMicrosoftTeams(guid string) (response MicrosoftTeamsAlertChannelResponseV2, err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateMicrosoftTeams updates a single instance of a MicrosoftTeams integration on the Lacework server
func (svc *AlertChannelsService) UpdateMicrosoftTeams(data AlertChannel) (response MicrosoftTeamsAlertChannelResponseV2, err error) {
	err = svc.update(data.ID(), data, &response)
	return
}

type MicrosoftTeamsData struct {
	TeamsURL string `json:"teamsUrl"`
}

type MicrosoftTeamsAlertChannelV2 struct {
	v2CommonIntegrationData
	Data MicrosoftTeamsData `json:"data"`
}

type MicrosoftTeamsAlertChannelResponseV2 struct {
	Data MicrosoftTeamsAlertChannelV2 `json:"data"`
}
