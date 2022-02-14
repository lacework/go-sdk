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
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
)

type ComplianceGcpReportConfig struct {
	OrganizationID string
	ProjectID      string
	Type           string
}

func (svc *ComplianceService) GetGcpReport(config ComplianceGcpReportConfig) (
	response complianceGcpReportResponse,
	err error,
) {
	if config.OrganizationID == "" || config.ProjectID == "" {
		err = errors.New("organization_id and project_id is required")
		return
	}
	apiPath := fmt.Sprintf(apiComplianceGcpLatestReport, config.OrganizationID, config.ProjectID)

	if config.Type != "" {
		apiPath = fmt.Sprintf("%s&REPORT_TYPE=%s", apiPath, config.Type)
	}

	// add JSON format, if not, the default is PDF
	apiPath = fmt.Sprintf("%s&FILE_FORMAT=json", apiPath)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *ComplianceService) DownloadGcpReportPDF(filepath string, config ComplianceGcpReportConfig) error {
	if config.OrganizationID == "" || config.ProjectID == "" {
		return errors.New("organization_id and project_id is required")
	}

	apiPath := fmt.Sprintf(apiComplianceGcpLatestReport, config.OrganizationID, config.ProjectID)

	if config.Type != "" {
		apiPath = fmt.Sprintf("%s&REPORT_TYPE=%s", apiPath, config.Type)
	}

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

func (svc *ComplianceService) RunGcpReport(projectID string) (
	response complianceRunGcpReportResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiRunReportGcp, projectID)
	err = svc.client.RequestDecoder("POST", apiPath, nil, &response)
	return
}

type complianceRunGcpReportResponse struct {
	IntgGuid             string                     `json:"intgGuid"`
	MultiContextEvalGuid string                     `json:"multiContextEvalGuid"`
	EvalRequests         []complianceGcpEvalRequest `json:"evalRequests"`
}

type complianceGcpEvalRequest struct {
	EvalGuid       string                       `json:"evalGuid"`
	EvalCtx        complianceRunReportGcpContex `json:"evalCtx"`
	SubmitErrorMsg interface{}                  `json:"submitErrorMsg"`
}

type complianceRunReportGcpContex struct {
	OrganizationID   string `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
	ProjectID        string `json:"projectId"`
	ProjectName      string `json:"projectName"`
}

type complianceGcpReportResponse struct {
	Data    []ComplianceGcpReport `json:"data"`
	Ok      bool                  `json:"ok"`
	Message string                `json:"message"`
}

type ComplianceGcpReport struct {
	ReportTitle      string                     `json:"reportTitle"`
	ReportType       string                     `json:"reportType"`
	ReportTime       time.Time                  `json:"reportTime"`
	OrganizationID   string                     `json:"organizationId"`
	OrganizationName string                     `json:"organizationName"`
	ProjectID        string                     `json:"projectId"`
	ProjectName      string                     `json:"projectName"`
	Summary          []ComplianceSummary        `json:"summary"`
	Recommendations  []ComplianceRecommendation `json:"recommendations"`
}

func (gcp ComplianceGcpReport) GetComplianceRecommendations() []ComplianceRecommendation {
	return gcp.Recommendations
}
