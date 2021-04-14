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

func TestLQLUpdateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQL,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PATCH", r.Method, "Update should be a PATCH method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.LQL.UpdateQuery(lqlQueryStr)
	assert.Nil(t, err)
}

func TestLQLUpdateBadInput(t *testing.T) {
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

	_, err = c.LQL.UpdateQuery("")
	assert.Equal(t, api.LQLQueryTranslateError, err.Error())
}

func TestLQLUpdateOK(t *testing.T) {
	mockResponse := mockLQLMessageResponse(
		fmt.Sprintf(`"lqlUpdated": "%s"`, lqlQueryID),
		"true",
	)

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

	updateExpected := api.LQLUpdateResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &updateExpected)

	var updateActual api.LQLUpdateResponse
	updateActual, err = c.LQL.UpdateQuery(lqlQueryStr)
	assert.Nil(t, err)

	assert.Equal(t, updateExpected, updateActual)
}

func TestLQLUpdateNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		api.ApiLQLCompile,
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

	_, err = c.LQL.UpdateQuery(lqlQueryStr)
	assert.NotNil(t, err)
}
