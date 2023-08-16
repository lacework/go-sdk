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

var alertsSearchRequest api.SearchFilter

func TestAlertsSearchMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Alerts/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search should be a POST method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.Search(alertsSearchRequest)
	assert.Nil(t, err)
}

func TestAlertsSearchOK(t *testing.T) {
	mockResponse := alertsJSON

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Alerts/search",
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
	listActual, err = c.V2.Alerts.Search(alertsSearchRequest)
	assert.Nil(t, err)
	assert.Equal(t, listExpected, listActual)
}

func TestAlertsSearchAllOK(t *testing.T) {
	fakeServer := lacework.MockServer()

	nextPage := fmt.Sprintf(
		"%s/api/v2/Alerts/search/nextPage",
		fakeServer.URL(),
	)
	mockResponse := fmt.Sprintf(alertsPage1JSON, nextPage)

	fakeServer.MockAPI(
		"Alerts/search",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, mockResponse)
		},
	)
	fakeServer.MockAPI(
		"Alerts/search/nextPage",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, alertsPage2JSON)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	var listActual api.AlertsResponse
	listActual, err = c.V2.Alerts.SearchAll(alertsSearchRequest)
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(listActual.Data))
}

func TestAlertsSearchError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Alerts/search",
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

	_, err = c.V2.Alerts.Search(alertsSearchRequest)
	assert.NotNil(t, err)
}
