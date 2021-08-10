// +build query

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryDescribeHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "query", "describe")
	assert.Contains(t, out.String(), "lacework query describe <data_source> [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryDescribeNoInput(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "describe")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryDescribeTable(t *testing.T) {
	t.Skip("skipping test due to unavailable api")

	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "describe", "CloudTrailRawEvents")
	assert.Contains(t, out.String(), "FIELD NAME")
	assert.Contains(t, out.String(), "INSERT_ID")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryDescribeJSON(t *testing.T) {
	t.Skip("skipping test due to unavailable api")

	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "describe", "CloudTrailRawEvents", "--json")
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
