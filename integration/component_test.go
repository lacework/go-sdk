//go:build component

// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

func TestComponentList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("component", "list")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(),
		"Loading components state...",
		"STDOUT changed, please check")

	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "VERSION",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "DESCRIPTION",
		"STDOUT table headers changed, please check")

	assert.Contains(t, out.String(), "Not Installed",
		"STDOUT our first component is not found, why?")
	assert.Contains(t, out.String(), "iac",
		"STDOUT our first component is not found, why?")
}

func TestComponentListJSON(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("component", "list", "--json")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	expectedJsonKeys := []string{
		"\"components\"",
		"\"artifacts\"",
		"\"breadcrumbs\"",
		"\"installationMessage\"",
		"\"description\"",
		"\"name\"",
		"\"type\"",
		"\"version\"",
		"\"arch\"",
		"\"os\"",
		"\"url\"",
		"\"signature\"",
	}
	t.Run("verify json keys", func(t *testing.T) {
		for _, header := range expectedJsonKeys {
			assert.Contains(t, out.String(), header,
				"STDOUT json keys changed, please check")
		}
	})

	assert.Contains(t, out.String(), "\"name\": \"iac\"",
		"missing IaC component in JSON output")
	assert.Contains(t, out.String(), "\"type\": \"CLI_COMMAND\"",
		"missing IaC component in JSON output")
}
