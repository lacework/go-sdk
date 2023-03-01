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

type querySourcesTableTest struct {
	Name     string
	Input    []api.Datasource
	Expected [][]string
}

var querySourcesTableTests = []querySourcesTableTest{
	querySourcesTableTest{
		Name:  "empty",
		Input: []api.Datasource{},
	},
	querySourcesTableTest{
		Name: "one",
		Input: []api.Datasource{
			api.Datasource{Name: "mySource", Description: "abc"},
		},
		Expected: [][]string{{"mySource\nabc"}},
	},
	querySourcesTableTest{
		Name: "sort",
		Input: []api.Datasource{
			api.Datasource{Name: "mySource", Description: "abc"},
			api.Datasource{Name: "aSource", Description: "xyz"},
		},
		Expected: [][]string{
			{"aSource\nxyz"},
			{"mySource\nabc"},
		},
	},
}

func TestQuerySourcesTable(t *testing.T) {
	for _, qtt := range querySourcesTableTests {
		t.Run(qtt.Name, func(t *testing.T) {
			out := querySourcesTable(qtt.Input)
			assert.Equal(t, out, qtt.Expected)
		})
	}
}
