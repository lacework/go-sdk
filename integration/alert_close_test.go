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
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/consul/sdk/testutil/retry"
	"github.com/stretchr/testify/assert"
)

func threeTimesWaitOneSecond() *retry.Counter {
	// ThreeTimes repeats an operation three times and waits 1s in between.
	return &retry.Counter{Count: 3, Wait: 1 * time.Second}
}

func alreadyClosedRetry(r *retry.R, s string) {
	if strings.Contains(s, "The action on this alert is not allowed") {
		r.Error("popAlert returned a closed alert")
	}
}

func TestAlertCloseMissingArg(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "close")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertCloseBadID(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "close", "me")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "alert ID must be a number")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertCloseReasonSurvey(t *testing.T) {
	retry.RunWith(threeTimesWaitOneSecond(), t, func(r *retry.R) {
		id, err := popAlert()
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		out, stderr, exitcode := LaceworkCLIWithTOMLConfig("alert", "close", id)
		alreadyClosedRetry(r, stderr.String())
		assert.Contains(t, out.String(), "[Use arrows to move, type to filter]")
		assert.Contains(t, stderr.String(), "unable to process alert close reason: EOF")
		assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
	})
}

func TestAlertCloseReasonInline(t *testing.T) {
	retry.RunWith(threeTimesWaitOneSecond(), t, func(r *retry.R) {
		id, err := popAlert()
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		out, stderr, exitcode := LaceworkCLIWithTOMLConfig("alert", "close", id, "-r", "1")
		alreadyClosedRetry(r, stderr.String())
		assert.Contains(t, out.String(), "Type a comment")
		assert.Contains(t, out.String(), "[Enter to launch editor]")
		assert.Contains(t, stderr.String(), "unable to process alert close comment: EOF")
		assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
	})
}

// unstable test disabled as part of GROW-1396
func _TestAlertCloseInline(t *testing.T) {
	retry.RunWith(threeTimesWaitOneSecond(), t, func(r *retry.R) {
		id, err := popAlert()
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		// verify
		out, stderr, exitcode := LaceworkCLIWithTOMLConfig(
			"alert", "close", id, "-r", "1", "-c", "everything is awesome")
		assert.Contains(t, out.String(), "Are you sure you want to close alert")
		assert.Contains(t, stderr.String(), "unable to confirm alert close attempt: EOF")
		assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")

		// close
		out, stderr, exitcode = LaceworkCLIWithTOMLConfig(
			"alert", "close", id, "-r", "1", "-c", "everything is awesome", "--noninteractive")
		alreadyClosedRetry(r, stderr.String())
		assert.Contains(t, out.String(), "was successfully closed.")
		assert.Empty(t, stderr.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

		// list closed
		out, stderr, exitcode = LaceworkCLIWithTOMLConfig("alert", "list", "--status", "Closed")
		assert.Contains(t, out.String(), id)
		assert.Empty(t, stderr.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	})
}

func TestAlertCloseDoesNotExist(t *testing.T) {
	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("alert", "close", "123456789101112")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, stderr.String(), "alert 123456789101112 does not exist")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}
