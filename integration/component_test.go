//go:build component

// Author:: Jon Stewart (<jon.stewart@lacework.net>)
// Copyright:: Copyright 2024, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package integration

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCDKComponentList(t *testing.T) {
	dir := setup()

	out := run(t, dir, "component", "list")

	found, err := regexp.Match("Not Installed\\s+component-example", []byte(out))
	assert.True(t, found)
	assert.Nil(t, err)

	cleanup(dir)
}

func TestCDKComponentListJSON(t *testing.T) {
	dir := setup()

	out := run(t, dir, "component", "list", "--json")

	expectedKeys := []string{
		"description",
		"latest_version",
		"name",
		"status",
		"type",
		"version",
	}

	for _, key := range expectedKeys {
		assert.Contains(t, out, fmt.Sprintf("\"%s\"", key),
			"STDOUT json keys changed")
	}

	cleanup(dir)
}

func TestCDKComponentShow(t *testing.T) {
	dir := setup()

	out := run(t, dir, "component", "show", "component-example")

	found, err := regexp.Match("Not Installed\\s+component-example", []byte(out))
	assert.True(t, found)
	assert.Nil(t, err)

	assert.Contains(t, out, "The following versions of this component are available to install")

	cleanup(dir)
}

func TestCDKComponentShowJSON(t *testing.T) {
	dir := setup()

	out := run(t, dir, "component", "show", "component-example", "--json")

	expectedKeys := []string{
		"description",
		"latest_version",
		"name",
		"status",
		"type",
		"version",
	}

	for _, key := range expectedKeys {
		assert.Contains(t, out, fmt.Sprintf("\"%s\"", key), "STDOUT json keys changed")
	}

	cleanup(dir)
}

func TestCDKComponentInstall(t *testing.T) {
	dir := setup()

	run(t, dir, "component", "install", "component-example")

	out := run(t, dir, "component", "list")

	found, err := regexp.Match("Installed\\s+component-example", []byte(out))
	assert.True(t, found)
	assert.Nil(t, err)

	out = run(t, dir, "component-example")
	assert.Contains(t, out, "component")

	cleanup(dir)
}

func TestCDKComponentUpdate(t *testing.T) {
	dir := setup()

	run(t, dir, "component", "install", "component-example", "--version", "0.9.0")

	out := run(t, dir, "component-example")
	assert.Contains(t, out, "component")

	t.Run("upgrade", func(t *testing.T) {
		run(t, dir, "component", "update", "component-example", "--version", "0.9.1")

		out := run(t, dir, "component-example")
		assert.Contains(t, out, "component")
	})

	t.Run("downgrade", func(t *testing.T) {
		run(t, dir, "component", "update", "component-example", "--version", "0.8.0")

		out := run(t, dir, "component-example")
		assert.Contains(t, out, "component")
	})

	cleanup(dir)
}

func TestCDKComponentUninstall(t *testing.T) {
	dir := setup()

	run(t, dir, "component", "install", "component-example")

	out := run(t, dir, "component-example")
	assert.Contains(t, out, "component")

	run(t, dir, "component", "uninstall", "component-example")

	_, err, _ := LaceworkCLIWithHome(dir, "component-example")
	assert.NotNil(t, err)

	out = run(t, dir, "component", "list")

	found, err2 := regexp.Match("Not Installed\\s+component-example", []byte(out))
	assert.True(t, found)
	assert.Nil(t, err2)

	cleanup(dir)
}

func TestCDKComponentDev(t *testing.T) {
	dir := setup()

	// GoLang scaffolding is too slow for the test
	run(t, dir, "component", "dev", "dev-component", "--description", "dev-component", "--type", "CLI_COMMAND", "--scaffolding", "\"No. Start from scratch\"")

	out := run(t, dir, "component", "list")

	found, err := regexp.Match("Development\\s+dev-component\\s+0.0.0-dev", []byte(out))
	assert.True(t, found)
	assert.Nil(t, err)

	run(t, dir, "component", "uninstall", "dev-component")

	out = run(t, dir, "component", "list")
	assert.NotContains(t, out, "dev-component")

	cleanup(dir)
}

func TestCDKComponentDevEnter(t *testing.T) {
	dir := setup()

	run(t, dir, "component", "install", "component-example")

	out := run(t, dir, "component", "list")

	found, err := regexp.Match("Installed\\s+component-example", []byte(out))
	assert.True(t, found)
	assert.Nil(t, err)

	out = run(t, dir, "component", "dev", "component-example", "--scaffolding", "\"No. Start from scratch\"")
	assert.Contains(t, out, "now in development mode")

	out = run(t, dir, "component", "list")

	found, err = regexp.Match("Development\\s+component-example", []byte(out))
	assert.True(t, found)
	assert.Nil(t, err)

	t.Run("uninstall", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithHome(dir, "component", "uninstall", "component-example")
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
		assert.Contains(t, out.String(), "- We will do better next time.\n\nDo you want to provide feedback?\n",
			"STDOUT json keys changed")
	})

	out = run(t, dir, "component", "list")

	found, err = regexp.Match("Not Installed\\s+component-example", []byte(out))
	assert.True(t, found)
	assert.Nil(t, err)

	cleanup(dir)
}

func setup() string {
	return createTOMLConfigFromCIvars()
}

func run(t *testing.T, dir string, command ...string) string {
	out, err, exitcode := LaceworkCLIWithHome(dir, command...)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	return out.String()
}

func cleanup(dir string) {
	os.RemoveAll(dir)
}
