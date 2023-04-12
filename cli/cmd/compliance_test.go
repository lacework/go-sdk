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

package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestComplianceRecommendationsFilterNoResults(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Category = []string{"monitoring"}
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 0)
	assert.Equal(t, output, "There are no recommendations with the specified filter(s).\n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnCategory(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Category = []string{"identity-and-access-management"}
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "1 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnService(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Service = []string{"aws:cloudtrail"}
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, output, "2 of 4 recommendations showing \n")
	clearFilters()
}

// Severity returns everything above the specified threshold eg. "low" returns low, medium, high, critical
func TestComplianceRecommendationsFilterOnSeverityLow(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Severity = "low"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 4)
	assert.Equal(t, output, "4 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnSeverityMedium(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Severity = "medium"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 3)
	assert.Equal(t, output, "3 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnSeverityCritical(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Severity = "critical"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "1 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnStatus(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Status = "non-compliant"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, output, "2 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterMultiple(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Severity = "high"
	compCmdState.Status = "non-compliant"

	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, output, "2 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterMultipleCategories(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Category = []string{"s3", "logging"}

	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 3)
	assert.Equal(t, output, "3 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterMultipleServices(t *testing.T) {
	mockRecommendations := []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Service = []string{"aws:s3", "aws:iam"}

	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, output, "2 of 4 recommendations showing \n")
	clearFilters()
}

func TestStatusInputToProperTransform(t *testing.T) {
	status := statusToProperTypes("non-compliant")
	assert.Equal(t, status, "NonCompliant")

	status = statusToProperTypes("compliant")
	assert.Equal(t, status, "Compliant")

	status = statusToProperTypes("suppressed")
	assert.Equal(t, status, "Suppressed")

	status = statusToProperTypes("requires-manual-assessment")
	assert.Equal(t, status, "RequiresManualAssessment")
}

func TestFiltersEnabled(t *testing.T) {
	NoneEnabled := complianceFiltersEnabled()
	assert.Equal(t, NoneEnabled, false)

	compCmdState.Category = []string{"s3"}
	compCmdState.Status = "non-compliant"
	compCmdState.Severity = "high"
	compCmdState.Service = []string{"aws:s3"}
	AllEnabled := complianceFiltersEnabled()
	assert.Equal(t, AllEnabled, true)

	compCmdState.Severity = ""
	compCmdState.Service = []string{}

	SomeEnabled := complianceFiltersEnabled()
	assert.Equal(t, SomeEnabled, true)

	clearFilters()
}

func TestComplianceRecommendationsRecommendationID(t *testing.T) {
	result, found := mockComplianceReport.GetComplianceRecommendation("LW_S3_1")

	if assert.True(t, found) {
		assert.NotNil(t, result)
		assert.Equal(t, result.Status, "NonCompliant")
		assert.Equal(t, result.AssessedResourceCount, 1)
		assert.Equal(t, result.ResourceCount, 1)
		assert.Equal(t, result.Category, "S3")
		assert.Equal(t, len(result.Violations), 1)
		assert.Equal(t, result.Violations[0].Resource, "arn:aws:s3:::resource-name")
		assert.Equal(t, result.Violations[0].Region, "us-east-1")
		assert.Equal(t, result.Violations[0].Reasons[0], "violation reason")
	}
}

func TestComplianceRecommendationsRecommendationIDNotFound(t *testing.T) {
	result, found := mockComplianceReport.GetComplianceRecommendation("1.2.3.4.5")
	assert.False(t, found)
	assert.Nil(t, result)
}

func clearFilters() {
	compCmdState.Category = []string{}
	compCmdState.Severity = ""
	compCmdState.Service = []string{}
	compCmdState.Status = ""
}

var (
	mockRecommendationOne = api.RecommendationV2{
		RecID:                 "LW_S3_1",
		AssessedResourceCount: 1,
		ResourceCount:         1,
		Category:              "S3",
		InfoLink:              "",
		Service:               "aws:s3",
		Severity:              2,
		Status:                "NonCompliant",
		Suppressions:          []string{},
		Title:                 "Mock S3",
		Violations: []api.ComplianceViolationV2{{
			Region:   "us-east-1",
			Resource: "arn:aws:s3:::resource-name",
			Reasons:  []string{"violation reason"},
		}},
	}
	mockRecommendationTwo = api.RecommendationV2{
		RecID:                 "AWS_CIS_1_7",
		AssessedResourceCount: 1,
		ResourceCount:         1,
		Category:              "Identity and Access Management",
		InfoLink:              "",
		Service:               "aws:iam",
		Severity:              1,
		Status:                "Compliant",
		Suppressions:          []string{},
		Title:                 "Mock IAM",
		Violations: []api.ComplianceViolationV2{{
			Region:   "us-east-2",
			Resource: "arn:aws:s3:::other-resource-name",
			Reasons:  []string{"other-violation reason"},
		}}}
	mockRecommendationThree = api.RecommendationV2{
		RecID:                 "AWS_CIS_2_2",
		AssessedResourceCount: 1,
		ResourceCount:         1,
		Category:              "Logging",
		InfoLink:              "",
		Service:               "aws:cloudtrail",
		Severity:              2,
		Status:                "NonCompliant",
		Suppressions:          []string{},
		Title:                 "Mock Log One",
		Violations:            []api.ComplianceViolationV2{},
	}

	mockRecommendationFour = api.RecommendationV2{
		RecID:                 "AWS_CIS_2_2",
		AssessedResourceCount: 1,
		ResourceCount:         1,
		Category:              "Logging",
		InfoLink:              "",
		Service:               "aws:cloudtrail",
		Severity:              4,
		Status:                "Compliant",
		Suppressions:          []string{},
		Title:                 "Mock Log Two",
		Violations:            []api.ComplianceViolationV2{},
	}

	mockComplianceReport = api.AwsReport{
		ReportTitle:     "This is a mock Report",
		ReportType:      "Report Type",
		ReportTime:      time.Time{},
		AccountID:       "12345679",
		AccountAlias:    "alias",
		Summary:         nil,
		Recommendations: []api.RecommendationV2{mockRecommendationOne, mockRecommendationTwo},
	}
)
