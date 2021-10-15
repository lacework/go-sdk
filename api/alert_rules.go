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

import (
	"fmt"

	"github.com/pkg/errors"
)

// AlertRulesService is the service that interacts with
// the AlertRules schema from the Lacework APIv2 Server
type AlertRulesService struct {
	client *Client
}

type alertRuleSeverity int

type AlertRuleSeverities []alertRuleSeverity

const alertRuleEventType = "Event"

func (sevs AlertRuleSeverities) toInt() []int {
	var res []int
	for _, i := range sevs {
		res = append(res, int(i))
	}
	return res
}

const (
	AlertRuleSeverityCritical alertRuleSeverity = 1
	AlertRuleSeverityHigh     alertRuleSeverity = 2
	AlertRuleSeverityMedium   alertRuleSeverity = 3
	AlertRuleSeverityLow      alertRuleSeverity = 4
	AlertRuleSeverityInfo     alertRuleSeverity = 5
)

// NewAlertRule returns an instance of the AlertRuleData struct
//
// Basic usage: Initialize a new AlertRuleData struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   alertRule := api.NewAlertRule(
//		"Foo",
//		api.AlertRule{
//		Description: "My Alert Rule"
//		Severities: api.AlertRuleSeverities{api.AlertRuleSeverityHigh,
//		Channels: []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
//		ResourceGroups: []string{"TECHALLY_111111111111AAAAAAAAAAAAAAAAAAAA"}
//       },
//     },
//   )
//
//   client.V2.AlertRules.Create(alertRule)
//
func NewAlertRule(name string, rule AlertRule) AlertRuleData {
	return AlertRuleData{
		Channels: rule.Channels,
		Type:     alertRuleEventType,
		Filter: AlertRuleFilter{
			Name:            name,
			Enabled:         1,
			Description:     rule.Description,
			Severity:        rule.Severities.toInt(),
			ResourceGroups:  rule.ResourceGroups,
			EventCategories: rule.EventCategories,
		},
	}
}

// List returns a list of Alert Rules
func (svc *AlertRulesService) List() (response AlertRulesResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2AlertRules, nil, &response)
	return
}

// Create creates a single Alert Rule
func (svc *AlertRulesService) Create(rule AlertRuleData) (
	response AlertRuleResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2AlertRules, rule, &response)
	return
}

// Delete deletes a Alert Rule that matches the provided guid
func (svc *AlertRulesService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2AlertRuleFromGUID, guid),
		nil,
		nil,
	)
}

// Update updates a single Alert Rule of the provided guid.
func (svc *AlertRulesService) Update(data AlertRuleData) (
	response AlertRuleResponse,
	err error,
) {
	if data.Guid == "" {
		err = errors.New("specify a Guid")
		return
	}
	apiPath := fmt.Sprintf(apiV2AlertRuleFromGUID, data.Guid)
	err = svc.client.RequestEncoderDecoder("PATCH", apiPath, data, &response)
	return
}

// Get returns a raw response of the Alert Rule with the matching guid.
func (svc *AlertRulesService) Get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify a Guid")
	}
	apiPath := fmt.Sprintf(apiV2AlertRuleFromGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, &response)
}

type AlertRule struct {
	Channels        []string
	Description     string
	Severities      AlertRuleSeverities
	ResourceGroups  []string
	EventCategories []string
}

type AlertRuleData struct {
	Guid     string          `json:"mcGuid,omitempty"`
	Type     string          `json:"type"`
	Channels []string        `json:"intgGuidList"`
	Filter   AlertRuleFilter `json:"filters"`
}

type AlertRuleFilter struct {
	Name                 string   `json:"name"`
	Enabled              int      `json:"enabled"`
	Description          string   `json:"description,omitempty"`
	Severity             []int    `json:"severity"`
	ResourceGroups       []string `json:"resourceGroups,omitempty"`
	EventCategories      []string `json:"eventCategory,omitempty"`
	CreatedOrUpdatedTime string   `json:"createdOrUpdatedTime,omitempty"`
	CreatedOrUpdatedBy   string   `json:"createdOrUpdatedBy,omitempty"`
}

type AlertRuleResponse struct {
	Data AlertRuleData `json:"data"`
}

type AlertRulesResponse struct {
	Data []AlertRuleData `json:"data"`
}
