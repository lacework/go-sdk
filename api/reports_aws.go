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

// v2AwsReportsService is a service that interacts with the APIv2
// vulnerabilities endpoints for hosts
type awsReportsService struct {
	client *Client
}

type AwsReportConfig struct {
	AccountID string
	Type      string
}

// Get returns a raw response of the Alert Profile with the matching guid.
func (svc *awsReportsService) Get(reportCfg AwsReportConfig) (response AwsReportResponse, err error) {
	if reportCfg.AccountID == "" {
		return AwsReportResponse{}, errors.New("specify an account id")
	}
	apiPath := fmt.Sprintf(apiV2Reports, reportCfg.AccountID, "json", reportCfg.Type)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *awsReportsService) DownloadPDF(filepath string, config AwsReportConfig) error {
	if config.AccountID == "" {
		return errors.New("account id is required")
	}

	apiPath := fmt.Sprintf(apiV2Reports, config.AccountID, "PDF", "")

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

type AwsReportResponse struct {
	Data    []AwsReport `json:"data"`
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
}

type AwsReport struct {
	ReportType      string             `json:"reportType"`
	ReportTitle     string             `json:"reportTitle"`
	Recommendations []RecommendationV2 `json:"recommendations"`
	Summary         []ReportSummary    `json:"summary"`
	AccountID       string             `json:"accountId"`
	AccountAlias    string             `json:"accountAlias"`
	ReportTime      time.Time          `json:"reportTime"`
}

func (aws AwsReport) GetComplianceRecommendation(recommendationID string) RecommendationV2 {
	for _, r := range aws.Recommendations {
		if r.RecID == recommendationID {
			return r
		}
	}
	return RecommendationV2{}
}
