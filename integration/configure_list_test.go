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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureListCommandWithConfig(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "list")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")

	expectedFields := []string{
		// headers
		"PROFILE",
		"ACCOUNT",
		"SUBACCOUNT",
		"API KEY",
		"API SECRET",
		"V",

		// column 1
		"> default",
		"dummy",
		"DUMMY_1234567890abcdefg",
		"*************cret",

		// column 2
		"dev",
		"dev.example",
		"DEVDEV_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC000",
		"*****************************1111",

		// column 3
		"v2",
		"v2.config",
		"subaccount.example",
		"V2_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00",
		"*****************************2222",

		// column 4
		"integration",
		"integration",
		"INTEGRATION_3DF1234AABBCCDD5678XXYYZZ1234ABC8BEC6500DC70",
		"*****************************4abc",

		// column 5
		"test",
		"test.account",
		"INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00",
		"*****************************0000",
	}
	t.Run("verify table fields", func(t *testing.T) {
		for _, field := range expectedFields {
			assert.Contains(t, out.String(), field,
				"STDOUT something intside the table is missing, please check")
		}
	})
}

func TestConfigureListCommandWithoutConfig(t *testing.T) {
	out, err, exitcode := LaceworkCLI("configure", "list")
	assert.Contains(t, err.String(), "ERROR unable to load profiles. No configuration file found.",
		"STDERR message changed, please check")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestConfigureListCommandWithConfigAndProfile(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithDummyConfig("configure", "list")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")

	t.Run("verify selected profile: > default", func(t *testing.T) {
		assert.Contains(t, out.String(), "> default",
			"STDOUT something intside the table is missing, please check")
	})

	out, err, exitcode = LaceworkCLIWithDummyConfig("configure", "list", "-p", "integration")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")

	t.Run("verify selected profile: > integration", func(t *testing.T) {
		assert.Contains(t, out.String(), "> integration",
			"STDOUT something intside the table is missing, please check")
	})

	out, err, exitcode = LaceworkCLIWithDummyConfig("configure", "list", "--profile", "dev")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")

	t.Run("verify selected profile: > dev", func(t *testing.T) {
		assert.Contains(t, out.String(), "> dev",
			"STDOUT something intside the table is missing, please check")
	})
}
