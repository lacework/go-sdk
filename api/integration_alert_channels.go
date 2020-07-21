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

package api

// Enum for Alert Severity Levels
type AlertLevel int

const (
	CriticalAlertLevel AlertLevel = 1 // Critical only
	HighAlertLevel     AlertLevel = 2 // High and above
	MediumAlertLevel   AlertLevel = 3 // Medium and above
	LowAlertLevel      AlertLevel = 4 // Low and above
	AllAlertLevel      AlertLevel = 5 // Info and above (which is All of them)
)

// AlertLevels is the list of available alert levels
var AlertLevels = map[AlertLevel]string{
	CriticalAlertLevel: "Critical",
	HighAlertLevel:     "High",
	MediumAlertLevel:   "Medium",
	LowAlertLevel:      "Low",
	AllAlertLevel:      "All",
}

// String returns the string representation of an alert level
func (i AlertLevel) String() string {
	return AlertLevels[i]
}

// Int returns the int representation of an alert level
func (i AlertLevel) Int() int {
	return int(i)
}

// Valid returns whether the AlertLevel is valid or not
func (i AlertLevel) Valid() bool {
	_, ok := AlertLevels[i]
	return ok
}
