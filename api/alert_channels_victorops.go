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

// GetVictorOps gets a single VictorOps alert channel matching the
// provided integration guid
func (svc *AlertChannelsService) GetVictorOps(guid string) (
	response VictorOpsAlertChannelResponseV2,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateVictorOps updates a single VictorOps integration on the Lacework Server
func (svc *AlertChannelsService) UpdateVictorOps(data AlertChannel) (
	response VictorOpsAlertChannelResponseV2,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type VictorOpsAlertChannelResponseV2 struct {
	Data VictorOpsAlertChannelV2 `json:"data"`
}

type VictorOpsAlertChannelV2 struct {
	v2CommonIntegrationData
	Data VictorOpsDataV2 `json:"data"`
}

type VictorOpsDataV2 struct {
	VictorOpsUrl string `json:"intgUrl"`
}
