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

type AlertIntegrationChannel struct {
	commonIntegrationData
	EnvironmentGUID string                 `json:"ENV_GUID"`
	Data            map[string]interface{} `json:"DATA"`
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
