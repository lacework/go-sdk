//
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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComplianceAwsGetReportFilter(t *testing.T) {
	account := os.Getenv("LW_INT_TEST_AWS_ACC")
	detailsOutput := "recommendations showing"
	out, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "aws", "get-report", account, "--details", "--status", "could-not-assess")
	assert.Contains(t, out.String(), detailsOutput, "Filtered detail output should contain filtered result")
	assert.Empty(t, err.String(),  "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}
