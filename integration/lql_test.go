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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	d, _       = time.ParseDuration("-24h")
	queryStart = time.Now().Add(d).Format(time.RFC3339)
	queryEnd   = time.Now().Format(time.RFC3339)
)

func cleanAndCreateQuery(id string, args ...string) (bytes.Buffer, bytes.Buffer, int) {
	LaceworkCLIWithTOMLConfig("query", "delete", id)
	return LaceworkCLIWithTOMLConfig(args...)
}

func TestQueryBase(t *testing.T) {
	out, err, exitcode := LaceworkCLI("query")
	assert.Contains(t, out.String(), "create")
	assert.Contains(t, out.String(), "delete")
	assert.Contains(t, out.String(), "list")
	assert.Contains(t, out.String(), "list-sources")
	assert.Contains(t, out.String(), "preview-source")
	assert.Contains(t, out.String(), "run")
	assert.Contains(t, out.String(), "show")
	assert.Contains(t, out.String(), "show-source")
	assert.Contains(t, out.String(), "update")
	assert.Contains(t, out.String(), "validate")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunEditor(t *testing.T) {
	// run
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "run")
	assert.Contains(t, out.String(), "Type a query to run")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to run query:")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")

	// validate_only
	out, err, exitcode = LaceworkCLIWithTOMLConfig("query", "run", "--validate_only")
	assert.Contains(t, out.String(), "Type a query to validate")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to validate query:")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunID(t *testing.T) {
	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryID)

	// run (default times)
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "run", queryID)
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// bad - file
	_, err, exitcode = LaceworkCLIWithTOMLConfig("query", "run", queryID, "--file", "foo")
	assert.Contains(t, err.String(), "ERROR flag --file not applicable when specifying query_id argument")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")

	// bad - url
	_, err, exitcode = LaceworkCLIWithTOMLConfig("query", "run", queryID, "--url", "foo")
	assert.Contains(t, err.String(), "ERROR flag --url not applicable when specifying query_id argument")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")

	// bad - validate_only
	_, err, exitcode = LaceworkCLIWithTOMLConfig("query", "run", queryID, "--validate_only")
	assert.Contains(t, err.String(), "ERROR flag --validate_only not applicable when specifying query_id argument")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFileJSONCrumb(t *testing.T) {
	// get temp file
	file, err := createTemporaryFile("TestQueryRunFile", "{")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// run
	_, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-f", file.Name(), "--start", queryStart, "--end", queryEnd)
	assert.Contains(t, stderr.String(), "query in JSON format")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFileYAMLCrumb(t *testing.T) {
	// get temp file
	file, err := createTemporaryFile("TestQueryRunFile", "tigerking")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// run
	_, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-f", file.Name(), "--start", queryStart, "--end", queryEnd)
	assert.Contains(t, stderr.String(), "query in YAML format")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFile(t *testing.T) {
	// get temp file
	file, err := createTemporaryFile(
		"TestQueryRunFile",
		fmt.Sprintf(queryJSONTemplate, queryID, queryText))
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// run (bad limit)
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-f", file.Name(), "--start", queryStart, "--end", queryEnd, "--limit", "-1")
	assert.Contains(t, stderr.String(), "limit must be at least 1")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")

	// run (explicit times / options)
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig(
		"query", "run", "-f", file.Name(), "--start", queryStart, "--end", queryEnd, "--limit", "1", "--json")

	// check limit
	var results []interface{}
	err = json.Unmarshal(out.Bytes(), &results)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, 1, len(results))

	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// validate_only
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig(
		"query", "run", "-f", file.Name(), "--validate_only")
	assert.Contains(t, out.String(), "Query validated successfully.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunURL(t *testing.T) {
	// run (natural time)
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", queryURL, "--range", "last week")
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// validate_only
	out, err, exitcode = LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", queryURL, "--validate_only")
	assert.Contains(t, out.String(), "Query validated successfully.")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFailOnBadInput(t *testing.T) {
	_, err, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", queryURL, "--fail_on_count", "!30")
	assert.Contains(t, err.String(), "ERROR count operation (!30) is invalid")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFailOnPos(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", queryURL, "--fail_on_count", ">1")
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Contains(t, err.String(), "ERROR (FAIL-ON): query matched fail_on_count expression")
	assert.Equal(t, 9, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFailOnNeg(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", queryURL, "--fail_on_count", "=0")
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
