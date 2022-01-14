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

// GetCiscoSparkWebhook gets a single instance of a Cisco Spark webhook alert channel
// with the corresponding integration guid
func (svc *AlertChannelsService) GetCiscoSparkWebhook(guid string) (response CiscoSparkWebhookAlertChannelResponseV2, err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateCiscoSparkWebhook updates a single instance of Cisco Spark webhook integration on the Lacework server
func (svc *AlertChannelsService) UpdateCiscoSparkWebhook(data AlertChannel) (response CiscoSparkWebhookAlertChannelResponseV2, err error) {
	err = svc.update(data.ID(), data, &response)
	return
}

type CiscoSparkWebhookDataV2 struct {
	Webhook string `json:"webhook"`
}

type CiscoSparkWebhookAlertChannelV2 struct {
	v2CommonIntegrationData
	Data CiscoSparkWebhookDataV2 `json:"data"`
}

type CiscoSparkWebhookAlertChannelResponseV2 struct {
	Data CiscoSparkWebhookAlertChannelV2 `json:"data"`
}
