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
	"time"

	"github.com/lacework/go-sdk/lwtime"

	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/internal/array"
)

// ContainerVulnerabilityService is a service that interacts with the vulnerabilities
// endpoints for the container space from the Lacework Server
type ContainerVulnerabilityService struct {
	client *Client
}

// Scan triggers a container vulnerability scan to the provider registry, repository,
// and tag provided. This function calls the underlaying API endpoint that assumes
// that the container repository has been already integrated with the platform.
func (svc *ContainerVulnerabilityService) Scan(registry, repository, tagOrHash string) (
	response vulnContainerScanResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST",
		apiVulnerabilitiesContainerScan,
		vulnContainerScanRequest{registry, repository, tagOrHash},
		&response,
	)
	return
}

func (svc *ContainerVulnerabilityService) ScanStatus(requestID string) (
	response vulnContainerScanStatusResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesContainerScanStatus, requestID)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *ContainerVulnerabilityService) AssessmentFromImageID(imageID string) (
	response VulnContainerAssessmentResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesAssessmentFromImageID, imageID)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

// ListAssessments leverages ListAssessmentsDateRange and returns a list of assessments from the last 7 days
func (svc *ContainerVulnerabilityService) AssessmentFromImageDigest(imageDigest string) (
	response VulnContainerAssessmentResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesAssessmentFromImageDigest, imageDigest)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

// ListAssessments leverages ListAssessmentsDateRange and returns a list of assessments from the last 7 days
func (svc *ContainerVulnerabilityService) ListAssessments() (VulnContainerAssessmentsResponse, error) {
	var (
		now = time.Now().UTC()

		// 7 days from now plus 2 minutes, why?
		// because our API has a limit of exactly 7 days
		from = now.AddDate(0, 0, -7).Add(time.Minute * time.Duration(2))
	)

	return svc.ListAssessmentsDateRange(from, now)
}

// ListAssessmentsDateRange returns a list of container assessments during the specified date range
func (svc *ContainerVulnerabilityService) ListAssessmentsDateRange(start, end time.Time) (
	response VulnContainerAssessmentsResponse,
	err error,
) {
	if start.After(end) {
		err = errors.New("data range should have a start time before the end time")
		return
	}

	apiPath := fmt.Sprintf(
		"%s?START_TIME=%s&END_TIME=%s",
		apiVulnContainerAssessmentsForDateRange,
		start.UTC().Format(time.RFC3339),
		end.UTC().Format(time.RFC3339),
	)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

type vulnContainerScanRequest struct {
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

type vulnContainerScanResponse struct {
	Data struct {
		Status    string `json:"status"`
		RequestID string `json:"requestId"`
	} `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type vulnContainerScanStatusResponse struct {
	Data    VulnContainerAssessment `json:"data"`
	Ok      bool                    `json:"ok"`
	Message string                  `json:"message"`
}

func (scan *vulnContainerScanStatusResponse) CheckStatus() string {
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

type VulnContainerAssessmentResponse struct {
	Data    VulnContainerAssessment `json:"data"`
	Ok      bool                    `json:"ok"`
	Message string                  `json:"message"`
}

func (res *VulnContainerAssessmentResponse) CheckStatus() string {
	if !res.Ok {
		return fmt.Sprintf("there is a problem with the vulnerability assessment: %s", res.Message)
	}

	if res.Data.ScanStatus != "" {
		return res.Data.ScanStatus
	}

	if res.Data.Status != "" {
		return res.Data.Status
	}

	return "Unknown"
}

type VulnContainerAssessment struct {
	TotalVulnerabilities    int32               `json:"total_vulnerabilities"`
	CriticalVulnerabilities int32               `json:"critical_vulnerabilities"`
	HighVulnerabilities     int32               `json:"high_vulnerabilities"`
	MediumVulnerabilities   int32               `json:"medium_vulnerabilities"`
	LowVulnerabilities      int32               `json:"low_vulnerabilities"`
	InfoVulnerabilities     int32               `json:"info_vulnerabilities"`
	FixableVulnerabilities  int32               `json:"fixable_vulnerabilities"`
	LastEvaluationTime      string              `json:"last_evaluation_time,omitempty"`
	Image                   *VulnContainerImage `json:"image,omitempty"`

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

func (report *VulnContainerAssessment) VulnFixableCount(severity string) int32 {
	if report.Image == nil || report.Image.ImageLayers == nil {
		return 0
	}

	severity = strings.ToLower(severity)
	if !array.ContainsStr(ValidVulnSeverities, severity) {
		return 0
	}

	if len(report.Image.ImageLayers) == 0 {
		return 0
	}

	var fixable int32 = 0
	for _, layer := range report.Image.ImageLayers {
		for _, pkg := range layer.Packages {
			for _, vuln := range pkg.Vulnerabilities {
				if vuln.Severity == severity && vuln.FixVersion != "" {
					fixable++
				}
			}
		}
	}
	return fixable
}

// HighestSeverity returns the highest severity level vulnerability in a VulnContainerAssessment
func (report *VulnContainerAssessment) HighestSeverity() string {
	if report.CriticalVulnerabilities != 0 {
		return "critical"
	}
	if report.HighVulnerabilities != 0 {
		return "high"
	}
	if report.MediumVulnerabilities != 0 {
		return "medium"
	}
	if report.LowVulnerabilities != 0 {
		return "low"
	}
	return "unknown"
}

// HighestFixableSeverity returns the highest fixable severity level vulnerability in a VulnContainerAssessment
func (report *VulnContainerAssessment) HighestFixableSeverity() string {
	if report.VulnFixableCount("critical") != 0 {
		return "critical"
	}
	if report.VulnFixableCount("high") != 0 {
		return "high"
	}
	if report.VulnFixableCount("medium") != 0 {
		return "medium"
	}
	if report.VulnFixableCount("low") != 0 {
		return "low"
	}
	return "unknown"
}

// TotalFixableVulnerabilities returns the total number of vulnerabilities that have a fix available
func (report *VulnContainerAssessment) TotalFixableVulnerabilities() int32 {
	return report.FixableVulnerabilities
}

type VulnContainerImage struct {
	ImageInfo   *vulnContainerImageInfo   `json:"image_info,omitempty"`
	ImageLayers []VulnContainerImageLayer `json:"image_layers,omitempty"`
}

type vulnContainerImageInfo struct {
	ImageDigest string   `json:"image_digest"`
	ImageID     string   `json:"image_id"`
	Registry    string   `json:"registry"`
	Repository  string   `json:"repository"`
	CreatedTime string   `json:"created_time"`
	Size        int64    `json:"size"`
	Tags        []string `json:"tags"`
}

type VulnContainerImageLayer struct {
	Hash      string                 `json:"hash"`
	CreatedBy string                 `json:"created_by"`
	Packages  []VulnContainerPackage `json:"packages"`
}

type VulnContainerPackage struct {
	Name            string                   `json:"name"`
	Namespace       string                   `json:"namespace"`
	Version         string                   `json:"version"`
	Vulnerabilities []ContainerVulnerability `json:"vulnerabilities"`

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

type ContainerVulnerability struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	Link        string                 `json:"link"`
	FixVersion  string                 `json:"fix_version"`
	Metadata    map[string]interface{} `json:"metadata"`
	Status      string                 `json:"status"`
}

// traverseMetadata will try to extract an interface from the nested tree of key
// values contain inside the Metadata field of a container vulnerability struct
//
// Example: extract the version 3 of the CVSS Score
//
//          score := v.traverseMetadata("NVD", "CVSSv3", "Score")
//
func (v *ContainerVulnerability) traverseMetadata(fields ...string) interface{} {

	var (
		metaInterface interface{}
		metaMap       = v.Metadata
	)
	for i, field := range fields {

		if i != 0 {
			if newMap, ok := metaInterface.(map[string]interface{}); ok {
				metaMap = newMap
			} else {
				return nil
			}
		}

		if found, ok := metaMap[field]; ok {
			metaInterface = found
		} else {
			return nil
		}

	}

	return metaInterface
}

func (v *ContainerVulnerability) CVSSv3Score() float64 {
	score := v.traverseMetadata("NVD", "CVSSv3", "Score")

	if f, ok := score.(float64); ok {
		return f
	}

	return 0
}

func (v *ContainerVulnerability) CVSSv2Score() float64 {
	score := v.traverseMetadata("NVD", "CVSSv2", "Score")

	if f, ok := score.(float64); ok {
		return f
	}

	return 0
}

type VulnContainerAssessmentsResponse struct {
	Assessments []VulnContainerAssessmentSummary `json:"data"`
	Ok          bool                             `json:"ok"`
	Message     string                           `json:"message"`
}

type VulnContainerAssessmentSummary struct {
	EvalGuid                    string          `json:"eval_guid"`
	EvalStatus                  string          `json:"eval_status"`
	EvalType                    string          `json:"eval_type"`
	ImageCreatedTime            lwtime.NanoTime `json:"image_created_time"`
	ImageDigest                 string          `json:"image_digest"`
	ImageID                     string          `json:"image_id"`
	ImageNamespace              string          `json:"image_namespace"`
	ImageRegistry               string          `json:"image_registry"`
	ImageRepo                   string          `json:"image_repo"`
	ImageScanErrorMsg           string          `json:"image_scan_error_msg"`
	ImageScanStatus             string          `json:"image_scan_status"`
	ImageScanTime               lwtime.NanoTime `json:"image_scan_time"`
	ImageSize                   string          `json:"image_size"`
	ImageTags                   []string        `json:"image_tags"`
	NdvContainers               string          `json:"ndv_containers"`
	NumFixes                    string          `json:"num_fixes"`
	NumVulnerabilitiesSeverity1 string          `json:"num_vulnerabilities_severity_1"`
	NumVulnerabilitiesSeverity2 string          `json:"num_vulnerabilities_severity_2"`
	NumVulnerabilitiesSeverity3 string          `json:"num_vulnerabilities_severity_3"`
	NumVulnerabilitiesSeverity4 string          `json:"num_vulnerabilities_severity_4"`
	NumVulnerabilitiesSeverity5 string          `json:"num_vulnerabilities_severity_5"`
	StartTime                   lwtime.NanoTime `json:"start_time"`
}
