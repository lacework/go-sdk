//go:build alert

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

func TestAlertCloseMissingArg(t *testing.T) {
	_, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "close")
	assert.Contains(t, err.String(), "accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertCloseBadID(t *testing.T) {
	_, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "close", "me")
	assert.Contains(t, err.String(), "alert ID must be a number")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertCloseReasonSurvey(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "close", "12345")
	assert.Contains(t, out.String(), "[Use arrows to move, type to filter]")
	assert.Contains(t, err.String(), "unable to process alert close reason: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertCloseReasonInline(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "close", "12345", "-r", "1")
	assert.Contains(t, out.String(), "Type a comment")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "unable to process alert close comment: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

/* need list to get a valid alert id
func TestAlertCloseAllInline(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"alert", "close", "12345", "-r", "1", "-c", "everything is awesome")
	assert.Contains(t, out.String(), "was successfully closed")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
*/
