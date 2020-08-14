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

import "fmt"

// HostVulnerabilityService is a service that interacts with the vulnerabilities
// endpoints for the host space from the Lacework Server
type HostVulnerabilityService struct {
	client *Client
}

// Scan requests an on-demand vulnerability assessment of your software packages
// to determine if the packages contain any common vulnerabilities and exposures
//
// NOTE: Only packages managed by a package manager for supported OS's are reported
func (svc *HostVulnerabilityService) Scan(manifest string) (
	response map[string]interface{},
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiVulnerabilitiesScanPkgManifest, manifest, &response)
	return
}

func (svc *HostVulnerabilityService) ListCves() (
	response hostVulnListCvesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiVulnerabilitiesHostListCves, nil, &response)
	return
}

func (svc *HostVulnerabilityService) ListHostsWithCVE(id string) (
	response hostVulnListHostsResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesListHostsWithCveID, id)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *HostVulnerabilityService) GetHostAssessment(id string) (
	response hostVulnHostResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesHostAssessmentFromMachineID, id)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

type hostVulnHostResponse struct {
	Assessment HostVulnHostAssessment `json:"data"`
	Ok         bool                   `json:"ok"`
	Message    string                 `json:"message"`
}

type HostVulnHostAssessment struct {
	Host hostVulnHostDetail `json:"host"`
	CVEs []HostVulnCVE      `json:"vulnerabilities"`
}

type hostVulnListHostsResponse struct {
	Hosts   []HostVulnDetail `json:"data"`
	Ok      bool             `json:"ok"`
	Message string           `json:"message"`
}

type HostVulnDetail struct {
	Details  hostVulnHostDetail `json:"host"`
	Packages []hostVulnPackage  `json:"packages"`
	Summary  hostVulnCveSummary `json:"summary"`
}

type hostVulnHostDetail struct {
	Hostname      string      `json:"hostname"`
	MachineID     string      `json:"machine_id"`
	MachineStatus string      `json:"machine_status"`
	Tags          hostVulnTag `json:"tags"`
}

type hostVulnTag struct {
	Account        string `json:"Account"`
	AmiID          string `json:"AmiId"`
	ExternalIP     string `json:"ExternalIp"`
	Hostname       string `json:"Hostname"`
	InstanceID     string `json:"InstanceId"`
	InternalIP     string `json:"InternalIp"`
	LwTokenShort   string `json:"LwTokenShort"`
	SubnetID       string `json:"SubnetId"`
	VmInstanceType string `json:"VmInstanceType"`
	VmProvider     string `json:"VmProvider"`
	VpcID          string `json:"VpcId"`
	Zone           string `json:"Zone"`
	Arch           string `json:"arch"`
	Os             string `json:"os"`
}

type hostVulnListCvesResponse struct {
	CVEs    []HostVulnCVE `json:"data"`
	Ok      bool          `json:"ok"`
	Message string        `json:"message"`
}

type HostVulnCVE struct {
	ID       string             `json:"cve_id"`
	Packages []hostVulnPackage  `json:"packages"`
	Summary  hostVulnCveSummary `json:"summary"`
}

type hostVulnPackage struct {
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	Severity      string `json:"severity"`
	Status        string `json:"status"`
	Version       string `json:"version"`
	HostCount     string `json:"host_count"`
	PackageStatus string `json:"package_status"`
	CveLink       string `json:"cve_link"`
	CvssScore     string `json:"cvss_score"`
	CvssV2Score   string `json:"cvss_v_2_score"`
	CvssV3Score   string `json:"cvss_v_3_score"`
	//FirstSeenTime time.Time `json:"first_seen_time"`
	FixAvailable string `json:"fix_available"`
	FixedVersion string `json:"fixed_version"`
}

type hostVulnCveSummary struct {
	Severity             map[string]interface{} `json:"severity"` // @afiune not sure if there is a defined structure for this field
	TotalVulnerabilities int                    `json:"total_vulnerabilities"`
	//LastEvaluationTime   time.Time              `json:"last_evaluation_time"`
}

// Severity above looks like this:
//
//"Medium": {
//"fixable": 0,
//"vulnerabilities": 2
//}
//"Low": {
//"fixable": 0,
//"vulnerabilities": 3
//},
//"Negligible": {
//"fixable": 0,
//"vulnerabilities": 1
//}
