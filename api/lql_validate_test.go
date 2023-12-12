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

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/lacework/go-sdk/internal/pointer"
)

var (
	validateQuery = api.ValidateQuery{
		QueryText: newQueryText,
	}
)

func TestQueryValidateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries/validate",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Compile should be a POST method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.Validate(api.ValidateQuery{})
	assert.Nil(t, err)
}

func testQueryValidateOKHelper(t *testing.T, expectedResponseData string, testQuery api.ValidateQuery) {
	mockResponse := mockQueryDataResponse(expectedResponseData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries/validate",
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

	validateExpected := api.QueryResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &validateExpected)

	var validateActual api.QueryResponse
	validateActual, err = c.V2.Query.Validate(testQuery)
	assert.Nil(t, err)
	assert.Equal(t, validateExpected, validateActual)
}

func TestLQLQueryValidateOK(t *testing.T) {
	testQueryValidateOKHelper(t, newQueryJSON, validateQuery)
}

func TestRegoQueryValidateOK(t *testing.T) {
	validateRegoQuery := api.ValidateQuery{
		QueryText:     newRegoQueryText,
		QueryLanguage: pointer.Of("Rego"),
	}
	testQueryValidateOKHelper(t, newRegoQueryJSON, validateRegoQuery)
}

func TestQueryValidateError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries/validate",
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

	_, err = c.V2.Query.Validate(api.ValidateQuery{})
	assert.NotNil(t, err)
}
