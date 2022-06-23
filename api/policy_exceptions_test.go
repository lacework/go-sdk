//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestPolicyExceptionsGet(t *testing.T) {
	var (
		exceptionID     = intgguid.New()
		policyID        = "my-policy-1"
		apiPath         = fmt.Sprintf("Exceptions/%s", exceptionID)
		policyException = singleMockPolicyException(exceptionID)
		fakeServer      = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generatePolicyExceptionResponse(policyException))
			}
		},
	)

	fakeServer.MockAPI("Exceptions/UNKNOWN_EXCEPTION_ID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when policy exception exists", func(t *testing.T) {
		var response api.PolicyExceptionResponse
		err := c.V2.Policy.Exceptions.Get(policyID, exceptionID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, exceptionID, response.Data.ExceptionID)
			assert.Equal(t, "exception description", response.Data.Description)
			assert.Equal(t, "accountIds", response.Data.Constraints[0].FieldKey)
			assert.Equal(t, "*", response.Data.Constraints[0].FieldValues[0])
			assert.Equal(t, "darren.murray@lacework.net", response.Data.LastUpdateUser)
			assert.Equal(t, "2022-06-17T12:50:43.591Z", response.Data.LastUpdateTime)
		}
	})

	t.Run("when policy exception does NOT exist", func(t *testing.T) {
		var response api.PolicyExceptionResponse
		err := c.V2.Policy.Exceptions.Get("UNKNOWN_POLICY_ID", "UNKNOWN_EXCEPTION_ID", response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/Exceptions/UNKNOWN_EXCEPTION_ID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestPolicyExceptionsDelete(t *testing.T) {
	var (
		exceptionID     = intgguid.New()
		policyID        = "my-policy-1"
		apiPath         = fmt.Sprintf("Exceptions/%s", exceptionID)
		policyException = singleMockPolicyException(exceptionID)
		getResponse     = generatePolicyExceptionResponse(policyException)
		fakeServer      = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if getResponse != "" {
				switch r.Method {
				case "GET":
					fmt.Fprintf(w, getResponse)
				case "DELETE":
					// once deleted, empty the getResponse so that
					// further GET requests return 404s
					getResponse = ""
				}
			} else {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	fakeServer.MockAPI("Exceptions/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify policy exception exists", func(t *testing.T) {
		var response api.PolicyExceptionResponse
		err := c.V2.Policy.Exceptions.Get(policyID, exceptionID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, "exception description", response.Data.Description)
			assert.Equal(t, "accountIds", response.Data.Constraints[0].FieldKey)
			assert.Equal(t, "*", response.Data.Constraints[0].FieldValues[0])
			assert.Equal(t, "darren.murray@lacework.net", response.Data.LastUpdateUser)
			assert.Equal(t, "2022-06-17T12:50:43.591Z", response.Data.LastUpdateTime)
		}
	})

	t.Run("when policy exception has been deleted", func(t *testing.T) {
		err := c.V2.Policy.Exceptions.Delete(policyID, exceptionID)
		assert.Nil(t, err)

		var response api.PolicyExceptionResponse
		err = c.V2.Policy.Exceptions.Get(policyID, exceptionID, &response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/Exceptions/")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestPolicyExceptionsList(t *testing.T) {
	var (
		policyID         = "my-policy-1"
		allGUIDs         []string
		policyExceptions = generateGuids(&allGUIDs, 3)
		expectedLen      = len(allGUIDs)
		fakeServer       = lacework.MockServer()
	)

	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Exceptions",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			policyExceptions := []string{
				generatePolicyExceptions(policyExceptions),
			}
			fmt.Fprintf(w,
				generatePolicyExceptionsResponse(
					strings.Join(policyExceptions, ", "),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Policy.Exceptions.List(policyID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.ExceptionID)
	}
}

func TestPolicyExceptionUpdate(t *testing.T) {
	var (
		exceptionID = intgguid.New()
		apiPath     = fmt.Sprintf("Exceptions/%s", exceptionID)
		fakeServer  = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Update() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "exception description", "policy description is missing")
			assert.Contains(t, body, "accountIds", "policy field key is missing")
			assert.Contains(t, body, "*", "missing field values")
		}

		fmt.Fprintf(w, generatePolicyExceptionResponse(singleMockPolicyException(exceptionID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	policyException := api.PolicyException{
		Description: "exception description",
		Constraints: []api.PolicyExceptionConstraint{{FieldKey: "accountIds", FieldValues: []string{"*"}}},
	}
	assert.Equal(t, "exception description", policyException.Description, "policy exception description mismatch")
	assert.Equal(t, "accountIds", policyException.Constraints[0].FieldKey, "policy exception field key mismatch")
	assert.Equal(t, "*", policyException.Constraints[0].FieldValues[0], "policy exception field values mismatch")
	policyException.ExceptionID = exceptionID

	response, err := c.V2.Policy.Exceptions.Update(policyID, policyException)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, exceptionID, response.Data.ExceptionID)
		assert.Equal(t, "exception description", response.Data.Description)
		assert.Equal(t, response.Data.Constraints[0].FieldKey, "accountIds")
		assert.Equal(t, response.Data.Constraints[0].FieldValues[0], "*")
	}
}

func generatePolicyExceptions(guids []string) string {
	policyExceptions := make([]string, len(guids))
	for i, guid := range guids {
		policyExceptions[i] = singleMockPolicyException(guid)
	}
	return strings.Join(policyExceptions, ", ")
}

func generatePolicyExceptionsResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generatePolicyExceptionResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleMockPolicyException(id string) string {
	return fmt.Sprintf(`
{
        "exceptionId": %q,
        "description": "exception description",
        "constraints": [
            {
                "fieldKey": "accountIds",
                "fieldValues": [
                    "*"
                ]
            }
        ],
        "lastUpdateTime": "2022-06-17T12:50:43.591Z",
        "lastUpdateUser": "darren.murray@lacework.net"
    }
	`, id)
}
