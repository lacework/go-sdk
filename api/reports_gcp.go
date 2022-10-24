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

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
)

// v2GcpReportsService is a service that interacts with the APIv2
type gcpReportsService struct {
	client *Client
}

type GcpReportConfig struct {
	OrganizationID string
	ProjectID      string
	Type           GcpReportType
}

type GcpReportType int

func (report GcpReportType) String() string {
	return gcpReportTypes[report]
}

func NewGcpReportType(report string) (GcpReportType, error) {
	for k, v := range gcpReportTypes {
		if v == report {
			return k, nil
		}
	}
	return NONE_GCP_REPORT, errors.Errorf("no report type found for %s", report)
}

var gcpReportTypes = map[GcpReportType]string{GCP_HIPAA: "GCP_HIPAA", GCP_CIS: "GCP_CIS", GCP_SOC: "GCP_SOC", GCP_CIS12: "GCP_CIS12",
	GCP_K8S: "GCP_K8S", GCP_PCI_Rev2: "GCP_PCI_Rev2", GCP_SOC_Rev2: "GCP_SOC_Rev2", GCP_HIPAA_Rev2: "GCP_HIPAA_Rev2", GCP_ISO_27001: "GCP_ISO_27001",
	GCP_NIST_CSF: "GCP_NIST_CSF", GCP_NIST_800_53_REV4: "GCP_NIST_800_53_REV4", GCP_NIST_800_171_REV2: "GCP_NIST_800_171_REV2", GCP_PCI: "GCP_PCI", GCP_CIS13: "GCP_CIS13"}

const (
	NONE_GCP_REPORT GcpReportType = iota
	GCP_HIPAA
	GCP_CIS
	GCP_SOC
	GCP_CIS12
	GCP_K8S
	GCP_PCI_Rev2
	GCP_SOC_Rev2
	GCP_HIPAA_Rev2
	GCP_ISO_27001
	GCP_NIST_CSF
	GCP_NIST_800_53_REV4
	GCP_NIST_800_171_REV2
	GCP_PCI
	GCP_CIS13
)

// Get returns a raw response of the Alert Profile with the matching guid.
func (svc *gcpReportsService) Get(reportCfg GcpReportConfig) (response GcpReportResponse, err error) {
	if reportCfg.ProjectID == "" || reportCfg.OrganizationID == "" {
		return GcpReportResponse{}, errors.New("project id and org id are required")
	}

	apiPath := fmt.Sprintf(apiV2ReportsSecondaryQuery, reportCfg.OrganizationID, reportCfg.ProjectID, "json", reportCfg.Type.String())
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *gcpReportsService) DownloadPDF(filepath string, config GcpReportConfig) error {
	if config.ProjectID == "" || config.OrganizationID == "" {
		return errors.New("project id and org id are required")
	}

	apiPath := fmt.Sprintf(apiV2Reports, config.ProjectID, "pdf", config.OrganizationID)

	request, err := svc.client.NewRequest("GET", apiPath, nil)
	if err != nil {
		return err
	}

	response, err := svc.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = checkErrorInResponse(response)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, response.Body)
	return err
}

func (gcp GcpReport) GetComplianceRecommendation(recommendationID string) RecommendationV2 {
	for _, r := range gcp.Recommendations {
		if r.RecID == recommendationID {
			return r
		}
	}
	return RecommendationV2{}
}

type GcpReportResponse struct {
	Data    []GcpReport `json:"data"`
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
}

type GcpReport struct {
	ReportType       string             `json:"reportType"`
	ReportTitle      string             `json:"reportTitle"`
	Recommendations  []RecommendationV2 `json:"recommendations"`
	Summary          []ReportSummary    `json:"summary"`
	ReportTime       time.Time          `json:"reportTime"`
	OrganizationName string             `json:"organizationName"`
	OrganizationID   string             `json:"organizationId"`
	ProjectName      string             `json:"projectName"`
	ProjectID        string             `json:"projectId"`
}
