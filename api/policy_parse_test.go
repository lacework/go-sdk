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

package api_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

var (
	newPolicyJSON = fmt.Sprintf(`{
	"evaluatorId": "%s",
	"policyId": "%s",
	"policyType": "%s",
	"queryId": "%s",
	"title": "%s",
	"enabled": %t,
	"description": "%s",
	"remediation": "%s",
	"severity": "%s",
	"evalFrequency": "%s",
	"alertEnabled": %t,
	"alertProfile": "%s"
}`, newPolicy.EvaluatorID, newPolicy.PolicyID, newPolicy.PolicyType, newPolicy.QueryID, newPolicy.Title,
		newPolicy.Enabled, newPolicy.Description, newPolicy.Remediation, newPolicy.Severity, newPolicy.EvalFrequency,
		newPolicy.AlertEnabled, newPolicy.AlertProfile)

	newPolicyYAML = fmt.Sprintf(`---
evaluatorId: %s
policyId: %s
policyType: %s
queryId: %s
title: %s
enabled: %t
description: %s
remediation: %s
severity: %s
evalFrequency: %s
alertEnabled: %t
alertProfile: %s
`, newPolicy.EvaluatorID, newPolicy.PolicyID, newPolicy.PolicyType, newPolicy.QueryID, newPolicy.Title,
		newPolicy.Enabled, newPolicy.Description, newPolicy.Remediation, newPolicy.Severity, newPolicy.EvalFrequency,
		newPolicy.AlertEnabled, newPolicy.AlertProfile)

	newPolicyNestedYAML = fmt.Sprintf(`---
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
  evalFrequency: %s
  alertEnabled: %t
  alertProfile: %s
`, newPolicy.EvaluatorID, newPolicy.PolicyID, newPolicy.PolicyType, newPolicy.QueryID, newPolicy.Title,
		newPolicy.Enabled, newPolicy.Description, newPolicy.Remediation, newPolicy.Severity, newPolicy.EvalFrequency,
		newPolicy.AlertEnabled, newPolicy.AlertProfile)

	falsePtr = false

	updatePolicy = api.UpdatePolicy{
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

	updatePolicyJSON = fmt.Sprintf(`{
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
}`, updatePolicy.EvaluatorID, updatePolicy.PolicyID, updatePolicy.PolicyType, updatePolicy.QueryID, updatePolicy.Title,
		false, updatePolicy.Description, updatePolicy.Remediation, updatePolicy.Severity, false,
		updatePolicy.AlertProfile)

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
`, updatePolicy.EvaluatorID, updatePolicy.PolicyID, updatePolicy.PolicyType, updatePolicy.QueryID, updatePolicy.Title,
		false, updatePolicy.Description, updatePolicy.Remediation, updatePolicy.Severity, false,
		updatePolicy.AlertProfile)

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
`, updatePolicy.EvaluatorID, updatePolicy.PolicyID, updatePolicy.PolicyType, updatePolicy.QueryID, updatePolicy.Title,
		false, updatePolicy.Description, updatePolicy.Remediation, updatePolicy.Severity, false,
		updatePolicy.AlertProfile)
)

type parseNewPolicyTest struct {
	Name     string
	Input    string
	Return   error
	Expected api.NewPolicy
}

var parseNewPolicyTests = []parseNewPolicyTest{
	parseNewPolicyTest{
		Name:     "empty-blob",
		Input:    "",
		Return:   errors.New("policy must be valid JSON or YAML"),
		Expected: api.NewPolicy{},
	},
	parseNewPolicyTest{
		Name:     "junk-blob",
		Input:    "this is junk",
		Return:   errors.New("policy must be valid JSON or YAML"),
		Expected: api.NewPolicy{},
	},
	parseNewPolicyTest{
		Name:     "partial-blob",
		Input:    "{",
		Return:   errors.New("policy must be valid JSON or YAML"),
		Expected: api.NewPolicy{},
	},
	parseNewPolicyTest{
		Name:     "json-blob",
		Input:    newPolicyJSON,
		Return:   nil,
		Expected: newPolicy,
	},
	parseNewPolicyTest{
		Name:     "yaml-blob",
		Input:    newPolicyYAML,
		Return:   nil,
		Expected: newPolicy,
	},
	parseNewPolicyTest{
		Name:     "yaml-nested-blob",
		Input:    newPolicyNestedYAML,
		Return:   nil,
		Expected: newPolicy,
	},
}

func TestParseNewPolicy(t *testing.T) {
	for _, pnpt := range parseNewPolicyTests {
		t.Run(pnpt.Name, func(t *testing.T) {
			actual, err := api.ParseNewPolicy(pnpt.Input)
			if pnpt.Return == nil {
				assert.Equal(t, pnpt.Return, err)
			} else {
				assert.Equal(t, pnpt.Return.Error(), err.Error())
			}
			assert.Equal(t, pnpt.Expected, actual)
		})
	}
}

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
		Input:    updatePolicyJSON,
		Return:   nil,
		Expected: updatePolicy,
	},
	parseUpdatePolicyTest{
		Name:     "yaml-blob",
		Input:    updatePolicyYAML,
		Return:   nil,
		Expected: updatePolicy,
	},
	parseUpdatePolicyTest{
		Name:     "yaml-nested-blob",
		Input:    updatePolicyNestedYAML,
		Return:   nil,
		Expected: updatePolicy,
	},
}

func TestParseUpdatePolicy(t *testing.T) {
	for _, pnpt := range parseUpdatePolicyTests {
		t.Run(pnpt.Name, func(t *testing.T) {
			actual, err := api.ParseUpdatePolicy(pnpt.Input)
			if pnpt.Return == nil {
				assert.Equal(t, pnpt.Return, err)
			} else {
				assert.Equal(t, pnpt.Return.Error(), err.Error())
			}
			assert.Equal(t, pnpt.Expected, actual)
		})
	}
}
