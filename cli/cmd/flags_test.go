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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFiltersFrom(t *testing.T) {
	cases := []struct {
		subjectStruct   interface{}
		expectedFilters []string
	}{
		// empty struct will return empty filters
		{struct{}{}, []string{}},

		// struct without 'json' tags will return empty filters
		{struct{ Foo string }{}, []string{}},

		// struct with lots of fields that has a 'json' tags
		{
			struct {
				Foo   string                 `json:"foo"`
				Bar   int                    `json:"bar"`
				Array []string               `json:"array"`
				Map   map[string]interface{} `json:"map"`
			}{},
			[]string{"foo", "bar", "array", "map"},
		},

		// struct mixed with and without 'json' tags
		{
			struct {
				FooWithTag    string `json:"foo_with_tag"`
				BarWithoutTag int
			}{},
			[]string{"foo_with_tag"},
		},

		// struct with deep nested structs
		{
			struct {
				Foo    string                 `json:"foo"`
				Bar    int                    `json:"bar"`
				Array  []string               `json:"array"`
				Map    map[string]interface{} `json:"map"`
				Struct struct {
					A    string `json:"a"`
					B    string `json:"b"`
					C    string
					Deep struct {
						X      string `json:"x"`
						Y      string
						Z      string `json:"z"`
						Nested struct {
							Thing string `json:"thing"`
						} `json:"nested"`
					} `json:"deep"`
				} `json:"struct"`
				BarWithoutTag int
			}{},
			[]string{
				"foo", "bar", "array", "map",
				"struct.a", "struct.b", "struct.deep.x", "struct.deep.z",
				"struct.deep.nested.thing"},
		},
	}

	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actualFilters := getFiltersFrom(kase.subjectStruct, "")
			assert.Equalf(t, actualFilters, kase.expectedFilters, "wrong filters")
		})
	}
}
