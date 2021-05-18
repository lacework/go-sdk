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
	lqlCompileData = `[{
		"name": "my_lql",
		"props": {
			"lql": true
		},
		"type": "Entity",
		"maxDuration": -1,
		"complexity": 3,
		"schema": [
			{
				"name": "INSERT_ID",
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

func TestLQLCompileMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/lql/compile",
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

	_, err = c.LQL.CompileQuery(lqlQueryStr)
	assert.Nil(t, err)
}

func TestLQLCompileBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/lql/compile",
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

	_, err = c.LQL.CompileQuery("")
	assert.Equal(t, api.LQLQueryTranslateError, err.Error())
}

func TestLQLCompileOK(t *testing.T) {
	mockResponse := mockLQLDataResponse(lqlCompileData)

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/lql/compile",
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

	compileExpected := api.LQLCompileResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &compileExpected)

	var compileActual api.LQLCompileResponse
	compileActual, err = c.LQL.CompileQuery(lqlQueryStr)
	assert.Nil(t, err)
	assert.Equal(t, compileExpected, compileActual)
}

func TestLQLCompileError(t *testing.T) {
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

	_, err = c.LQL.CompileQuery(lqlQueryStr)
	assert.NotNil(t, err)
}
