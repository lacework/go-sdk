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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryUpdateEditor(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "update")
	assert.Contains(t, out.String(), "Type a query to update")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to update query:")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryUpdateFile(t *testing.T) {
	// get temp file
	file, err := createTemporaryFile(
		"TestQueryUpdateFile",
		fmt.Sprintf(queryJSONTemplate, queryID, queryUpdateText),
	)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryID)

	// update
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("query", "update", "-f", file.Name())
	assert.Contains(t, out.String(), fmt.Sprintf("The query %s was updated.", queryID))
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("query", "show", queryID)
	assert.Contains(t, out.String(), "INSERT_TIME")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryUpdateURL(t *testing.T) {
	// get temp file
	file, err := createTemporaryFile(
		"TestQueryUpdateFile",
		fmt.Sprintf(queryJSONTemplate, queryID, queryUpdateText),
	)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryID)

	// update
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("query", "update", "-f", file.Name())
	assert.Contains(t, out.String(), fmt.Sprintf("The query %s was updated.", queryID))
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("query", "show", queryID)
	assert.Contains(t, out.String(), "INSERT_TIME")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryUpdateFromIDNotFound(t *testing.T) {
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("query", "update", "ID_NOT_FOUND", "--noninteractive")

	assert.Empty(t, out.String(), "STDOUT should be empty") // added --noninteractive to avoid polluting STDOUT
	assert.Contains(t, stderr.String(), "unable to load query from your account")
	assert.Contains(t, stderr.String(), "/api/v2/Queries/ID_NOT_FOUND")
	assert.Contains(t, stderr.String(), "Query id ID_NOT_FOUND not found")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}
