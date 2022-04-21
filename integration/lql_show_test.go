//go:build query

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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryShowHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "query", "show")
	assert.Contains(t, out.String(), "lacework query show <query_id> [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryShowNoInput(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "show")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryShowYAML(t *testing.T) {
	queryIdFromPlatform := "LW_Global_AWS_CTA_AccessKeyDeleted"
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "show", queryIdFromPlatform, "--yaml")
	assert.Contains(t, out.String(), fmt.Sprintf("queryId: %s", queryIdFromPlatform))
	assert.Contains(t, out.String(), "queryText: |-")
	assert.Contains(t, out.String(), "INSERT_ID")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
