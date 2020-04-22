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
)

// VulnerabilitiesService is a service that interacts with the vulnerabilities
// endpoints from the Lacework Server
type VulnerabilitiesService struct {
	client *Client
}

// Scan triggers a vulnerability scan to the provider registry, repository, and
// tag provided. This function calls the underlaying API endpoint that assumes
// that the container repository has been already integrated with the platform.
func (svc *VulnerabilitiesService) Scan(registry, repository, tagOrHash string) (
	response vulScanResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST",
		apiVulnerabilitiesScan,
		vulScanRequest{registry, repository, tagOrHash},
		&response,
	)
	return
}

type vulScanRequest struct {
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

type vulScanResponse struct {
	Data struct {
		Status    string `json:"status"`
		RequestID string `json:"requestId"`
	} `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (svc *VulnerabilitiesService) ScanStatus(requestID string) (
	response vulScanStatusResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesScanStatus, requestID)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

type vulScanStatusResponse struct {
	Data    VulContainerReport `json:"data"`
	Ok      bool               `json:"ok"`
	Message string             `json:"message"`
}

func (scan *vulScanStatusResponse) CheckStatus() string {
	if !scan.Ok {
		return fmt.Sprintf("there is a problem with the vulnerability scan: %s", scan.Message)
	}

	if scan.Data.Status != "" {
		// @afiune as far as I can see, the three status we could have are:
		// * Scanning
		// * Failed
		// * NotFound
		//
		// Where is "Success"? Not here. Why?
		return scan.Data.Status
	}

	// If the scan is not running, that means, it completed running, now the
	// status of the scan changes to be stored in "ScanStatus" :sadpanda:
	if scan.Data.ScanStatus != "" {
		return scan.Data.ScanStatus
	}

	return "Unknown"
}

func (svc *VulnerabilitiesService) ReportFromID(imageID string) (
	response VulContainerReportResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesReportFromID, imageID)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *VulnerabilitiesService) ReportFromDigest(imageDigest string) (
	response VulContainerReportResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesReportFromDigest, imageDigest)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

type VulContainerReportResponse struct {
	Data    VulContainerReport `json:"data"`
	Ok      bool               `json:"ok"`
	Message string             `json:"message"`
}

func (res *VulContainerReportResponse) CheckStatus() string {
	if !res.Ok {
		return fmt.Sprintf("there is a problem with the vulnerability report: %s", res.Message)
	}

	if res.Data.ScanStatus != "" {
		return res.Data.ScanStatus
	}

	if res.Data.Status != "" {
		return res.Data.Status
	}

	return "Unknown"
}

type VulContainerReport struct {
	TotalVulnerabilities    int32             `json:"total_vulnerabilities"`
	CriticalVulnerabilities int32             `json:"critical_vulnerabilities"`
	HighVulnerabilities     int32             `json:"high_vulnerabilities"`
	MediumVulnerabilities   int32             `json:"medium_vulnerabilities"`
	LowVulnerabilities      int32             `json:"low_vulnerabilities"`
	InfoVulnerabilities     int32             `json:"info_vulnerabilities"`
	FixableVulnerabilities  int32             `json:"fixable_vulnerabilities"`
	LastEvaluationTime      string            `json:"last_evaluation_time"`
	Image                   vulContainerImage `json:"image"`

	// @afiune these two parameters, Status and Message will appear when
	// the vulnerability scan is still running. ugh. why?
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`

	// ScanStatus is a property that will appear when the vulnerability scan finished
	// running, this status indicates whether the scan finished successfully or not
	ScanStatus string `json:"scan_status,omitempty"`

	// @afiune why we can't parse the time?
	//LastEvaluationTime      time.Time         `json:"last_evaluation_time"`
}

func (report *VulContainerReport) VulCountsTable() [][]string {
	return [][]string{
		[]string{"Critical", fmt.Sprint(report.CriticalVulnerabilities),
			fmt.Sprint(report.VulFixableCount("critical"))},
		[]string{"High", fmt.Sprint(report.HighVulnerabilities),
			fmt.Sprint(report.VulFixableCount("high"))},
		[]string{"Medium", fmt.Sprint(report.MediumVulnerabilities),
			fmt.Sprint(report.VulFixableCount("medium"))},
		[]string{"Low", fmt.Sprint(report.LowVulnerabilities),
			fmt.Sprint(report.VulFixableCount("low"))},
		[]string{"Info", fmt.Sprint(report.InfoVulnerabilities),
			fmt.Sprint(report.VulFixableCount("info"))},
	}
}

func (report *VulContainerReport) VulFixableCount(severity string) int32 {
	// @afiune check valid severity
	severity = strings.ToLower(severity)

	if len(report.Image.ImageLayers) == 0 {
		return 0
	}

	var fixable int32 = 0
	for _, layer := range report.Image.ImageLayers {
		for _, pkg := range layer.Packages {
			for _, vul := range pkg.Vulnerabilities {
				if vul.Severity == severity && vul.FixVersion != "" {
					fixable++
				}
			}
		}
	}
	return fixable
}

type vulContainerImage struct {
	ImageInfo   vulContainerImageInfo    `json:"image_info"`
	ImageLayers []vulContainerImageLayer `json:"image_layers"`
}

func (image *vulContainerImage) Table() [][]string {
	info := image.ImageInfo
	return [][]string{
		[]string{"ID", info.ImageID},
		[]string{"Digest", info.ImageDigest},
		[]string{"Registry", info.Registry},
		[]string{"Repository", info.Repository},
		[]string{"Size", ByteCountBinary(info.Size)},
		[]string{"Created At", info.CreatedTime},
		[]string{"Tags", strings.Join(info.Tags, ",")},
	}
}

func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

type vulContainerImageInfo struct {
	ImageDigest string   `json:"image_digest"`
	ImageID     string   `json:"image_id"`
	Registry    string   `json:"registry"`
	Repository  string   `json:"repository"`
	CreatedTime string   `json:"created_time"`
	Size        int64    `json:"size"`
	Tags        []string `json:"tags"`
}

type vulContainerImageLayer struct {
	Hash      string                `json:"hash"`
	CreatedBy string                `json:"created_by"`
	Packages  []vulContainerPackage `json:"packages"`
}

type vulContainerPackage struct {
	Name            string                   `json:"name"`
	Namespace       string                   `json:"namescape"`
	Version         string                   `json:"version"`
	Vulnerabilities []containerVulnerability `json:"vulnerabilities"`

	// @afiune maybe these fields are host related information and not container
	FixAvailable  string `json:"fix_available,omitempty"`
	FixedVersion  string `json:"fixed_version,omitempty"`
	HostCount     string `json:"host_count,omitempty"`
	Severity      string `json:"severity,omitempty"`
	Status        string `json:"status,omitempty"`
	CveLink       string `json:"cve_link,omitempty"`
	CveScore      string `json:"cve_score,omitempty"`
	CvssV3Score   string `json:"cvss_v3_score,omitempty"`
	CvssV2Score   string `json:"cvss_v2_score,omitempty"`
	FirstSeenTime string `json:"first_seen_time,omitempty"`
}

type containerVulnerability struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	Link        string                 `json:"link"`
	FixVersion  string                 `json:"fix_version"`
	Metadata    map[string]interface{} `json:"metadata"`
}
