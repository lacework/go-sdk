//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

var (
	falsePtr         = false
	updateTestPolicy = api.UpdatePolicy{
		EvaluatorID:  "Cloudtrail",
		PolicyType:   "Violation",
		QueryID:      "LW_CLI_AWS_CTA_IntegrationTest",
		Title:        "My Policy Title",
		Enabled:      &falsePtr,
		Description:  "My Policy Description",
		Remediation:  "Check yourself...",
		Severity:     "high",
		Limit:        nil,
		AlertEnabled: &falsePtr,
		AlertProfile: "LW_CloudTrail_Alerts",
	}
	updateTestPolicyJSON = fmt.Sprintf(`{
	"evaluatorId": "%s",
	"policyId": "%s",
	"policyType": "%s",
	"queryId": "%s",
	"title": "%s",
	"enabled": %t,
	"description": "%s",
	"remediation": "%s",
	"severity": "%s",
	"alertEnabled": %t,
	"alertProfile": "%s"
}`, updateTestPolicy.EvaluatorID, updateTestPolicy.PolicyID, updateTestPolicy.PolicyType, updateTestPolicy.QueryID, updateTestPolicy.Title,
		false, updateTestPolicy.Description, updateTestPolicy.Remediation, updateTestPolicy.Severity, false,
		updateTestPolicy.AlertProfile)
	updatePolicyYAML = fmt.Sprintf(`---
evaluatorId: %s
policyId: %s
policyType: %s
queryId: %s
title: %s
enabled: %t
description: %s
remediation: %s
severity: %s
alertEnabled: %t
alertProfile: %s
`, updateTestPolicy.EvaluatorID, updateTestPolicy.PolicyID, updateTestPolicy.PolicyType, updateTestPolicy.QueryID, updateTestPolicy.Title,
		false, updateTestPolicy.Description, updateTestPolicy.Remediation, updateTestPolicy.Severity, false,
		updateTestPolicy.AlertProfile)
	updatePolicyNestedYAML = fmt.Sprintf(`---
policies:
- evaluatorId: %s
  policyId: %s
  policyType: %s
  queryId: %s
  title: %s
  enabled: %t
  description: %s
  remediation: %s
  severity: %s
  alertEnabled: %t
  alertProfile: %s
`, updateTestPolicy.EvaluatorID, updateTestPolicy.PolicyID, updateTestPolicy.PolicyType, updateTestPolicy.QueryID, updateTestPolicy.Title,
		false, updateTestPolicy.Description, updateTestPolicy.Remediation, updateTestPolicy.Severity, false,
		updateTestPolicy.AlertProfile)
)

type parseUpdatePolicyTest struct {
	Name     string
	Input    string
	Return   error
	Expected api.UpdatePolicy
}

var parseUpdatePolicyTests = []parseUpdatePolicyTest{
	parseUpdatePolicyTest{
		Name:     "empty-blob",
		Input:    "",
		Return:   errors.New("policy must be valid JSON or YAML"),
		Expected: api.UpdatePolicy{},
	},
	parseUpdatePolicyTest{
		Name:     "junk-blob",
		Input:    "this is junk",
		Return:   errors.New("policy must be valid JSON or YAML"),
		Expected: api.UpdatePolicy{},
	},
	parseUpdatePolicyTest{
		Name:     "partial-blob",
		Input:    "{",
		Return:   errors.New("policy must be valid JSON or YAML"),
		Expected: api.UpdatePolicy{},
	},
	parseUpdatePolicyTest{
		Name:     "json-blob",
		Input:    updateTestPolicyJSON,
		Return:   nil,
		Expected: updateTestPolicy,
	},
	parseUpdatePolicyTest{
		Name:     "yaml-blob",
		Input:    updatePolicyYAML,
		Return:   nil,
		Expected: updateTestPolicy,
	},
	parseUpdatePolicyTest{
		Name:     "yaml-nested-blob",
		Input:    updatePolicyNestedYAML,
		Return:   nil,
		Expected: updateTestPolicy,
	},
}

func TestParseUpdatePolicy(t *testing.T) {
	for _, pnpt := range parseUpdatePolicyTests {
		t.Run(pnpt.Name, func(t *testing.T) {
			actual, err := parseUpdatePolicy(pnpt.Input)
			if pnpt.Return == nil {
				assert.Equal(t, pnpt.Return, err)
			} else {
				assert.Equal(t, pnpt.Return.Error(), err.Error())
			}
			assert.Equal(t, pnpt.Expected, actual)
		})
	}
}
