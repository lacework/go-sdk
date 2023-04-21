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
	"fmt"
	"os"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestComplianceAwsList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "aws", "list",
	)
	assert.Empty(t, err.String(), "STDERR should be empty")

	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "AWS ACCOUNT", "STDOUT changed, please check")
	assert.Contains(t, out.String(), "STATUS", "STDOUT changed, please check")
}

func TestComplianceAwsGetReportFilter(t *testing.T) {
	account := os.Getenv("LW_INT_TEST_AWS_ACC")
	detailsOutput := "recommendations showing"
	out, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "aws", "get-report", account, "--status", "compliant", "--type", "AWS_CIS_S3")

	assert.Contains(t, out.String(), detailsOutput, "Filtered detail output should contain filtered result")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "COMPLIANCE REPORT DETAILS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), account,
		"Account ID in compliance report is not correct")
	assert.Contains(t, out.String(), "NON-COMPLIANT RECOMMENDATIONS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "ID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "RECOMMENDATION",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "SEVERITY",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "SERVICE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "AFFECTED",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "ASSESSED",
		"STDOUT table headers changed, please check")
}

func TestComplianceAwsGetReportDetails(t *testing.T) {
	account := os.Getenv("LW_INT_TEST_AWS_ACC")
	detailsOutput := "recommendations showing"
	out, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "aws", "get-report", account, "--details")

	assert.NotContains(t, out.String(), detailsOutput,
		"Details table without filter should not contain filtered output")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "COMPLIANCE REPORT DETAILS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), account,
		"Account ID in compliance report is not correct")
	assert.Contains(t, out.String(), "NON-COMPLIANT RECOMMENDATIONS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "ID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "RECOMMENDATION",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "SEVERITY",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "SERVICE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "AFFECTED",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "ASSESSED",
		"STDOUT table headers changed, please check")
}

func TestComplianceAwsGetReportFiltersWithJsonOutput(t *testing.T) {
	account := os.Getenv("LW_INT_TEST_AWS_ACC")
	out, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "aws", "get-report", account, "--severity", "critical", "--json")
	severities := []string{"\"severity\": 2", "\"severity\": 3", "\"severity\": 4", "\"severity\": 5"}
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	// When critical severity filter is set, other severities should not be returned in json result
	assert.NotContains(t, severities, out.String(),
		"Json output does not adhere to severity filter")
}

func TestComplianceAwsGetReportAccountIDWithAlias(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "aws", "get-report", "account-id (account-alias)",
	)
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(),
		"Getting compliance report...",
		"STDOUT changed, please check")
	assert.Contains(t, err.String(),
		"no data found in the report",
		"STDERR changed, please check")
}

func TestComplianceAwsGetReportTypeAWS_SOC_Rev2(t *testing.T) {
	account := os.Getenv("LW_INT_TEST_AWS_ACC")
	out, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "aws", "get-report", account, "--type", "AWS_SOC_Rev2")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "AWS SOC 2 Report Rev2",
		"STDOUT report type missing or something else is going on, please check")
	assert.Contains(t, out.String(), "Report Type",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), account,
		"Account ID in compliance report is not correct")
}

func TestComplianceAwsGetReportByName(t *testing.T) {
	account := os.Getenv("LW_INT_TEST_AWS_ACC")
	out, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "aws", "get-report", account, "--report_name", "AWS CSA CCM 4.0.5")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "AWS Cloud Security Alliance",
		"STDOUT report type missing or something else is going on, please check")
	assert.Contains(t, out.String(), "Report Title",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), account,
		"Account ID in compliance report is not correct")
}

func TestComplianceAwsGetAllReportType(t *testing.T) {
	account := os.Getenv("LW_INT_TEST_AWS_ACC")
	for _, reportType := range api.AwsReportTypes() {
		t.Run(reportType, func(t *testing.T) {
			out, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "aws", "get-report", account, "--type", reportType)
			assert.Empty(t, err.String(), "STDERR should be empty")
			assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
			assert.Contains(t, out.String(), "COMPLIANCE REPORT DETAILS",
				"STDOUT table headers changed, please check")
		})
	}
}

func TestComplianceAwsGetReportRecommendationID(t *testing.T) {
	account := os.Getenv("LW_INT_TEST_AWS_ACC")
	out, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "aws", "get-report", account, "2.1.2")

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "RECOMMENDATION DETAILS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "SEVERITY",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "SERVICE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "CATEGORY",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "ASSESSED RESOURCES ",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "AFFECTED RESOURCES",
		"STDOUT table headers changed, please check")
}

func TestComplianceAwsGetReportRecommendationIDNotFound(t *testing.T) {
	account := os.Getenv("LW_INT_TEST_AWS_ACC")
	_, err, exitcode := LaceworkCLIWithTOMLConfig("compliance", "aws", "get-report", account, "rec-not-found")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, err.String(), "recommendation id 'rec-not-found' not found.",
		"STDERR changed?, please check")
}

func TestComplianceAwsSearchEmpty(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "aws", "search", "example",
	)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "Resource 'example' not found.", "STDOUT changed, please check")
}

func TestComplianceAwsScan(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "aws", "scan",
	)

	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "STATUS")
	assert.Contains(t, out.String(), "DETAILS")
}

func _TestComplianceAwsSearch(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig(
		"compliance", "aws", "search", "arn:aws:s3:::tech-ally-test",
	)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	fmt.Println(out.String())

	assert.Contains(t, out.String(), "RECOMMENDATION ID", "table headers missing")
	assert.Contains(t, out.String(), "ACCOUNT ID", "table headers missing")
	assert.Contains(t, out.String(), "REASON", "table headers missing")
	assert.Contains(t, out.String(), "SEVERITY", "table headers missing")
	assert.Contains(t, out.String(), "STATUS", "table headers missing")

	assert.Contains(t, out.String(), "LW_S3_12", "recommendation id missing")
	assert.Contains(t, out.String(), "S3 bucket does not have MFA", "reason missing")
}
