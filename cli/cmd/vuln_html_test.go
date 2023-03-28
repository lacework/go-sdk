//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

func TestVulContainerImageLayersToHTML(t *testing.T) {
	var (
		expectedVulnIDs = []string{"CVE-2021-24215", "CVE-2020-24215"}
		resp            = mockVulnerabilityAssessment()
		subject         = vulContainerImageLayersToHTML(resp.Data)
	)
	if assert.Equal(t, 2, len(subject), "wrong number of vulnerabilities") {
		for _, vuln := range subject {
			assert.NotEmpty(t, vuln.Layer, "the layer should not be empty, did response change?")
			assert.True(t, vuln.UseNoScore, "do we have CVSS scores? update, please")
			assert.Contains(t, expectedVulnIDs, vuln.CVE,
				"missing CVE, check HTML output")
		}
	}
}
