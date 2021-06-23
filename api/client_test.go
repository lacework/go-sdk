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
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestNewClient(t *testing.T) {
	c, err := api.NewClient("test")
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "default API version should be v1")
	}
}

func TestNewClientAccountEmptyError(t *testing.T) {
	c, err := api.NewClient("")
	assert.Nil(t, c)
	if assert.NotNil(t, err) {
		assert.Equal(t, "account cannot be empty", err.Error(),
			"we cannot create an api client without a Lacework account")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithURL(fakeServer.URL()),
		api.WithExpirationTime(1800),
		api.WithApiV2(),
		api.WithTimeout(time.Minute*5),
		api.WithLogLevel("DEBUG"),
		api.WithHeader("User-Agent", "test-agent"),
		api.WithTokenFromKeys("KEY", "SECRET"), // this option has to be the last one
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "v2", c.ApiVersion(), "modified API version should be v2")
	}
}

func TestCopyClientWithOptions(t *testing.T) {
	var v interface{}
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI(
		"endpoint-org-access",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "true", r.Header.Get("Org-Access"))
			fmt.Fprintf(w, "{}")
		},
	)
	fakeServer.MockAPI(
		"endpoint-NO-org-access",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "", r.Header.Get("Org-Access"))
			fmt.Fprintf(w, "{}")
		},
	)

	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithURL(fakeServer.URL()),
		api.WithExpirationTime(1800),
		api.WithTimeout(time.Minute*5),
		api.WithLogLevel("DEBUG"),
		api.WithHeader("User-Agent", "test-agent"),
		api.WithTokenFromKeys("KEY", "SECRET"), // this option has to be the last one
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "default API version should be v1")
		assert.Contains(t, c.URL(), "http://127.0.0.1:", "wrong URL")
		assert.True(t, c.ValidAuth())
	}

	err = c.RequestDecoder("GET", "endpoint-NO-org-access", nil, v)
	assert.Nil(t, err)

	newExactClient, err := api.CopyClient(c)
	if assert.Nil(t, err) {
		assert.Equal(t, c.ApiVersion(), newExactClient.ApiVersion(), "copy client mismatch")
		assert.Equal(t, c.URL(), newExactClient.URL(), "copy client mismatch")
		assert.True(t, newExactClient.ValidAuth())
	}

	td, err := newExactClient.GenerateToken()
	if assert.Nil(t, err) {
		assert.Equal(t, "TOKEN", td.Token)
	}

	newModifiedClient, err := api.CopyClient(c,
		api.WithURL("https://new.lacework.net/"),
		api.WithExpirationTime(3600),
		api.WithApiV2(),
		api.WithTimeout(time.Minute*60), // LOL!
		api.WithLogLevel("INFO"),
		api.WithOrgAccess(),
	)
	if assert.Nil(t, err) {
		assert.NotEqual(t, c.ApiVersion(), newModifiedClient.ApiVersion(), "copy modified client mismatch")
		assert.NotEqual(t, c.URL(), newModifiedClient.URL(), "copy modified client mismatch")
		assert.Equal(t, "v2", newModifiedClient.ApiVersion(), "copy modified API version should be v2")
		assert.Equal(t, "https://new.lacework.net/", newModifiedClient.URL(), "copy modified client mismatch")
		assert.True(t, newExactClient.ValidAuth())
	}

	err = c.RequestDecoder("GET", "endpoint-org-access", nil, v)
	assert.Nil(t, err)
}
