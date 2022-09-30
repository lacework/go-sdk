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

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

var alertDetailsJSON = `{
    "data": {
        "alertId": 967272,
        "startTime": "2022-09-29T15:00:00.000Z",
        "alertType": "ChangedFile",
        "severity": "Medium",
        "reachability": "UnknownReachability",
        "derivedFields": {
            "category": "Policy",
            "sub_category": "File",
            "source": "Agent"
        },
        "endTime": "2022-09-29T16:00:00.000Z",
        "lastUserUpdatedTime": "0",
        "status": "Open",
        "alertName": "Clone of Files Changed ",
        "alertInfo": {
            "subject": "Clone of Files Changed",
            "description": "Custom Policy Violation"
        },
        "policyId": "CUSTOM_FIM_276",
        "alertSource": "UnknownAlertSource",
        "entityMap": {
            "CustomRule": [
                {
                    "KEY": {
                        "rule_guid": "OBJ_503B388C004E428E23D86397CD6DE818F307B3306143D7BB13E0"
                    },
                    "PROPS": {
                        "display_filter": "'[{\"operator\":\"include\",\"field\":\"PATH\",\"values\":[\"*\"]},{\"operator\":\"include\",\"field\":\"HOSTNAME\",\"values\":[\"*\"]}]'",
                        "lastupdated_time": "1654547359267",
                        "lastupdated_user": "'Amanpreet Dhindsa'"
                    }
                }
            ]
        }
    }
}`

func TestAlertsGetDetailsMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetDetails should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.GetDetails(alertID)
	assert.Nil(t, err)
}

func TestAlertsGetDetailsOK(t *testing.T) {
	mockResponse := alertDetailsJSON

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
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

	detailsExpected := api.AlertDetailsResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &detailsExpected)

	var detailsActual api.AlertDetailsResponse
	detailsActual, err = c.V2.Alerts.GetDetails(alertID)
	assert.Nil(t, err)
	assert.Equal(t, detailsExpected, detailsActual)
}

func TestAlertsGetDetailsError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
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

	_, err = c.V2.Alerts.GetDetails(alertID)
	assert.NotNil(t, err)
}
