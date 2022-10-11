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
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// v2GcpReportsService is a service that interacts with the APIv2
// vulnerabilities endpoints for hosts
type gcpReportsService struct {
	client *Client
}

type GcpReportConfig struct {
	OrganizationID string
	ProjectID      string
	Type           string
}

// Get returns a raw response of the Alert Profile with the matching guid.
func (svc *gcpReportsService) Get(reportCfg GcpReportConfig) (response GcpReportResponse, err error) {
	if reportCfg.ProjectID == "" || reportCfg.OrganizationID == "" {
		return GcpReportResponse{}, errors.New("project id and org id are required")
	}

	apiPath := fmt.Sprintf(apiV2Reports, reportCfg.ProjectID, reportCfg.Type, reportCfg.ProjectID)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
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

func (gcp GcpReport) GetComplianceRecommendation(recommendationID string) RecommendationV2 {
	for _, r := range gcp.Recommendations {
		if r.RecID == recommendationID {
			return r
		}
	}
	return RecommendationV2{}
}

func (svc *gcpReportsService) DownloadPDF(filepath string, config GcpReportConfig) error {
	if config.ProjectID == "" || config.OrganizationID == "" {
		return errors.New("project id and org id are required")
	}

	apiPath := fmt.Sprintf(apiV2Reports, config.ProjectID, "PDF", config.OrganizationID)

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
