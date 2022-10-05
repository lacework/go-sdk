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

type RelatedAlert struct {
	ID        string    `json:"eventId"`
	Name      string    `json:"eventName"`
	Type      string    `json:"eventType"`
	Severity  string    `json:"severity"`
	Rank      int       `json:"rank"`
	Info      AlertInfo `json:"eventInfo"`
	StartTime string    `json:"startTime"`
	EndTime   string    `json:"endTime"`
}

type RelatedAlertsResponse struct {
	Data []RelatedAlert `json:"data"`
}

func (svc *AlertsService) GetRelatedAlerts(id int) (
	response RelatedAlertsResponse,
	err error,
) {
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf(apiV2AlertsDetails, id, AlertRelatedAlertsScope),
		nil,
		&response,
	)
	return
}
