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

// GetNewRelicInsights gets a single NewRelic alert channel matching the
// provided integration guid
func (svc *AlertChannelsService) GetNewRelicInsights(guid string) (
	response NewRelicInsightsAlertChannelResponseV2,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateNewRelicInsights updates a single NewRelic integration on the Lacework Server
func (svc *AlertChannelsService) UpdateNewRelicInsights(data AlertChannel) (
	response NewRelicInsightsAlertChannelResponseV2,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type NewRelicInsightsAlertChannelResponseV2 struct {
	Data NewRelicInsightsAlertChannelV2 `json:"data"`
}

type NewRelicInsightsAlertChannelV2 struct {
	v2CommonIntegrationData
	Data NewRelicInsightsDataV2 `json:"data"`
}

type NewRelicInsightsDataV2 struct {
	AccountID int    `json:"accountId"`
	InsertKey string `json:"insertKey"`
}
