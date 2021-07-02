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
	policyText string = `{
	"policy_id": "lacework-clitest-1",
	"title": "My Policy Title",
	"enabled": false,
	"alert_enabled": false,
	"lql_id": "LW_CLI_AWS_CTA_IntegrationTest",
	"severity": "low",
	"description": "My Policy Description",
	"remediation": "Check yourself..."
}`
	policyURL string = "https://raw.githubusercontent.com/lacework/go-sdk/main/integration/test_resources/policy/lacework-clitest-1.json"
)

func TestPolicyHelp(t *testing.T) {
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
	assert.Contains(t, out.String(), "list")
	assert.Contains(t, out.String(), "show")
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
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryID)

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

	// update-url (output human)
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("policy", "update", "-u", policyURL)
	assert.Contains(t, out.String(), "Policy (lacework-clitest-1) updated successfully.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// list enabled-only
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("policy", "list", "--enabled")
	assert.Contains(t, out.String(), "lacework-global-1")
	assert.NotContains(t, out.String(), "lacework-clitest-1")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// list alert_enabled-only
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("policy", "list", "--alert_enabled")
	assert.Contains(t, out.String(), "lacework-global-1")
	assert.NotContains(t, out.String(), "lacework-clitest-1")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// delete
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("policy", "delete", "lacework-clitest-1")
	assert.Contains(t, out.String(), "Policy (lacework-clitest-1) deleted successfully.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateURL(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryID)

	// create (output human)
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("policy", "create", "-u", policyURL)
	// teardown policy
	defer LaceworkCLIWithTOMLConfig("policy", "delete", "lacework-clitest-1")

	assert.Contains(t, out.String(), "Policy (lacework-clitest-1) created successfully.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// update-file (output json)
	// get temp file
	file, err := ioutil.TempFile("", "TestPolicyUpdateFile")
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

	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("policy", "update", "-f", file.Name(), "--json")
	assert.Contains(t, out.String(), `"policy_id"`)
	assert.Contains(t, out.String(), `"low"`)
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyListHelp(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLI("help", "policy", "list")
	assert.Contains(t, out.String(), "lacework policy list [flags]")
	assert.Contains(t, out.String(), "--alert_enabled     only show alert_enabled policies")
	assert.Contains(t, out.String(), "--enabled           only show enabled policies")
	assert.Contains(t, out.String(), "--severity string   filter policies by severity threshold (critical, high, medium, low, info)")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyList(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	// list (output human)
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "list")
	assert.Contains(t, out.String(), "POLICY ID")
	assert.Contains(t, out.String(), "lacework-global-1")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// list (output json)
	out, err, exitcode = LaceworkCLIWithTOMLConfig("policy", "list", "--json")
	assert.Contains(t, out.String(), `"policy_id"`)
	assert.Contains(t, out.String(), `"lacework-global-1"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyBadSeverity(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	_, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "list", "--severity", "superhigh")
	assert.Contains(t, err.String(), "the severity superhigh is not valid, use one of critical, high, medium, low, info")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")

}

func TestPolicySeverityCritical(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "list", "--severity", "critical")
	assert.Contains(t, out.String(), "lacework-global-8")
	assert.NotContains(t, out.String(), "lacework-global-1")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyShowHelp(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	out, err, exitcode := LaceworkCLI("help", "policy", "show")
	assert.Contains(t, out.String(), "lacework policy show <policy_id> [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyShow(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	// show (output)
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "show", "lacework-global-1")
	assert.Contains(t, out.String(), "POLICY ID")
	assert.Contains(t, out.String(), "lacework-global-1")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// show (output json)
	out, err, exitcode = LaceworkCLIWithTOMLConfig("policy", "show", "lacework-global-1", "--json")
	assert.Contains(t, out.String(), `"policy_id"`)
	assert.Contains(t, out.String(), `"lacework-global-1"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyUpdateHelp(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLI("help", "policy", "update")
	assert.Contains(t, out.String(), "lacework policy update [policy_id] [flags]")
	assert.Contains(t, out.String(), "-f, --file string")
	assert.Contains(t, out.String(), "-u, --url string")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyUpdateEditor(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// update
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "update")

	assert.Contains(t, out.String(), "Type a policy to update")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to update policy: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyUpdateFile(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// update file tested by virtue of TestPolicyCreateURL
}

func TestPolicyUpdateURL(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	// update URL tested by virtue of TestPolicyCreateFile
}

func TestPolicyDeleteHelp(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}
	out, err, exitcode := LaceworkCLI("help", "policy", "delete")
	assert.Contains(t, out.String(), "lacework policy delete <policy_id> [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyDeleteNoInput(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "delete")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyDelete(t *testing.T) {
	if os.Getenv("CI_BETA") == "" {
		t.Skip("skipping test in production mode")
	}

	// setup query
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)
	// teardown query
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryID)
	// setup policy
	LaceworkCLIWithTOMLConfig("policy", "create", "-u", policyURL)

	// human delete tested by virtue of TestPolicyCreateFile

	// json
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "delete", "lacework-clitest-1", "--json")
	assert.Contains(t, out.String(), `"lacework-clitest-1"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
