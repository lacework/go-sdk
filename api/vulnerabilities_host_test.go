//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	subject "github.com/lacework/go-sdk/api"
)

func TestHostVulnHostAssessmentVulnerabilityCounts(t *testing.T) {

	assessment := subject.HostVulnHostAssessment{
		CVEs: []subject.HostVulnCVE{
			subject.HostVulnCVE{
				Packages: []subject.HostVulnPackage{
					subject.HostVulnPackage{
						Severity:     "Critical",
						FixedVersion: "1.2.3-mocked-fix",
					},
					subject.HostVulnPackage{Severity: "High"},
					subject.HostVulnPackage{Severity: "High"},
					subject.HostVulnPackage{Severity: "High"},
					subject.HostVulnPackage{Severity: "High"},
					subject.HostVulnPackage{Severity: "High"},
					subject.HostVulnPackage{Severity: "Low"},
				},
			},
		},
	}

	assessmentCounts := assessment.VulnerabilityCounts()
	assert.Equal(t, int32(1), assessmentCounts.Critical, "wrong critical vuln")
	assert.Equal(t, int32(1), assessmentCounts.CritFixable, "wrong critical fixable vuln")
	assert.Equal(t, int32(5), assessmentCounts.High, "wrong high vuln")
	assert.Equal(t, int32(0), assessmentCounts.HighFixable, "wrong high fixable vuln")
	assert.Equal(t, int32(0), assessmentCounts.Medium, "wrong medium vuln")
	assert.Equal(t, int32(0), assessmentCounts.MedFixable, "wrong medium fixable vuln")
	assert.Equal(t, int32(1), assessmentCounts.Low, "wrong low vuln")
	assert.Equal(t, int32(0), assessmentCounts.LowFixable, "wrong low fixable vuln")
	assert.Equal(t, int32(0), assessmentCounts.Info, "wrong info vuln")
	assert.Equal(t, int32(0), assessmentCounts.InfoFixable, "wrong info fixable vuln")

	assert.Equal(t, int32(7), assessmentCounts.Total, "wrong total vuln")
	assert.Equal(t, int32(1), assessmentCounts.TotalFixable, "wrong total fixable vuln")
}

func TestHostVulnSeverityCountsVulnerabilityCounts(t *testing.T) {

	assessment := subject.HostVulnSeverityCounts{
		High:   &subject.HostVulnSeverityCountsDetails{32, 45},
		Medium: &subject.HostVulnSeverityCountsDetails{1, 432},
		Info:   &subject.HostVulnSeverityCountsDetails{0, 923},
	}

	assessmentCounts := assessment.VulnerabilityCounts()
	assert.Equal(t, int32(0), assessmentCounts.Critical, "wrong critical vuln")
	assert.Equal(t, int32(0), assessmentCounts.CritFixable, "wrong critical fixable vuln")
	assert.Equal(t, int32(45), assessmentCounts.High, "wrong high vuln")
	assert.Equal(t, int32(32), assessmentCounts.HighFixable, "wrong high fixable vuln")
	assert.Equal(t, int32(432), assessmentCounts.Medium, "wrong medium vuln")
	assert.Equal(t, int32(1), assessmentCounts.MedFixable, "wrong medium fixable vuln")
	assert.Equal(t, int32(0), assessmentCounts.Low, "wrong low vuln")
	assert.Equal(t, int32(0), assessmentCounts.LowFixable, "wrong low fixable vuln")
	assert.Equal(t, int32(923), assessmentCounts.Info, "wrong info vuln")
	assert.Equal(t, int32(0), assessmentCounts.InfoFixable, "wrong info fixable vuln")

	assert.Equal(t, int32(1400), assessmentCounts.Total, "wrong total vuln")
	assert.Equal(t, int32(33), assessmentCounts.TotalFixable, "wrong total fixable vuln")
}

func TestHostVulnScanPkgManifestResponseVulnerabilityCounts(t *testing.T) {

	assessment := subject.HostVulnScanPkgManifestResponse{
		Vulns: []subject.HostScanPackageVulnDetails{
			subject.HostScanPackageVulnDetails{
				Severity: "Medium",
				FixInfo:  subject.HostScanPackageVulnFixInfo{FixAvailable: 1},
			},
			subject.HostScanPackageVulnDetails{Severity: "Medium"},
			subject.HostScanPackageVulnDetails{Severity: "Low"},
			subject.HostScanPackageVulnDetails{Severity: "Low"},
			subject.HostScanPackageVulnDetails{Severity: "Low"},
			subject.HostScanPackageVulnDetails{Severity: "Low"},
			subject.HostScanPackageVulnDetails{Severity: "Low"},
			subject.HostScanPackageVulnDetails{Severity: "Low"},
			subject.HostScanPackageVulnDetails{
				Severity: "Info",
				FixInfo:  subject.HostScanPackageVulnFixInfo{FixAvailable: 1},
			},
		},
	}

	assessmentCounts := assessment.VulnerabilityCounts()
	assert.Equal(t, int32(0), assessmentCounts.Critical, "wrong critical vuln")
	assert.Equal(t, int32(0), assessmentCounts.CritFixable, "wrong critical fixable vuln")
	assert.Equal(t, int32(0), assessmentCounts.High, "wrong high vuln")
	assert.Equal(t, int32(0), assessmentCounts.HighFixable, "wrong high fixable vuln")
	assert.Equal(t, int32(2), assessmentCounts.Medium, "wrong medium vuln")
	assert.Equal(t, int32(1), assessmentCounts.MedFixable, "wrong medium fixable vuln")
	assert.Equal(t, int32(6), assessmentCounts.Low, "wrong low vuln")
	assert.Equal(t, int32(0), assessmentCounts.LowFixable, "wrong low fixable vuln")
	assert.Equal(t, int32(1), assessmentCounts.Info, "wrong info vuln")
	assert.Equal(t, int32(1), assessmentCounts.InfoFixable, "wrong info fixable vuln")

	assert.Equal(t, int32(9), assessmentCounts.Total, "wrong total vuln")
	assert.Equal(t, int32(2), assessmentCounts.TotalFixable, "wrong total fixable vuln")
}
