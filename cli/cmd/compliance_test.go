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

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestComplianceRecommendationsFilterNoResults(t *testing.T) {
	mockRecommendations := []api.ComplianceRecommendation{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Category = "monitoring"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 0)
	assert.Equal(t, output, "There are no recommendations with the specified filters.\n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnCategory(t *testing.T) {
	mockRecommendations := []api.ComplianceRecommendation{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Category = "identity-and-access-management"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "1 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnService(t *testing.T) {
	mockRecommendations := []api.ComplianceRecommendation{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Service = "aws:cloudtrail"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, output, "2 of 4 recommendations showing \n")
	clearFilters()
}

// Severity returns everything above the specified threshold eg. "low" returns low, medium, high, critical
func TestComplianceRecommendationsFilterOnSeverityLow(t *testing.T) {
	mockRecommendations := []api.ComplianceRecommendation{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Severity = "low"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 4)
	assert.Equal(t, output, "4 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnSeverityMedium(t *testing.T) {
	mockRecommendations := []api.ComplianceRecommendation{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Severity = "medium"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 3)
	assert.Equal(t, output, "3 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnSeverityCritical(t *testing.T) {
	mockRecommendations := []api.ComplianceRecommendation{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Severity = "critical"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "1 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterOnStatus(t *testing.T) {
	mockRecommendations := []api.ComplianceRecommendation{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Status = "non-compliant"
	result, output := filterRecommendations(mockRecommendations)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, output, "2 of 4 recommendations showing \n")
	clearFilters()
}

func TestComplianceRecommendationsFilterMultiple(t *testing.T) {
	mockRecommendations := []api.ComplianceRecommendation{mockRecommendationOne, mockRecommendationTwo,
		mockRecommendationThree, mockRecommendationFour}
	compCmdState.Severity = "high"
	compCmdState.Status = "non-compliant"

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
	NoneEnabled := filtersEnabled()
	assert.Equal(t, NoneEnabled, false)

	compCmdState.Category = "s3"
	compCmdState.Status = "non-compliant"
	compCmdState.Severity = "high"
	compCmdState.Service = "aws:s3"
	AllEnabled := filtersEnabled()
	assert.Equal(t, AllEnabled, true)

	compCmdState.Severity = ""
	compCmdState.Service = ""

	SomeEnabled := filtersEnabled()
	assert.Equal(t, SomeEnabled, true)

	clearFilters()
}

func clearFilters() {
	compCmdState.Category = ""
	compCmdState.Severity = ""
	compCmdState.Service = ""
	compCmdState.Status = ""
}

var (
	mockRecommendationOne = api.ComplianceRecommendation{
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
		Violations:            []api.ComplianceViolation{},
	}
	mockRecommendationTwo = api.ComplianceRecommendation{
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
		Violations:            []api.ComplianceViolation{},
	}
	mockRecommendationThree = api.ComplianceRecommendation{
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
		Violations:            []api.ComplianceViolation{},
	}

	mockRecommendationFour = api.ComplianceRecommendation{
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
		Violations:            []api.ComplianceViolation{},
	}
)
