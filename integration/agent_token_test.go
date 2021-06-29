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

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testTokenAlias = "codefresh-int-test-token"
	testTokenDesc  = "this token is used for our ci/cd tests (do-not-update)"
	// since we are wrapping the table output, we need to check for different strings
	testTokenDescWrap1 = "this token is used for our"
	testTokenDescWrap2 = "ci/cd tests (do-not-update)"
)

func TestAgentTokenCommandAliases(t *testing.T) {
	// lacework agent token
	out, err, exitcode := LaceworkCLI("help", "agent", "token")
	assert.Contains(t, out.String(), "lacework agent token [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework agent tokens
	out, err, exitcode = LaceworkCLI("help", "agent", "tokens")
	assert.Contains(t, out.String(), "lacework agent token [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAgentTokenCommandList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("agent", "token", "list")
	assert.Contains(t, out.String(), "TOKEN",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATE",
		"STDOUT table headers changed, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestAgentTokenCommandEndToEnd(t *testing.T) {
	var (
		out         bytes.Buffer
		err         bytes.Buffer
		exitcode    int
		agentTokens []agentToken
		// @afiune last_updated_time doesn't seem to exist in API v2 ....
		//tokenUpdatedTime time.Time
	)
	t.Run("list agent tokens", func(t *testing.T) {
		out, err, exitcode = LaceworkCLIWithTOMLConfig("agent", "token", "list", "--json")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
	})

	t.Run("inspecting json output", func(t *testing.T) {
		errJson := json.Unmarshal(out.Bytes(), &agentTokens)
		assert.Nil(t, errJson)
		assert.NotEmpty(t, agentTokens, "check JSON token list response")
	})

	var tokenID string
	for _, token := range agentTokens {
		if token.TokenAlias == testTokenAlias {
			tokenID = token.AccessToken
		}
	}

	if tokenID == "" {
		panic(fmt.Sprintf("something happened to the test token '%s'. :sadpanda:", testTokenAlias))
	}

	t.Run(fmt.Sprintf("token show command using alias '%s'", testTokenAlias), func(t *testing.T) {
		out, err, exitcode = LaceworkCLIWithTOMLConfig("agent", "token", "show", tokenID)
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")

		expectedOutput := []string{
			// headers
			"AGENT ACCESS TOKEN DETAILS",
			"NAME",
			"DESCRIPTION",
			"VERSION",
			"STATE",
			"CREATED AT",
			// @afiune last_updated_time doesn't seem to exist in API v2 ....
			//"UPDATED AT",
		}
		for _, str := range expectedOutput {
			assert.Contains(t, out.String(), str,
				"STDOUT table does not contain the '"+str+"' output")
		}
	})

	t.Run("storing last update time", func(t *testing.T) {
		out, err, exitcode = LaceworkCLIWithTOMLConfig("agent", "token", "show", tokenID, "--json")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")

		var token agentToken
		errJson := json.Unmarshal(out.Bytes(), &token)
		assert.Nil(t, errJson)
		assert.NotEmpty(t, token, "check JSON token show response while parsing update time")
		//tokenUpdatedTime = token.LastUpdatedTime
	})

	t.Run("token update: disable", func(t *testing.T) {
		out, err, exitcode = LaceworkCLIWithTOMLConfig(
			"agent", "token", "update", tokenID, "--disable")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Contains(t, out.String(), "Disabled",
			"STDOUT token should be disabled")
	})

	t.Run("token update: enable", func(t *testing.T) {
		out, err, exitcode = LaceworkCLIWithTOMLConfig(
			"agent", "token", "update", tokenID, "--enable")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Contains(t, out.String(), "Enabled",
			"STDOUT token should be enabled")
	})

	t.Run("token update: description", func(t *testing.T) {
		newDescription := "this shall not be seen! :)"
		out, err, exitcode = LaceworkCLIWithTOMLConfig(
			"agent", "token", "update", tokenID, "--description", newDescription)
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Contains(t, out.String(), newDescription,
			"STDOUT token description mismatch")

		out, err, exitcode = LaceworkCLIWithTOMLConfig(
			"agent", "token", "update", tokenID, "--description", testTokenDesc)
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Contains(t, out.String(), testTokenDescWrap1,
			"STDOUT token description mismatch")
		assert.Contains(t, out.String(), testTokenDescWrap2,
			"STDOUT token description mismatch")
	})

	t.Run("verify token was updated", func(t *testing.T) {
		out, err, exitcode = LaceworkCLIWithTOMLConfig("agent", "token", "show", tokenID, "--json")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")

		var token agentToken
		errJson := json.Unmarshal(out.Bytes(), &token)
		assert.Nil(t, errJson)
		assert.NotEmpty(t, token, "check JSON token show response while parsing update time")
		//assert.True(t,
		//tokenUpdatedTime.Before(token.LastUpdatedTime),
		//fmt.Sprintf("the agent access token last_update_time was NOT updated! check please. (old:%s) (new:%s)",
		//tokenUpdatedTime, token.LastUpdatedTime,
		//),
		//)
	})
}

type agentToken struct {
	AccessToken string          `json:"accessToken"`
	Props       agentTokenProps `json:"props,omitempty"`
	TokenAlias  string          `json:"tokenAlias"`
	Enabled     int             `json:"tokenEnabled"`
	Version     string          `json:"version"`
}

type agentTokenProps struct {
	CreatedTime time.Time `json:"createdTime,omitempty"`
	Description string    `json:"description,omitempty"`
}
