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

// NewSlackChannelIntegration returns an instance of SlackChanIntegration
// with the provided name and data.
//
// Basic usage: Initialize a new SlackChanIntegration struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   slackChannel := api.NewSlackChannelIntegration("foo",
//     api.SlackChannelData{
//       SlackUrl: "https://hooks.slack.com/services/ABCD/12345/abcd1234",
//       MinAlertSeverity: 3,
//     },
//   )
//
//   client.Integrations.CreateSlackChannel(slackChannel)
//
func NewSlackChannelIntegration(name string, data SlackChannelData) SlackChanIntegration {
	return SlackChanIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    SlackChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

type SlackAlertLevel int

const (
	CriticalSlackAlertLevel SlackAlertLevel = 1
	HighSlackAlertLevel     SlackAlertLevel = 2
	MediumSlackAlertLevel   SlackAlertLevel = 3
	LowSlackAlertLevel      SlackAlertLevel = 4
	AllSlackAlertLevel      SlackAlertLevel = 5
)

// SlackAlertLevels is the list of available slack alert levels
var SlackAlertLevels = map[SlackAlertLevel]string{
	CriticalSlackAlertLevel: "Critical",
	HighSlackAlertLevel:     "High",
	MediumSlackAlertLevel:   "Medium",
	LowSlackAlertLevel:      "Low",
	AllSlackAlertLevel:      "All",
}

// String returns the string representation of a slack alert level
func (i SlackAlertLevel) String() string {
	return SlackAlertLevels[i]
}

// Int returns the int representation of a slack alert level
func (i SlackAlertLevel) Int() int {
	return int(i)
}

// CreateSlackChannel creates a slack channel alert integration on the Lacework Server
func (svc *IntegrationsService) CreateSlackChannel(integration SlackChanIntegration) (
	response SlackChanIntResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetSlackChannel gets a slack channel alert integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetSlackChannel(guid string) (
	response SlackChanIntResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateSlackChannel updates a single slack channel alert integration
func (svc *IntegrationsService) UpdateSlackChannel(data SlackChanIntegration) (
	response SlackChanIntResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListSlackChannel lists the SLACK_CHANNEL external integrations available on the Lacework Server
func (svc *IntegrationsService) ListSlackChannel() (response SlackChanIntResponse, err error) {
	err = svc.listByType(SlackChannelIntegration, &response)
	return
}

type SlackChanIntResponse struct {
	Data    []SlackChanIntegration `json:"data"`
	Ok      bool                   `json:"ok"`
	Message string                 `json:"message"`
}

type SlackChanIntegration struct {
	commonIntegrationData
	Data SlackChannelData `json:"DATA"`
}

type SlackChannelData struct {
	IssueGrouping    string          `json:"ISSUE_GROUPING,omitempty" mapstructure:"ISSUE_GROUPING"`
	SlackUrl         string          `json:"SLACK_URL" mapstructure:"SLACK_URL"`
	MinAlertSeverity SlackAlertLevel `json:"MIN_ALERT_SEVERITY" mapstructure:"MIN_ALERT_SEVERITY"`
}
