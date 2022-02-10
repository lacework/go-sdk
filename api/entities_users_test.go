//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestEntities_Users_List(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/Users/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockUserEntityResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Entities.ListUsers()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "systemd-network", response.Data[0].PrimaryGroupName)
	}
}

func TestEntities_Users_List_All(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/Users/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockUserEntityResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Entities.ListAllUsers()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "systemd-network", response.Data[0].PrimaryGroupName)
	}
}

func mockUserEntityResponse() string {
	return `
{
  "data": [
    {
      "createdTime": "2022-02-08T10:03:08.459Z",
      "mid": 51,
      "otherGroupNames": [
        "systemd-network"
      ],
      "primaryGroupName": "systemd-network",
      "uid": 100,
      "username": "systemd-network"
    }
  ],
  "paging": {
    "rows": 1,
    "totalRows": 1,
    "urls": {
      "nextPage": null
    }
  }
}
	`
}
