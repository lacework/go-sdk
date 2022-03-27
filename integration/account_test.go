//go:build account

// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

func TestAccountCommandAliases(t *testing.T) {
	// lacework account
	out, err, exitcode := LaceworkCLI("help", "account")
	assert.Contains(t, out.String(), "lacework account")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework account list
	out, err, exitcode = LaceworkCLI("help", "acc", "ls")
	assert.Contains(t, out.String(), "lacework account list [flags]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAccountCommandList(t *testing.T) {
	if os.Getenv("CI_STANDALONE_ACCOUNT") != "" {
		t.Skip("skipping organizational account test")
	}
	out, err, exitcode := LaceworkCLIWithTOMLConfig("account", "list")
	assert.Contains(t, out.String(), "ACCOUNTS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), os.Getenv("CI_ACCOUNT"),
		"STDOUT unable to find account, please check")
	assert.Contains(t, out.String(), os.Getenv("CI_SUBACCOUNT"),
		"STDOUT unable to find account, please check")
	assert.Contains(t, out.String(),
		"Use '--subaccount <name>' to switch any command to a different account.",
		"STDOUT wrong breadcrumbs, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}
