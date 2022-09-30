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

var alertTimelineJSON = `{
    "data": [
        {
            "id": 129,
            "alertId": 168155,
            "entryType": "Comment",
            "entryAuthorType": "UserUpdate",
            "message": {
                "value": "this is fine"
            },
            "externalTime": "0",
            "user": {
                "userGuid": "DEV7CE4D_76A94FE2C8D3EEEEEEEEED62EA6745D4A42DF777",
                "username": "david.hazekamp@lacework.net"
            },
            "updateContext": {}
        }
    ]
}`

func TestAlertsGetTimelineMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetTimeline should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.GetTimeline(alertID)
	assert.Nil(t, err)
}

func TestAlertsGetTimelineOK(t *testing.T) {
	mockResponse := alertInvestigationJSON

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

	timelineExpected := api.AlertTimelineResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &timelineExpected)

	var timelineActual api.AlertTimelineResponse
	timelineActual, err = c.V2.Alerts.GetTimeline(alertID)
	assert.Nil(t, err)
	assert.Equal(t, timelineExpected, timelineActual)
}

func TestAlertsGetTimelineError(t *testing.T) {
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

	_, err = c.V2.Alerts.GetTimeline(alertID)
	assert.NotNil(t, err)
}
