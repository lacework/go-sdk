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

	"github.com/stretchr/testify/assert"
)

const (
	policyText string = `
{
	"policy_id": "my-policy-2",
	"title": "My Policy Title",
	"enabled": false,
	"lql_id": "MyLQL",
	"severity": "high",
	"description": "My Policy Description",
	"remediation": "Check yourself..."
}`
	policyURL string = "https://raw.githubusercontent.com/lacework/go-sdk/main/integration/test_resources/policy/my-policy-1.json"
)

func TestPolicyAliases(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// lacework policy
	out, err, exitcode := LaceworkCLI("help", "policy")
	assert.Contains(t, out.String(), "lacework policy [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyBase(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLI("policy")
	assert.Contains(t, out.String(), "create")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateHelp(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLI("help", "policy", "create")
	assert.Contains(t, out.String(), "lacework policy create [flags]")
	assert.Contains(t, out.String(), "-f, --file string")
	assert.Contains(t, out.String(), "-u, --url string")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateEditor(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// create
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "create")
	assert.Contains(t, out.String(), "Type a policy to create")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to create policy: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateFile(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", lqlQueryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", lqlQueryID)

	// get temp file
	file, err := ioutil.TempFile("", "TestPolicyCreateFile")
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// write-to and close file
	_, err = file.Write([]byte(policyText))
	if err != nil {
		t.FailNow()
	}
	file.Close()

	// create (output json)
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("policy", "create", "-f", file.Name(), "--json")
	assert.Contains(t, out.String(), `"policy_id"`)
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateURL(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", lqlQueryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", lqlQueryID)

	// create (output human)
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "create", "-u", policyURL)
	assert.Contains(t, out.String(), "Policy (my-policy-1) created successfully.")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
