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

// ReportsService is a service that interacts with the Reports
// endpoints from the Lacework APIv2 Server
type ReportsService struct {
	client *Client
	Aws    *awsReportsService
	Azure  *azureReportsService
	Gcp    *gcpReportsService
}

func NewReportsService(c *Client) *ReportsService {
	return &ReportsService{c,
		&awsReportsService{c},
		&azureReportsService{c},
		&gcpReportsService{c},
	}
}

// The method by which a report can be retrieved from v2/Reports/ api
// can be 'reportName' or 'reportType'
type reportFilter int

const (
	ReportFilterType reportFilter = iota
	ReportFilterName
)

// reportFilterTypes is the list of available report filter types
var reportFilterTypes = map[reportFilter]string{
	ReportFilterType: "reportType",
	ReportFilterName: "reportName",
}

func (r reportFilter) String() string {
	return reportFilterTypes[r]
}

type ReportSummary struct {
	NumRecommendations        int `json:"NUM_RECOMMENDATIONS"`
	NumSeverity2NonCompliance int `json:"NUM_SEVERITY_2_NON_COMPLIANCE"`
	NumSeverity4NonCompliance int `json:"NUM_SEVERITY_4_NON_COMPLIANCE"`
	NumSeverity1NonCompliance int `json:"NUM_SEVERITY_1_NON_COMPLIANCE"`
	NumCompliant              int `json:"NUM_COMPLIANT"`
	NumSeverity3NonCompliance int `json:"NUM_SEVERITY_3_NON_COMPLIANCE"`
	AssessedResourceCount     int `json:"ASSESSED_RESOURCE_COUNT"`
	NumSuppressed             int `json:"NUM_SUPPRESSED"`
	NumSeverity5NonCompliance int `json:"NUM_SEVERITY_5_NON_COMPLIANCE"`
	NumNotComplinace          int `json:"NUM_NOT_COMPLIANT"`
	ViolatedResourceCount     int `json:"VIOLATED_RESOURCE_COUNT"`
	SuppressedResourceCount   int `json:"SUPPRESSED_RESOURCE_COUNT"`
}

type RecommendationV2 struct {
	AccountID             string                  `json:"ACCOUNT_ID"`
	AccountAlias          string                  `json:"ACCOUNT_ALIAS"`
	Service               string                  `json:"SERVICE"`
	StartTime             int64                   `json:"START_TIME"`
	Suppressions          []string                `json:"SUPPRESSIONS"`
	InfoLink              string                  `json:"INFO_LINK"`
	AssessedResourceCount int                     `json:"ASSESSED_RESOURCE_COUNT"`
	Status                string                  `json:"STATUS"`
	RecID                 string                  `json:"REC_ID"`
	Category              string                  `json:"CATEGORY"`
	Title                 string                  `json:"TITLE"`
	Violations            []ComplianceViolationV2 `json:"VIOLATIONS"`
	ResourceCount         int                     `json:"RESOURCE_COUNT"`
	Severity              int                     `json:"SEVERITY"`
}

type ComplianceViolationV2 struct {
	Region   string   `json:"region"`
	Resource string   `json:"resource"`
	Reasons  []string `json:"reasons"`
}

func (r *RecommendationV2) SeverityString() string {
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

type CloudComplianceReportV2 interface {
	GetComplianceRecommendation(recommendationID string) (*RecommendationV2, bool)
}

// ValidComplianceStatus is a list of all valid compliance status
var ValidComplianceStatus = []string{"non-compliant", "requires-manual-assessment", "suppressed", "compliant", "could-not-assess"}
