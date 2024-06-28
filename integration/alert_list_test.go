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
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/circleci/fork-lacework-go-sdk/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func popAlert() (string, error) {
	var alerts api.Alerts

	out, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"alert", "list", "--status", "Open", "--json", "--range", "last 7 days", "--nocache",
	)
	if stderr.String() != "" {
		return "-1", errors.New(stderr.String())
	}
	if exitcode != 0 {
		return "-1", errors.New("non-zero exit code")
	}

	err := json.Unmarshal(out.Bytes(), &alerts)
	if err != nil {
		return "-1", err
	}

	alerts.SortByID()
	for _, a := range alerts {
		return strconv.Itoa(a.ID), nil
	}
	return "-1", errors.New("no open alerts found")
}

func TestAlertListBadRange(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--range", "bad")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR unable to list alerts: natural time (bad) is invalid")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListBadStart(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--start", "bad")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR unable to list alerts: unable to parse time (bad)")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListBadEnd(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--end", "bad")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR unable to list alerts: unable to parse time (bad)")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListHuman(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list")
	assert.Contains(t, out.String(), "ALERT ID")
	// breadcrumb
	assert.Contains(t, out.String(), "lacework alert show <alert_id>")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListJSON(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--json")
	assert.Contains(t, out.String(), `"alertId"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListNone(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--start", "-10s", "--end", "now")
	assert.Contains(t, out.String(), "There are no alerts in the specified time range.")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListSeverityBad(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--severity", "foo")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR unable to list alerts: the severity (foo) is not valid, use one of (critical, high, medium, low, info)")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListSeverityHighAndCritical(t *testing.T) {
	// Sometimes tech-ally doesn't have critical or high alerts, so those get skipped
	// this is a strange but required workaround since these tests run against real data
	// we expect to always find medium low and info
	t.Run("should be able to find critical alerts", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--severity", "critical", "--range", "last week")
		if strings.Contains(out.String(), "Critical") {
			assert.Contains(t, out.String(), "Critical")
			assert.Empty(t, err.String(), "STDERR should be empty")
			assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
		} else {
			t.Skip()
		}
	})
	t.Run("should be able to find high alerts", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--severity", "high", "--range", "last week")
		if strings.Contains(out.String(), "high") {
			assert.Contains(t, out.String(), "high")
			assert.Empty(t, err.String(), "STDERR should be empty")
			assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
		} else {
			t.Skip()
		}
	})
	t.Run("should be able to find low,medium,info alerts", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--severity", "high", "--range", "last week")
		assert.NotContains(t, out.String(), "Low")
		assert.NotContains(t, out.String(), "Medium")
		assert.NotContains(t, out.String(), "Info")
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	})
}

func TestAlertListStatusBad(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--status", "foo")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR unable to list alerts: the status (foo) is not valid, use one of (Open, Closed)")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListStatusOpen(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--status", "Open", "--range", "last 7 days")
	assert.Contains(t, out.String(), "Open")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListTypeBad(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--type", "foo")
	assert.Contains(t, out.String(), "No alerts match the specified filters within the given time range. Try removing filters or expanding the time range.")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
