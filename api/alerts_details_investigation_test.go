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

var alertInvestigationJSON = `{
    "data": [
        {
            "question": "Has a new user been involved in the event in the last 60 days?",
            "answer": "No"
        },
        {
            "question": "Have the users involved in the event authenticated without MFA in the last 60 days?",
            "answer": "No"
        },
        {
            "question": "Have any of the users involved in the event used the Root account in the last 60 days?",
            "answer": "No"
        }
    ]
}`

func TestAlertsGetInvestigationMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetInvestigation should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.GetInvestigation(alertID)
	assert.Nil(t, err)
}

func TestAlertsGetInvestigationOK(t *testing.T) {
	mockResponse := alertInvestigationJSON

	fakeServer := lacework.MockServer()
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

	investigationExpected := api.AlertInvestigationResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &investigationExpected)

	var investigationActual api.AlertInvestigationResponse
	investigationActual, err = c.V2.Alerts.GetInvestigation(alertID)
	assert.Nil(t, err)
	assert.Equal(t, investigationExpected, investigationActual)
}

func TestAlertsGetInvestigationError(t *testing.T) {
	fakeServer := lacework.MockServer()
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

	_, err = c.V2.Alerts.GetInvestigation(alertID)
	assert.NotNil(t, err)
}
