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
	executeQueryArguments = []api.ExecuteQueryArgument{
		api.ExecuteQueryArgument{
			Name:  "StartTimeRange",
			Value: "2021-07-11T00:00:00.000Z",
		},
		api.ExecuteQueryArgument{
			Name:  "EndTimeRange",
			Value: "2021-07-12T00:00:00.000Z",
		},
	}
	executeQuery = api.ExecuteQueryRequest{
		Query: api.ExecuteQuery{
			QueryText: newQueryText,
		},
		Arguments: executeQueryArguments,
	}
	executeQueryByID = api.ExecuteQueryByIDRequest{
		QueryID:   queryID,
		Arguments: executeQueryArguments,
	}
	executeQueryData = `[
	{
		"INSERT_ID": "35308423"
	}
]`
)

func TestQueryExecuteMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/execute",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Execute should be a POST method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.Execute(executeQuery)
	assert.Nil(t, err)
}

func TestQueryExecuteOK(t *testing.T) {
	mockResponse := mockQueryDataResponse(executeQueryData)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/execute",
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
	runActual, err = c.V2.Query.Execute(executeQuery)
	assert.Nil(t, err)

	assert.Equal(t, runExpected, runActual)
}

func TestQueryExecuteError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/execute",
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

	_, err = c.V2.Query.Execute(api.ExecuteQueryRequest{})
	assert.NotNil(t, err)
}

func TestQueryExecuteByIDMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Queries/%s/execute", queryID),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Execute should be a POST method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.ExecuteByID(executeQueryByID)
	assert.Nil(t, err)
}

func TestQueryExecuteByIDOK(t *testing.T) {
	mockResponse := mockQueryDataResponse(executeQueryData)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Queries/%s/execute", queryID),
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
	runActual, err = c.V2.Query.ExecuteByID(executeQueryByID)
	assert.Nil(t, err)

	assert.Equal(t, runExpected, runActual)
}

func TestQueryExecuteByIDError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Queries/%s/execute", queryID),
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

	_, err = c.V2.Query.ExecuteByID(api.ExecuteQueryByIDRequest{})
	assert.NotNil(t, err)
}
