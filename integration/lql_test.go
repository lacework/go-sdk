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
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	lqlQueryID     string = "MyLQL"
	lqlQueryText   string = "MyLQL(CloudTrailRawEvents e) {SELECT INSERT_ID LIMIT 1}"
	lqlQueryUpdate string = "MyLQL(CloudTrailRawEvents e) {SELECT INSERT_ID, INSERT_TIME LIMIT 1}"
	lqlQueryURL    string = "https://raw.githubusercontent.com/lacework/go-sdk/main/integration/test_resources/lql/my.lql"
)

var (
	d, _          = time.ParseDuration("-24h")
	lqlQueryStart = time.Now().Add(d).Format(time.RFC3339)
	lqlQueryEnd   = time.Now().Format(time.RFC3339)
)

func TestQueryAliases(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// lacework query
	out, err, exitcode := LaceworkCLI("help", "query")
	assert.Contains(t, out.String(), "lacework query [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework lql
	out, err, exitcode = LaceworkCLI("help", "lql")
	assert.Contains(t, out.String(), "lacework query [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryBase(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLI("query")
	assert.Contains(t, out.String(), "create")
	assert.Contains(t, out.String(), "delete")
	assert.Contains(t, out.String(), "describe")
	assert.Contains(t, out.String(), "list")
	assert.Contains(t, out.String(), "list-sources")
	assert.Contains(t, out.String(), "run")
	assert.Contains(t, out.String(), "show")
	assert.Contains(t, out.String(), "update")
	assert.Contains(t, out.String(), "validate")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunHelp(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLI("help", "query", "run")
	assert.Contains(t, out.String(), "lacework query run [query_id] [flags]")
	assert.Contains(t, out.String(), "--end string")
	assert.Contains(t, out.String(), `end time for LQL query (default "now")`)
	assert.Contains(t, out.String(), "-f, --file string")
	assert.Contains(t, out.String(), "--start string")
	assert.Contains(t, out.String(), `start time for LQL query (default "@d")`)
	assert.Contains(t, out.String(), "--range string")
	assert.Contains(t, out.String(), "-u, --url string")
	assert.Contains(t, out.String(), "--validate_only")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunEditor(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	// run
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "run")
	assert.Contains(t, out.String(), "Type a query to run")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to run LQL query: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")

	// validate_only
	out, err, exitcode = LaceworkCLIWithTOMLConfig("query", "run", "--validate_only")
	assert.Contains(t, out.String(), "Type a query to validate")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to run LQL query: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunID(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", lqlQueryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", lqlQueryID)

	// run (default times)
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "run", lqlQueryID)
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// validate_only
	out, err, exitcode = LaceworkCLIWithTOMLConfig("query", "run", lqlQueryID, "--validate_only")
	assert.Contains(t, out.String(), "LQL query validated successfully.")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFileJSONCrumb(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// get temp file
	file, err := ioutil.TempFile("", "TestQueryRunFile")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// write-to and close file
	_, err = file.Write([]byte("{"))
	if err != nil {
		t.FailNow()
	}
	file.Close()

	// run
	_, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-f", file.Name(), "--start", lqlQueryStart, "--end", lqlQueryEnd)
	assert.Contains(t, stderr.String(), "LQL query in JSON format")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFilePlainCrumb(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// get temp file
	file, err := ioutil.TempFile("", "TestQueryRunFile")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// write-to and close file
	_, err = file.Write([]byte("tigerking"))
	if err != nil {
		t.FailNow()
	}
	file.Close()

	// run
	_, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-f", file.Name(), "--start", lqlQueryStart, "--end", lqlQueryEnd)
	assert.Contains(t, stderr.String(), "LQL query in plain text format")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFile(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// get temp file
	file, err := ioutil.TempFile("", "TestQueryRunFile")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// write-to and close file
	_, err = file.Write([]byte(lqlQueryText))
	if err != nil {
		t.FailNow()
	}
	file.Close()

	// run (explicit times)
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-f", file.Name(), "--start", lqlQueryStart, "--end", lqlQueryEnd)
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// validate_only
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig(
		"query", "run", "-f", file.Name(), "--validate_only")
	assert.Contains(t, out.String(), "LQL query validated successfully.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunURL(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// run (natural time)
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", lqlQueryURL, "--range", "last week")
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// validate_only
	out, err, exitcode = LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", lqlQueryURL, "--validate_only")
	assert.Contains(t, out.String(), "LQL query validated successfully.")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFailOnBadInput(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	_, err, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", lqlQueryURL, "--fail_on_count", "!30")
	assert.Contains(t, err.String(), "ERROR count operation (!30) is invalid")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFailOnPos(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", lqlQueryURL, "--fail_on_count", "=1")
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 9, exitcode, "EXITCODE is not the expected one")
}

func TestQueryRunFailOnNeg(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "run", "-u", lqlQueryURL, "--fail_on_count", ">1")
	assert.Contains(t, out.String(), `"INSERT_ID"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
