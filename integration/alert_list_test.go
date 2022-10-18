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
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func popAlert() (int, error) {
	var alerts []api.Alert

	out, stderr, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--json")
	if stderr.String() != "" {
		return 0, errors.New(stderr.String())
	}
	if exitcode != 0 {
		return 0, errors.New("non-zero exit code")
	}

	err := json.Unmarshal(out.Bytes(), &alerts)
	if err != nil {
		return 0, err
	}

	for _, a := range alerts {
		if a.Status == "Open" {
			return a.ID, nil
		}
	}
	return 0, errors.New("no open alerts found")
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

func TestAlertList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list")
	assert.Contains(t, out.String(), "ALERT ID")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertListJSON(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "list", "--json")
	assert.Contains(t, out.String(), `"alertId"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
