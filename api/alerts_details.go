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
	"net/http"

	"github.com/pkg/errors"
)

type alertScope int

const (
	AlertDetailsScope alertScope = iota
	AlertInvestigationScope
	AlertEventsScope
	AlertRelatedAlertsScope
	AlertIntegrationsScope
	AlertTimelineScope
)

var AlertScopes = map[alertScope]string{
	AlertDetailsScope:       "Details",
	AlertInvestigationScope: "Investigation",
	AlertEventsScope:        "Events",
	AlertRelatedAlertsScope: "RelatedAlerts",
	AlertIntegrationsScope:  "Integrations",
	AlertTimelineScope:      "Timeline",
}

func (i alertScope) String() string {
	return AlertScopes[i]
}

type AlertDetails struct {
	Alert
	EntityMap map[string]interface{} `json:"entityMap"` // @dhazekamp: this needs to be built out properly
}

type AlertDetailsResponse struct {
	Data AlertDetails `json:"data"`
}

func (svc *AlertsService) Get(id int, scope alertScope) (interface{}, error) {
	switch scope {
	case AlertDetailsScope:
		return svc.GetDetails(id)
	case AlertInvestigationScope:
		return svc.GetInvestigation(id)
	case AlertEventsScope:
		return svc.GetEvents(id)
	case AlertRelatedAlertsScope:
		return svc.GetRelatedAlerts(id)
	case AlertIntegrationsScope:
		return svc.GetIntegrations(id)
	case AlertTimelineScope:
		return svc.GetTimeline(id)
	default:
		return nil, errors.New(fmt.Sprintf("alert scope (%s) not recognized", scope))
	}
}

func (svc *AlertsService) GetDetails(id int) (
	response AlertDetailsResponse,
	err error,
) {
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf(apiV2AlertsDetails, id, AlertDetailsScope),
		nil,
		&response,
	)
	return
}

func (svc *AlertsService) Exists(id int) (bool, error) {
	var response AlertDetailsResponse
	err := svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf(apiV2AlertsDetails, id, AlertDetailsScope),
		nil,
		&response,
	)

	if err == nil {
		return true, nil
	}
	errResponse, ok := err.(*errorResponse)
	if ok && errResponse.Response.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, err
}
