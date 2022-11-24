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

// v2AwsReportsService is a service that interacts with the APIv2
type awsReportsService struct {
	client *Client
}

type AwsReportConfig struct {
	AccountID string
	Type      AwsReportType
}

type AwsReportType int

func (report AwsReportType) String() string {
	return awsReportTypes[report]
}

func NewAwsReportType(report string) (AwsReportType, error) {
	for k, v := range awsReportTypes {
		if v == report {
			return k, nil
		}
	}
	return NONE_AWS_REPORT, errors.Errorf("no report type found for %s", report)
}

var awsReportTypes = map[AwsReportType]string{AWS_CIS_S3: "AWS_CIS_S3", NIST_800_53_Rev4: "NIST_800-53_Rev4",
	NIST_800_171_Rev2: "NIST_800-171_Rev2", ISO_2700: "ISO_2700", HIPAA: "HIPAA", SOC: "SOC",
	AWS_SOC_Rev2: "AWS_SOC_Rev2", PCI: "PCI", AWS_CIS_14: "AWS_CIS_14", AWS_CMMC_1_02: "AWS_CMMC_1.02",
	AWS_ISO_27001_2013: "AWS_ISO_27001:2013", AWS_NIST_CSF: "AWS_NIST_CSF", AWS_HIPAA: "AWS_HIPAA",
	AWS_NIST_800_53_rev5: "AWS_NIST_800-53_rev5", AWS_NIST_800_171_rev2: "AWS_NIST_800-171_rev2",
	AWS_PCI_DSS_3_2_1: "AWS_PCI_DSS_3.2.1", AWS_SOC_2: "AWS_SOC_2", LW_AWS_SEC_ADD_1_0: "LW_AWS_SEC_ADD_1_0"}

const (
	NONE_AWS_REPORT AwsReportType = iota
	AWS_CIS_S3
	NIST_800_53_Rev4
	NIST_800_171_Rev2
	ISO_2700
	HIPAA
	SOC
	AWS_SOC_Rev2
	PCI
	AWS_CIS_14
	AWS_CMMC_1_02
	AWS_HIPAA
	AWS_ISO_27001_2013
	AWS_NIST_CSF
	AWS_NIST_800_171_rev2
	AWS_NIST_800_53_rev5
	AWS_PCI_DSS_3_2_1
	AWS_SOC_2
	LW_AWS_SEC_ADD_1_0
)

// Get returns an AwsReportResponse
func (svc *awsReportsService) Get(reportCfg AwsReportConfig) (response AwsReportResponse, err error) {
	if reportCfg.AccountID == "" {
		return AwsReportResponse{}, errors.New("specify an account id")
	}
	apiPath := fmt.Sprintf(apiV2Reports, reportCfg.AccountID, "json", reportCfg.Type.String())
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *awsReportsService) DownloadPDF(filepath string, config AwsReportConfig) error {
	if config.AccountID == "" {
		return errors.New("account id is required")
	}

	apiPath := fmt.Sprintf(apiV2Reports, config.AccountID, "pdf", "")

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

func (aws AwsReport) GetComplianceRecommendation(recommendationID string) RecommendationV2 {
	for _, r := range aws.Recommendations {
		if r.RecID == recommendationID {
			return r
		}
	}
	return RecommendationV2{}
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
