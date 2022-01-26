//go:build policy

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
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	newPolicyYAML string = `---
evaluatorId: Cloudtrail
policyId: clitest-1
policyType: Violation
queryId: LW_CLI_AWS_CTA_IntegrationTest
title: My Policy Title
enabled: false
description: My Policy Description
remediation: Check yourself...
severity: high
alertEnabled: false
alertProfile: LW_CloudTrail_Alerts
`
	newHostPolicyYAML string = `---
evaluatorId:
policyType: Violation
queryId: LW_CLI_Host_Files_IntegrationTest
title: My Policy Title
enabled: false
description: My Policy Description
remediation: Check yourself...
severity: high
alertEnabled: false
alertProfile: LW_HE_FILES_DEFAULT_PROFILE.HE_File_NewViolation
`
	// nested
	updatePolicyYAML string = `---
policies:
  - severity: low
`
	policyURL string = "https://raw.githubusercontent.com/lacework/go-sdk/main/integration/test_resources/policy/account-clitest-1.json"
)

var (
	policyIDRE *regexp.Regexp = regexp.MustCompile(
		`([\w-]+-(?:(?:cli.*?test-1)|(?:default-\d+)))`)
)

func getPolicyIdFromStdout(s string) (string, error) {
	if match := policyIDRE.FindStringSubmatch(s); match != nil {
		return match[1], nil
	}
	return "", errors.New("unable to extract policy identifier")
}

func TestPolicyHelp(t *testing.T) {
	// lacework policy
	out, err, exitcode := LaceworkCLI("help", "policy")
	assert.Contains(t, out.String(), "lacework policy [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyBase(t *testing.T) {
	out, err, exitcode := LaceworkCLI("policy")
	assert.Contains(t, out.String(), "create")
	assert.Contains(t, out.String(), "list")
	assert.Contains(t, out.String(), "show")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "policy", "create")
	assert.Contains(t, out.String(), "lacework policy create [flags]")
	assert.Contains(t, out.String(), "-f, --file string")
	assert.Contains(t, out.String(), "-u, --url string")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateEditor(t *testing.T) {
	// create
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "create")
	assert.Contains(t, out.String(), "Type a policy to create")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to create policy:")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateFile(t *testing.T) {
	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)

	// get temp file
	file, err := createTemporaryFile("TestPolicyCreateFile", newPolicyYAML)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// create (output json)
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("policy", "create", "-f", file.Name(), "--json")
	assert.Contains(t, out.String(), `"policyId"`)
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	policyID, err := getPolicyIdFromStdout(out.String())
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	// update-url (output human)
	// ideally we wouldn't specify a policyID here since it's in policyURL
	// but the policy API is completely jacked up with respect to persisting a policyID as code
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("policy", "update", policyID, "-u", policyURL)
	assert.Contains(t, out.String(),
		fmt.Sprintf("The policy %s was updated.", policyID))
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// list enabled-only
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("policy", "list", "--enabled")
	assert.Contains(t, out.String(), "lacework-global-1")
	assert.NotContains(t, out.String(), policyID)
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// list alert_enabled-only
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("policy", "list", "--alert_enabled")
	assert.Contains(t, out.String(), "lacework-global-1")
	assert.NotContains(t, out.String(), policyID)
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// force delete
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig("policy", "delete", policyID, "--force")
	assert.Contains(t, out.String(),
		fmt.Sprintf("The policy %s was deleted.", policyID))
	assert.Contains(t, out.String(),
		fmt.Sprintf("The query %s was deleted.", queryID))
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateURL(t *testing.T) {
	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryID)

	// create (output human)
	policyID := "$account-clitest-1"
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("policy", "create", "-u", policyURL)
	// teardown policy
	defer LaceworkCLIWithTOMLConfig("policy", "delete", policyID)

	assert.Contains(t, out.String(), "clitest-1 was created.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// update-file (output json)
	// get temp file
	file, err := createTemporaryFile("TestPolicyUpdateFile", updatePolicyYAML)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// this test intentionally (vs. being required to) specifies the policyID inline
	out, stderr, exitcode = LaceworkCLIWithTOMLConfig(
		"policy", "update", policyID, "-f", file.Name(), "--json") // specify id inline
	assert.Contains(t, out.String(), `"policyId"`)
	assert.Contains(t, out.String(), `"low"`)
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyCreateStdin(t *testing.T) {
	var out, stderr bytes.Buffer

	// get CLI
	dir := createTOMLConfigFromCIvars()
	defer os.RemoveAll(dir)

	// setup query
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)
	// teardown query
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryID)

	// get temp policy file
	file, err := createTemporaryFile("TestPolicyCreateFile", newPolicyYAML)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// open file by name
	stdin, err := os.Open(file.Name())
	if err != nil {
		t.FailNow()
	}

	// setup command
	cmd := NewLaceworkCLI(dir, stdin, "policy", "create")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// run command
	exitcode, err := runLaceworkCLIFromCmd(cmd)
	if exitcode == 999 {
		fmt.Println(stderr)
		if _, err := stderr.WriteString(err.Error()); err != nil {
			// @afiune we should never get here but if we do, lets print the error
			fmt.Println(err)
		}
	}
	assert.Contains(t, out.String(), "was created.")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	policyID, err := getPolicyIdFromStdout(out.String())
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer LaceworkCLIWithTOMLConfig("policy", "delete", policyID)
}

func TestPolicyCreateHost(t *testing.T) {
	// setup
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryHostURL)
	// teardown
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryHostID)

	// get temp file
	file, err := createTemporaryFile("TestPolicyCreateHost", newHostPolicyYAML)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// create (output json)
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("policy", "create", "-f", file.Name(), "--json")
	assert.Contains(t, out.String(), `"policyId"`)
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	policyID, err := getPolicyIdFromStdout(out.String())
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer LaceworkCLIWithTOMLConfig("policy", "delete", policyID)
}

func TestPolicyListHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "policy", "list")
	assert.Contains(t, out.String(), "lacework policy list [flags]")
	assert.Contains(t, out.String(), "--alert_enabled     only show alert_enabled policies")
	assert.Contains(t, out.String(), "--enabled           only show enabled policies")
	assert.Contains(t, out.String(), "--severity string   filter policies by severity threshold (critical, high, medium, low, info)")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyList(t *testing.T) {
	// list (output human)
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "list")
	assert.Contains(t, out.String(), "POLICY ID")
	assert.Contains(t, out.String(), "lacework-global-1")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// list (output json)
	out, err, exitcode = LaceworkCLIWithTOMLConfig("policy", "list", "--json")
	assert.Contains(t, out.String(), `"policyId"`)
	assert.Contains(t, out.String(), `"lacework-global-1"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyBadSeverity(t *testing.T) {
	_, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "list", "--severity", "superhigh")
	assert.Contains(t, err.String(), "the severity superhigh is not valid, use one of critical, high, medium, low, info")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")

}

func TestPolicySeverityCritical(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "list", "--severity", "critical")
	assert.Contains(t, out.String(), "lacework-global-8")
	assert.NotContains(t, out.String(), "lacework-global-1")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyShowHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "policy", "show")
	assert.Contains(t, out.String(), "lacework policy show <policy_id> [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyShow(t *testing.T) {
	// show (output)
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "show", "lacework-global-1")
	assert.Contains(t, out.String(), "POLICY ID")
	assert.Contains(t, out.String(), "lacework-global-1")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// show (output json)
	out, err, exitcode = LaceworkCLIWithTOMLConfig("policy", "show", "lacework-global-1", "--json")
	assert.Contains(t, out.String(), `"policyId"`)
	assert.Contains(t, out.String(), `"lacework-global-1"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyUpdateHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "policy", "update")
	assert.Contains(t, out.String(), "lacework policy update [policy_id] [flags]")
	assert.Contains(t, out.String(), "-f, --file string")
	assert.Contains(t, out.String(), "-u, --url string")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyUpdateEditor(t *testing.T) {
	// update
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "update")

	assert.Contains(t, out.String(), "Type a policy to update")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "ERROR unable to update policy:")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyUpdateFile(t *testing.T) {
	// update file tested by virtue of TestPolicyCreateURL
}

func TestPolicyUpdateURL(t *testing.T) {
	// update URL tested by virtue of TestPolicyCreateFile
}

func TestPolicyDeleteHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "policy", "delete")
	assert.Contains(t, out.String(), "lacework policy delete <policy_id> [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyDeleteNoInput(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("policy", "delete")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyDelete(t *testing.T) {
	// setup query
	LaceworkCLIWithTOMLConfig("query", "create", "-u", queryURL)
	// teardown query
	defer LaceworkCLIWithTOMLConfig("query", "delete", queryID)

	// setup policy
	out, _, _ := LaceworkCLIWithTOMLConfig("policy", "create", "-u", policyURL)

	policyID, err := getPolicyIdFromStdout(out.String())
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	// human delete tested by virtue of TestPolicyCreateFile

	// json
	_, stderr, exitcode := LaceworkCLIWithTOMLConfig("policy", "delete", policyID, "--json")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
