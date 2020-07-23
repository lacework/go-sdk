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

// NewPagerDutyAlertChannel returns an instance of PagerDutyAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new PagerDutyAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   pagerduty := api.NewPagerDutyAlertChannel("foo",
//     api.PagerDutyData{
//       IntegrationKey:   "1234abc8901abc567abc123abc78e012",
//       MinAlertSeverity: api.AllAlertLevel,
//     },
//   )
//
//   client.Integrations.CreatePagerDutyAlertChannel(pagerduty)
//
func NewPagerDutyAlertChannel(name string, data PagerDutyData) PagerDutyAlertChannel {
	return PagerDutyAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    PagerDutyIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreatePagerDutyAlertChannel creates a pager duty alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreatePagerDutyAlertChannel(integration PagerDutyAlertChannel) (
	response PagerDutyAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetPagerDutyAlertChannel gets a pager duty alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetPagerDutyAlertChannel(guid string) (
	response PagerDutyAlertChannelResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdatePagerDutyAlertChannel updates a single pager duty alert channel integration
func (svc *IntegrationsService) UpdatePagerDutyAlertChannel(data PagerDutyAlertChannel) (
	response PagerDutyAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListPagerDutyAlertChannel lists the PAGER_DUTY_API external integrations available on the Lacework Server
func (svc *IntegrationsService) ListPagerDutyAlertChannel() (response PagerDutyAlertChannelResponse, err error) {
	err = svc.listByType(PagerDutyIntegration, &response)
	return
}

type PagerDutyAlertChannelResponse struct {
	Data    []PagerDutyAlertChannel `json:"data"`
	Ok      bool                    `json:"ok"`
	Message string                  `json:"message"`
}

type PagerDutyAlertChannel struct {
	commonIntegrationData
	Data PagerDutyData `json:"DATA"`
}

type PagerDutyData struct {
	IssueGrouping    string     `json:"ISSUE_GROUPING,omitempty" mapstructure:"ISSUE_GROUPING"`
	IntegrationKey   string     `json:"API_INTG_KEY" mapstructure:"API_INTG_KEY"`
	MinAlertSeverity AlertLevel `json:"MIN_ALERT_SEVERITY,omitempty" mapstructure:"MIN_ALERT_SEVERITY"`
}
