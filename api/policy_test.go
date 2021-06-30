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
	policyURI = "Policies"
	policyID  = "my-policy-1"
	policyStr = `{
	"policyId": "my-policy-1",
	"title": "My Policy Title",
	"enabled": false,
	"queryId": "MyExampleQuery",
	"severity": "high",
	"description": "My Policy Description",
	"remediation": "Check yourself..."
}`
	policyCreateData = `{
	"policyId": "my-policy-1",
	"title": "My Policy Title",
	"enabled": false,
	"description": "My Policy Description",
	"remediation": "Check yourself...",
	"severity": "2",
	"evalFrequency": "Hourly",
	"limit": 1000,
	"alertEnabled": false,
	"alertProfile": "LW_CloudTrail_Alerts",
	"policyType": "Violation"
}`
	mockPolicyErrorResponse = `{
	"message": "This is an error message"
}`
)

func mockPolicyDataResponse(data string) string {
	return `{
	"data": ` + data + `,
	"message": "SUCCESS"
}`
}

func TestPolicyCreateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		policyURI,
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

func TestPolicyCreateBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		policyURI,
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
	assert.Equal(t, "policy must be valid JSON: unexpected end of JSON input", err.Error())
}

func TestPolicyCreateOK(t *testing.T) {
	mockResponse := mockPolicyDataResponse(policyCreateData)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		policyURI,
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
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		policyURI,
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, mockPolicyErrorResponse, http.StatusBadRequest)
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
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
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
	mockResponse := mockPolicyDataResponse(policyCreateData)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
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

func TestPolicyGetByIDNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, mockPolicyErrorResponse, http.StatusBadRequest)
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

func TestPolicyUpdateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PATCH", r.Method, "Update should be a PATCH method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.Policy.Update(policyID, policyStr)
	assert.Nil(t, err)
}

func TestPolicyUpdateNoPolicyID(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
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

	_, err = c.Policy.Update("", "{}")
	assert.Equal(t, "policy ID must be provided", err.Error())
}

func TestPolicyUpdatePolicyIDFromPayload(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
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

	_, err = c.Policy.Update("",
		fmt.Sprintf(`{"policy_id": "%s"}`, policyID))
	assert.Nil(t, err)
}

func TestPolicyUpdateOK(t *testing.T) {
	// policy create and update data are same-same
	mockResponse := mockPolicyDataResponse(policyCreateData)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
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

	updateExpected := api.PolicyResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &updateExpected)

	updateActual, err := c.Policy.Update(policyID, `{"severity": "high"}`)
	assert.Nil(t, err)
	assert.Equal(t, updateExpected, updateActual)
}

func TestPolicyUpdateError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, mockPolicyErrorResponse, http.StatusBadRequest)
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

func TestPolicyDeleteMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method, "Delete should be a DELETE method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.Policy.Delete(policyID)
	assert.Nil(t, err)
}

func TestPolicyDeleteBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
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

	_, err = c.Policy.Delete("")
	assert.Equal(t, "policy ID must be provided", err.Error())
}

func TestPolicyDeleteOK(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
		func(w http.ResponseWriter, r *http.Request) {
			// send the headers with a 204 response code.
			w.WriteHeader(http.StatusNoContent)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.Policy.Delete(policyID)
	assert.Nil(t, err)
}

func TestPolicyDeleteError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, policyID),
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, mockPolicyErrorResponse, http.StatusBadRequest)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.Policy.Delete("no-such-policy")
	assert.NotNil(t, err)
}
