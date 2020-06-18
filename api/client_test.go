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
	fakeServer.ApiVersion = "v2"
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithURL(fakeServer.URL()),
		api.WithExpirationTime(1800),
		api.WithApiV2(),
		api.WithLogLevel("DEBUG"),
		api.WithHeader("User-Agent", "test-agent"),
		api.WithTokenFromKeys("KEY", "SECRET"), // this option has to be the last one
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "v2", c.ApiVersion(), "modified API version should be v2")
	}
}
