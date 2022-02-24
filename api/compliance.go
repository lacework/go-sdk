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

// ComplianceService is a service that interacts with the compliance
// endpoints from the Lacework Server
type ComplianceService struct {
	client *Client
}

type CloudComplianceReport interface {
	GetComplianceRecommendation(recommendationID string) ComplianceRecommendation
}

func (svc *ComplianceService) ListGcpProjects(orgID string) (
	response compGcpProjectsResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiComplianceGcpListProjects, orgID)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *ComplianceService) RunIntegrationReport(intgGuid string) (
	response map[string]interface{},
	err error,
) {
	apiPath := fmt.Sprintf(apiRunReportIntegration, intgGuid)
	err = svc.client.RequestDecoder("POST", apiPath, nil, &response)
	return
}

type compGcpProjectsResponse struct {
	Data    []CompGcpProjects `json:"data"`
	Ok      bool              `json:"ok"`
	Message string            `json:"message"`
}

type CompGcpProjects struct {
	Organization string   `json:"organization"`
	Projects     []string `json:"projects"`
}

type ComplianceSummary struct {
	AssessedResourceCount     int `json:"assessed_resource_count"`
	NumCompliant              int `json:"num_compliant"`
	NumNotCompliant           int `json:"num_not_compliant"`
	NumRecommendations        int `json:"num_recommendations"`
	NumSeverity1NonCompliance int `json:"num_severity_1_non_compliance"`
	NumSeverity2NonCompliance int `json:"num_severity_2_non_compliance"`
	NumSeverity3NonCompliance int `json:"num_severity_3_non_compliance"`
	NumSeverity4NonCompliance int `json:"num_severity_4_non_compliance"`
	NumSeverity5NonCompliance int `json:"num_severity_5_non_compliance"`
	NumSuppressed             int `json:"num_suppressed"`
	SuppressedResourceCount   int `json:"suppressed_resource_count"`
	ViolatedResourceCount     int `json:"violated_resource_count"`
}

type ComplianceRecommendation struct {
	RecID                 string                `json:"rec_id"`
	AssessedResourceCount int                   `json:"assessed_resource_count"`
	ResourceCount         int                   `json:"resource_count"`
	Category              string                `json:"category"`
	InfoLink              string                `json:"info_link"`
	Service               string                `json:"service"`
	Severity              int                   `json:"severity"`
	Status                string                `json:"status"`
	Suppressions          []string              `json:"suppressions"`
	Title                 string                `json:"title"`
	Violations            []ComplianceViolation `json:"violations"`
}

func (r *ComplianceRecommendation) SeverityString() string {
	switch r.Severity {
	case 1:
		return "Critical"
	case 2:
		return "High"
	case 3:
		return "Medium"
	case 4:
		return "Low"
	case 5:
		return "Info"
	default:
		return "Unknown"
	}
}

type ComplianceViolation struct {
	Region   string   `json:"region"`
	Resource string   `json:"resource"`
	Reasons  []string `json:"reasons"`
}

// ValidComplianceStatus is a list of all valid compliance status
var ValidComplianceStatus = []string{"non-compliant", "requires-manual-assessment", "suppressed", "compliant", "could-not-assess"}
