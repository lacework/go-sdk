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

type ComplianceAzureReportConfig struct {
	TenantID       string
	SubscriptionID string
	Type           string
}

func (svc *ComplianceService) ListAzureTenants() ([]string, error) {
	var response AzureIntegrationsResponse
	var tenants []string

	err := svc.client.RequestDecoder("GET", fmt.Sprintf(apiIntegrationsByType, "AZURE_CFG"), nil, &response)
	if err != nil {
		return nil, err
	}
	for _, azure := range response.Data {
		tenants = append(tenants, azure.Data.TenantID)
	}
	return tenants, err
}

func (svc *ComplianceService) ListAzureSubscriptions(tenantID string) (
	response compAzureSubsResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiComplianceAzureListSubscriptions, tenantID)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *ComplianceService) GetAzureReport(config ComplianceAzureReportConfig) (
	response complianceAzureReportResponse,
	err error,
) {
	if config.TenantID == "" || config.SubscriptionID == "" {
		err = errors.New("tenant_id and subscription_id are required")
		return
	}
	apiPath := fmt.Sprintf(apiComplianceAzureLatestReport, config.TenantID, config.SubscriptionID)

	if config.Type != "" {
		apiPath = fmt.Sprintf("%s&REPORT_TYPE=%s", apiPath, config.Type)
	}

	// add JSON format, if not, the default is PDF
	apiPath = fmt.Sprintf("%s&FILE_FORMAT=json", apiPath)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *ComplianceService) DownloadAzureReportPDF(filepath string, config ComplianceAzureReportConfig) error {
	if config.TenantID == "" || config.SubscriptionID == "" {
		return errors.New("tenant_id and subscription_id are required")
	}

	apiPath := fmt.Sprintf(apiComplianceAzureLatestReport, config.TenantID, config.SubscriptionID)

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

func (svc *ComplianceService) RunAzureReport(tenantID string) (
	response complianceRunAzureReportResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiRunReportAzure, tenantID)
	err = svc.client.RequestDecoder("POST", apiPath, nil, &response)
	return
}

type complianceRunAzureReportResponse struct {
	IntgGuid             string                       `json:"intgGuid"`
	MultiContextEvalGuid string                       `json:"multiContextEvalGuid"`
	EvalRequests         []complianceAzureEvalRequest `json:"evalRequests"`
}

type complianceAzureEvalRequest struct {
	EvalGuid       string                         `json:"evalGuid"`
	EvalCtx        complianceRunReportAzureContex `json:"evalCtx"`
	SubmitErrorMsg interface{}                    `json:"submitErrorMsg"`
}

type complianceRunReportAzureContex struct {
	SubscriptionID   string `json:"subscriptionId"`
	SubscriptionName string `json:"subscriptionName"`
	TenantID         string `json:"tenantId"`
	TenantName       string `json:"tenantName"`
}

type compAzureSubsResponse struct {
	Data    []CompAzureSubscriptions `json:"data"`
	Ok      bool                     `json:"ok"`
	Message string                   `json:"message"`
}

type CompAzureSubscriptions struct {
	Tenant        string   `json:"tenant"`
	Subscriptions []string `json:"subscriptions"`
}

type complianceAzureReportResponse struct {
	Data    []ComplianceAzureReport `json:"data"`
	Ok      bool                    `json:"ok"`
	Message string                  `json:"message"`
}

type ComplianceAzureReport struct {
	ReportTitle      string                     `json:"reportTitle"`
	ReportType       string                     `json:"reportType"`
	ReportTime       time.Time                  `json:"reportTime"`
	TenantID         string                     `json:"tenantId"`
	TenantName       string                     `json:"tenantName"`
	SubscriptionID   string                     `json:"subscriptionId"`
	SubscriptionName string                     `json:"subscriptionName"`
	Summary          []ComplianceSummary        `json:"summary"`
	Recommendations  []ComplianceRecommendation `json:"recommendations"`
}
