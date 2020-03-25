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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestWithApiV2(t *testing.T) {
	c, err := api.NewClient("test", api.WithApiV2())
	if assert.Nil(t, err) {
		assert.Equal(t, "v2", c.ApiVersion(), "API version should be v2")
	}
}

func TestWithToken(t *testing.T) {
	c, err := api.NewClient("test", api.WithToken("TOKEN"))
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "API version should be v2")
	}
}

func TestApiVersion(t *testing.T) {
	c, err := api.NewClient("foo")
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "wrong default API version")
	}
}

func TestWithApiKeys(t *testing.T) {
	c, err := api.NewClient("foo", api.WithApiKeys("KEY", "SECRET"))
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "wrong default API version")
	}
}

func TestWithTokenFromKeys(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	c, err := api.NewClient("foo",
		api.WithURL(fakeServer.URL()),
		api.WithTokenFromKeys("KEY", "SECRET"), // this option has to be the last one
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "v1", c.ApiVersion(), "wrong default API version")
	}
}

func TestGenerateToken(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	c, err := api.NewClient("foo",
		api.WithURL(fakeServer.URL()),
		api.WithApiKeys("KEY", "SECRET"),
	)
	if assert.Nil(t, err) {
		response, err := c.GenerateToken()
		assert.Nil(t, err)
		assert.Equal(t, "TOKEN", response.Token(), "token mismatch")
	}
}

func TestGenerateTokenWithKeys(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	c, err := api.NewClient("foo", api.WithURL(fakeServer.URL()))
	if assert.Nil(t, err) {
		response, err := c.GenerateTokenWithKeys("KEY", "SECRET")
		assert.Nil(t, err)
		assert.Equal(t, "TOKEN", response.Token(), "token mismatch")
	}
}

func TestGenerateTokenErrorKeysMissing(t *testing.T) {
	c, err := api.NewClient("where-are-my-keys")
	if assert.Nil(t, err) {
		response, err := c.GenerateToken()
		if assert.NotNil(t, err) {
			assert.Empty(t, response, "token must be empty")
			assert.Equal(t,
				"unable to generate access token: auth keys missing",
				err.Error(),
				"error message mismatch",
			)
		}
	}
}
