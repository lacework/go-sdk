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

type AlertTimelineMessage struct {
	Format string `json:"format"`
	Value  string `json:"value"`
}

type AlertTimelineUser struct {
	UserGUID string `json:"userGuid"`
	Name     string `json:"username"`
}

type AlertTimelineNewIntegrationContext struct {
	AlertID                string `json:"alertId"`
	LastSyncTime           string `json:"lastSyncTime"`
	AlertIntegrationStatus string `json:"alertIntegrationStatus"`
	Status                 string `json:"status"`
	Bidirectional          bool   `json:"isBidirectional"`
}

type AlertTimelineUpdateContext struct {
	NewIntegration AlertTimelineNewIntegrationContext `json:"newIntegration"`
}

type AlertTimeline struct {
	ID              int                        `json:"id"`
	AlertID         int                        `json:"alertId"`
	EntryType       string                     `json:"entryType"`
	EntryAuthorType string                     `json:"entryAuthorType"`
	IntgGUID        string                     `json:"intgGuid"`
	Message         AlertTimelineMessage       `json:"message"`
	ExternalTime    string                     `json:"externalTime"`
	User            AlertTimelineUser          `json:"user"`
	UpdateContext   AlertTimelineUpdateContext `json:"updateContext"`
	Channel         AlertIntegrationChannel    `json:"alertChannel"`
}

type AlertTimelineResponse struct {
	Data []AlertTimeline `json:"data"`
}

func (svc *AlertsService) GetTimeline(id int) (
	response AlertTimelineResponse,
	err error,
) {
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf(apiV2AlertsDetails, id, AlertTimelineScope),
		nil,
		&response,
	)
	return
}
