//go:build configure

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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureSwitchProfileNoConfigFails(t *testing.T) {
	out, err, exitcode := LaceworkCLI("configure", "switch-profile", "foo")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Contains(t,
		err.String(),
		"ERROR unable to load profiles. No configuration file found.",
		"STDERR message is not correct")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestConfigureSwitchProfileNotFound(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "switch-profile", "bar")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Contains(t,
		err.String(),
		"ERROR Profile 'bar' not found. Try 'lacework configure list' to see all configured profiles.",
		"STDERR message is not correct")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestConfigureSwitchProfileWithConfig(t *testing.T) {
	// @afiune we store the temp directory since we need to run more
	// commands after switching a profile
	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)

	out, err, exitcode := LaceworkCLIWithHome(dir, "configure", "switch-profile", "dev")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Contains(t,
		out.String(), "Profile switched to 'dev'.",
		"STDOUT message is not correct")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")

	out, err, exitcode = LaceworkCLIWithHome(dir, "configure", "show", "account")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Contains(t,
		out.String(), "dev.example",
		"STDOUT message is not correct")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")

	t.Run("re-running should pass", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "switch-profile", "dev")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Contains(t,
			out.String(), "Profile switched to 'dev'.",
			"STDOUT message is not correct")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")

		out, err, exitcode = LaceworkCLIWithHome(dir, "configure", "show", "account")
		assert.Contains(t, out.String(), "dev.example", "STDOUT message is not correct")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
	})

	t.Run("switch back to the default profile", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithHome(dir, "configure", "switch-profile", "default")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Contains(t,
			out.String(), "Profile switched back to default.",
			"STDOUT message is not correct")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")

		out, err, exitcode = LaceworkCLIWithHome(dir, "configure", "show", "account")
		assert.Contains(t, out.String(), "dummy", "STDOUT message is not correct")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")

		t.Run("re-run should pass", func(t *testing.T) {
			out, err, exitcode := LaceworkCLIWithHome(dir, "configure", "switch-profile", "default")
			assert.Empty(t,
				err.String(),
				"STDERR should be empty")
			assert.Contains(t,
				out.String(), "Profile switched back to default.",
				"STDOUT message is not correct")
			assert.Equal(t, 0, exitcode,
				"EXITCODE is not the expected one")

			out, err, exitcode = LaceworkCLIWithHome(dir, "configure", "show", "account")
			assert.Contains(t, out.String(), "dummy", "STDOUT message is not correct")
			assert.Empty(t,
				err.String(),
				"STDERR should be empty")
			assert.Equal(t, 0, exitcode,
				"EXITCODE is not the expected one")
		})
	})
}
