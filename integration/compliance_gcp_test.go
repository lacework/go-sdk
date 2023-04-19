//go:build compliance

// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

func TestComplianceGoogleGetReportOrgAndProjectWithAlias(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "gcp", "get-report", "org-id (org-alias)", "proj-id (proj-alias)",
	)
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(),
		"Getting compliance report...",
		"STDOUT changed, please check")
	assert.Contains(t, err.String(),
		"no data found in the report",
		"STDERR changed, please check")
}

func TestComplianceGoogleList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "gcp", "list",
	)
	assert.Empty(t, err.String(), "STDERR should be empty")

	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "PROJECT ID", "STDOUT changed, please check")
	assert.Contains(t, out.String(), "ORGANIZATION ID", "STDOUT changed, please check")
	assert.Contains(t, out.String(), "STATUS", "STDOUT changed, please check")
}

func TestComplianceGcpScan(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "gcp", "scan",
	)

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "STATUS    scanning")
	assert.Contains(t, out.String(), "DETAILS   Scan has been requested")
}
