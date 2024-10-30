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

import (
	"fmt"
)

type AlertIntegrationChannelState struct {
	Ok                 bool                   `json:"ok"`
	LastUpdatedTime    int                    `json:"lastUpdatedTime"`
	LastSuccessfulTime int                    `json:"lastSuccessfulTime"`
	Details            map[string]interface{} `json:"details,omitempty"`
}

type AlertIntegrationChannel struct {
	IntgGuid             string                       `json:"INTG_GUID,omitempty"`
	Name                 string                       `json:"NAME"`
	CreatedOrUpdatedTime string                       `json:"CREATED_OR_UPDATED_TIME,omitempty"`
	CreatedOrUpdatedBy   string                       `json:"CREATED_OR_UPDATED_BY,omitempty"`
	Type                 string                       `json:"TYPE"`
	Enabled              int                          `json:"ENABLED"`
	State                AlertIntegrationChannelState `json:"STATE,omitempty"`
	IsOrg                int                          `json:"IS_ORG,omitempty"`
	TypeName             string                       `json:"TYPE_NAME,omitempty"`
	EnvironmentGUID      string                       `json:"ENV_GUID"`
	Data                 map[string]interface{}       `json:"DATA"`
}

func (c AlertIntegrationChannel) Status() string {
	if c.Enabled == 1 {
		return "Enabled"
	}
	return "Disabled"
}

func (c AlertIntegrationChannel) StateString() string {
	if c.State.Ok {
		return "Ok"
	}
	return "Pending"
}

type AlertIntegrationContext struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

type AlertIntegration struct {
	ID            string                  `json:"alertIntegrationId"`
	AlertID       int                     `json:"alertId"`
	Type          string                  `json:"integrationType"`
	Channel       AlertIntegrationChannel `json:"alertChannel"`
	Context       AlertIntegrationContext `json:"integrationContext"`
	IntgGUID      string                  `json:"intgGuid"`
	LastSyncTime  string                  `json:"lastSyncTime"`
	Status        string                  `json:"status"`
	Bidirectional bool                    `json:"isBidirectional"`
}

type AlertIntegrationsResponse struct {
	Data []AlertIntegration `json:"data"`
}

func (svc *AlertsService) GetIntegrations(id int) (
	response AlertIntegrationsResponse,
	err error,
) {
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf(apiV2AlertsDetails, id, AlertIntegrationsScope),
		nil,
		&response,
	)
	return
}
