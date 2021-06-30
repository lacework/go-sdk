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

func TestLQLDeleteMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/my_lql",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method, "Delete should be a DELETE method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.LQL.Delete(lqlQueryID)
	fmt.Println(err)
	assert.Nil(t, err)
}

func TestLQLDeleteBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/my_lql",
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

	_, err = c.LQL.Delete("")
	assert.NotNil(t, err)
}

func TestLQLDeleteOK(t *testing.T) {

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/my_lql",
		func(w http.ResponseWriter, r *http.Request) {
			// send the headers with a 204 response code.
			w.WriteHeader(http.StatusNoContent)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	deleteExpected := api.LQLDeleteResponse{}
	_ = json.Unmarshal([]byte(""), &deleteExpected)

	var deleteActual api.LQLDeleteResponse
	deleteActual, err = c.LQL.Delete(lqlQueryID)
	assert.Nil(t, err)

	assert.Equal(t, deleteExpected, deleteActual)
}

func TestLQLDeleteNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/my_lql",
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

	_, err = c.LQL.Delete(lqlQueryID)
	assert.NotNil(t, err)
}
