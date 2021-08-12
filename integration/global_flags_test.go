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
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalFlagToken(t *testing.T) {
	// generating a token with toml config
	token, err, exitcode := LaceworkCLIWithTOMLConfig("access-token")
	assert.Contains(t, token.String(), "_", // @afiune tokens start with "_secret123"
		"STDOUT table headers changed, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	if assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one") {

		// running the Lacework CLI without toml config but with the token flag
		out, err, exitcode := LaceworkCLI("int", "list",
			"--token", strings.Trim(token.String(), "\n"), "--account", os.Getenv("CI_ACCOUNT"))
		assert.Contains(t, out.String(), "INTEGRATION GUID")
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	}
}
