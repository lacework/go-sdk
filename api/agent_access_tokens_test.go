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
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestAgentAccessTokenPrettyState(t *testing.T) {
	subject := api.AgentAccessToken{Enabled: 1}
	assert.Equal(t, "Enabled", subject.PrettyState())

	subject.Enabled = 0
	assert.Equal(t, "Disabled", subject.PrettyState())
	// anything else
	subject.Enabled = 9
	assert.Equal(t, "Disabled", subject.PrettyState())
}

func TestAgentAccessTokenState(t *testing.T) {
	subject := api.AgentAccessToken{Enabled: 1}
	assert.True(t, subject.State())

	subject.Enabled = 0
	assert.False(t, subject.State())
	// anything else
	subject.Enabled = 9
	assert.False(t, subject.State())
}

func TestAgentAccessTokenGet(t *testing.T) {
	var (
		expectedTokenID    = mockToken()
		expectedTokenAlias = "test-alias"
		apiPath            = fmt.Sprintf("AgentAccessTokens/%s", expectedTokenID)
		mockedAccessToken  = singleAgentAccessToken(expectedTokenID, expectedTokenAlias, "")
		fakeServer         = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateV2Response(mockedAccessToken))
			}
		},
	)

	fakeServer.MockAPI("AgentAccessTokens/UNKNOWN_TOKEN",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when agent access token exists", func(t *testing.T) {
		response, err := c.V2.AgentAccessTokens.Get(expectedTokenID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, expectedTokenID, response.Data.AccessToken)
			assert.Equal(t, expectedTokenAlias, response.Data.TokenAlias)
			assert.Equal(t, 1, response.Data.Enabled)
		}
	})

	t.Run("when agent access token does NOT exist", func(t *testing.T) {
		response, err := c.V2.AgentAccessTokens.Get("UNKNOWN_TOKEN")
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/AgentAccessTokens/UNKNOWN_TOKEN")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestAgentAccessTokenList(t *testing.T) {
	var (
		expectedAgentAccTokens = []string{mockToken(), mockToken(), mockToken()}
		expectedLen            = len(expectedAgentAccTokens)
		fakeServer             = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("AgentAccessTokens",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			fmt.Fprintf(w,
				generateV2ArrayResponse(
					generateAgentAccessTokens(expectedAgentAccTokens),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AgentAccessTokens.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, expectedAgentAccTokens, d.AccessToken)
	}
}

func generateAgentAccessTokens(ids []string) string {
	agentAccessTokens := make([]string, len(ids))
	for i, id := range ids {
		agentAccessTokens[i] = singleAgentAccessToken(id, fmt.Sprintf("test-%d", i), "")
	}
	return strings.Join(agentAccessTokens, ", ")
}

func generateV2ArrayResponse(data string) string {
	return `{ "data": [` + data + `] }`
}

func generateV2Response(data string) string {
	return `{ "data": ` + data + ` }`
}

func singleAgentAccessToken(token, alias, desc string) string {
	if token == "" {
		return "{}"
	}

	return `
  {
    "accessToken": "` + token + `",
    "createdTime": "2017-06-01T16:25:07.928Z",
    "props": {
      "description": "` + desc + `",
      "createdTime": "2017-06-01T16:25:07.928Z"
    },
    "tokenAlias": "` + alias + `",
    "tokenEnabled": 1,
    "version": "0.1"
  }
`
}

// mockToken generates a mocked agent access token for test purposes
func mockToken() string {
	now := time.Now().UTC().UnixNano()
	seed := rand.New(rand.NewSource(now))
	return strconv.FormatInt(seed.Int63(), 16)
}
