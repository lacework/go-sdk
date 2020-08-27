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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventCommandAliases(t *testing.T) {
	// lacework event
	out, err, exitcode := LaceworkCLI("help", "event")
	assert.Contains(t, out.String(), "lacework event [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// lacework events
	out, err, exitcode = LaceworkCLI("help", "events")
	assert.Contains(t, out.String(), "lacework event [command]")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestEventCommandList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("event", "list")
	assert.Contains(t, out.String(), "EVENT ID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "SEVERITY",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "START TIME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "END TIME",
		"STDOUT table headers changed, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestEventCommandList3Days(t *testing.T) {
	// @afiune could we find a way to generate a consistent event? but if we do
	// wouldn't the ML learn it and then become a known behavior? uhmmm
	// for now we will just check that we have the headers :wink:
	out, err, exitcode := LaceworkCLIWithTOMLConfig("event", "list", "--days", "3")
	assert.Contains(t, out.String(), "EVENT ID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "SEVERITY",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "START TIME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "END TIME",
		"STDOUT table headers changed, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestEventCommandListSeverityError(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("event", "list", "--severity", "foo")
	assert.Contains(t, err.String(), "ERROR the severity foo is not valid, use one of critical, high, medium, low, info",
		"STDERR the message to the user has changed, update please")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestEventCommandListTimeRange(t *testing.T) {
	var (
		now  = time.Now().UTC()
		from = now.AddDate(0, 0, -4) // 4 days from now
	)

	out, err, exitcode := LaceworkCLIWithTOMLConfig("event", "list", "--start", from.Format(time.RFC3339), "--end", now.Format(time.RFC3339))
	assert.Contains(t, out.String(), "EVENT ID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "SEVERITY",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "START TIME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "END TIME",
		"STDOUT table headers changed, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestEventCommandOpenError(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("event", "open", "123abc")
	assert.Contains(t, err.String(),
		"ERROR invalid event id 123abc. Event id should be a numeric value",
		"STDERR the error message changed, update")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}
