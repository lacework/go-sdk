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
	lqlDataSource   string = "CloudTrailRawEvents"
	lqlDescribeData string = `[{
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

func TestLQLDescribeMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", api.ApiLQLDescribe, url.QueryEscape(lqlDataSource)),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Describe should be a GET method")
			assert.Subset(
				t,
				[]byte(r.RequestURI),
				[]byte(lqlDataSource),
				"Describe should specify lqlDataSource route",
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

	_, err = c.LQL.Describe(lqlDataSource)
	assert.Nil(t, err)
}

func TestLQLDescribeOK(t *testing.T) {
	mockResponse := mockLQLDataResponse(lqlDescribeData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", api.ApiLQLDescribe, url.QueryEscape(lqlDataSource)),
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

	describeExpected := api.LQLDescribeResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &describeExpected)

	var describeActual api.LQLDescribeResponse
	describeActual, err = c.LQL.Describe(lqlDataSource)
	assert.Nil(t, err)

	assert.Equal(t, describeExpected, describeActual)
}

func TestLQLDescribeNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("%s/%s", api.ApiLQLDescribe, url.QueryEscape(lqlDataSource)),
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

	_, err = c.LQL.Describe("NoSuchDataSource")
	assert.NotNil(t, err)
}

func TestLQLDescribeError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQLCompile,
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

	_, err = c.LQL.Describe(lqlDataSource)
	assert.NotNil(t, err)
}
