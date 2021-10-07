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

// GetDatadog gets a single instance of a Datadog alert channel
// with the corresponding integration guid
func (svc *AlertChannelsService) GetDatadog(guid string) (response DatadogAlertChannelResponseV2, err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateDatadog updates a single instance of a Datadog integration on the Lacework server
func (svc *AlertChannelsService) UpdateDatadog(data AlertChannel) (response DatadogAlertChannelResponseV2, err error) {
	err = svc.update(data.ID(), data, &response)
	return
}

type DatadogDataV2 struct {
	ApiKey      string         `json:"apiKey"`
	DatadogSite datadogSite    `json:"datadogSite,omitempty"`
	DatadogType datadogService `json:"datadogType,omitempty"`
}

type DatadogAlertChannelV2 struct {
	v2CommonIntegrationData
	Data DatadogDataV2 `json:"data"`
}

type DatadogAlertChannelResponseV2 struct {
	Data DatadogAlertChannelV2 `json:"data"`
}
