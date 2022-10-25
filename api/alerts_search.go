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

type alertsFilterField string

const (
	AlertsFilterFieldType     alertsFilterField = "alertType"
	AlertsFilterFieldSeverity alertsFilterField = "severity"
	AlertsFilterFieldStatus   alertsFilterField = "status"
)

func (svc *AlertsService) Search(filter SearchFilter) (
	response AlertsResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder(
		"POST",
		apiV2AlertsSearch,
		filter,
		&response,
	)
	return
}

func (svc *AlertsService) SearchAll(filter SearchFilter) (
	response AlertsResponse,
	err error,
) {
	response, err = svc.Search(filter)
	if err != nil {
		return
	}

	var (
		all    Alerts
		pageOk bool
	)
	for {
		all = append(all, response.Data...)

		pageOk, err = svc.client.NextPage(&response)
		if err == nil && pageOk {
			continue
		}
		break
	}

	response.Data = all
	response.ResetPaging()
	return
}
