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
	"encoding/csv"
	"fmt"
	"strings"
)

// Used to clean CSV inputs prior to rendering
func csvCleanData(input []string) []string {
	var data []string
	for _, h := range input {
		data = append(data, strings.Replace(h, "\n", "", -1))
	}
	return data
}

// Used to produce CSV output
func renderAsCSV(headers []string, data [][]string) string {
	csvOut := &strings.Builder{}
	csv := csv.NewWriter(csvOut)

	if len(headers) > 0 {
		csv.Write(csvCleanData(headers))
	}

	for _, record := range data {
		if err := csv.Write(csvCleanData(record)); err != nil {
			fmt.Printf("Failed to build csv output, got error: %s", err.Error())
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	csv.Flush()
	return csvOut.String()
}
