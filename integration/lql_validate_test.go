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
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryValidateHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "query", "validate")
	assert.Contains(t, out.String(), "lacework query validate [flags]")
	assert.Contains(t, out.String(), "-f, --file string")
	assert.Contains(t, out.String(), "-u, --url string")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryValidateEditor(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "validate")
	assert.Contains(t, out.String(), "Type a query to validate")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to validate query: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryValidateFile(t *testing.T) {
	// get temp file
	file, err := createTemporaryFile(
		"TestQueryValidateFile",
		fmt.Sprintf(queryJSONTemplate, evaluatorID, queryID, queryText),
	)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// validate
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("query", "validate", "-f", file.Name())
	assert.Contains(t, out.String(), "Query validated successfully.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryValidateHost(t *testing.T) {
	// get temp file
	file, err := createTemporaryFile(
		"TestQueryValidateHost",
		fmt.Sprintf(queryJSONTemplate, "", queryHostID, queryHostText))
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// validate_only
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"query", "validate", "-f", file.Name())
	assert.Contains(t, out.String(), "Query validated successfully.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryValidateURL(t *testing.T) {
	// validate
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "validate", "-u", queryURL)
	assert.Contains(t, out.String(), "Query validated successfully.")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryValidateStdin(t *testing.T) {
	var out, stderr bytes.Buffer

	// get CLI
	dir := createTOMLConfigFromCIvars()
	defer os.RemoveAll(dir)

	// get temp file
	file, err := createTemporaryFile(
		"TestQueryValidateFile",
		fmt.Sprintf(queryJSONTemplate, evaluatorID, queryID, queryText),
	)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// open file by name
	stdin, err := os.Open(file.Name())
	if err != nil {
		t.FailNow()
	}

	cmd := NewLaceworkCLI(dir, stdin, "query", "validate")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	exitcode, err := runLaceworkCLIFromCmd(cmd)
	if exitcode == 999 {
		fmt.Println(stderr)
		if _, err := stderr.WriteString(err.Error()); err != nil {
			// @afiune we should never get here but if we do, lets print the error
			fmt.Println(err)
		}
	}

	assert.Contains(t, out.String(), "Query validated successfully.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
