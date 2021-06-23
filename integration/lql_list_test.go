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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryListHelp(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLI("help", "query", "list")
	assert.Contains(t, out.String(), "lacework query list [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryList(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", lqlQueryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", lqlQueryID)

	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "list")
	assert.Contains(t, out.String(), "QUERY ID")
	assert.Contains(t, out.String(), "LW_CLI_AWS_CTA_IntegrationTest")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	out, err, exitcode = LaceworkCLIWithTOMLConfig("query", "list", "--json")
	assert.Contains(t, out.String(), `"LW_CLI_AWS_CTA_IntegrationTest"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
