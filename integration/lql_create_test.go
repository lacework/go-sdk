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
	"fmt"
	"io/ioutil"
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
	assert.Contains(t, out.String(), "Type a query to create")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to create query: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestQueryCreateFile(t *testing.T) {
	// get temp file
	file, err := ioutil.TempFile("", "TestCreateFile")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// write-to and close file
	query := fmt.Sprintf(queryJSONTemplate, evaluatorID, queryID, queryText)
	_, err = file.Write([]byte(query))
	if err != nil {
		t.FailNow()
	}
	file.Close()

	// create
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("query", "create", "-f", file.Name())
	assert.Contains(t, out.String(), fmt.Sprintf("Query (%s) created successfully.", queryID))
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
	assert.Contains(t, out.String(), fmt.Sprintf("Query (%s) deleted successfully.", queryID))
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestQueryCreateURL(t *testing.T) {
	// This is tested by virtue of setup in other tests
	return
}
