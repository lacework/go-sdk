//go:build integration

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationCommandAliases(t *testing.T) {
	// lacework integration
	out, err, exitcode := LaceworkCLI("help", "integration")
	assert.Contains(t, out.String(), "lacework integration [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework integrations
	out, err, exitcode = LaceworkCLI("help", "integrations")
	assert.Contains(t, out.String(), "lacework integration [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework int
	out, err, exitcode = LaceworkCLI("help", "int")
	assert.Contains(t, out.String(), "lacework integration [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestIntegrationCommandList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("integration", "list")
	assert.Contains(t, out.String(), "INTEGRATION GUID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestIntegrationCommandListWithTypeFlag(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("integration", "list", "--type", "AWS_CFG")
	assert.Contains(t, out.String(), "INTEGRATION GUID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")

	// TODO @afiune lets try to create an environment where we can be 100% sure that
	// integrations will exist and assert against it

	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestIntegrationCommandListWithTypeFlagErrorUnknownType(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("integration", "list", "--type", "FOO_BAR")
	assert.Emptyf(t, out.String(),
		"STDOUT should be empty")
	assert.Contains(t, err.String(),
		"ERROR unknown integration type 'FOO_BAR'",
		"STDERR should contain the unknown type")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestIntegrationCommandShowNonExistingError(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("integration", "show", "EXXAMPLE_123")
	assert.Emptyf(t, out.String(),
		"STDOUT should be empty")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
	assert.Contains(t, err.String(),
		"ERROR the provided integration GUID was not found",
		"STDERR should contain a helpful message")
	assert.Contains(t, err.String(),
		"To list the available integrations in your account run 'lacework integrations list'",
		"STDERR should contain a helpful message")
}
