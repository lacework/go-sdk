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
	malformedLCL LaceworkContentLibrary = LaceworkContentLibrary{
		Queries: map[string]LCLQuery{
			"my_query": LCLQuery{References: []LCLReference{}},
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
	}
)

type getRefTest struct {
	Name      string
	Library   LaceworkContentLibrary
	QueryID   string
	Reference LCLReference
	Error     error
}

var getRefTests = []getRefTest{
	getRefTest{
		Name:  "NoQueryID",
		Error: errors.New("query ID must be provided"),
	},
	getRefTest{
		Name:    "QueryNotFound",
		QueryID: "my_query",
		Error:   errors.New("query does not exist in library"),
	},
	getRefTest{
		Name:    "QueryMalformed",
		Library: malformedLCL,
		QueryID: "my_query",
		Error:   errors.New("query exists but is malformed"),
	},
	getRefTest{
		Name:    "QueryOK",
		Library: mockLCL,
		QueryID: "my_query",
		Error:   nil,
	},
}

func TestGetRef(t *testing.T) {
	for _, grt := range getRefTests {
		t.Run(grt.Name, func(t *testing.T) {
			actualRef, actualError := grt.Library.getReferenceForQuery(grt.QueryID)

			if grt.Error != nil {
				assert.Equal(t, grt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
				assert.Equal(t, grt.QueryID, actualRef.ID)
			}
		})
	}
}
