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

package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

var alertsJSON = `{
    "paging": {
        "rows": 412,
        "totalRows": 412,
        "urls": {
            "nextPage": null
        }
    },
    "data": [
        {
            "alertId": 967278,
            "startTime": "2022-09-29T16:00:00.000Z",
            "alertType": "CloudActivityLogIngestionFailed",
            "severity": "High",
            "reachability": "UnknownReachability",
            "derivedFields": {
                "category": "Policy",
                "sub_category": "Platform",
                "source": ""
            },
            "endTime": "2022-09-29T17:00:00.000Z",
            "lastUserUpdatedTime": "0",
            "status": "Open",
            "alertName": "Clone of Cloud Activity log ingestion failure detected",
            "alertInfo": {
                "subject": "Clone of Cloud Activity log ingestion failure detected",
                "description": "New integration failure detected for kiki-intg-2 (and 3 more)"
            },
            "policyId": "CUSTOM_PLATFORM_130",
            "alertSource": "UnknownAlertSource"
        }
	]
}`

func TestAlertsListMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Alerts",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.List()
	assert.Nil(t, err)
}

func TestAlertsListOK(t *testing.T) {
	mockResponse := alertsJSON

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Alerts",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, mockResponse)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	listExpected := api.AlertsResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &listExpected)

	var listActual api.AlertsResponse
	listActual, err = c.V2.Alerts.List()
	assert.Nil(t, err)
	assert.Equal(t, listExpected, listActual)
}

func TestAlertsListError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Alerts",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, lqlErrorReponse, http.StatusInternalServerError)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.List()
	assert.NotNil(t, err)
}

func TestAlertsListByTimeMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Alerts",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListByTime should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.ListByTime(time.Now(), time.Now())
	assert.Nil(t, err)
}

func TestAlertsListByTimeOK(t *testing.T) {
	mockResponse := alertsJSON

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Alerts",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, mockResponse)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	listExpected := api.AlertsResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &listExpected)

	var listActual api.AlertsResponse
	listActual, err = c.V2.Alerts.ListByTime(time.Now(), time.Now())
	assert.Nil(t, err)
	assert.Equal(t, listExpected, listActual)
}

func TestAlertsListByTimeError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Alerts",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, lqlErrorReponse, http.StatusInternalServerError)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.ListByTime(time.Now(), time.Now())
	assert.NotNil(t, err)
}
