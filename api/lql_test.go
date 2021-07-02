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
	"github.com/lacework/go-sdk/lwtime"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	queryEvaluator = "Cloudtrail"
	queryID        = "my_lql"
	queryText      = `my_lql { source { CloudTrailRawEvents } return { INSERT_ID } }`
	queryJSON      = fmt.Sprintf(`{
	"evaluatorId": "%s",
	"queryId": "%s",
	"queryText": "%s"
}`, queryEvaluator, queryID, queryText)
	queryYAML = fmt.Sprintf(`---
evaluatorId: %s
queryId: %s
queryText: %s
`, queryEvaluator, queryID, queryText)
	lqlRunData = `[
	{
		"INSERT_ID": "35308423"
	}
]`
	lqlErrorReponse = `{ "message": "This is an error message" }`
)

type LQLParseQueryTimeTest struct {
	Name       string
	Input      string
	ReturnTime string
	ReturnErr  error
}

var (
	atDay, _               = lwtime.ParseRelative("@d")
	lqlParseQueryTimeTests = []LQLParseQueryTimeTest{
		LQLParseQueryTimeTest{
			Name:       "valid-rfc-utc",
			Input:      "2021-03-31T00:00:00Z",
			ReturnTime: "2021-03-31T00:00:00Z",
			ReturnErr:  nil,
		},
		LQLParseQueryTimeTest{
			Name:       "valid-rfc-central",
			Input:      "2021-03-31T00:00:00-05:00",
			ReturnTime: "2021-03-31T05:00:00Z",
			ReturnErr:  nil,
		},
		LQLParseQueryTimeTest{
			Name:       "valid-milli",
			Input:      "1617230464000",
			ReturnTime: "2021-03-31T22:41:04Z",
			ReturnErr:  nil,
		},
		LQLParseQueryTimeTest{
			Name:       "valid-relative",
			Input:      "@d",
			ReturnTime: atDay.UTC().Format(time.RFC3339),
			ReturnErr:  nil,
		},
		LQLParseQueryTimeTest{
			Name:       "empty",
			Input:      "",
			ReturnTime: "0001-01-01T00:00:00Z",
			ReturnErr:  errors.New("unable to parse time ()"),
		},
		LQLParseQueryTimeTest{
			Name:       "invalid",
			Input:      "jweaver",
			ReturnTime: "0001-01-01T00:00:00Z",
			ReturnErr:  errors.New("unable to parse time (jweaver)"),
		},
	}
)

func TestLQLParseQueryTime(t *testing.T) {
	for _, lqlPQTT := range lqlParseQueryTimeTests {
		t.Run(lqlPQTT.Name, func(t *testing.T) {
			outTime, err := api.ParseQueryTime(lqlPQTT.Input)
			if err == nil {
				assert.Equal(t, lqlPQTT.ReturnErr, err)
			} else {
				assert.Equal(t, lqlPQTT.ReturnErr.Error(), err.Error())
			}
			assert.Equal(t, lqlPQTT.ReturnTime, outTime.UTC().Format(time.RFC3339))
		})
	}
}

type LQLValidateQueryRangeTest struct {
	Name           string
	StartTimeRange time.Time
	EndTimeRange   time.Time
	Return         error
}

var lqlValidateQueryRangeTests = []LQLValidateQueryRangeTest{
	LQLValidateQueryRangeTest{
		Name:           "ok",
		StartTimeRange: time.Unix(0, 0),
		EndTimeRange:   time.Unix(1, 0),
		Return:         nil,
	},
	LQLValidateQueryRangeTest{
		Name:           "empty-start",
		StartTimeRange: time.Time{},
		EndTimeRange:   time.Unix(1, 0),
		Return:         nil,
	},
	LQLValidateQueryRangeTest{
		Name:           "empty-end",
		StartTimeRange: time.Unix(1, 0),
		EndTimeRange:   time.Time{},
		Return:         errors.New("date range should have a start time before the end time"),
	},
	LQLValidateQueryRangeTest{
		Name:           "start-after-end",
		StartTimeRange: time.Unix(1717333947, 0),
		EndTimeRange:   time.Unix(1617333947, 0),
		Return:         errors.New("date range should have a start time before the end time"),
	},
	LQLValidateQueryRangeTest{
		Name:           "start-equal-end",
		StartTimeRange: time.Unix(1617333947, 0),
		EndTimeRange:   time.Unix(1617333947, 0),
		Return:         nil,
	},
}

func TestLQLValidateQueryRange(t *testing.T) {
	for _, lqlVQRT := range lqlValidateQueryRangeTests {
		t.Run(lqlVQRT.Name, func(t *testing.T) {
			err := api.ValidateQueryRange(lqlVQRT.StartTimeRange, lqlVQRT.EndTimeRange)
			if err == nil {
				assert.Equal(t, lqlVQRT.Return, err)
			} else {
				assert.Equal(t, lqlVQRT.Return.Error(), err.Error())
			}
		})
	}
}

type LQLParseQueryTest struct {
	Name     string
	Input    string
	Return   error
	Expected api.Query
}

var lqlParseQueryTests = []LQLParseQueryTest{
	LQLParseQueryTest{
		Name:     "empty-blob",
		Input:    "",
		Return:   errors.New("query must be valid JSON or YAML"),
		Expected: api.Query{},
	},
	LQLParseQueryTest{
		Name:     "junk-blob",
		Input:    "this is junk",
		Return:   errors.New("query must be valid JSON or YAML"),
		Expected: api.Query{},
	},
	LQLParseQueryTest{
		Name:     "partial-blob",
		Input:    "{",
		Return:   errors.New("query must be valid JSON or YAML"),
		Expected: api.Query{},
	},
	LQLParseQueryTest{
		Name:   "json-blob",
		Input:  queryJSON,
		Return: nil,
		Expected: api.Query{
			ID:          queryID,
			QueryText:   queryText,
			EvaluatorID: queryEvaluator,
		},
	},
	LQLParseQueryTest{
		Name:   "yaml-blob",
		Input:  queryYAML,
		Return: nil,
		Expected: api.Query{
			ID:          queryID,
			QueryText:   queryText,
			EvaluatorID: queryEvaluator,
		},
	},
}

func TestLQLParseQuery(t *testing.T) {
	for _, lqlPQT := range lqlParseQueryTests {
		t.Run(lqlPQT.Name, func(t *testing.T) {
			actual, err := api.ParseQuery(lqlPQT.Input)
			if err == nil {
				assert.Equal(t, lqlPQT.Return, err)
			} else {
				assert.Equal(t, lqlPQT.Return.Error(), err.Error())
			}
			assert.Equal(t, lqlPQT.Expected, actual)
		})
	}
}

func mockLQLDataResponse(data string) string {
	return `{
	"data": ` + data + `
}`
}

func TestLQLCreateMethod(t *testing.T) {
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

	_, err = c.LQL.Create(queryJSON)
	assert.Nil(t, err)
}

func TestLQLCreateBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"Queries",
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

	_, err = c.LQL.Create("")
	assert.Equal(t, "query must be valid JSON or YAML", err.Error())
}

func TestLQLCreateOK(t *testing.T) {
	mockResponse := mockLQLDataResponse(queryJSON)

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
	createActual, err = c.LQL.Create(queryJSON)
	assert.Nil(t, err)

	assert.Equal(t, createExpected, createActual)
}

func TestLQLCreateError(t *testing.T) {
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

	_, err = c.LQL.Create(queryJSON)
	assert.NotNil(t, err)
}

func TestLQLGetQueriesMethod(t *testing.T) {
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

	_, err = c.LQL.GetQueries()
	assert.Nil(t, err)
}

func TestLQLGetQueryByIDOK(t *testing.T) {
	mockResponse := mockLQLDataResponse(queryJSON)

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
	getActual, err = c.LQL.GetByID(queryID)
	assert.Nil(t, err)

	assert.Equal(t, getExpected, getActual)
}

func TestLQLGetByIDNotFound(t *testing.T) {
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

	_, err = c.LQL.GetByID("NoSuchQuery")
	assert.NotNil(t, err)
}

func TestLQLRunMethod(t *testing.T) {
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

	_, err = c.LQL.Run(queryJSON, "0", "1")
	assert.Nil(t, err)
}

func TestLQLRunBadInput(t *testing.T) {
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

	_, err = c.LQL.Run("", "", "")
	assert.Equal(t, "query must be valid JSON or YAML", err.Error())
}

func TestLQLRunOK(t *testing.T) {
	mockResponse := mockLQLDataResponse(lqlRunData)

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
	runActual, err = c.LQL.Run(queryJSON, "0", "1")
	assert.Nil(t, err)

	assert.Equal(t, runExpected, runActual)
}

func TestLQLRunError(t *testing.T) {
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

	_, err = c.LQL.Run(queryJSON, "0", "1")
	assert.NotNil(t, err)
}
