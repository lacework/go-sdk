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

// GetGcpPubSub gets a single instance of a GCP Pub Sub alert channel with the corresponding guid
func (svc *AlertChannelsService) GetGcpPubSub(guid string) (response GcpPubSubAlertChannelResponseV2, err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateGcpPubSub updates a single instance of GCP Pub Sub integration on the Lacework server
func (svc *AlertChannelsService) UpdateGcpPubSub(data AlertChannel) (response GcpPubSubAlertChannelResponseV2, err error) {
	err = svc.update(data.ID(), data, &response)
	return
}

type GcpPubSubDataV2 struct {
	Credentials   GcpPubSubCredentials `json:"credentials"`
	IssueGrouping string               `json:"issueGrouping"`
	ProjectID     string               `json:"projectId"`
	TopicID       string               `json:"topicId"`
}

type GcpPubSubAlertChannelV2 struct {
	v2CommonIntegrationData
	Data GcpPubSubDataV2 `json:"data"`
}

type GcpPubSubAlertChannelResponseV2 struct {
	Data GcpPubSubAlertChannelV2 `json:"data"`
}

type GcpPubSubCredentials struct {
	ClientEmail  string `json:"clientEmail"`
	ClientID     string `json:"clientId"`
	PrivateKey   string `json:"privateKey"`
	PrivateKeyID string `json:"privateKeyId"`
}
