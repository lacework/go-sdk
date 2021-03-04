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

package array

import "strings"

func ContainsStr(array []string, expected string) bool {
	for _, value := range array {
		if expected == value {
			return true
		}
	}
	return false
}

func ContainsStrCaseInsensitive(array []string, expected string) bool {
	for _, value := range array {
		if strings.EqualFold(expected, value) {
			return true
		}
	}
	return false
}

func ContainsPartialStr(array []string, expected string) bool {
	for _, value := range array {
		if strings.Contains(value, expected) {
			return true
		}
	}
	return false
}

func ContainsInt(array []int, expected int) bool {
	for _, value := range array {
		if expected == value {
			return true
		}
	}
	return false
}

func ContainsBool(array []bool, expected bool) bool {
	for _, value := range array {
		if expected == value {
			return true
		}
	}
	return false
}
