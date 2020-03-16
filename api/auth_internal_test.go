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

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiPath(t *testing.T) {
	c1 := &Client{apiVersion: "v1"}
	assert.Equal(t, "/api/v1/foo", c1.apiPath("foo"), "api path mismatch")
	assert.Equal(t,
		"/api/v1/access/tokens",
		c1.apiPath(apiTokens),
		"token api path mismatch")

	c2 := &Client{apiVersion: "v2"}
	assert.Equal(t, "/api/v2/bar", c2.apiPath("bar"), "api path mismatch")
	assert.Equal(t,
		"/api/v2/external/integrations",
		c2.apiPath(apiIntegrations),
		"integrations api path mismatch")
}
