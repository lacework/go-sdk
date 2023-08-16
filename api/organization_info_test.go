//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestOrganizationInfoGet(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("OrganizationInfo",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetOrganizationInfo() should be a GET method")
			fmt.Fprintf(w, mockOrganizationInfoResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.OrganizationInfo.Get()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "myAccountName.lacework.net", response.Data[0].OrgAccountURL)
		assert.True(t, response.Data[0].OrgAccount)
		assert.Equal(t, "myAccountName", response.Data[0].AccountName())
	}
}

func mockOrganizationInfoResponse() string {
	return `
{
  "data": [
    { "orgAccount": true,
      "orgAccountUrl": "myAccountName.lacework.net"
	}
  ]
}`
}
