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
)

func TestAgentTokenPrettyStatus(t *testing.T) {
	subject := api.AgentToken{Enabled: "true"}
	assert.Equal(t, "Enabled", subject.PrettyStatus())

	subject.Enabled = "false"
	assert.Equal(t, "Disabled", subject.PrettyStatus())
	subject.Enabled = "anything else"
	assert.Equal(t, "Disabled", subject.PrettyStatus())
}

func TestAgentTokenStatus(t *testing.T) {
	subject := api.AgentToken{Enabled: "true"}
	assert.True(t, subject.Status())

	subject.Enabled = "false"
	assert.False(t, subject.Status())
	subject.Enabled = "anything else"
	assert.False(t, subject.Status())
}

func TestAgentTokenEnabledInt(t *testing.T) {
	subject := api.AgentToken{Enabled: "true"}
	assert.Equal(t, 1, subject.EnabledInt())

	subject.Enabled = "false"
	assert.Equal(t, 0, subject.EnabledInt())
	subject.Enabled = "anything else"
	assert.Equal(t, 0, subject.EnabledInt())
}
