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

var (
	policyID  = "my-policy-1"
	policyStr = `{
	"policy_id": "my-policy-1",
	"title": "My Policy Title",
	"enabled": false,
	"lql_id": "MyExampleQuery",
	"severity": "high",
	"description": "My Policy Description",
	"remediation": "Check yourself..."
}`
	policyCreateData = `[
	{
		"policy_id": "my-policy-1",
		"title": "My Policy Title",
		"enabled": false,
		"description": "My Policy Description",
		"remediation": "Check yourself...",
		"severity": "2",
		"eval_frequency": "Hourly",
		"limit": 1000,
		"alert_enabled": false,
		"alert_profile": "LW_CloudTrail_Alerts",
		"policy_type": "Violation"
	}
]`
	policyAlreadyExistsError = mockPolicyDataResponse(
		`{ "message": "{\"error\":\"Error: Cannot create rule my-policy-1 because it already exists in database.\"}" }`,
		"false",
	)
	policyUnableToLocateError = mockPolicyDataResponse(
		`{ "message": "{\"error\":\"Error: Unable to locate policy foo, please double check the policy exists and has not already been deleted.\"}" }`,
		"false",
	)
)

func mockPolicyDataResponse(data string, ok string) string {
	return `{
	"data": ` + data + `,
	"ok": ` + ok + `,
	"message": "SUCCESS"
}`
}

func TestPolicyCreateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiPolicy,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Create should be a POST method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.Policy.Create(policyStr)
	assert.Nil(t, err)
}

func TestPolicyCreateBadJSON(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiPolicy,
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.Policy.Create("")
	assert.Equal(t, "unexpected end of JSON input", err.Error())
}

func TestPolicyCreateOK(t *testing.T) {
	mockResponse := mockPolicyDataResponse(policyCreateData, "true")
	fmt.Println(mockResponse)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiPolicy,
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

	createExpected := api.PolicyResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &createExpected)

	createActual, err := c.Policy.Create("{}")
	assert.Nil(t, err)
	assert.Equal(t, createExpected, createActual)
}

func TestPolicyCreateError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiPolicy,
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, policyAlreadyExistsError, http.StatusBadRequest)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.Policy.Create(policyStr)
	assert.NotNil(t, err)
}

func TestPolicyGetByIDMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiPolicy,
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

	_, err = c.Policy.GetByID(policyID)
	assert.Nil(t, err)
}

func TestPolicyGetByIDOK(t *testing.T) {
	mockResponse := mockPolicyDataResponse(policyCreateData, "true")

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiPolicy,
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

	getExpected := api.PolicyResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &getExpected)

	var getActual api.PolicyResponse
	getActual, err = c.Policy.GetByID(policyID)
	assert.Nil(t, err)

	assert.Equal(t, getExpected, getActual)
}

func TestPolicyByIDNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, policyUnableToLocateError, http.StatusBadRequest)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.Policy.GetByID("NoSuchPolicy")
	assert.NotNil(t, err)
}
