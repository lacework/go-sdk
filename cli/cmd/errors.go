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

package cmd

import (
	"fmt"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwseverity"
	"github.com/pkg/errors"
)

type vulnerabilityPolicyError struct {
	SeverityRating        string
	FixableSeverityRating string
	FixableVulnCount      int32
	FailOnSeverity        string
	FailOnFixable         bool
	ExitCode              int
	Message               string
	Err                   error
}

func NewVulnerabilityPolicyErrorV2(
	assessment api.VulnerabilitiesContainersResponse,
	failOnSeverity string, failOnFixable bool,
) *vulnerabilityPolicyError {
	return &vulnerabilityPolicyError{
		SeverityRating:        assessment.HighestSeverity(),
		FixableSeverityRating: assessment.HighestFixableSeverity(),
		FixableVulnCount:      assessment.TotalFixableVulnerabilities(),
		FailOnSeverity:        failOnSeverity,
		FailOnFixable:         failOnFixable,
		// we use a default exit code that might change
		// during NonCompliant() or Compliant()
		ExitCode: 9,
	}
}

func NewVulnerabilityPolicyError(
	assessment api.VulnerabilityAssessment,
	failOnSeverity string, failOnFixable bool,
) *vulnerabilityPolicyError {
	return &vulnerabilityPolicyError{
		SeverityRating:        assessment.HighestSeverity(),
		FixableSeverityRating: assessment.HighestFixableSeverity(),
		FixableVulnCount:      assessment.TotalFixableVulnerabilities(),
		FailOnSeverity:        failOnSeverity,
		FailOnFixable:         failOnFixable,
		// we use a default exit code that might change
		// during NonCompliant() or Compliant()
		ExitCode: 9,
	}
}

// Example of an error message sent to the end-user:
//
// ERROR (FAIL-ON): fixable vulnerabilities found with threshold 'critical' (exit code: 9)
func (e *vulnerabilityPolicyError) Error() string {
	if e.ExitCode == 0 {
		return ""
	}
	return fmt.Sprintf("(FAIL-ON): %s (exit code: %d)", e.Message, e.ExitCode)
}

func (e *vulnerabilityPolicyError) Unwrap() error {
	return e.Err
}

func (e *vulnerabilityPolicyError) NonCompliant() bool {
	return !e.validate()
}

func (e *vulnerabilityPolicyError) Compliant() bool {
	return e.validate()
}

// validate returns true if the error policy is compliant,
// that is, when the provided assessment doesn't meet the
// thresholds. It returns false if the policy is NOT compliant
func (e *vulnerabilityPolicyError) validate() bool {
	severityRating, _ := lwseverity.Normalize(e.SeverityRating)
	fixableSeverityRating, _ := lwseverity.Normalize(e.FixableSeverityRating)
	threshold, _ := lwseverity.Normalize(e.FailOnSeverity)

	cli.Log.Debugw("validating policy",
		"severity_rating", severityRating,
		"fixable_severity_rating", fixableSeverityRating,
		"threshold", threshold,
		"fixable_vuln_count", e.FixableVulnCount,
	)
	if e.FailOnSeverity == "" && e.FailOnFixable && e.FixableVulnCount > 0 {
		e.Message = "fixable vulnerabilities found"
		return false
	}

	if e.FailOnFixable && e.FixableVulnCount > 0 && threshold >= fixableSeverityRating {
		e.Message = fmt.Sprintf(
			"fixable vulnerabilities found with threshold '%s'",
			e.FailOnSeverity)
		return false
	}

	if !e.FailOnFixable && (severityRating <= threshold && severityRating != 0) {
		e.Message = fmt.Sprintf(
			"vulnerabilities found with threshold '%s'",
			e.FailOnSeverity)
		return false
	}

	e.Message = "Compliant policy"
	e.ExitCode = 0
	return true
}

func yikes(msg string) error {
	return errors.Wrap(
		errors.New("something went pretty wrong here, contact support@lacework.net"),
		msg,
	)
}
