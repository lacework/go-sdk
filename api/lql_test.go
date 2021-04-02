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
	lqlTranslateTimeTests []LQLTranslateTimeTest = []LQLTranslateTimeTest{
		LQLTranslateTimeTest{
			Name:       "valid-rfc-utc",
			Input:      "2021-03-31T00:00:00Z",
			ReturnTime: "2021-03-31T00:00:00Z",
			ReturnErr:  nil,
		},
		LQLTranslateTimeTest{
			Name:       "valid-rfc-central",
			Input:      "2021-03-31T00:00:00-05:00",
			ReturnTime: "2021-03-31T05:00:00Z",
			ReturnErr:  nil,
		},
		LQLTranslateTimeTest{
			Name:       "valid-milli",
			Input:      "1617230464000",
			ReturnTime: "2021-03-31T22:41:04Z",
			ReturnErr:  nil,
		},
		LQLTranslateTimeTest{
			Name:       "empty",
			Input:      "",
			ReturnTime: "",
			ReturnErr:  nil,
		},
		LQLTranslateTimeTest{
			Name:       "invalid",
			Input:      "jweaver",
			ReturnTime: "",
			ReturnErr:  "unable to parse time (jweaver)",
		},
	}
	lqlValidateRangeTests []LQLValidateRangeTest = []LQLValidateRangeTest{
		LQLValidateRangeTest{
			Name: "ok",
			Input: api.LQLQuery{
				StartTimeRange: "0",
				EndTimeRange:   "1",
				QueryBlob:      lqlQueryStr,
			},
			AllowEmpty: false,
			Return:     nil,
		},
		LQLValidateRangeTest{
			Name: "empty-start-allowed",
			Input: api.LQLQuery{
				StartTimeRange: "",
				EndTimeRange:   "1",
				QueryBlob:      lqlQueryStr,
			},
			AllowEmpty: true,
			Return:     nil,
		},
		LQLValidateRangeTest{
			Name: "empty-start-disallowed",
			Input: api.LQLQuery{
				StartTimeRange: "",
				EndTimeRange:   "1",
				QueryBlob:      lqlQueryStr,
			},
			AllowEmpty: false,
			Return:     "start time must not be empty",
		},
		LQLValidateRangeTest{
			Name: "empty-end-allowed",
			Input: api.LQLQuery{
				StartTimeRange: "0",
				EndTimeRange:   "",
				QueryBlob:      lqlQueryStr,
			},
			AllowEmpty: true,
			Return:     nil,
		},
		LQLValidateRangeTest{
			Name: "empty-end-disallowed",
			Input: api.LQLQuery{
				StartTimeRange: "0",
				EndTimeRange:   "",
				QueryBlob:      lqlQueryStr,
			},
			AllowEmpty: false,
			Return:     "end time must not be empty",
		},
		LQLValidateRangeTest{
			Name: "empty-both-allowed",
			Input: api.LQLQuery{
				StartTimeRange: "",
				EndTimeRange:   "",
				QueryBlob:      lqlQueryStr,
			},
			AllowEmpty: true,
			Return:     nil,
		},
		LQLValidateRangeTest{
			Name: "start-after-end",
			Input: api.LQLQuery{
				StartTimeRange: "1717333947000",
				EndTimeRange:   "1617333947000",
				QueryBlob:      lqlQueryStr,
			},
			AllowEmpty: false,
			Return:     "date range should have a start time before the end time",
		},
		LQLValidateRangeTest{
			Name: "start-equal-end",
			Input: api.LQLQuery{
				StartTimeRange: "1617333947000",
				EndTimeRange:   "1617333947000",
				QueryBlob:      lqlQueryStr,
			},
			AllowEmpty: false,
			Return:     nil,
		},
	}
	lqlValidateTests []LQLValidateTest = []LQLValidateTest{
		LQLValidateTest{
			Name: "empty",
			Input: &api.LQLQuery{
				StartTimeRange: "0",
				EndTimeRange:   "1",
				QueryText:      lqlQueryStr,
			},
			Return: nil,
		},
	}
	lqlQueryTypeTests []LQLQueryTest = []LQLQueryTest{
		LQLQueryTest{
			Name: "empty-blob",
			Input: &api.LQLQuery{
				QueryBlob: ``,
			},
			Return: api.LQLQueryTranslateError,
			Expected: &api.LQLQuery{
				QueryText: ``,
				QueryBlob: ``,
			},
		},
		LQLQueryTest{
			Name: "junk-blob",
			Input: &api.LQLQuery{
				QueryBlob: `this is junk`,
			},
			Return: api.LQLQueryTranslateError,
			Expected: &api.LQLQuery{
				QueryText: ``,
				QueryBlob: `this is junk`,
			},
		},
		LQLQueryTest{
			Name: "json-blob",
			Input: &api.LQLQuery{
				QueryBlob: `{
	"START_TIME_RANGE": "678910",
	"END_TIME_RANGE": "111213141516",
	"QUERY_TEXT": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
			Return: nil,
			Expected: &api.LQLQuery{
				StartTimeRange: "1970-01-01T00:11:18Z",
				EndTimeRange:   "1973-07-11T04:32:21Z",
				QueryText:      "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
				QueryBlob: `{
	"START_TIME_RANGE": "678910",
	"END_TIME_RANGE": "111213141516",
	"QUERY_TEXT": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
		},
		LQLQueryTest{
			Name: "json-blob-lower",
			Input: &api.LQLQuery{
				QueryBlob: `{
	"start_time_range": "678910",
	"end_time_range": "111213141516",
	"query_text": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
			Return: nil,
			Expected: &api.LQLQuery{
				StartTimeRange: "1970-01-01T00:11:18Z",
				EndTimeRange:   "1973-07-11T04:32:21Z",
				QueryText:      "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
				QueryBlob: `{
	"start_time_range": "678910",
	"end_time_range": "111213141516",
	"query_text": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
		},
		LQLQueryTest{
			Name: "lql-blob",
			Input: &api.LQLQuery{
				QueryBlob: "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
			},
			Return: nil,
			Expected: &api.LQLQuery{
				QueryText: "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
				QueryBlob: "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }",
			},
		},
		LQLQueryTest{
			Name: "overwrite-blob",
			Input: &api.LQLQuery{
				StartTimeRange: "0",
				EndTimeRange:   "1",
				QueryText:      "should not overwrite",
				QueryBlob: `{
	"START_TIME_RANGE": "678910",
	"END_TIME_RANGE": "111213141516",
	"QUERY_TEXT": "my_lql(CloudTrailRawEvents e) { SELECT INSERT_ID LIMIT 10 }"
}`,
			},
			Return: nil,
			Expected: &api.LQLQuery{
				StartTimeRange: "1970-01-01T00:00:00Z",
				EndTimeRange:   "1970-01-01T00:00:00Z",
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

type LQLTranslateTimeTest struct {
	Name       string
	Input      string
	ReturnTime string
	ReturnErr  interface{}
}

type LQLValidateRangeTest struct {
	Name       string
	Input      api.LQLQuery
	AllowEmpty bool
	Return     interface{}
}

type LQLValidateTest struct {
	Name   string
	Input  *api.LQLQuery
	Return interface{}
}

type LQLQueryTest struct {
	Name     string
	Input    *api.LQLQuery
	Return   interface{}
	Expected *api.LQLQuery
}

func TestTranslateTime(t *testing.T) {
	for _, lqlTranslateTimeTest := range lqlTranslateTimeTests {
		t.Run(lqlTranslateTimeTest.Name, func(t *testing.T) {
			outTime, err := api.LQLQuery{}.TranslateTime(lqlTranslateTimeTest.Input)
			if err == nil {
				assert.Equal(t, lqlTranslateTimeTest.ReturnTime, outTime)
				assert.Equal(t, lqlTranslateTimeTest.ReturnErr, err)
			} else {
				assert.Equal(t, lqlTranslateTimeTest.ReturnErr, err.Error())
			}
		})
	}
}

func TestValidateRange(t *testing.T) {
	for _, lqlValidateRangeTest := range lqlValidateRangeTests {
		t.Run(lqlValidateRangeTest.Name, func(t *testing.T) {
			err := lqlValidateRangeTest.Input.Translate()
			assert.Nil(t, err)
			err = lqlValidateRangeTest.Input.ValidateRange(lqlValidateRangeTest.AllowEmpty)
			if err == nil {
				assert.Equal(t, lqlValidateRangeTest.Return, err)
			} else {
				assert.Equal(t, lqlValidateRangeTest.Return, err.Error())
			}
		})
	}
}

func TestValidate(t *testing.T) {
	for _, lqlValidateTest := range lqlValidateTests {
		t.Run(lqlValidateTest.Name, func(t *testing.T) {
			err := lqlValidateTest.Input.Translate()
			assert.Nil(t, err)
			err = lqlValidateTest.Input.ValidateRange(true)
			if err == nil {
				assert.Equal(t, lqlValidateTest.Return, err)
			} else {
				assert.Equal(t, lqlValidateTest.Return, err.Error())
			}
		})
	}
}

func TestLQLQueryTranslate(t *testing.T) {
	for _, lqlQueryTest := range lqlQueryTypeTests {
		t.Run(lqlQueryTest.Name, func(t *testing.T) {
			if err := lqlQueryTest.Input.Translate(); err == nil {
				assert.Equal(t, lqlQueryTest.Return, err)
			} else {
				assert.Equal(t, lqlQueryTest.Return, err.Error())
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

	_, err = c.LQL.RunQuery(lqlQueryStr, "0", "1")
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
	runActual, err = c.LQL.RunQuery(lqlQueryStr, "0", "1")
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
