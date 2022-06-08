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

	"github.com/stretchr/testify/assert"
)

func TestVulnerabilityPolicyErrorFailOnFixable(t *testing.T) {
	mockAssessment := mockVulnAssessment{"high", "high", 2}
	mockPolicy := NewVulnerabilityPolicyError(&mockAssessment, "", true)

	if assert.Truef(t, mockPolicy.NonCompliant(), "policy should not be compliant") {
		assert.Equal(t, 9, mockPolicy.ExitCode)
		assert.Equal(t,
			"(FAIL-ON): fixable vulnerabilities found (exit code: 9)",
			mockPolicy.Error(),
		)
		assert.Falsef(t, mockPolicy.Compliant(), "policy should not be compliant")
	}
}

func TestVulnerabilityPolicyErrorFailOnSeverityWithFixable(t *testing.T) {
	mockAssessment := mockVulnAssessment{"high", "high", 2}
	mockPolicy := NewVulnerabilityPolicyError(&mockAssessment, "high", true)

	if assert.Truef(t, mockPolicy.NonCompliant(), "policy should not be compliant") {
		assert.Equal(t, 9, mockPolicy.ExitCode)
		assert.Equal(t,
			"(FAIL-ON): fixable vulnerabilities found with threshold 'high' (exit code: 9)",
			mockPolicy.Error(),
		)
		assert.Falsef(t, mockPolicy.Compliant(), "policy should not be compliant")
	}
}

func TestVulnerabilityPolicyErrorFailOnSeverityHigh(t *testing.T) {
	mockAssessment := mockVulnAssessment{"high", "high", 2}
	mockPolicy := NewVulnerabilityPolicyError(&mockAssessment, "high", false)

	if assert.Truef(t, mockPolicy.NonCompliant(), "policy should not be compliant") {
		assert.Equal(t, 9, mockPolicy.ExitCode)
		assert.Equal(t,
			"(FAIL-ON): vulnerabilities found with threshold 'high' (exit code: 9)",
			mockPolicy.Error(),
		)
		assert.Falsef(t, mockPolicy.Compliant(), "policy should not be compliant")
	}
}

func TestVulnerabilityPolicyErrorShouldNotFailOnSeverityCritical(t *testing.T) {
	mockAssessment := mockVulnAssessment{"medium", "medium", 1}
	mockPolicy := NewVulnerabilityPolicyError(&mockAssessment, "critical", false)

	assert.False(t, mockPolicy.NonCompliant(), "policy should be compliant")
	assert.True(t, mockPolicy.Compliant(), "policy should be compliant")
}

func TestVulnerabilityPolicyErrorShouldNotFailOnSeverityCriticalFailOnFixable(t *testing.T) {
	mockAssessment := mockVulnAssessment{"medium", "medium", 1}
	mockPolicy := NewVulnerabilityPolicyError(&mockAssessment, "critical", true)

	assert.False(t, mockPolicy.NonCompliant(), "policy should be compliant")
	assert.True(t, mockPolicy.Compliant(), "policy should be compliant")
}

func TestVulnerabilityPolicyErrorShouldNotFailOnSeverityCriticalFailWithNoVulns(t *testing.T) {
	mockAssessment := mockVulnAssessment{"unknown", "unknown", 0}
	mockPolicy := NewVulnerabilityPolicyError(&mockAssessment, "critical", false)

	assert.False(t, mockPolicy.NonCompliant(), "policy should be compliant")
	assert.True(t, mockPolicy.Compliant(), "policy should be compliant")
}

func TestVulnerabilityPolicyErrorShouldNotFailOnFixable(t *testing.T) {
	mockAssessment := mockVulnAssessment{"medium", "", 0}
	mockPolicy := NewVulnerabilityPolicyError(&mockAssessment, "", true)

	assert.False(t, mockPolicy.NonCompliant(), "policy should be compliant")
	assert.True(t, mockPolicy.Compliant(), "policy should be compliant")
}

type mockVulnAssessment struct {
	highestSeverity             string
	highestFixableSeverity      string
	totalFixableVulnerabilities int32
}

func (m *mockVulnAssessment) HighestSeverity() string {
	return m.highestSeverity
}
func (m *mockVulnAssessment) HighestFixableSeverity() string {
	return m.highestFixableSeverity
}
func (m *mockVulnAssessment) TotalFixableVulnerabilities() int32 {
	return m.totalFixableVulnerabilities
}
