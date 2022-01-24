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

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	mockPolicyReferences []LCLReference = []LCLReference{
		LCLReference{
			ID:   "my_query",
			Type: "query",
			Path: "queries/my_query",
		},
		LCLReference{
			ID:   "my_policy",
			Type: "policy",
			Path: "policies/my_policy",
		},
	}
	malformedLCL LaceworkContentLibrary = LaceworkContentLibrary{
		Queries: map[string]LCLQuery{
			"my_query": LCLQuery{References: []LCLReference{}},
		},
		Policies: map[string]LCLPolicy{
			"my_policy": LCLPolicy{References: []LCLReference{}},
		},
	}
	mockLCL LaceworkContentLibrary = LaceworkContentLibrary{
		Queries: map[string]LCLQuery{
			"my_query": LCLQuery{
				References: []LCLReference{
					LCLReference{
						ID:   "my_query",
						Type: "query",
						Path: "queries/my_query",
					},
				},
			},
		},
		Policies: map[string]LCLPolicy{
			"my_policy": LCLPolicy{
				References: mockPolicyReferences,
			},
		},
	}
)

func TestGetPolicyReference(t *testing.T) {
	ref, err := getPolicyReference([]LCLReference{})
	assert.NotNil(t, err)
	ref, _ = getPolicyReference(mockPolicyReferences)
	assert.Equal(t, mockPolicyReferences[1], ref)
}

type getQueryRefTest struct {
	Name      string
	Library   LaceworkContentLibrary
	QueryID   string
	Reference LCLReference
	Error     error
}

var getQueryRefTests = []getQueryRefTest{
	getQueryRefTest{
		Name:  "NoQueryID",
		Error: errors.New("query ID must be provided"),
	},
	getQueryRefTest{
		Name:    "QueryNotFound",
		QueryID: "my_query",
		Error:   errors.New("query does not exist in library"),
	},
	getQueryRefTest{
		Name:    "QueryMalformed",
		Library: malformedLCL,
		QueryID: "my_query",
		Error:   errors.New("query exists but is malformed"),
	},
	getQueryRefTest{
		Name:    "QueryOK",
		Library: mockLCL,
		QueryID: "my_query",
		Error:   nil,
	},
}

func TestGetQueryRef(t *testing.T) {
	for _, gqrt := range getQueryRefTests {
		t.Run(gqrt.Name, func(t *testing.T) {
			actualRef, actualError := gqrt.Library.getReferenceForQuery(gqrt.QueryID)

			if gqrt.Error != nil {
				assert.Equal(t, gqrt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
				assert.Equal(t, gqrt.QueryID, actualRef.ID)
			}
		})
	}
}

type getPolicyRefsTest struct {
	Name       string
	Library    LaceworkContentLibrary
	PolicyID   string
	References []LCLReference
	Error      error
}

var getPolicyRefsTests = []getPolicyRefsTest{
	getPolicyRefsTest{
		Name:  "NoPolicyID",
		Error: errors.New("policy ID must be provided"),
	},
	getPolicyRefsTest{
		Name:     "PolicyNotFound",
		PolicyID: "my_policy",
		Error:    errors.New("policy does not exist in library"),
	},
	getPolicyRefsTest{
		Name:     "PolicyMalformed",
		Library:  malformedLCL,
		PolicyID: "my_policy",
		Error:    errors.New("policy exists but is malformed"),
	},
	getPolicyRefsTest{
		Name:     "PolicyOK",
		Library:  mockLCL,
		PolicyID: "my_policy",
		Error:    nil,
	},
}

func TestGetPolicyRefs(t *testing.T) {
	for _, gprt := range getPolicyRefsTests {
		t.Run(gprt.Name, func(t *testing.T) {
			actualRefs, actualError := gprt.Library.getReferencesForPolicy(gprt.PolicyID)

			if gprt.Error != nil {
				assert.Equal(t, gprt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
				assert.Equal(t, mockPolicyReferences, actualRefs)
			}
		})
	}
}
