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

import (
	"fmt"
	"strings"

	"github.com/lacework/go-sdk/lwtime"
)

// HostVulnerabilityService is a service that interacts with the vulnerabilities
// endpoints for the host space from the Lacework Server
type HostVulnerabilityService struct {
	client *Client
}

// Scan requests an on-demand vulnerability assessment of your software packages
// to determine if the packages contain any common vulnerabilities and exposures
//
// NOTE: Only packages managed by a package manager for supported OS's are reported
func (svc *HostVulnerabilityService) Scan(manifest *PackageManifest) (
	response HostVulnScanPkgManifestResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST",
		apiVulnerabilitiesScanPkgManifest,
		manifest,
		&response,
	)

	if err == nil {
		// the API response coming from the Lacework server contains too much
		// information that could confuse our users, this function will parse
		// all the vulnerabilities and remove the non-matching ones
		response.CleanResponse()
	}

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
	Packages []HostVulnPackage  `json:"packages"`
	Summary  HostVulnCveSummary `json:"summary"`
}

type hostVulnHostDetail struct {
	Hostname      string      `json:"hostname"`
	MachineID     string      `json:"machine_id"`
	MachineStatus string      `json:"machine_status,omitempty"`
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
	Packages []HostVulnPackage  `json:"packages"`
	Summary  HostVulnCveSummary `json:"summary"`
}

type HostVulnPackage struct {
	Name                string          `json:"name"`
	Namespace           string          `json:"namespace"`
	Severity            string          `json:"severity"`
	Status              string          `json:"status,omitempty"`
	VulnerabilityStatus string          `json:"vulnerability_status,omitempty"`
	Version             string          `json:"version"`
	HostCount           string          `json:"host_count"`
	PackageStatus       string          `json:"package_status"`
	Description         string          `json:"description"`
	CveLink             string          `json:"cve_link"`
	CvssScore           string          `json:"cvss_score"`
	CvssV2Score         string          `json:"cvss_v_2_score"`
	CvssV3Score         string          `json:"cvss_v_3_score"`
	FirstSeenTime       lwtime.RFC1123Z `json:"first_seen_time"`
	FixAvailable        string          `json:"fix_available"`
	FixedVersion        string          `json:"fixed_version"`
}

func (assessment *HostVulnHostAssessment) VulnerabilityCounts() HostVulnCounts {
	var hostCounts = HostVulnCounts{}

	for _, cve := range assessment.CVEs {
		for _, pkg := range cve.Packages {

			switch strings.ToLower(pkg.Severity) {
			case "critical":
				hostCounts.Critical++
				hostCounts.Total++
				if pkg.FixedVersion != "" {
					hostCounts.CritFixable++
					hostCounts.TotalFixable++
				}
			case "high":
				hostCounts.High++
				hostCounts.Total++
				if pkg.FixedVersion != "" {
					hostCounts.HighFixable++
					hostCounts.TotalFixable++
				}
			case "medium":
				hostCounts.Medium++
				hostCounts.Total++
				if pkg.FixedVersion != "" {
					hostCounts.MedFixable++
					hostCounts.TotalFixable++
				}
			case "low":
				hostCounts.Low++
				hostCounts.Total++
				if pkg.FixedVersion != "" {
					hostCounts.LowFixable++
					hostCounts.TotalFixable++
				}
			default:
				hostCounts.Info++
				hostCounts.Total++
				if pkg.FixedVersion != "" {
					hostCounts.InfoFixable++
					hostCounts.TotalFixable++
				}
			}
		}
	}

	return hostCounts
}

// HighestSeverity returns the highest severity level vulnerability
func (h *HostVulnCounts) HighestSeverity() string {
	if h.Critical != 0 {
		return "critical"
	}
	if h.High != 0 {
		return "high"
	}
	if h.Medium != 0 {
		return "medium"
	}
	if h.Low != 0 {
		return "low"
	}
	return "unknown"
}

// HighestFixableSeverity returns the highest fixable severity level vulnerability
func (h *HostVulnCounts) HighestFixableSeverity() string {
	if h.CritFixable != 0 {
		return "critical"
	}
	if h.HighFixable != 0 {
		return "high"
	}
	if h.MedFixable != 0 {
		return "medium"
	}
	if h.LowFixable != 0 {
		return "low"
	}
	return "unknown"
}

// TotalFixableVulnerabilities returns the total number of vulnerabilities that have a fix available
func (h *HostVulnCounts) TotalFixableVulnerabilities() int32 {
	return h.TotalFixable
}

type HostVulnCounts struct {
	Critical     int32
	CritFixable  int32
	High         int32
	HighFixable  int32
	Medium       int32
	MedFixable   int32
	Low          int32
	LowFixable   int32
	Info         int32
	InfoFixable  int32
	Total        int32
	TotalFixable int32
}

type HostVulnSeverityCounts struct {
	Critical *HostVulnSeverityCountsDetails `json:"Critical"`
	High     *HostVulnSeverityCountsDetails `json:"High"`
	Medium   *HostVulnSeverityCountsDetails `json:"Medium"`
	Low      *HostVulnSeverityCountsDetails `json:"Low"`
	Info     *HostVulnSeverityCountsDetails `json:"Info"`
}

func (counts *HostVulnSeverityCounts) VulnerabilityCounts() HostVulnCounts {
	var hostCounts = HostVulnCounts{}

	if counts.Critical != nil {
		hostCounts.Critical += counts.Critical.Vulnerabilities
		hostCounts.CritFixable += counts.Critical.Fixable
		hostCounts.Total += counts.Critical.Vulnerabilities
		hostCounts.TotalFixable += counts.Critical.Fixable
	}

	if counts.High != nil {
		hostCounts.High += counts.High.Vulnerabilities
		hostCounts.HighFixable += counts.High.Fixable
		hostCounts.Total += counts.High.Vulnerabilities
		hostCounts.TotalFixable += counts.High.Fixable
	}

	if counts.Medium != nil {
		hostCounts.Medium += counts.Medium.Vulnerabilities
		hostCounts.MedFixable += counts.Medium.Fixable
		hostCounts.Total += counts.Medium.Vulnerabilities
		hostCounts.TotalFixable += counts.Medium.Fixable
	}

	if counts.Low != nil {
		hostCounts.Low += counts.Low.Vulnerabilities
		hostCounts.LowFixable += counts.Low.Fixable
		hostCounts.Total += counts.Low.Vulnerabilities
		hostCounts.TotalFixable += counts.Low.Fixable
	}

	if counts.Info != nil {
		hostCounts.Info += counts.Info.Vulnerabilities
		hostCounts.InfoFixable += counts.Info.Fixable
		hostCounts.Total += counts.Info.Vulnerabilities
		hostCounts.TotalFixable += counts.Info.Fixable
	}

	return hostCounts
}

type HostVulnSeverityCountsDetails struct {
	Fixable         int32 `json:"fixable"`
	Vulnerabilities int32 `json:"vulnerabilities"`
}

type HostVulnCveSummary struct {
	Severity             HostVulnSeverityCounts `json:"severity"`
	TotalVulnerabilities int                    `json:"total_vulnerabilities"`
	LastEvaluationTime   lwtime.EpochString     `json:"last_evaluation_time"`
}

type HostVulnScanPkgManifestResponse struct {
	Vulns   []HostScanPackageVulnDetails `json:"data"`
	Ok      bool                         `json:"ok"`
	Message string                       `json:"message"`
}

// CleanResponse will go over all the vulnerabilities from a package-manifest
// scan and remove the non-matching ones, leaving only the vulnerabilities
// that matter
func (scanPkgManifest *HostVulnScanPkgManifestResponse) CleanResponse() {
	filteredVulns := make([]HostScanPackageVulnDetails, 0)

	for _, vuln := range scanPkgManifest.Vulns {
		if !vuln.Match() {
			continue
		}
		filteredVulns = append(filteredVulns, vuln)
	}

	scanPkgManifest.Vulns = filteredVulns
}

func (scanPkgManifest *HostVulnScanPkgManifestResponse) VulnerabilityCounts() HostVulnCounts {
	var hostCounts = HostVulnCounts{}

	for _, vuln := range scanPkgManifest.Vulns {
		switch vuln.Severity {
		case "Critical":
			hostCounts.Critical++
			hostCounts.Total++
			if vuln.HasFix() {
				hostCounts.CritFixable++
				hostCounts.TotalFixable++
			}
		case "High":
			hostCounts.High++
			hostCounts.Total++
			if vuln.HasFix() {
				hostCounts.HighFixable++
				hostCounts.TotalFixable++
			}
		case "Medium":
			hostCounts.Medium++
			hostCounts.Total++
			if vuln.HasFix() {
				hostCounts.MedFixable++
				hostCounts.TotalFixable++
			}
		case "Low":
			hostCounts.Low++
			hostCounts.Total++
			if vuln.HasFix() {
				hostCounts.LowFixable++
				hostCounts.TotalFixable++
			}
		default:
			hostCounts.Info++
			hostCounts.Total++
			if vuln.HasFix() {
				hostCounts.InfoFixable++
				hostCounts.TotalFixable++
			}
		}
	}

	return hostCounts
}

func (scanPkg *HostScanPackageVulnDetails) ScoreString() string {
	if scanPkg.CVEProps.Metadata.NVD.CVSSv3.Score != 0 {
		return fmt.Sprintf("%.1f", scanPkg.CVEProps.Metadata.NVD.CVSSv3.Score)
	}

	if scanPkg.CVEProps.Metadata.NVD.CVSSv2.Score != 0 {
		return fmt.Sprintf("%.1f", scanPkg.CVEProps.Metadata.NVD.CVSSv2.Score)
	}

	return ""
}

type HostScanPackageVulnDetails struct {
	CVEProps struct {
		CveBatchID  string `json:"cve_batch_id"`
		Description string `json:"description"`
		Link        string `json:"link"`
		Metadata    struct {
			NVD struct {
				CVSSv2 struct {
					PublishedDateTime string  `json:"PublishedDateTime"`
					Score             float64 `json:"Score"`
					Vectors           string  `json:"Vectors"`
				} `json:"CVSSv2"`
				CVSSv3 struct {
					ExploitabilityScore float64 `json:"ExploitabilityScore"`
					ImpactScore         float64 `json:"ImpactScore"`
					Score               float64 `json:"Score"`
					Vectors             string  `json:"Vectors"`
				} `json:"CVSSv3"`
			} `json:"NVD"`
		} `json:"metadata"`
	} `json:"CVE_PROPS"`
	FeatureKey struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"FEATURE_KEY"`
	FixInfo   HostScanPackageVulnFixInfo `json:"FIX_INFO"`
	OsPkgInfo struct {
		Namespace     string `json:"namespace"`
		Os            string `json:"os"`
		OsVer         string `json:"os_ver"`
		Pkg           string `json:"pkg"`
		PkgVer        string `json:"pkg_ver"`
		VersionFormat string `json:"version_format"`
	} `json:"OS_PKG_INFO"`
	Props struct {
		EvalAlgo string `json:"eval_algo"`
	} `json:"PROPS"`
	Severity string `json:"SEVERITY"`
	Summary  struct {
		EvalCreatedTime          string `json:"eval_created_time"`
		EvalStatus               string `json:"eval_status"`
		NumFixableVuln           int    `json:"num_fixable_vuln"`
		NumFixableVulnBySeverity struct {
			Num1 int `json:"1"`
			Num2 int `json:"2"`
			Num3 int `json:"3"`
			Num4 int `json:"4"`
			Num5 int `json:"5"`
		} `json:"num_fixable_vuln_by_severity"`
		NumTotal          int `json:"num_total"`
		NumVuln           int `json:"num_vuln"`
		NumVulnBySeverity struct {
			Num1 int `json:"1"`
			Num2 int `json:"2"`
			Num3 int `json:"3"`
			Num4 int `json:"4"`
			Num5 int `json:"5"`
		} `json:"num_vuln_by_severity"`
	} `json:"SUMMARY"`
	VulnID string `json:"VULN_ID"`
}

func (v *HostScanPackageVulnDetails) Match() bool {
	if v.Summary.EvalStatus != "MATCH_VULN" {
		return false
	}
	if v.FixInfo.EvalStatus != "VULNERABLE" {
		return false
	}
	return true
}

func (v *HostScanPackageVulnDetails) HasFix() bool {
	return v.FixInfo.FixAvailable == 1
}

// PackageManifest is the representation of a package manifest
// that the Lacework API server expects when executing a scan
//
// {
//     "os_pkg_info_list": [
//         {
//             "os":"Ubuntu",
//             "os_ver":"18.04",
//             "pkg": "openssl",
//             "pkg_ver": "1.1.1-1ubuntu2.1~18.04.6"
//         }
//     ]
// }
type PackageManifest struct {
	OsPkgInfoList []OsPkgInfo `json:"os_pkg_info_list"`
}

type OsPkgInfo struct {
	Os     string `json:"os"`
	OsVer  string `json:"os_ver"`
	Pkg    string `json:"pkg"`
	PkgVer string `json:"pkg_ver"`
}

type HostScanPackageVulnFixInfo struct {
	CompareResult               int    `json:"compare_result"`
	EvalStatus                  string `json:"eval_status"`
	FixAvailable                int    `json:"fix_available"`
	FixedVersion                string `json:"fixed_version"`
	FixedVersionComparisonInfos []struct {
		CurrFixVer                         string `json:"curr_fix_ver"`
		IsCurrFixVerGreaterThanOtherFixVer string `json:"is_curr_fix_ver_greater_than_other_fix_ver"`
		OtherFixVer                        string `json:"other_fix_ver"`
	} `json:"fixed_version_comparison_infos"`
	FixedVersionComparisonScore int    `json:"fixed_version_comparison_score"`
	MaxPrefixMatchingLenScore   int    `json:"max_prefix_matching_len_score"`
	VersionInstalled            string `json:"version_installed"`
}
