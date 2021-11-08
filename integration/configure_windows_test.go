//go:build configure || (configure && !linux)

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

func TestConfigureCommandWithJSONFileFlagError(t *testing.T) {
	out, err, exitcode := LaceworkCLI("configure", "--json_file", "foo")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Contains(t,
		err.String(),
		"ERROR unable to load keys from the provided json file: open foo: The system cannot find the file specified.",
		"STDERR should be empty")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestConfigureListHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("configure", "list", "--help")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Contains(t,
		out.String(),
		`C:\> $env:LW_PROFILE = 'my-profile'`,
		"STDOUT the environment variable in the help message is not correct")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}
