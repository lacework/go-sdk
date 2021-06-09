//
// Author:: Matt Cadorette (<matthew.cadorette@lacework.net>)
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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderSimpleCSV(t *testing.T) {
	expectedCsv := strings.TrimPrefix(`
KEY,VALUE
key1,value1
key2,value2
key3,value3
`, "\n")

	csv, _ := renderAsCSV(
		[]string{"KEY", "VALUE"},
		[][]string{
			{"key1", "value1"},
			{"key2", "value2"},
			{"key3", "value3"},
		},
	)
	assert.Equal(t,
		csv,
		expectedCsv,
		"csv is not being formatted correctly")
}

func TestRenderComplexCSV(t *testing.T) {
	expectedCsv := strings.TrimPrefix(`
KEY HEADER VALUE,"VALUE,TEST"
key1,"this is a value, from [a, b, c]"
key2,value2
key3,value3
`, "\n")

	csv, _ := renderAsCSV(
		[]string{"KEY\n HEADER VALUE", "VALUE,TEST"},
		[][]string{
			{"key1", "this is a value, from [a, b, c]"},
			{"key2", "value2"},
			{"key3", "value3"},
		},
	)
	assert.Equal(t,
		csv,
		expectedCsv,
		"csv is not being formatted correctly")
}

func TestCSVDataCleanup(t *testing.T) {
	data := csvCleanData([]string{"KEY\n HEADER\n VALUE", "VALUE,TEST\n"})
	assert.NotContains(t, strings.Join(data, ""), "\n", "data is not being cleaned up properly")
}
