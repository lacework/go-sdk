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

func TestAlertCommentMissingArg(t *testing.T) {
	_, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "comment")
	assert.Contains(t, err.String(), "accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertCommentBadID(t *testing.T) {
	_, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "comment", "me")
	assert.Contains(t, err.String(), "alert ID must be a number")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertComment(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "comment", "12345")
	assert.Contains(t, out.String(), "Type a comment")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, err.String(), "unable to process alert comment: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

/* need list to get a valid alert id
func TestAlertCommentInline(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"alert", "comment", "12345", "-c", "everything is awesome")
	assert.Contains(t, out.String(), "was successfully closed")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
*/
