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
	"fmt"
	"strings"
	"testing"

	"github.com/lacework/go-sdk/v2/api"
	"github.com/stretchr/testify/assert"
)

func TestGenerateContainerVulnListCacheKey(t *testing.T) {
	cases := []struct {
		filterFlagsToHash cacheFiltersToBuildVulnContainerHash
		expectedCacheKey  string
	}{
		{cacheFiltersToBuildVulnContainerHash{
			"", "", "", []string{}, []string{}},
			"vulnerability/container/v2_3285545029616131935"},
		{cacheFiltersToBuildVulnContainerHash{
			"@d", "now", "", []string{}, []string{}},
			"vulnerability/container/v2_8666301743654077811"},
		{cacheFiltersToBuildVulnContainerHash{
			"@d", "now", "", []string{"repo1", "repo2"}, []string{"reg1"}},
			"vulnerability/container/v2_2929007791209551587"},
		{cacheFiltersToBuildVulnContainerHash{
			"", "now", "", []string{}, []string{"reg1"}},
			"vulnerability/container/v2_5320155942991519168"},
		// note, this is just like the first case
		{cacheFiltersToBuildVulnContainerHash{
			"", "", "", []string{}, []string{}},
			"vulnerability/container/v2_3285545029616131935"},
	}

	// first time we test all the the test cases
	for i, kase := range cases {
		t.Run(fmt.Sprintf("first case %d", i), func(t *testing.T) {
			vulCmdState.Start = kase.filterFlagsToHash.Start
			vulCmdState.End = kase.filterFlagsToHash.End
			vulCmdState.Range = kase.filterFlagsToHash.Range
			vulCmdState.Repositories = kase.filterFlagsToHash.Repositories
			vulCmdState.Registries = kase.filterFlagsToHash.Registries

			assert.Equal(t, kase.expectedCacheKey, generateContainerVulnListCacheKey())
		})
	}

	// second time should generate the same hashes
	for i, kase := range cases {
		t.Run(fmt.Sprintf("second case %d", i), func(t *testing.T) {
			vulCmdState.Start = kase.filterFlagsToHash.Start
			vulCmdState.End = kase.filterFlagsToHash.End
			vulCmdState.Range = kase.filterFlagsToHash.Range
			vulCmdState.Repositories = kase.filterFlagsToHash.Repositories
			vulCmdState.Registries = kase.filterFlagsToHash.Registries

			assert.Equal(t, kase.expectedCacheKey, generateContainerVulnListCacheKey())
		})
	}
}

func TestSeveritySummary(t *testing.T) {

	assessments := buildVulnCtrAssessmentSummary(mockVulnerabilityObservationsImageSummary())
	summaryString := severityCtrSummary(assessments[0].VulnCount, assessments[0].FixableCount)
	assert.Equal(t, "8 High 9 Fixable", summaryString)
}

func TestBuildCSVVulnCtrReportVulnerabilitiesListing(t *testing.T) {
	cli.EnableCSVOutput()
	defer func() { cli.csvOutput = false }()

	headers := []string{"Registry", "Repository", "Tag", "Last Scan", "Status", "Containers", "Vulnerabilities", "Image Digest"}
	assessments := buildVulnCtrAssessmentSummary(mockVulnerabilityObservationsImageSummary())
	filteredAssessments := applyVulnCtrFilters(assessments)
	assessmentOutput := assessmentSummaryToOutputFormat(filteredAssessments)
	rows := vulAssessmentsToTable(assessmentOutput)
	csv, err := renderAsCSV(headers, rows)
	if err != nil {
		panic(err)
	}

	expected := `
Registry,Repository,Tag,Last Scan,Status,Containers,Vulnerabilities,Image Digest
docker.io,lacework/jre,alpine-test,2025-08-06T13:05:11Z,Success,0,8 High 9 Fixable,sha256:a41ec54e6450ccc66d9f2ff975a0004d889349f3e8f5b086ebe8704e7ae962ac
`

	assert.Equal(t, strings.TrimPrefix(expected, "\n"), csv)
}

func mockVulnerabilityObservationsImageSummary() []api.VulnerabilityObservationsImageSummary {
	return []api.VulnerabilityObservationsImageSummary{
		{
			ContainerCount:           0,
			Digest:                   "sha256:a41ec54e6450ccc66d9f2ff975a0004d889349f3e8f5b086ebe8704e7ae962ac",
			ImageId:                  "sha256:b167326fa5f713a3cf7d742852967303b1b9301a147f84a0132ae58c47086fb4",
			LastScanTime:             "2025-08-06 13:05:11.000 Z",
			Registry:                 "docker.io",
			Repository:               "lacework/jre",
			Tag:                      "alpine-test",
			ScanStatus:               "Success",
			VulnCountCritical:        0,
			VulnCountCriticalFixable: 0,
			VulnCountHigh:            8,
			VulnCountHighFixable:     8,
			VulnCountMedium:          0,
			VulnCountMediumFixable:   0,
			VulnCountLow:             1,
			VulnCountLowFixable:      1,
			VulnCountInfo:            0,
			VulnCountInfoFixable:     0,
		},
	}
}
