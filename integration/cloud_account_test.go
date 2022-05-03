//go:build cloudAccount

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

package cloudAccount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudAccountCommandAliases(t *testing.T) {
	// lacework cloud-account
	out, err, exitcode := LaceworkCLI("help", "cloud-account")
	assert.Contains(t, out.String(), "lacework cloud-account [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework cloud-accounts
	out, err, exitcode = LaceworkCLI("help", "cloud-accounts")
	assert.Contains(t, out.String(), "lacework cloud-accounts [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework cloud
	out, err, exitcode = LaceworkCLI("help", "cloud")
	assert.Contains(t, out.String(), "lacework cloud [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework ca
	out, err, exitcode = LaceworkCLI("help", "ca")
	assert.Contains(t, out.String(), "lacework ca [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func _TestCloudAccountCommandList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("cloud-account", "list")
	assert.Contains(t, out.String(), "CLOUD ACCOUNT INTEGRATION GUID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATE",
		"STDOUT table headers changed, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func _TestCloudAccountCommandListWithTypeFlag(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("cloud-account", "list", "--type", "AWS_CFG")
	assert.Contains(t, out.String(), "CLOUD ACCOUNT INTEGRATION GUID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")

	// TODO @afiune lets try to create an environment where we can be 100% sure that
	// cloud account integrations will exist and assert against it

	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func _TestCloudAccountCommandListWithTypeFlagErrorUnknownType(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("cloud-account", "list", "--type", "FOO_BAR")
	assert.Emptyf(t, out.String(),
		"STDOUT should be empty")
	assert.Contains(t, err.String(),
		"ERROR unknown cloud account type 'FOO_BAR'",
		"STDERR should contain the unknown type")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}
