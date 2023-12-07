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

	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/internal/pointer"
	"github.com/lacework/go-sdk/lwseverity"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

var (
	policyURI = "Policies"
	policyID  = "my-policy-1"

	policy = api.Policy{
		Tags: []string{"fhqwhgads"},
	}
	newPolicy = api.NewPolicy{
		PolicyID:      policyID,
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
	}
	policyCreateData = fmt.Sprintf(`{
	"policyId": "%s",
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
}`, policyID)

	updatePolicyMinimal = api.UpdatePolicy{
		PolicyID: policyID,
		Title:    "My New Policy Title",
	}
	policyUpdateData = fmt.Sprintf(`{
	"policyId": "%s",
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
}`, policyID)

	updateBulkPolicy = api.BulkUpdatePolicies{{
		PolicyID: "lacework-global-1",
		Enabled:  ptr.Bool(true),
		Severity: lwseverity.Medium.String(),
	}}

	mockPolicyErrorResponse = `{
	"message": "This is an error message"
}`
	bulkPolicyUpdateResponse = `
{
"data": [
{
"evaluatorId": "Cloudtrail",
"policyId": "lacework-global-1",
"policyType": "Violation",
"queryId": "LW_Global_AWS_Example",
"queryText": "{ source { CloudTrailRawEvents }",
"title": "Test Example",
"enabled": true,
"description": "An example description",
"remediation": "Example remediation",
"severity": "medium",
"limit": 1000,
"evalFrequency": "Hourly",
"alertEnabled": true,
"alertProfile": "LW_CloudTrail_Alerts.Example",
"owner": "Lacework",
"lastUpdateTime": "2023-01-30T14:00:00.000Z",
"lastUpdateUser": "test.user@lacework.net",
"tags": [
"domain:AWS",
"subdomain:Cloudtrail"]
}
]
}
`

	regoPolicyId = "rego-policy-1"
	regoPolicy   = api.NewPolicy{
		PolicyID:      regoPolicyId,
		PolicyType:    "Violation",
		QueryID:       "MyRegoQuery",
		QueryLanguage: pointer.Of("Rego"),
		Title:         "My Rego Policy Title",
		Enabled:       false,
		Description:   "My Policy Description",
		Remediation:   "Check yourself...",
		Severity:      "high",
		EvalFrequency: "Hourly",
		Limit:         1000,
		AlertEnabled:  false,
		AlertProfile:  "LW_CloudTrail_Alerts",
	}
	regoPolicyCreateData = fmt.Sprintf(`{
	"policyId": "%s",
	"policyType": "%s",
	"queryId": "%s",
	"queryLanguage": "%s",
	"title": "%s",
	"enabled": %v,
	"description": "%s",
	"remediation": "%s",
	"severity": "%s",
	"evalFrequency": "%s",
	"limit": %d,
	"alertEnabled": %v,
	"alertProfile": "%s"
}`, regoPolicy.PolicyID, regoPolicy.PolicyType, regoPolicy.QueryID, *regoPolicy.QueryLanguage, regoPolicy.Title,
		regoPolicy.Enabled, regoPolicy.Description, regoPolicy.Remediation, regoPolicy.Severity,
		regoPolicy.EvalFrequency, regoPolicy.Limit, regoPolicy.AlertEnabled, regoPolicy.AlertProfile)
)

func mockPolicyDataResponse(data string) string {
	return `{
	"data": ` + data + `,
	"message": "SUCCESS"
}`
}

func TestPolicyHasTags(t *testing.T) {
	assert.Equal(t, true, policy.HasTag("fhqwhgads"))
	assert.Equal(t, false, policy.HasTag(""))
}

func TestPolicyCreateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
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

func testPolicyCreateOKHelper(t *testing.T, expectedPolicyData string, testPolicy api.NewPolicy) {
	mockResponse := mockPolicyDataResponse(expectedPolicyData)

	fakeServer := lacework.MockServer()
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

	createActual, err := c.V2.Policy.Create(testPolicy)
	assert.Nil(t, err)
	assert.Equal(t, createExpected, createActual)
}

func TestLqlPolicyCreateOK(t *testing.T) {
	testPolicyCreateOKHelper(t, policyCreateData, newPolicy)
}

func TestRegoPolicyCreateOK(t *testing.T) {
	testPolicyCreateOKHelper(t, regoPolicyCreateData, regoPolicy)
}

func TestPolicyCreateError(t *testing.T) {
	fakeServer := lacework.MockServer()
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

func testPolicyGetOKHelper(t *testing.T, expectedPolicyData string, testPolicyId string) {
	mockResponse := mockPolicyDataResponse(expectedPolicyData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", policyURI, testPolicyId),
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
	getActual, err = c.V2.Policy.Get(testPolicyId)
	assert.Nil(t, err)

	assert.Equal(t, getExpected, getActual)
}

func TestLQLPolicyGetOK(t *testing.T) {
	testPolicyGetOKHelper(t, policyCreateData, policyID)
}

func TestRegoPolicyGetOK(t *testing.T) {
	testPolicyGetOKHelper(t, regoPolicyCreateData, regoPolicyId)
}

func TestPolicyGetNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
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

func TestPolicyListMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Policies",
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

	_, err = c.V2.Policy.List()
	assert.Nil(t, err)
}

func TestPolicyListTagsMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Policies/Tags",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListTags should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Policy.ListTags()
	assert.Nil(t, err)
}

func TestPolicyUpdateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
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

	_, err = c.V2.Policy.Update(updatePolicyMinimal)
	assert.Nil(t, err)
}

func TestPolicyUpdateManyMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		policyURI,
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

	_, err = c.V2.Policy.UpdateMany(updateBulkPolicy)
	assert.Nil(t, err)
}

func TestPolicyUpdateManyOK(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		policyURI,
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, bulkPolicyUpdateResponse)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	updateExpected := api.BulkPolicyUpdateResponse{}
	_ = json.Unmarshal([]byte(bulkPolicyUpdateResponse), &updateExpected)

	updateActual, err := c.V2.Policy.UpdateMany(updateBulkPolicy)
	assert.Nil(t, err)
	assert.Equal(t, updateExpected, updateActual)
}

func TestPolicyUpdateBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
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

	_, err = c.V2.Policy.Update(api.UpdatePolicy{})
	assert.Equal(t, "policy ID must be provided", err.Error())
}

func TestPolicyUpdateOK(t *testing.T) {
	// policy create and update data are same-same
	mockResponse := mockPolicyDataResponse(policyUpdateData)

	fakeServer := lacework.MockServer()
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

	updateActual, err := c.V2.Policy.Update(updatePolicyMinimal)
	assert.Nil(t, err)
	assert.Equal(t, updateExpected, updateActual)
}

func TestPolicyUpdateError(t *testing.T) {
	fakeServer := lacework.MockServer()
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

	_, err = c.V2.Policy.Update(api.UpdatePolicy{})
	assert.NotNil(t, err)
}

func TestPolicyDeleteMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
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
