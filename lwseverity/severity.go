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

// A package for Lacework severities
package lwseverity

import (
	"sort"
	"strings"
)

type Severity interface {
	GetSeverity() string
}

// Take a string representation of Lacework severity and
// return it's normalized integer and string values
// Critical Severity      => 1, "Critical"
// High Severity          => 2, "High"
// Medium Severity        => 3, "Medium"
// Low Severity           => 4, "Low"
// Informational Severity => 5, "Info"
// Unknown Severity       => 0, "Unknown"
func SeverityToProperTypes(severity string) (int, string) {
	switch strings.ToLower(severity) {
	case "1", "critical":
		return 1, "Critical"
	case "2", "high":
		return 2, "High"
	case "3", "medium":
		return 3, "Medium"
	case "4", "low":
		return 4, "Low"
	case "5", "info":
		return 5, "Info"
	default:
		return 0, "Unknown"
	}
}

// Returns true if the first severity not as critical as the second severity
//
// For instance:
//
// "info" is not as crtical as "low" (true)
// "medium" is as critical as "medium" (false)
// "high" is more critical than "medium" (false)
// "unknown" is more critical than "medium" (false)
// "medium" is not as critical as "unknown" (true)
func NotAsCritical(first, second string) bool {
	sevFirst, _ := SeverityToProperTypes(first)
	sevSecond, _ := SeverityToProperTypes(second)
	return sevFirst > sevSecond
}

// Returns true if the threshold is proper and the severity is
// greater than or equal to the threshold
//
// For instance:
//
// "medium" severity should be filtered for "high" threshold
// "critical" severity should NOT be filtered for "high" threshold
// "info" severity should NOT be filtered for "info" threshold
// invalid (unknown) severity should NOT be filtered for * threshold
// all severities should NOT be filtered for an invalid (unknown) threshold
func ShouldFilter(severity, threshold string) bool {
	sevThreshold, _ := SeverityToProperTypes(threshold)
	if sevThreshold == 0 {
		return false
	}
	return NotAsCritical(severity, threshold)
}

// Sort a slice of Severity interfaces from critical -> info
func SortSlice[S Severity](s []S) {
	sort.Slice(s, func(i, j int) bool {
		sevI, _ := SeverityToProperTypes(s[i].GetSeverity())
		sevJ, _ := SeverityToProperTypes(s[j].GetSeverity())
		return sevI < sevJ
	})
}

// Sort a slice of Severity interfaces from info -> critical
func SortSliceA[S Severity](s []S) {
	sort.Slice(s, func(i, j int) bool {
		sevI, _ := SeverityToProperTypes(s[i].GetSeverity())
		sevJ, _ := SeverityToProperTypes(s[j].GetSeverity())
		return sevI > sevJ
	})
}
