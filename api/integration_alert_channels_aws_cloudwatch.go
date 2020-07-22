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

// NewAwsCloudWatchAlertChannel returns an instance of AwsCloudWatchAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new AwsCloudWatchAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   awsCloudWatch := api.NewAwsCloudWatchAlertChannel("foo",
//     api.AwsCloudWatchData{
//       EventBusArn: "arn:aws:events:us-west-2:1234567890:event-bus/default",
//       MinAlertSeverity: api.MediumAlertLevel,
//     },
//   )
//
//   client.Integrations.CreateAwsCloudWatchAlertChannel(awsCloudWatch)
//
func NewAwsCloudWatchAlertChannel(name string, data AwsCloudWatchData) AwsCloudWatchAlertChannel {
	return AwsCloudWatchAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    AwsCloudWatchIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateAwsCloudWatchAlertChannel creates a AWS CloudWatch alert channel on the Lacework Server
func (svc *IntegrationsService) CreateAwsCloudWatchAlertChannel(integration AwsCloudWatchAlertChannel) (
	response AwsCloudWatchResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAwsCloudWatchAlertChannel gets a AWS CloudWatch alert channel that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetAwsCloudWatchAlertChannel(guid string) (
	response AwsCloudWatchResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsCloudWatchAlertChannel updates a single AWS CloudWatch alert channel
func (svc *IntegrationsService) UpdateAwsCloudWatchAlertChannel(data AwsCloudWatchAlertChannel) (
	response AwsCloudWatchResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListAwsCloudWatchAlertChannel lists the CLOUDWATCH_EB external integrations available on the Lacework Server
func (svc *IntegrationsService) ListAwsCloudWatchAlertChannel() (response AwsCloudWatchResponse, err error) {
	err = svc.listByType(AwsCloudWatchIntegration, &response)
	return
}

type AwsCloudWatchResponse struct {
	Data    []AwsCloudWatchAlertChannel `json:"data"`
	Ok      bool                        `json:"ok"`
	Message string                      `json:"message"`
}

type AwsCloudWatchAlertChannel struct {
	commonIntegrationData
	Data AwsCloudWatchData `json:"DATA"`
}

type AwsCloudWatchData struct {
	IssueGrouping    string     `json:"ISSUE_GROUPING,omitempty" mapstructure:"ISSUE_GROUPING"`
	EventBusArn      string     `json:"EVENT_BUS_ARN" mapstructure:"EVENT_BUS_ARN"`
	MinAlertSeverity AlertLevel `json:"MIN_ALERT_SEVERITY,omitempty" mapstructure:"MIN_ALERT_SEVERITY"`
}
