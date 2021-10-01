//
// Author:: Vatasha White (<vatasha.white@lacework.net>)
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

// GetCloudwatchEb gets a single instance of an AWS Cloudwatch alert channel
// with the corresponding integration guid
func (svc *AlertChannelsService) GetCloudwatchEb(guid string) (response CloudwatchEbAlertChannelResponseV2, err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateCloudwatchEb Update AWSCloudWatch updates a single instance of an AWS cloudwatch integration on the Lacework server
func (svc *AlertChannelsService) UpdateCloudwatchEb(data AlertChannel) (response CloudwatchEbAlertChannelResponseV2, err error) {
	err = svc.update(data.ID(), data, &response)
	return
}

type CloudwatchEbDataV2 struct {
	EventBusArn   string `json:"eventBusArn"`
	IssueGrouping string `json:"issueGrouping,omitempty"`
}

type CloudwatchEbAlertChannelV2 struct {
	v2CommonIntegrationData
	Data CloudwatchEbDataV2 `json:"data"`
}

type CloudwatchEbAlertChannelResponseV2 struct {
	Data CloudwatchEbAlertChannelV2 `json:"data"`
}
