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

func TestDeepKeyValueExtract(t *testing.T) {
	tests := []struct {
		name     string
		mmap     map[string]interface{}
		expected map[string]string
	}{
		{"empty",
			map[string]interface{}{},
			map[string]string{},
		},
		{"all strings",
			map[string]interface{}{
				"str1": interface{}("foo"),
				"str2": interface{}("bar"),
			},
			map[string]string{
				"str1": "foo",
				"str2": "bar",
			},
		},
		{"numbers",
			map[string]interface{}{
				"k1": interface{}(1),
				"k2": interface{}(342634423432843943),
				"k3": interface{}(2.0),
			},
			map[string]string{
				"k1": "1",
				"k2": "342634423432843943",
				"k3": "2",
			},
		},
		{"bool",
			map[string]interface{}{
				"bool1": interface{}(true),
				"bool2": interface{}(false),
			},
			map[string]string{
				"bool1": "ENABLE",
				"bool2": "DISABLE",
			},
		},
		{"nested maps",
			map[string]interface{}{
				"map1": map[string]interface{}{
					"k1": interface{}(1),
					"k2": interface{}(342634423432843943),
					"k3": interface{}(2.0),
				},
				"map2": map[string]interface{}{
					"str1": interface{}("foo"),
					"str2": interface{}("bar"),
				},
				"bool2": interface{}(false),
			},
			map[string]string{
				"str1":  "foo",
				"str2":  "bar",
				"k1":    "1",
				"k2":    "342634423432843943",
				"k3":    "2",
				"bool2": "DISABLE",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, deepKeyValueExtract(tt.mmap), tt.expected)
		})
	}
}
