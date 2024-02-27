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

func TestQueryCreateHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "query", "create")
	assert.Contains(t, out.String(), "lacework query create [flags]")
	assert.Contains(t, out.String(), "-f, --file string")
	assert.Contains(t, out.String(), "-u, --url string")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryCreateEditor(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("query", "create")
	assert.Contains(t, out.String(), "Choose query language")
	assert.Contains(t, err.String(), "ERROR unable to create query:")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryCreateFile(t *testing.T) {
	// get temp file
	file, err := createTemporaryFile(
		"TestQueryCreateFile",
		fmt.Sprintf(queryJSONTemplate, queryID, queryText),
	)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// create
	out, stderr, exitcode := cleanAndCreateQuery(queryID, "query", "create", "-f", file.Name())
	assert.Contains(t, out.String(), fmt.Sprintf("The query %s was created.", queryID))
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// show
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("query", "show", queryID)
	assert.Contains(t, out.String(), "INSERT_ID")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("query", "show", queryID, "--json")
	assert.Contains(t, out.String(), "INSERT_ID")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// delete
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("query", "delete", queryID)
	assert.Contains(t, out.String(), fmt.Sprintf("The query %s was deleted.", queryID))
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
