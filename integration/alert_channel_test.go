//go:build alert_channel

// Author:: Darren Murray (<darren.murray@lacework.net>)
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

func TestAlertChannelShow(t *testing.T) {
	t.Parallel()
	if !alertChannelExists() {
		t.Skip("Alert Channel does not exist")
	}
	out, err, exitcode := LaceworkCLIWithTOMLConfig("ac", "show", "TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A")
	// Summary Table
	assert.Contains(t, out.String(), "ALERT CHANNEL GUID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATE",
		"STDOUT table headers changed, please check")

	// Details Table
	assert.Contains(t, out.String(), "EXTERNAL ID",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "ROLE ARN",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "BUCKET ARN",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "UPDATED AT",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "UPDATED BY",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "STATE UPDATED AT",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "LAST SUCCESSFUL STATE",
		"STDOUT details headers changed, please check")

	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestAlertChannelList(t *testing.T) {
	t.Parallel()
	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert-channel", "list")
	assert.Contains(t, out.String(), "ALERT CHANNEL GUID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATE",
		"STDOUT table headers changed, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func alertChannelExists() bool {
	lw, err := laceworkIntegrationTestClient()
	resp, err := lw.V2.AlertChannels.GetAwsS3("TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A")
	if err != nil {
		return false
	}
	name := resp.Data.Name
	return name != ""
}
