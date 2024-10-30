//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

package api

import "fmt"

// v2SoftwarePackagesVulnerabilityService is a service that interacts with the APIv2
// vulnerabilities endpoints for software packages
type v2SoftwarePackagesVulnerabilityService struct {
	client *Client
}

// Scan on-demand vulnerability assessment of your software packages
func (svc *v2SoftwarePackagesVulnerabilityService) Scan(manifest VulnerabilitiesPackageManifest) (
	response VulnerabilitySoftwarePackagesResponse, err error,
) {
	err = svc.client.RequestEncoderDecoder(
		"POST", apiV2VulnerabilitiesSoftwarePackagesScan,
		manifest, &response,
	)
	return
}

func (v *VulnerabilitySoftwarePackage) HasFix() bool {
	return v.FixInfo.FixAvailable == 1
}

func (v *VulnerabilitySoftwarePackage) IsVulnerable() bool {
	return v.FixInfo.EvalStatus == "VULNERABLE"
}

func (v *VulnerabilitySoftwarePackage) ScoreString() string {
	if v.CveProps.Metadata.Nvd.Cvssv3.Score != 0 {
		return fmt.Sprintf("%.1f", v.CveProps.Metadata.Nvd.Cvssv3.Score)
	}

	if v.CveProps.Metadata.Nvd.Cvssv2.Score != 0 {
		return fmt.Sprintf("%.1f", v.CveProps.Metadata.Nvd.Cvssv2.Score)
	}

	return ""
}

func (v *VulnerabilitySoftwarePackagesResponse) VulnerabilityCounts() HostVulnCounts {
	var hostCounts = HostVulnCounts{}

	for _, vuln := range v.Data {
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

type VulnerabilitySoftwarePackagesResponse struct {
	Data []VulnerabilitySoftwarePackage `json:"data"`
}

type VulnerabilitySoftwarePackage struct {
	OsPkgInfo struct {
		Namespace     string `json:"namespace"`
		Os            string `json:"os"`
		OsVer         string `json:"osVer"`
		Pkg           string `json:"pkg"`
		PkgVer        string `json:"pkgVer"`
		VersionFormat string `json:"versionFormat"`
	} `json:"osPkgInfo"`
	VulnID     string `json:"vulnId"`
	Severity   string `json:"severity"`
	FeatureKey struct {
		AffectedRange struct {
			End struct {
				Inclusive bool   `json:"inclusive"`
				Value     string `json:"value"`
			} `json:"end"`
			FixVersion string `json:"fixVersion"`
			Start      struct {
				Inclusive bool   `json:"inclusive"`
				Value     string `json:"value"`
			} `json:"start"`
		} `json:"affectedRange"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"featureKey"`
	CveProps struct {
		CveBatchId  string `json:"cveBatchId"`
		Description string `json:"description"`
		Link        string `json:"link"`
		Metadata    struct {
			Nvd struct {
				Cvssv2 struct {
					Publisheddatetime string  `json:"publisheddatetime"`
					Score             float64 `json:"score"`
					Vectors           string  `json:"vectors"`
				} `json:"cvssv2"`
				Cvssv3 struct {
					Exploitabilityscore float64 `json:"exploitabilityscore"`
					Impactscore         float64 `json:"impactscore"`
					Score               float64 `json:"score"`
					Vectors             string  `json:"vectors"`
				} `json:"cvssv3"`
			} `json:"nvd"`
		} `json:"metadata"`
	} `json:"cveProps"`
	FixInfo struct {
		CompareResult               int    `json:"compareResult"`
		EvalStatus                  string `json:"evalStatus"`
		FixAvailable                int    `json:"fixAvailable"`
		FixedVersion                string `json:"fixedVersion"`
		FixedVersionComparisonInfos []struct {
			CurrFixVer                         string `json:"currFixVer"`
			IsCurrFixVerGreaterThanOtherFixVer string `json:"isCurrFixVerGreaterThanOtherFixVer"`
			OtherFixVer                        string `json:"otherFixVer"`
		} `json:"fixedVersionComparisonInfos"`
		FixedVersionComparisonScore int    `json:"fixedVersionComparisonScore"`
		MaxPrefixMatchingLenScore   int    `json:"maxPrefixMatchingLenScore"`
		VersionInstalled            string `json:"versionInstalled"`
	} `json:"fixInfo"`
	Summary struct {
		EvalCreatedTime          string `json:"evalCreatedTime"`
		EvalStatus               string `json:"evalStatus"`
		NumFixableVuln           int    `json:"numFixableVuln"`
		NumFixableVulnBySeverity struct {
			Critical int `json:"1"`
			High     int `json:"2"`
			Medium   int `json:"3"`
			Low      int `json:"4"`
			Info     int `json:"5"`
		} `json:"numFixableVulnBySeverity"`
		NumTotal          int `json:"numTotal"`
		NumVuln           int `json:"numVuln"`
		NumVulnBySeverity struct {
			Critical int `json:"1"`
			High     int `json:"2"`
			Field3   int `json:"3"`
			Medium   int `json:"4"`
			Info     int `json:"5"`
		} `json:"numVulnBySeverity"`
	} `json:"summary"`
	Props struct {
		EvalAlgo string `json:"evalAlgo"`
	} `json:"props"`
}

type VulnerabilitiesPackageManifest struct {
	OsPkgInfoList []VulnerabilitiesOsPkgInfo `json:"osPkgInfoList"`
}

type VulnerabilitiesOsPkgInfo struct {
	Os     string `json:"os"`
	OsVer  string `json:"osVer"`
	Pkg    string `json:"pkg"`
	PkgVer string `json:"pkgVer"`
}
