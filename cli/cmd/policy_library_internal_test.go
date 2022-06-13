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

	"github.com/stretchr/testify/assert"
)

type policySyncOpsSummaryTest struct {
	Name     string
	Psos     []PolicySyncOperation
	Expected string
}

var policySyncOpsSummaryTests = []policySyncOpsSummaryTest{
	policySyncOpsSummaryTest{
		Name:     "empty",
		Expected: "Policy sync-library will create 0 policies, update 0 policies, create 0 queries, update 0 queries.",
	},
	policySyncOpsSummaryTest{
		Name: "full",
		Psos: []PolicySyncOperation{
			PolicySyncOperation{
				ID:          "larry",
				ContentType: "policy",
				Operation:   "create",
			},
			PolicySyncOperation{
				ID:          "curly",
				ContentType: "policy",
				Operation:   "update",
			},
			PolicySyncOperation{
				ID:          "moe",
				ContentType: "query",
				Operation:   "create",
			},
			PolicySyncOperation{
				ID:          "shremp",
				ContentType: "query",
				Operation:   "update",
			},
			PolicySyncOperation{
				ID:          "invalidcontenttype",
				ContentType: "invalid",
				Operation:   "create",
			},
			PolicySyncOperation{
				ID:          "invalidoperationtype",
				ContentType: "policy",
				Operation:   "invalid",
			},
		},
		Expected: "Policy sync-library will create 1 policies, update 1 policies, create 1 queries, update 1 queries.",
	},
}

func TestPolicySyncOpsSummary(t *testing.T) {
	for _, psost := range policySyncOpsSummaryTests {
		t.Run(psost.Name, func(t *testing.T) {
			actual := policySyncOpsSummary(psost.Psos)
			assert.Equal(t, psost.Expected, actual)
		})
	}
}

type policySyncOpsDetailsTest struct {
	Name     string
	Psos     []PolicySyncOperation
	Expected string
}

var policySyncOpsDetailsTests = []policySyncOpsDetailsTest{
	policySyncOpsDetailsTest{
		Name:     "empty",
		Expected: "Operation details:\n\n",
	},
	policySyncOpsDetailsTest{
		Name: "full",
		Psos: []PolicySyncOperation{
			PolicySyncOperation{
				ID:          "larry",
				ContentType: "policy",
				Operation:   "create",
			},
			PolicySyncOperation{
				ID:          "curly",
				ContentType: "policy",
				Operation:   "update",
			},
			PolicySyncOperation{
				ID:          "moe",
				ContentType: "query",
				Operation:   "create",
			},
			PolicySyncOperation{
				ID:          "shremp",
				ContentType: "query",
				Operation:   "update",
			},
			PolicySyncOperation{
				ID:          "invalidcontenttype",
				ContentType: "invalid",
				Operation:   "create",
			},
			PolicySyncOperation{
				ID:          "invalidoperationtype",
				ContentType: "policy",
				Operation:   "invalid",
			},
		},
		Expected: `Operation details:
  Policy larry will be created.
  Policy curly will be updated.
  Query moe will be created.
  Query shremp will be updated.

`,
	},
}

func TestPolicySyncOpsDetails(t *testing.T) {
	for _, psost := range policySyncOpsDetailsTests {
		t.Run(psost.Name, func(t *testing.T) {
			actual := policySyncOpsDetails(psost.Psos)
			assert.Equal(t, psost.Expected, actual)
		})
	}
}
