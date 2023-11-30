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
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/lacework/go-sdk/internal/pointer"
)

var (
	queryID      = "my_lql"
	newQueryText = `my_lql { source { CloudTrailRawEvents } return { INSERT_ID } }`
	newQuery     = api.NewQuery{
		QueryID:   queryID,
		QueryText: newQueryText,
	}
	newQueryJSON = fmt.Sprintf(`{
	"queryId": "%s",
	"queryText": "%s"
}`, queryID, newQueryText)
	newQueryYAML = fmt.Sprintf(`---
queryId: %s
queryText: %s`, newQuery.QueryID, newQuery.QueryText)
	lqlErrorReponse = `{ "message": "This is an error message" }`

	regoQueryID      = "my_rego"
	newRegoQueryText = "package clitest\n" +
		"import future.keywords\n" +
		"import data.lacework\n" +
		"source := lacework.aws.cfg.list(\"s3\", \"list-buckets\")\n" +
		"assess := assessment.violation(input, \"just because\")"
	newRegoQuery = api.NewQuery{
		QueryID:       regoQueryID,
		QueryLanguage: pointer.Of("Rego"),
		QueryText:     newRegoQueryText,
	}
	newRegoQueryJSON = fmt.Sprintf(`{
	"queryId": "%s",
	"queryLanguage": "Rego",
	"queryText": %#v
}`, queryID, newRegoQueryText)
)

func mockQueryDataResponse(data string) string {
	return `{
	"data": ` + data + `
}`
}

type parseNewQueryTest struct {
	Name     string
	Input    string
	Expected api.NewQuery
	Error    error
}

var parseNewQueryTests = []parseNewQueryTest{
	parseNewQueryTest{
		Name:     "empty-blob",
		Input:    "",
		Expected: api.NewQuery{},
		Error:    errors.New("unable to parse query"),
	},
	parseNewQueryTest{
		Name:     "junk-blob",
		Input:    "this is junk",
		Expected: api.NewQuery{},
		Error:    errors.New("unable to parse query"),
	},
	parseNewQueryTest{
		Name:     "partial-blob",
		Input:    "{",
		Expected: api.NewQuery{},
		Error:    errors.New("unable to parse query"),
	},
	parseNewQueryTest{
		Name:     "json-blob",
		Input:    newQueryJSON,
		Expected: newQuery,
		Error:    nil,
	},
	parseNewQueryTest{
		Name:     "yaml-blob",
		Input:    newQueryYAML,
		Expected: newQuery,
		Error:    nil,
	},
}

func TestParseNewQuery(t *testing.T) {
	for _, pnqt := range parseNewQueryTests {
		t.Run(pnqt.Name, func(t *testing.T) {
			actual, err := api.ParseNewQuery(pnqt.Input)
			if pnqt.Error == nil {
				assert.Equal(t, pnqt.Error, err)
			} else {
				assert.Equal(t, pnqt.Error.Error(), err.Error())
			}
			assert.Equal(t, pnqt.Expected, actual)
		})
	}
}

func TestQueryCreateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries",
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

	_, err = c.V2.Query.Create(newQuery)
	assert.Nil(t, err)
}

func createQueryOKTestHelper(t *testing.T, expectedResponseData string, testQuery api.NewQuery) {
	mockResponse := mockQueryDataResponse(expectedResponseData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries",
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

	createExpected := api.QueryResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &createExpected)

	var createActual api.QueryResponse
	createActual, err = c.V2.Query.Create(testQuery)
	assert.Nil(t, err)

	assert.Equal(t, createExpected, createActual)

	if strings.Contains(expectedResponseData, "queryLanguage") {
		assert.Equal(t, "Rego", *createActual.Data.QueryLanguage)
	} else {
		assert.Nil(t, createActual.Data.QueryLanguage)
	}
}

func TestLQLQueryCreateOK(t *testing.T) {
	// queryLanguage is not available
	createQueryOKTestHelper(t, newQueryJSON, newQuery)
}

func TestRegoQueryCreateOK(t *testing.T) {
	// queryLanguage is available
	createQueryOKTestHelper(t, newRegoQueryJSON, newRegoQuery)
}

func TestQueryCreateError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries",
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

	_, err = c.V2.Query.Create(newQuery)
	assert.NotNil(t, err)
}

func TestQueryListMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries",
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

	_, err = c.V2.Query.List()
	assert.Nil(t, err)
}

func getQueryByIDTestHelper(t *testing.T, expectedResponseData string, queryId string) {
	mockResponse := mockQueryDataResponse(expectedResponseData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries/my_lql",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Get should be a GET method")
			fmt.Fprint(w, mockResponse)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	getExpected := api.QueryResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &getExpected)

	var getActual api.QueryResponse
	getActual, err = c.V2.Query.Get(queryID)
	assert.Nil(t, err)

	assert.Equal(t, getExpected, getActual)

	if strings.Contains(expectedResponseData, "queryLanguage") {
		assert.Equal(t, "Rego", *getActual.Data.QueryLanguage)
	} else {
		assert.Nil(t, getActual.Data.QueryLanguage)
	}
}

func TestLQLQueryGetQueryByIDOK(t *testing.T) {
	// queryLanguage is not available
	getQueryByIDTestHelper(t, newQueryJSON, queryID)
}

func TestRegoQueryGetQueryByIDOK(t *testing.T) {
	// queryLanguage is available
	getQueryByIDTestHelper(t, newRegoQueryJSON, regoQueryID)
}

func TestQueryGetNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, lqlErrorReponse, http.StatusBadRequest)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.Get("NoSuchQuery")
	assert.NotNil(t, err)
}
