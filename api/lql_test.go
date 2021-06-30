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
	lqlQueryID    = "my_lql"
	lqlQueryStr   = "my_lql { source { CloudTrailRawEvents } return { INSERT_ID } }"
	lqlCreateData = fmt.Sprintf(`{
	"queryID": "my_lql",
	"queryText": "%s"
}`, lqlQueryStr)
	lqlRunData = `[
	{
		"INSERT_ID": "35308423"
	}
]`
	lqlErrorReponse = `{ "message": "This is an error message" }`
)

type LQLTranslateTimeTest struct {
	Name       string
	Input      string
	ReturnTime string
	ReturnErr  interface{}
}

var (
	atDay, _              = lwtime.ParseRelative("@d")
	lqlTranslateTimeTests = []LQLTranslateTimeTest{
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
			Name:       "valid-relative",
			Input:      "@d",
			ReturnTime: atDay.UTC().Format(time.RFC3339),
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
)

func TestLQLTranslateTime(t *testing.T) {
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

type LQLValidateRangeTest struct {
	Name       string
	Input      api.LQLQuery
	AllowEmpty bool
	Return     interface{}
}

var lqlValidateRangeTests = []LQLValidateRangeTest{
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

func TestLQLValidateRange(t *testing.T) {
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

type LQLValidateTest struct {
	Name   string
	Input  *api.LQLQuery
	Return interface{}
}

var lqlValidateTests = []LQLValidateTest{
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

func TestLQLValidate(t *testing.T) {
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

type LQLQueryTest struct {
	Name     string
	Input    *api.LQLQuery
	Return   error
	Expected *api.LQLQuery
}

var lqlQueryPopulateIDTests = []LQLQueryTest{
	LQLQueryTest{
		Name:     "empty",
		Input:    &api.LQLQuery{},
		Return:   errors.New("unable to extract ID from query text"),
		Expected: &api.LQLQuery{},
	},
	LQLQueryTest{
		Name: "junk",
		Input: &api.LQLQuery{
			QueryText: "this is junk",
		},
		Return: errors.New("unable to extract ID from query text"),
		Expected: &api.LQLQuery{
			QueryText: "this is junk",
		},
	},
	LQLQueryTest{
		Name: "simple",
		Input: &api.LQLQuery{
			QueryText: lqlQueryStr,
		},
		Return: nil,
		Expected: &api.LQLQuery{
			ID:        lqlQueryID,
			QueryText: lqlQueryStr,
		},
	},
	LQLQueryTest{
		Name: "newlines",
		Input: &api.LQLQuery{
			QueryText: `
-- a comment
my query {
}`,
		},
		Return: nil,
		Expected: &api.LQLQuery{
			ID: "my query",
			QueryText: `
-- a comment
my query {
}`,
		},
	},
}

func TestLQLQueryPopulateID(t *testing.T) {
	for _, lqlQueryTest := range lqlQueryPopulateIDTests {
		t.Run(lqlQueryTest.Name, func(t *testing.T) {
			if err := lqlQueryTest.Input.PopulateID(); err == nil {
				assert.Equal(t, lqlQueryTest.Return, err)
			} else {
				assert.Equal(t, lqlQueryTest.Return.Error(), err.Error())
			}
			assert.Equal(t, lqlQueryTest.Expected, lqlQueryTest.Input)
		})
	}
}

var lqlQueryTypeTests = []LQLQueryTest{
	LQLQueryTest{
		Name: "empty-blob",
		Input: &api.LQLQuery{
			QueryBlob: ``,
		},
		Return: errors.New(api.LQLQueryTranslateError),
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
		Return: errors.New(api.LQLQueryTranslateError),
		Expected: &api.LQLQuery{
			QueryText: ``,
			QueryBlob: `this is junk`,
		},
	},
	LQLQueryTest{
		Name: "partial-blob",
		Input: &api.LQLQuery{
			QueryBlob: `{`,
		},
		Return: errors.New(api.LQLQueryTranslateError),
		Expected: &api.LQLQuery{
			QueryText: ``,
			QueryBlob: `{`,
		},
	},
	LQLQueryTest{
		Name: "json-blob",
		Input: &api.LQLQuery{
			QueryBlob: fmt.Sprintf(`{
"start_time_range": "678910",
"end_time_range": "111213141516",
"queryText": "%s"
}`, lqlQueryStr),
		},
		Return: nil,
		Expected: &api.LQLQuery{
			StartTimeRange: "1970-01-01T00:11:18Z",
			EndTimeRange:   "1973-07-11T04:32:21Z",
			ID:             lqlQueryID,
			QueryText:      lqlQueryStr,
			QueryBlob: fmt.Sprintf(`{
"start_time_range": "678910",
"end_time_range": "111213141516",
"queryText": "%s"
}`, lqlQueryStr),
		},
	},
	LQLQueryTest{
		Name: "json-blob-lower",
		Input: &api.LQLQuery{
			QueryBlob: fmt.Sprintf(`{
"start_time_range": "678910",
"end_time_range": "111213141516",
"queryText": "%s"
}`, lqlQueryStr),
		},
		Return: nil,
		Expected: &api.LQLQuery{
			StartTimeRange: "1970-01-01T00:11:18Z",
			EndTimeRange:   "1973-07-11T04:32:21Z",
			ID:             lqlQueryID,
			QueryText:      lqlQueryStr,
			QueryBlob: fmt.Sprintf(`{
"start_time_range": "678910",
"end_time_range": "111213141516",
"queryText": "%s"
}`, lqlQueryStr),
		},
	},
	LQLQueryTest{
		Name: "lql-blob",
		Input: &api.LQLQuery{
			QueryBlob: fmt.Sprintf("--a comment\n%s", lqlQueryStr),
		},
		Return: nil,
		Expected: &api.LQLQuery{
			ID:        lqlQueryID,
			QueryText: fmt.Sprintf("--a comment\n%s", lqlQueryStr),
			QueryBlob: fmt.Sprintf("--a comment\n%s", lqlQueryStr),
		},
	},
	LQLQueryTest{
		Name: "overwrite-blob",
		Input: &api.LQLQuery{
			StartTimeRange: "0",
			EndTimeRange:   "1",
			QueryText:      "should not overwrite",
			QueryBlob: fmt.Sprintf(`{
"start_time_range": "678910",
"end_time_range": "111213141516",
"queryText": "%s"
}`, lqlQueryStr),
		},
		Return: errors.New("unable to extract ID from query text"),
		Expected: &api.LQLQuery{
			StartTimeRange: "0",
			EndTimeRange:   "1",
			ID:             "",
			QueryText:      "should not overwrite",
			QueryBlob: fmt.Sprintf(`{
"start_time_range": "678910",
"end_time_range": "111213141516",
"queryText": "%s"
}`, lqlQueryStr),
		},
	},
}

func TestLQLQueryTranslate(t *testing.T) {
	for _, lqlQueryTest := range lqlQueryTypeTests {
		t.Run(lqlQueryTest.Name, func(t *testing.T) {
			if err := lqlQueryTest.Input.Translate(); err == nil {
				assert.Equal(t, lqlQueryTest.Return, err)
			} else {
				assert.Equal(t, lqlQueryTest.Return.Error(), err.Error())
			}
			assert.Equal(t, lqlQueryTest.Expected, lqlQueryTest.Input)
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

	_, err = c.LQL.Create(lqlQueryStr)
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
	assert.Equal(t, api.LQLQueryTranslateError, err.Error())
}

func TestLQLCreateOK(t *testing.T) {
	mockResponse := mockLQLDataResponse(lqlCreateData)

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

	createExpected := api.LQLQueryResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &createExpected)

	var createActual api.LQLQueryResponse
	createActual, err = c.LQL.Create(lqlQueryStr)
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

	_, err = c.LQL.Create(lqlQueryStr)
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
	mockResponse := mockLQLDataResponse(lqlCreateData)

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

	getExpected := api.LQLQueryResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &getExpected)

	var getActual api.LQLQueryResponse
	getActual, err = c.LQL.GetByID(lqlQueryID)
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

	_, err = c.LQL.Run(lqlQueryStr, "0", "1")
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
	assert.Equal(t, api.LQLQueryTranslateError, err.Error())
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
	runActual, err = c.LQL.Run(lqlQueryStr, "0", "1")
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

	_, err = c.LQL.Run(lqlQueryStr, "0", "1")
	assert.NotNil(t, err)
}
