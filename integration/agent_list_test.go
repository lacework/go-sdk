//go:build agent

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

func TestAgentList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("agent", "list")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")

	expectedOutput := []string{
		// headers
		"AGENT VERSION",
		"NAME",
		"STATUS",
		"HOSTNAME",
		"MID",
		"LAST CHECK-IN",
	}
	for _, str := range expectedOutput {
		assert.Contains(t, out.String(), str,
			"STDOUT table does not contain the '"+str+"' output")
	}
}
