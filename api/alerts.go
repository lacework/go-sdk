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
	"sort"
	"time"

	"github.com/lacework/go-sdk/lwtime"
)

// AlertsService is a service that interacts with the Alerts
// endpoints from the Lacework Server
type AlertsService struct {
	client *Client
}

type AlertInfo struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type AlertSpec struct {
	Profile string `json:"alertProfile"`
	Name    string `json:"name"`
}

type AlertDerivedFields struct {
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
	Source      string `json:"source"`
}

type Alert struct {
	ID            int                `json:"alertId"`
	Name          string             `json:"alertName"`
	Type          string             `json:"alertType"`
	Severity      string             `json:"severity"`
	Info          AlertInfo          `json:"alertInfo"`
	Spec          AlertSpec          `json:"alertSpec"`
	Status        string             `json:"status"`
	StartTime     string             `json:"startTime"`
	EndTime       string             `json:"endTime"`
	UpdateTime    string             `json:"lastUserUpdateTime"`
	PolicyID      string             `json:"policyId"`
	DerivedFields AlertDerivedFields `json:"derivedFields"`
	Reachability  string             `json:"reachability"`
}

type Alerts []Alert

func (a Alerts) SortDescending() []Alert {
	sort.Slice(a, func(i, j int) bool {
		return a[i].ID > a[j].ID
	})
	return a
}

type AlertsResponse struct {
	Data Alerts `json:"data"`
}

func (svc *AlertsService) List() (
	response AlertsResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiV2Alerts, nil, &response)
	return
}

func (svc *AlertsService) ListByTime(start, end time.Time) (
	response AlertsResponse,
	err error,
) {
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf(
			apiV2AlertsByTime,
			start.UTC().Format(lwtime.RFC3339Milli),
			end.UTC().Format(lwtime.RFC3339Milli),
		),
		nil,
		&response,
	)
	return
}
