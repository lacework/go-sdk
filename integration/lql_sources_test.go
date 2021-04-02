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

func TestQuerySourcesAliases(t *testing.T) {
	if env := os.Getenv("CI_BETA"); env == "" {
		t.Skip("skipping test in production mode")
	}
	// lacework query sources
	out, err, exitcode := LaceworkCLI("help", "query", "sources")
	assert.Contains(t, out.String(), "lacework query list-sources [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework query list-sources
	out, err, exitcode = LaceworkCLI("help", "query", "list-sources")
	assert.Contains(t, out.String(), "lacework query list-sources [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQuerySourcesTable(t *testing.T) {
	if env := os.Getenv("CI_BETA"); env == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "sources")
	assert.Contains(t, out.String(), "DATA SOURCE")
	assert.Contains(t, out.String(), "CloudTrailRawEvents")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQuerySourcesJSON(t *testing.T) {
	if env := os.Getenv("CI_BETA"); env == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "sources", "--json")
	assert.Contains(t, out.String(), "[")
	assert.Contains(t, out.String(), `"CloudTrailRawEvents"`)
	assert.Contains(t, out.String(), "]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
