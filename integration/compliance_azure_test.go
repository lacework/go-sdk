//go:build compliance

// Author:: Darren Murray (<darren.murray@lacework.net>)
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

func TestComplianceAzureListTenants(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "az", "list-tenants")
	assert.Contains(t, out.String(),
		"Fetching list of configured Azure tenants",
		"STDOUT changed, please check")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Empty(t, err.String(), "STDERR should be empty")
}

func TestComplianceAzureGetReportTenantAndSubscriptionWithAlias(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "azure", "get-report", "tenant-id (tenant-alias)", "subscription-id (subscription-alias)",
	)
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(),
		"Getting compliance report...",
		"STDOUT changed, please check")
	assert.Contains(t, err.String(),
		"no data found in the report",
		"STDERR changed, please check")
}

func TestComplianceAzureScan(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "azure", "scan",
	)

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "STATUS    scanning")
	assert.Contains(t, out.String(), "DETAILS   Scan has been requested")
}
