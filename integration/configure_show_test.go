//go:build configure

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

func TestConfigureShowCommandWrongKey(t *testing.T) {
	t.Parallel()
	out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "show", "foo")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Contains(t, err.String(), "unknown configuration key.",
		"STDERR is not correct, please update")
	assert.Contains(t, err.String(), "(available: profile, account, subaccount, api_secret, api_key, version)",
		"STDERR is not correct, please update")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestConfigureShowCommandWithConfig(t *testing.T) {
	t.Parallel()
	out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "show", "profile")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "default",
		"STDOUT wrong computed profile")
}

func TestConfigureShowCommandWithoutConfig(t *testing.T) {
	t.Parallel()
	out, err, exitcode := LaceworkCLI("configure", "show", "account")
	assert.Empty(t,
		err.String(),
		"STDOUT should be empty")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestConfigureShowCommandWithConfigAndProfile(t *testing.T) {
	t.Parallel()
	t.Run("dev.account", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "show", "account", "-p", "dev")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Equal(t, "dev.example\n", out.String(),
			"STDOUT does not match with the correct value")
	})

	t.Run("integration.api_secret", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "show", "api_secret", "-p", "integration")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Equal(t, "_1234abdc00ff11vv22zz33xyz1234abc\n", out.String(),
			"STDOUT does not match with the correct value")
	})

	t.Run("test.api_key", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "show", "api_key", "-p", "test")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Equal(t, "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00\n", out.String(),
			"STDOUT does not match with the correct value")
	})

	t.Run("v2.subaccount", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "show", "subaccount", "--profile", "v2")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Equal(t, "subaccount.example\n", out.String(),
			"STDOUT does not match with the correct value")
	})

	t.Run("v2.version", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "show", "version", "-p", "v2")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Equal(t, "2\n", out.String(),
			"STDOUT does not match with the correct value")
	})

	t.Run("foo.unknown", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "show", "account", "-p", "foo")
		assert.Empty(t,
			err.String(),
			"STDERR should be empty")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
		assert.Empty(t,
			out.String(),
			"STDERR should be empty")
	})
}
