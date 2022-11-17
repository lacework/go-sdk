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
	"fmt"
	"sort"
	"strings"
)

type severity int

func (s severity) GetSeverity() string {
	return s.String()
}

const (
	// Unknown severity
	Unknown severity = iota
	// Critical severity
	Critical
	// High severity
	High
	// Medium severity
	Medium
	// Low severity
	Low
	// Informational severity
	Info
)

var severities = map[severity]string{
	Unknown:  "Unknown",
	Critical: "Critical",
	High:     "High",
	Medium:   "Medium",
	Low:      "Low",
	Info:     "Info",
}

// Get severity as a string type
func (s severity) String() string {
	return severities[s]
}

type validSeverities []severity

// A list of valid Lacework severities (critical, high, medium, low, info)
var ValidSeverities = validSeverities{Critical, High, Medium, Low, Info}

// Return a string representation of valid severities
// "critical, high, medium, low, info"
func (v validSeverities) String() string {
	s := ""

	for _, severity := range v {
		s += fmt.Sprintf("%s, ", strings.ToLower(severities[severity]))
	}

	return strings.TrimRight(s, ", ")
}

// Initialize a severity from string
func NewSeverity(s string) severity {
	switch strings.ToLower(s) {
	case "1", "critical":
		return Critical
	case "2", "high":
		return High
	case "3", "medium":
		return Medium
	case "4", "low":
		return Low
	case "5", "info":
		return Info
	default:
		return Unknown
	}
}

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
func Normalize(s string) (int, string) {
	severity := NewSeverity(s)
	return int(severity), severity.String()
}

// Take a string representation of Lacework severity and
// return whether it properly maps to a valid severity (not unknown)
func IsValid(s string) bool {
	return NewSeverity(s) != Unknown
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
	sevFirst, _ := Normalize(first)
	sevSecond, _ := Normalize(second)
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
	sevThreshold, _ := Normalize(threshold)
	if sevThreshold == 0 {
		return false
	}
	return NotAsCritical(severity, threshold)
}

// Sort a slice of Severity interfaces from critical -> info
func SortSlice[S Severity](s []S) {
	sort.Slice(s, func(i, j int) bool {
		sevI, _ := Normalize(s[i].GetSeverity())
		sevJ, _ := Normalize(s[j].GetSeverity())
		return sevI < sevJ
	})
}

// Sort a slice of Severity interfaces from info -> critical
func SortSliceA[S Severity](s []S) {
	sort.Slice(s, func(i, j int) bool {
		sevI, _ := Normalize(s[i].GetSeverity())
		sevJ, _ := Normalize(s[j].GetSeverity())
		return sevI > sevJ
	})
}
