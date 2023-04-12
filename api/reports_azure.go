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
	"sort"
	"time"

	"github.com/pkg/errors"
)

// v2AzureReportsService is a service that interacts with the APIv2
type azureReportsService struct {
	client *Client
}

type AzureReportConfig struct {
	TenantID       string
	SubscriptionID string
	Value          string
	Parameter      reportFilter
}

type AzureReportType int

func (report AzureReportType) String() string {
	return azureReportTypes[report]
}

func NewAzureReportType(report string) (AzureReportType, error) {
	for k, v := range azureReportTypes {
		if v == report {
			return k, nil
		}
	}
	return NONE_AZURE_REPORT, errors.Errorf("no report type found for %s", report)
}

func AzureReportTypes() []string {
	reportTypes := make([]string, 0, len(azureReportTypes))

	for _, report := range azureReportTypes {
		reportTypes = append(reportTypes, report)
	}

	sort.Strings(reportTypes)
	return reportTypes
}

var azureReportTypes = map[AzureReportType]string{
	AZURE_CIS: "AZURE_CIS", AZURE_CIS_131: "AZURE_CIS_131", AZURE_SOC: "AZURE_SOC", AZURE_SOC_Rev2: "AZURE_SOC_Rev2",
	AZURE_PCI: "AZURE_PCI", AZURE_PCI_Rev2: "AZURE_PCI_Rev2", AZURE_ISO_27001: "AZURE_ISO_27001", AZURE_NIST_CSF: "AZURE_NIST_CSF",
	AZURE_NIST_800_53_REV5: "AZURE_NIST_800_53_REV5", AZURE_NIST_800_171_REV2: "AZURE_NIST_800_171_REV2", AZURE_HIPAA: "AZURE_HIPAA"}

const (
	NONE_AZURE_REPORT AzureReportType = iota
	AZURE_CIS
	AZURE_CIS_131
	AZURE_SOC
	AZURE_SOC_Rev2
	AZURE_PCI
	AZURE_PCI_Rev2
	AZURE_ISO_27001
	AZURE_NIST_CSF
	AZURE_NIST_800_53_REV5
	AZURE_NIST_800_171_REV2
	AZURE_HIPAA
)

// Get returns an AzureReportResponse
func (svc *azureReportsService) Get(reportCfg AzureReportConfig) (response AzureReportResponse, err error) {
	if reportCfg.SubscriptionID == "" {
		return AzureReportResponse{}, errors.New("specify an account id")
	}

	apiPath := fmt.Sprintf(apiV2ReportsSecondaryQuery, reportCfg.TenantID, reportCfg.SubscriptionID, "json", reportCfg.Parameter.String(), reportCfg.Value)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *azureReportsService) DownloadPDF(filepath string, config AzureReportConfig) error {
	if config.TenantID == "" || config.SubscriptionID == "" {
		return errors.New("tenant_id and subscription_id are required")
	}

	apiPath := fmt.Sprintf(apiV2ReportsSecondaryQuery, config.TenantID, config.SubscriptionID, "pdf", config.Parameter.String(), config.Value)

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

func (azure AzureReport) GetComplianceRecommendation(recommendationID string) (*RecommendationV2, bool) {
	for _, r := range azure.Recommendations {
		if r.RecID == recommendationID {
			return &r, true
		}
	}
	return nil, false
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
