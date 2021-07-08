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
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

var (
	queryEvaluator = "Cloudtrail"
	queryID        = "my_lql"
	newQueryText   = `my_lql { source { CloudTrailRawEvents } return { INSERT_ID } }`
	newQuery       = api.NewQuery{
		EvaluatorID: queryEvaluator,
		QueryID:     queryID,
		QueryText:   newQueryText,
	}
	newQueryJSON = fmt.Sprintf(`{
	"evaluatorId": "%s",
	"queryId": "%s",
	"queryText": "%s"
}`, queryEvaluator, queryID, newQueryText)
	queryRunData = `[
	{
		"INSERT_ID": "35308423"
	}
]`
	lqlErrorReponse = `{ "message": "This is an error message" }`
)

func mockQueryDataResponse(data string) string {
	return `{
	"data": ` + data + `
}`
}

func TestQueryCreateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
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

func TestQueryCreateOK(t *testing.T) {
	mockResponse := mockQueryDataResponse(newQueryJSON)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
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
	createActual, err = c.V2.Query.Create(newQuery)
	assert.Nil(t, err)

	assert.Equal(t, createExpected, createActual)
}

func TestQueryCreateError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
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
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Get should be a GET method")
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

func TestQueryGetQueryByIDOK(t *testing.T) {
	mockResponse := mockQueryDataResponse(newQueryJSON)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
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
	getActual, err = c.V2.Query.GetByID(queryID)
	assert.Nil(t, err)

	assert.Equal(t, getExpected, getActual)
}

func TestQueryGetByIDNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
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

	_, err = c.V2.Query.GetByID("NoSuchQuery")
	assert.NotNil(t, err)
}

func TestQueryExecuteMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/lql/query",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Run should be a POST method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.Execute(newQueryText, time.Unix(0, 0), time.Unix(1, 0))
	assert.Nil(t, err)
}

func TestQueryExecuteBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/lql/query",
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

	_, err = c.V2.Query.Execute("", time.Unix(0, 0), time.Unix(1, 0))
	assert.Equal(t, "query text must be provided", err.Error())
}

func TestQueryExecuteOK(t *testing.T) {
	mockResponse := mockQueryDataResponse(queryRunData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/lql/query",
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

	var runExpected map[string]interface{}
	_ = json.Unmarshal([]byte(mockResponse), &runExpected)

	var runActual map[string]interface{}
	runActual, err = c.V2.Query.Execute(newQueryText, time.Unix(0, 0), time.Unix(1, 0))
	assert.Nil(t, err)

	assert.Equal(t, runExpected, runActual)
}

func TestQueryExecuteError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/lql/query",
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

	_, err = c.V2.Query.Execute(newQueryText, time.Unix(0, 0), time.Unix(1, 0))
	assert.NotNil(t, err)
}
