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
	"net/url"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

var (
	QueryDataSource   = "CloudTrailRawEvents"
	QueryDescribeData = `[{
		"props": {
			"cachable": true
		},
		"type": "Entity",
		"maxDuration": -1,
		"complexity": 3,
		"schema": [
			{
				"name": "INSERT_ID",
				"type": "Number",
				"props": {}
			},
			{
				"name": "INSERT_TIME",
				"type": "Timestamp",
				"props": {}
			},
			{
				"name": "EVENT_TIME",
				"type": "Timestamp",
				"props": {}
			},
			{
				"name": "EVENT",
				"type": "String",
				"props": {
					"jsonObject": true
				}
			},
			{
				"name": "EVENT_SOURCE",
				"type": "String",
				"props": {}
			},
			{
				"name": "EVENT_NAME",
				"type": "String",
				"props": {}
			},
			{
				"name": "ERROR_CODE",
				"type": "String",
				"props": {}
			}
		],
		"parameters": [
			{
				"required": false,
				"name": "StartTimeRange",
				"type": "Timestamp",
				"default": null,
				"props": null
			},
			{
				"required": true,
				"name": "EventRawTable",
				"type": "String",
				"default": "CLOUD_TRAIL_INTERNAL.EVENT_RAW_T",
				"props": null
			},
			{
				"required": false,
				"name": "BATCH_ID",
				"type": "Number",
				"default": null,
				"props": null
			},
			{
				"required": false,
				"name": "EndTimeRange",
				"type": "Timestamp",
				"default": null,
				"props": null
			}
		],
		"primaryKey": []
	}]`
)

func TestQueryDescribeMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("external/lql/describe/%s", url.QueryEscape(QueryDataSource)),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Describe should be a GET method")
			assert.Subset(
				t,
				[]byte(r.RequestURI),
				[]byte(QueryDataSource),
				"Describe should specify QueryDataSource route",
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

	_, err = c.Query.Describe(QueryDataSource)
	assert.Nil(t, err)
}

func TestQueryDescribeOK(t *testing.T) {
	mockResponse := mockQueryDataResponse(QueryDescribeData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("external/lql/describe/%s", url.QueryEscape(QueryDataSource)),
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

	describeExpected := api.QueryDescribeResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &describeExpected)

	var describeActual api.QueryDescribeResponse
	describeActual, err = c.Query.Describe(QueryDataSource)
	assert.Nil(t, err)

	assert.Equal(t, describeExpected, describeActual)
}

func TestQueryDescribeNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("external/lql/describe/%s", url.QueryEscape(QueryDataSource)),
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

	_, err = c.Query.Describe("NoSuchDataSource")
	assert.NotNil(t, err)
}

func TestQueryDescribeError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/lql/compile",
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

	_, err = c.Query.Describe(QueryDataSource)
	assert.NotNil(t, err)
}
