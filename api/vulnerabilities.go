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

// VulnerabilitiesService is a service that interacts with the vulnerabilities
// endpoints from the Lacework Server
type VulnerabilitiesService struct {
	client    *Client
	Host      *HostVulnerabilityService
	Container *ContainerVulnerabilityService
}

// ValidVulnSeverities is a list of all valid severities in a vulnerability report
var ValidVulnSeverities = []string{"critical", "high", "medium", "low", "info"}

func NewVulnerabilityService(c *Client) *VulnerabilitiesService {
	return &VulnerabilitiesService{c,
		&HostVulnerabilityService{c},
		&ContainerVulnerabilityService{c},
	}
}

// VulnerabilityAssessment is used to provide common functions that are
// required by host or container vulnerability assessments, this is used
// to treat them both as equal
type VulnerabilityAssessment interface {
	HighestSeverity() string
	HighestFixableSeverity() string
	TotalFixableVulnerabilities() int32
}
