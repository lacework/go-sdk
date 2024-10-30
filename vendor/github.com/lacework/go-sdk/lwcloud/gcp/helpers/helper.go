//
// Author:: Nicholas Schmeller (<nick.schmeller@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

package helpers

import (
	"fmt"
	"os"
	"strings"
)

const (
	GCP_PROJECT_TYPE      = "PROJECT"
	GCP_ORGANIZATION_TYPE = "ORGANIZATION"
)

func SkipEntry(Entry string, skipList, allowList map[string]bool) bool {
	if skipList != nil {
		if _, skip := skipList[Entry]; skip {
			return true
		}
	}

	if allowList != nil {
		if _, allow := allowList[Entry]; allow {
			return false
		} else {
			// skip all other entries
			return true
		}
	}

	return false
}

func GetGcpFormatedLabel(in string) string {
	lower := strings.ToLower(in)
	out := strings.ReplaceAll(lower, ":", "-")
	out = strings.ReplaceAll(out, ".", "-")
	out = strings.ReplaceAll(out, "\"", "-")
	out = strings.ReplaceAll(out, "{", "-")
	out = strings.ReplaceAll(out, "}", "-")
	return out
}

func CombineErrors(old error, new error) error {
	if new == nil {
		return old
	}

	if old == nil {
		return new
	}

	return fmt.Errorf("%s, %s", old.Error(), new.Error())
}

func IsProjectScanScope() bool {
	return os.Getenv("GCP_SCAN_SCOPE") == GCP_PROJECT_TYPE
}

func IsOrgScanScope() bool {
	return os.Getenv("GCP_SCAN_SCOPE") == GCP_ORGANIZATION_TYPE
}
