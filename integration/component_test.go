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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/internal/file"
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

func TestComponentDevModeGolang(t *testing.T) {
	cName := "go-component"
	dir := createTOMLConfigFromCIvars()
	defer os.RemoveAll(dir)

	t.Run("component not found", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithHome(dir, "component", "list")
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
		assert.NotContains(t, out.String(), cName,
			"the test component should not be here already, check!")

		out, err, exitcode = LaceworkCLIWithHome(dir, cName)
		assert.Empty(t, out.String(), "STDOUT should be empty")
		assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
		assert.NotContains(t, out.String(), cName,
			fmt.Sprintf("ERROR unknown command \"%s\"", cName))
	})

	t.Run("enter dev-mode", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithHome(
			dir, "component", "dev", cName, "--type", "CLI_COMMAND", "--nocolor",
			"--description", "A Go component for testing", "--noninteractive",
		)
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
		assert.Contains(t, out.String(),
			fmt.Sprintf("Component '%s' in now in development mode.", cName),
			"the test component should not be here already, check!")
		assert.Contains(t, out.String(),
			fmt.Sprintf("lacework/components/%s/.dev", cName),
			"the test component should not be here already, check!")
		assert.Contains(t, out.String(),
			"Deploy your dev component at:",
			"the test component should not be here already, check!")
	})

	t.Run(fmt.Sprintf("deploy %s component", cName), func(t *testing.T) {
		osBin := fmt.Sprintf("%s-%s-%s", cName, runtime.GOOS, runtime.GOARCH)
		if runtime.GOOS == "windows" {
			osBin += ".exe"
		}
		fromBin := filepath.Join(
			"test_resources", "cdk", cName, "bin", osBin,
		)
		toBin := filepath.Join(
			dir, ".config", "lacework", "components", cName, cName,
		)
		assert.Nil(t, file.Copy(fromBin, toBin))
	})

	t.Run("component should be found and installed", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithHome(dir, "component", "list")
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
		assert.Contains(t, out.String(), cName,
			"the test component SHOULD be here already, check!")
		assert.Contains(t, out.String(), "Installed",
			"the test component SHOULD be installed, check!")
		assert.Contains(t, out.String(), "0.0.0-dev",
			"the test component SHOULD be installed, check!")
		assert.Contains(t, out.String(), "(dev-mode) A Go component for testing",
			"the test component SHOULD match descriptiohn, check!")
	})

	t.Run("execute component", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithHome(dir, cName)
		assert.Empty(t, out.String(), "STDOUT should be empty")
		assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
		assert.Contains(t, err.String(),
			fmt.Sprintf("ERROR This is a dummy %s for testing", cName),
			"go component changed? Or something is wrong, check!")

		out, err, exitcode = LaceworkCLIWithHome(dir, cName, "run")
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
		assert.Contains(t, out.String(), "Running...",
			"go component changed? Or something is wrong, check!")
		assert.Contains(t, out.String(), "Highest Severity:",
			"go component changed? Or something is wrong, check!")

		out, err, exitcode = LaceworkCLIWithHome(dir, cName, "fail")
		assert.Empty(t, out.String(), "STDOUT should be empty")
		assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
		assert.Contains(t, err.String(), "ERROR Purposely failing...",
			"go component changed? Or something is wrong, check!")
	})

	t.Run("global help should show component", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithHome(dir, "help")
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
		assert.Contains(t, out.String(), cName,
			"the test component SHOULD be shown in help message, check!")
		assert.Contains(t, out.String(), "(dev-mode) A Go component for testing",
			"the test component SHOULD be shown in help message, check!")
	})

	t.Run("component should be displayed in version command", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithHome(dir, "version")
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
		assert.Contains(t, out.String(), fmt.Sprintf("> %s v0.0.0-dev", cName),
			"the test component SHOULD be shown in version command, check!")
	})
}
