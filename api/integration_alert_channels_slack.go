//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

// NewSlackAlertChannel returns an instance of SlackAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new SlackAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   slackChannel := api.NewSlackAlertChannel("foo",
//     api.SlackChannelData{
//       SlackUrl: "https://hooks.slack.com/services/ABCD/12345/abcd1234",
//       MinAlertSeverity: api.CriticalAlertLevel,
//     },
//   )
//
//   client.Integrations.CreateSlackAlertChannel(slackChannel)
//
func NewSlackAlertChannel(name string, data SlackChannelData) SlackAlertChannel {
	return SlackAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    SlackChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateSlackAlertChannel creates a slack alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateSlackAlertChannel(integration SlackAlertChannel) (
	response SlackAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetSlackAlertChannel gets a slack alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetSlackAlertChannel(guid string) (
	response SlackAlertChannelResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateSlackAlertChannel updates a single slack alert channel integration
func (svc *IntegrationsService) UpdateSlackAlertChannel(data SlackAlertChannel) (
	response SlackAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListSlackAlertChannel lists the SLACK_CHANNEL external integrations available on the Lacework Server
func (svc *IntegrationsService) ListSlackAlertChannel() (response SlackAlertChannelResponse, err error) {
	err = svc.listByType(SlackChannelIntegration, &response)
	return
}

type SlackAlertChannelResponse struct {
	Data    []SlackAlertChannel `json:"data"`
	Ok      bool                `json:"ok"`
	Message string              `json:"message"`
}

type SlackAlertChannel struct {
	commonIntegrationData
	Data SlackChannelData `json:"DATA"`
}

type SlackChannelData struct {
	IssueGrouping    string     `json:"ISSUE_GROUPING,omitempty" mapstructure:"ISSUE_GROUPING"`
	SlackUrl         string     `json:"SLACK_URL" mapstructure:"SLACK_URL"`
	MinAlertSeverity AlertLevel `json:"MIN_ALERT_SEVERITY,omitempty" mapstructure:"MIN_ALERT_SEVERITY"`
}
