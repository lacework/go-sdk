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
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeComment(id string) (bytes.Buffer, bytes.Buffer, int) {
	return LaceworkCLIWithTOMLConfig(
		"alert", "comment", id, "-c", "everything is awesome...cause go-sdk is here", "-f", "plaintext")
}

func TestAlertCommentMissingArg(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "comment")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertCommentBadID(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "comment", "me")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "alert ID must be a number")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertComment(t *testing.T) {
	id, err := popAlert()
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("alert", "comment", id)
	assert.Contains(t, out.String(), "Type a comment")
	assert.Contains(t, out.String(), "[Enter to launch editor]")
	assert.Contains(t, stderr.String(), "unable to process alert comment: EOF")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertCommentInline(t *testing.T) {
	id, err := popAlert()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	out, stderr, exitcode := makeComment(id)
	assert.Contains(t, out.String(), "Comment added successfully")
	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertCommentDoesNotExist(t *testing.T) {
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("alert", "comment", "123456789101112")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, stderr.String(), "alert 123456789101112 does not exist")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}
