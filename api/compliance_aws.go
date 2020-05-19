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

type ComplianceAwsReportConfig struct {
	AccountID string
	Type      string
}

func (svc *ComplianceService) GetAwsReport(config ComplianceAwsReportConfig) (
	response complianceAwsReportResponse,
	err error,
) {
	if config.AccountID == "" {
		err = errors.New("account_id is required")
		return
	}
	apiPath := fmt.Sprintf(apiComplianceAwsLatestReport, config.AccountID)

	if config.Type != "" {
		apiPath = fmt.Sprintf("%s&REPORT_TYPE=%s", apiPath, config.Type)
	}

	// add JSON format, if not, the default is PDF
	apiPath = fmt.Sprintf("%s&FILE_FORMAT=json", apiPath)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *ComplianceService) DownloadAwsReportPDF(filepath string, config ComplianceAwsReportConfig) error {
	if config.AccountID == "" {
		return errors.New("account_id is required")
	}

	apiPath := fmt.Sprintf(apiComplianceAwsLatestReport, config.AccountID)

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

func (svc *ComplianceService) RunAwsReport(accountID string) (
	response map[string]interface{}, // @afiune not consistent with the other cloud providers
	err error,
) {
	apiPath := fmt.Sprintf(apiRunReportAws, accountID)
	err = svc.client.RequestDecoder("POST", apiPath, nil, &response)
	return
}

type complianceAwsReportResponse struct {
	Data    []ComplianceAwsReport `json:"data"`
	Ok      bool                  `json:"ok"`
	Message string                `json:"message"`
}

type ComplianceAwsReport struct {
	ReportTitle     string                     `json:"reportTitle"`
	ReportType      string                     `json:"reportType"`
	ReportTime      time.Time                  `json:"reportTime"`
	AccountID       string                     `json:"accountId"`
	AccountAlias    string                     `json:"accountAlias"`
	Summary         []ComplianceSummary        `json:"summary"`
	Recommendations []ComplianceRecommendation `json:"recommendations"`
}
