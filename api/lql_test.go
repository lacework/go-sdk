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
	lqlQueryID    string = "my_lql"
	lqlQueryStr   string = "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
	lqlCreateData string = `[
	{
		"lql_id": "my_lql",
		"query_text": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
	}
]`
	lqlRunData string = `[
	{
		"INSERT_ID": "35308423"
	}
]`
	lqlErrorReponse string = mockLQLDataResponse(
		`{ "message": "Error Serving Request" }`,
	)
	lqlUnableResponse string = mockLQLMessageResponse(
		`"message": "{\"error\":\"Error: Unable to locate lql query NoSuchQuery, please double check the query exists and has not already been updated.\"}"`,
		"false",
	)
	lqlQueryTypeTests []TestLQLQuery = []TestLQLQuery{
		TestLQLQuery{
			Name: "empty-blob",
			Input: &api.LQLQuery{
				QueryBlob: ``,
			},
			Output: api.LQLQueryTranslateError,
			Expected: &api.LQLQuery{
				QueryText: ``,
				QueryBlob: ``,
			},
		},
		TestLQLQuery{
			Name: "junk-blob",
			Input: &api.LQLQuery{
				QueryBlob: `this is junk`,
			},
			Output: api.LQLQueryTranslateError,
			Expected: &api.LQLQuery{
				QueryText: ``,
				QueryBlob: `this is junk`,
			},
		},
		TestLQLQuery{
			Name: "json-blob",
			Input: &api.LQLQuery{
				QueryBlob: `{
	"START_TIME_RANGE": "678910",
	"END_TIME_RANGE": "111213141516",
	"QUERY_TEXT": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
			Output: nil,
			Expected: &api.LQLQuery{
				StartTimeRange: "678910",
				EndTimeRange:   "111213141516",
				QueryText:      "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
				QueryBlob: `{
	"START_TIME_RANGE": "678910",
	"END_TIME_RANGE": "111213141516",
	"QUERY_TEXT": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
		},
		TestLQLQuery{
			Name: "json-blob-lower",
			Input: &api.LQLQuery{
				QueryBlob: `{
	"start_time_range": "678910",
	"end_time_range": "111213141516",
	"query_text": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
			Output: nil,
			Expected: &api.LQLQuery{
				StartTimeRange: "678910",
				EndTimeRange:   "111213141516",
				QueryText:      "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
				QueryBlob: `{
	"start_time_range": "678910",
	"end_time_range": "111213141516",
	"query_text": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
		},
		TestLQLQuery{
			Name: "lql-blob",
			Input: &api.LQLQuery{
				QueryBlob: "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
			},
			Output: nil,
			Expected: &api.LQLQuery{
				QueryText: "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
				QueryBlob: "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
			},
		},
		TestLQLQuery{
			Name: "overwrite-blob",
			Input: &api.LQLQuery{
				StartTimeRange: "should not overwrite",
				EndTimeRange:   "should not overwrite",
				QueryText:      "should not overwrite",
				QueryBlob: `{
	"START_TIME_RANGE": "678910",
	"END_TIME_RANGE": "111213141516",
	"QUERY_TEXT": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
			Output: nil,
			Expected: &api.LQLQuery{
				StartTimeRange: "should not overwrite",
				EndTimeRange:   "should not overwrite",
				QueryText:      "should not overwrite",
				QueryBlob: `{
	"START_TIME_RANGE": "678910",
	"END_TIME_RANGE": "111213141516",
	"QUERY_TEXT": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
		},
	}
)

type TestLQLQuery struct {
	Name     string
	Input    *api.LQLQuery
	Output   interface{}
	Expected *api.LQLQuery
}

func TestLQLQueryTranslate(t *testing.T) {
	for _, lqlQueryTest := range lqlQueryTypeTests {
		t.Run(lqlQueryTest.Name, func(t *testing.T) {
			if err := lqlQueryTest.Input.Translate(); err == nil {
				assert.Equal(t, lqlQueryTest.Output, err)
			} else {
				assert.Equal(t, lqlQueryTest.Output, err.Error())
			}
			assert.Equal(t, lqlQueryTest.Expected, lqlQueryTest.Input)
		})
	}
}

func mockLQLDataResponse(data string) string {
	return `{
	"data": ` + data + `,
	"ok": true,
	"message": "SUCCESS"
}`
}

func mockLQLMessageResponse(message string, ok string) string {
	return `{
	"ok": ` + ok + `,
	"message": {` + message + `
	}
}`
}

func TestCreateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
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

	_, err = c.LQL.CreateQuery(lqlQueryStr)
	assert.Nil(t, err)
}

func TestCreateBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
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

	_, err = c.LQL.CreateQuery("")
	assert.Equal(t, api.LQLQueryTranslateError, err.Error())
}

func TestCreateOK(t *testing.T) {
	mockResponse := mockLQLDataResponse(lqlCreateData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
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

	createExpected := api.LQLQueryResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &createExpected)

	var createActual api.LQLQueryResponse
	createActual, err = c.LQL.CreateQuery(lqlQueryStr)
	assert.Nil(t, err)

	assert.Equal(t, createExpected, createActual)
}

func TestCreateError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
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

	_, err = c.LQL.CreateQuery(lqlQueryStr)
	assert.NotNil(t, err)
}

func TestGetQueriesMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Get should be a GET method")
			assert.NotSubset(
				t,
				[]byte(r.RequestURI),
				[]byte(api.ApiLQL+"?LQL_ID"),
				"GetQueries should not specify LQL_ID argument",
			)
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.LQL.GetQueries()
	assert.Nil(t, err)
}

func TestGetQueryByIDOK(t *testing.T) {
	mockResponse := mockLQLDataResponse(lqlCreateData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
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

	createExpected := api.LQLQueryResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &createExpected)

	var createActual api.LQLQueryResponse
	createActual, err = c.LQL.GetQueryByID(lqlQueryStr)
	assert.Nil(t, err)

	assert.Equal(t, createExpected, createActual)
}

func TestGetQueryByIDNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, lqlUnableResponse, http.StatusBadRequest)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.LQL.GetQueryByID("NoSuchQuery")
	assert.NotNil(t, err)
}

func TestRunQueryMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQLQuery,
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

	_, err = c.LQL.RunQuery(lqlQueryStr, "", "")
	assert.Nil(t, err)
}

func TestRunQueryBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQLQuery,
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

	_, err = c.LQL.RunQuery("", "", "")
	assert.Equal(t, api.LQLQueryTranslateError, err.Error())
}

func TestRunQueryOK(t *testing.T) {
	mockResponse := mockLQLDataResponse(lqlRunData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQLQuery,
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
	runActual, err = c.LQL.RunQuery(lqlQueryStr, "", "")
	assert.Nil(t, err)

	assert.Equal(t, runExpected, runActual)
}

func TestRunQueryError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
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

	_, err = c.LQL.CreateQuery(lqlQueryStr)
	assert.NotNil(t, err)
}
