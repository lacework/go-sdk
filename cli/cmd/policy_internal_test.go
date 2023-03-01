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
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

type policyTagsTableTest struct {
	Name     string
	Input    []string
	Expected [][]string
}

var policyTagsTableTests = []policyTagsTableTest{
	policyTagsTableTest{
		Name:  "empty",
		Input: []string{},
	},
	policyTagsTableTest{
		Name:     "one",
		Input:    []string{"myTag"},
		Expected: [][]string{{"myTag"}},
	},
	policyTagsTableTest{
		Name:  "sort",
		Input: []string{"myTag", "aTag"},
		Expected: [][]string{
			{"aTag"},
			{"myTag"},
		},
	},
}

func TestPolicyTagsTable(t *testing.T) {
	for _, qtt := range policyTagsTableTests {
		t.Run(qtt.Name, func(t *testing.T) {
			out := policyTagsTable(qtt.Input)
			assert.Equal(t, out, qtt.Expected)
		})
	}
}

type policyTableTest struct {
	Name     string
	Input    []api.Policy
	Expected [][]string
}

var policyTableTests = []policyTableTest{
	policyTableTest{
		Name:  "empty",
		Input: []api.Policy{},
	},
	policyTableTest{
		Name: "one",
		Input: []api.Policy{
			api.Policy{PolicyID: "my-policy-1"},
		},
		Expected: [][]string{
			[]string{"my-policy-1", "", "", "Disabled", "Disabled", "", ""},
		},
	},
	policyTableTest{
		Name: "sort",
		Input: []api.Policy{
			api.Policy{PolicyID: "my-policy-10"},
			api.Policy{PolicyID: "my-policy-3"},
		},
		Expected: [][]string{
			[]string{"my-policy-3", "", "", "Disabled", "Disabled", "", ""},
			[]string{"my-policy-10", "", "", "Disabled", "Disabled", "", ""},
		},
	},
}

func TestPolicyTable(t *testing.T) {
	for _, qtt := range policyTableTests {
		t.Run(qtt.Name, func(t *testing.T) {
			out := policyTable(qtt.Input)
			assert.Equal(t, out, qtt.Expected)
		})
	}
}
