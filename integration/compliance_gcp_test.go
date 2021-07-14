//
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
		"unable to get gcp compliance report",
		"STDERR changed, please check")
	assert.Contains(t, err.String(),
		"GCP_ORG_ID=org-id&GCP_PROJ_ID=proj-id&",
		"STDERR changed, please check")
}

func TestComplianceGoogleList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "gcp", "list",
	)
	assert.Empty(t, err.String(), "STDERR should be empty")

	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(),
		"PROJECT ID   ORGANIZATION ID",
		"STDOUT changed, please check")
}
