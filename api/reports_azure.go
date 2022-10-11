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

// v2AzureReportsService is a service that interacts with the APIv2
// vulnerabilities endpoints for hosts
type azureReportsService struct {
	client *Client
}

type AzureReportConfig struct {
	TenantID       string
	SubscriptionID string
	Type           string
	Format         string
}

// Get returns a raw response of the Alert Profile with the matching guid.
func (svc *azureReportsService) Get(reportCfg AzureReportConfig) (response AzureReportResponse, err error) {
	format := reportCfg.Format
	if reportCfg.SubscriptionID == "" {
		return AzureReportResponse{}, errors.New("specify an account id")
	}

	if format == "" {
		format = "json"
	}

	apiPath := fmt.Sprintf(apiV2Reports, reportCfg.SubscriptionID, format, reportCfg.Type)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *azureReportsService) DownloadPDF(filepath string, config AzureReportConfig) error {
	if config.TenantID == "" || config.SubscriptionID == "" {
		return errors.New("tenant_id and subscription_id are required")
	}

	apiPath := fmt.Sprintf(apiV2Reports, config.TenantID, "PDF", config.SubscriptionID)

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

type AzureReportResponse struct {
	Data    []AzureReport `json:"data"`
	Ok      bool          `json:"ok"`
	Message string        `json:"message"`
}

type AzureReport struct {
	ReportType       string             `json:"reportType"`
	ReportTitle      string             `json:"reportTitle"`
	Recommendations  []RecommendationV2 `json:"recommendations"`
	Summary          []ReportSummary    `json:"summary"`
	ReportTime       time.Time          `json:"reportTime"`
	SubscriptionName string             `json:"subscriptionName"`
	SubscriptionID   string             `json:"SubscriptionID"`
	TenantName       string             `json:"tenantName"`
	TenantID         string             `json:"tenantId"`
}

func (azure AzureReport) GetComplianceRecommendation(recommendationID string) RecommendationV2 {
	for _, r := range azure.Recommendations {
		if r.RecID == recommendationID {
			return r
		}
	}
	return RecommendationV2{}
}
