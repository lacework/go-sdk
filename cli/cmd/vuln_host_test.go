//
// Author:: Darren Murray (<dmurray-lacework@lacework.net>)
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

package cmd

import (
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestListCvesFilterSeverity(t *testing.T) {
	vulCmdState.Severity = "critical"
	defer clearVulnFilters()

	mockCves := []api.HostVulnCVE{mockCveOne}
	result, output := hostVulnCVEs(mockCves)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "\n1 of 2 cve(s) showing \n")
}

func TestShowAssessmentFilterSeverity(t *testing.T) {
	vulCmdState.Severity = "critical"
	defer clearVulnFilters()

	mockCves := []api.HostVulnCVE{mockCveOne}
	result, output := hostVulnCVEs(mockCves)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "\n1 of 2 cve(s) showing \n")
}

func TestShowAssessmentFilterSeverityWithPackages(t *testing.T) {
	vulCmdState.Severity = "critical"
	vulCmdState.Packages = true
	defer clearVulnFilters()

	mockCves := []api.HostVulnCVE{mockCveOne}
	result, output := hostVulnPackagesTable(mockCves, true)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, output, "1 of 2 package(s) showing \n")
}

func clearVulnFilters() {
	vulCmdState.Severity = ""
}

var mockCveOne = api.HostVulnCVE{
	ID:       "TestID",
	Packages: []api.HostVulnPackage{mockPackageOne, mockPackageTwo},
	Summary: api.HostVulnCveSummary{
		Severity: api.HostVulnSeverityCounts{
			Critical: &api.HostVulnSeverityCountsDetails{
				Fixable:         1,
				Vulnerabilities: 1,
			},
			High: &api.HostVulnSeverityCountsDetails{
				Fixable:         1,
				Vulnerabilities: 1,
			},
		},
		TotalVulnerabilities: 2,
		LastEvaluationTime:   api.Json16DigitTime{},
	},
}

var mockPackageOne = api.HostVulnPackage{
	Name:         "test",
	Namespace:    "rhel:8",
	Severity:     "High",
	Status:       "Active",
	HostCount:    "1",
	FixAvailable: "1",
}

var mockPackageTwo = api.HostVulnPackage{
	Name:         "test2",
	Namespace:    "rhel:8",
	Severity:     "Critical",
	Status:       "Active",
	HostCount:    "1",
	FixAvailable: "1",
}
