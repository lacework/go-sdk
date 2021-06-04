//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

func TestAccountOrganizationInfoForStandalone(t *testing.T) {
	var fakeServer = lacework.MockServer()
	fakeServer.MockAPI("external/account/organizationInfo", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetOrganizationInfo should be a GET method")
		fmt.Fprintf(w, accountOrganizationInfoResponseStandalone())
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Account.GetOrganizationInfo()
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.False(t, response.OrgAccount)
		assert.Empty(t, response.OrgAccountURL)
		assert.Empty(t, response.AccountName())
	}
}

func TestAccountOrganizationInfoForOrganizational(t *testing.T) {
	var fakeServer = lacework.MockServer()
	fakeServer.MockAPI("external/account/organizationInfo", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetOrganizationInfo should be a GET method")
		fmt.Fprintf(w, accountOrganizationInfoResponseOrganizational("test-org"))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Account.GetOrganizationInfo()
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.True(t, response.OrgAccount)
		assert.Equal(t, "test-org", response.AccountName())
		assert.Equal(t, "test-org.lacework.net", response.OrgAccountURL)
	}
}

func accountOrganizationInfoResponseStandalone() string {
	return `{ "orgAccount": false }`
}

func accountOrganizationInfoResponseOrganizational(name string) string {
	return `{
  "orgAccount": true,
  "orgAccountUrl": "` + name + `.lacework.net"
}`
}
