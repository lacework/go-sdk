// Author:: Darren Murray (<dmurray-lacework@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

// Test that correct input prompts are shown for each restraint data type
func TestPolicyExceptionCreateConstraintPrompts(t *testing.T) {
	policyExceptionTableTests := []struct {
		Name     string
		Input    []api.PolicyExceptionConfigurationConstraints
		Expected []PolicyExceptionSurveyQuestion
	}{
		{
			Name: "String type constraint",
			Input: []api.PolicyExceptionConfigurationConstraints{{
				DataType:   "String",
				FieldKey:   "resourceNames",
				MultiValue: false,
			}},
			Expected: []PolicyExceptionSurveyQuestion{{questions: []*survey.Question{{Name: "resourceNames", Prompt: &survey.Input{Message: "resourceNames value:"}, Validate: survey.Required}}}},
		},
		{
			Name: "String list type constraint",
			Input: []api.PolicyExceptionConfigurationConstraints{{
				DataType:   "String",
				FieldKey:   "resourceNames",
				MultiValue: true,
			}},
			Expected: []PolicyExceptionSurveyQuestion{{questions: []*survey.Question{{Name: "resourceNames", Prompt: &survey.Multiline{Message: "resourceNames values:"}, Validate: survey.Required}}}},
		},
		{
			Name: "Key value type constraint",
			Input: []api.PolicyExceptionConfigurationConstraints{{
				DataType:   "KVTagPair",
				FieldKey:   "resourceNames",
				MultiValue: false,
			}},
			Expected: []PolicyExceptionSurveyQuestion{{questions: []*survey.Question{{Name: "resourceNames-key", Prompt: &survey.Input{Message: "key:"}},
				{Name: "resourceNames-value", Prompt: &survey.Input{Message: "value:"}}}}},
		},
		{
			Name: "Key value list type constraint",
			Input: []api.PolicyExceptionConfigurationConstraints{{
				DataType:   "KVTagPair",
				FieldKey:   "resourceNames",
				MultiValue: true,
			}},
			Expected: []PolicyExceptionSurveyQuestion{{questions: []*survey.Question{{Name: "resourceNames-key", Prompt: &survey.Input{Message: "key:"}},
				{Name: "resourceNames-value", Prompt: &survey.Input{Message: "value:"}}}}},
		},
		{
			Name: "Invalid type constraint",
			Input: []api.PolicyExceptionConfigurationConstraints{{
				DataType:   "INVALID",
				FieldKey:   "resourceNames",
				MultiValue: true,
			}},
			Expected: []PolicyExceptionSurveyQuestion{},
		},
		{
			Name:     "Empty type constraint",
			Input:    []api.PolicyExceptionConfigurationConstraints{},
			Expected: []PolicyExceptionSurveyQuestion{},
		},
	}

	for _, ptt := range policyExceptionTableTests {
		t.Run(ptt.Name, func(t *testing.T) {
			questions := buildPromptAddExceptionConstraintListQuestions(ptt.Input)
			if len(questions) > 0 {
				assert.EqualValuesf(t, questions[0].questions[0].Prompt, ptt.Expected[0].questions[0].Prompt, fmt.Sprintf("%q showing incorrect prompt", ptt.Name))
			} else {
				assert.Equal(t, questions, ptt.Expected)
			}
		})
	}
}
