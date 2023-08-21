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
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var alertShowID, _ = popAlert()

func TestAlertShowBadID(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "show", "bad")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "alert ID must be a number")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertShowBadScope(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "bad")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR scope (bad) is not recognized")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestAlertShowDetails(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "show", alertShowID)
	assert.Contains(t, out.String(), "SUBJECT")
	assert.Contains(t, out.String(), "For further investigation")

	accountName := os.Getenv("CI_ACCOUNT")
	if os.Getenv("CI_SUBACCOUNT") != "" {
		accountName = os.Getenv("CI_SUBACCOUNT")
	}
	assert.Contains(
		t,
		out.String(),
		fmt.Sprintf(
			"https://%s.lacework.net/ui/investigation/monitor/AlertInbox/%s/details?accountName=%s",
			os.Getenv("CI_ACCOUNT"),
			alertShowID,
			accountName,
		),
	)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertShowDetailsJSON(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--json")
	assert.Contains(t, out.String(), `"alertId"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertShowInvestigation(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "Investigation")

	// RAIN-38816: Alerts: Valid Scopes Return 404 When Entries Are Not Available
	if strings.Contains(err.String(), "[404] Not found") {
		return
	}
	assert.Contains(t, out.String(), "QUESTION")
	assert.Contains(t, out.String(), "For further investigation")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	out, err, exitcode = LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "Investigation", "--json")
	assert.Contains(t, out.String(), `"question"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertShowEvents(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "Events")

	// RAIN-38816: Alerts: Valid Scopes Return 404 When Entries Are Not Available
	if strings.Contains(err.String(), "[404] Not found") {
		return
	}
	assert.Contains(t, out.String(), `"additionalEventInfo"`)
	assert.Contains(t, out.String(), "For further investigation")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	out, err, exitcode = LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "Events", "--json")
	assert.Contains(t, out.String(), `"additionalEventInfo"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertShowRelatedAlerts(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "RelatedAlerts")

	// RAIN-38816: Alerts: Valid Scopes Return 404 When Entries Are Not Available
	if strings.Contains(err.String(), "[404] Not found") {
		return
	}
	assert.Contains(t, out.String(), "RANK")
	assert.Contains(t, out.String(), "For further investigation")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	out, err, exitcode = LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "RelatedAlerts", "--json")
	assert.Contains(t, out.String(), `"rank"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertShowIntegrations(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "Integrations")

	// RAIN-38816: Alerts: Valid Scopes Return 404 When Entries Are Not Available
	if strings.Contains(err.String(), "[404] Not found") {
		return
	}
	assert.Contains(t, out.String(), "GUID")
	assert.Contains(t, out.String(), "For further investigation")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	out, err, exitcode = LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "Integrations", "--json")
	assert.Contains(t, out.String(), `"intgGuid"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertShowTimeline(t *testing.T) {
	makeComment(alertShowID)
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "Timeline")

	assert.Contains(t, out.String(), "@")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	out, err, exitcode = LaceworkCLIWithTOMLConfig("alert", "show", alertShowID, "--scope", "Timeline", "--json")
	assert.Contains(t, out.String(), `"entryType"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
