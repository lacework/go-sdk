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
	newPolicy = api.NewPolicy{
		EvaluatorID:   "Cloudtrail",
		PolicyID:      "my-policy-1",
		PolicyType:    "Violation",
		QueryID:       "MyExampleQuery",
		Title:         "My Policy Title",
		Enabled:       false,
		Description:   "My Policy Description",
		Remediation:   "Check yourself...",
		Severity:      "high",
		EvalFrequency: "Hourly",
		AlertEnabled:  false,
		AlertProfile:  "LW_CloudTrail_Alerts",
		PolicyUI:      map[string]string{"domain": "AWS", "subdomain": "Cloudtrail"},
	}
	updatePolicy = api.UpdatePolicy{
		Title: "My New Policy Title",
	}
	policyCreateData = `{
	"evaluatorId": "Cloudtrail",
	"policyId": "my-policy-1",
	"title": "My Policy Title",
	"enabled": false,
	"description": "My Policy Description",
	"remediation": "Check yourself...",
	"severity": "high",
	"evalFrequency": "Hourly",
	"limit": 1000,
	"alertEnabled": false,
	"alertProfile": "LW_CloudTrail_Alerts",
	"policyType": "Violation"
}`
	policyUpdateData = `{
	"evaluatorId": "Cloudtrail",
	"policyId": "my-policy-1",
	"title": "My New Policy Title",
	"enabled": false,
	"description": "My Policy Description",
	"remediation": "Check yourself...",
	"severity": "high",
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

	_, err = c.V2.Policy.Create(api.NewPolicy{})
	assert.Nil(t, err)
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

	createActual, err := c.V2.Policy.Create(newPolicy)
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

	_, err = c.V2.Policy.Create(api.NewPolicy{})
	assert.NotNil(t, err)
}

func TestPolicyGetMethod(t *testing.T) {
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

	_, err = c.V2.Policy.Get(policyID)
	assert.Nil(t, err)
}

func TestPolicyGetOK(t *testing.T) {
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
	getActual, err = c.V2.Policy.Get(policyID)
	assert.Nil(t, err)

	assert.Equal(t, getExpected, getActual)
}

func TestPolicyGetNotFound(t *testing.T) {
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

	_, err = c.V2.Policy.Get("NoSuchPolicy")
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

	_, err = c.V2.Policy.Update(policyID, api.UpdatePolicy{})
	assert.Nil(t, err)
}

func TestPolicyUpdateBadInput(t *testing.T) {
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

	_, err = c.V2.Policy.Update("", api.UpdatePolicy{})
	assert.Equal(t, "policy ID must be provided", err.Error())
}

func TestPolicyUpdateOK(t *testing.T) {
	// policy create and update data are same-same
	mockResponse := mockPolicyDataResponse(policyUpdateData)

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

	updateActual, err := c.V2.Policy.Update(policyID, updatePolicy)
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

	_, err = c.V2.Policy.Update(policyID, api.UpdatePolicy{})
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

	_, err = c.V2.Policy.Delete(policyID)
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

	_, err = c.V2.Policy.Delete("")
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

	_, err = c.V2.Policy.Delete(policyID)
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

	_, err = c.V2.Policy.Delete("no-such-policy")
	assert.NotNil(t, err)
}
