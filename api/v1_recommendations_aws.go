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
	"github.com/pkg/errors"
)

// AwsRecommendationsV1 is a service that interacts with the V1 Recommendations
// endpoints from the Lacework Server
type AwsRecommendationsV1 struct {
	client *Client
}

func (svc *AwsRecommendationsV1) List() ([]RecommendationV1, error) {
	return svc.client.Recommendations.list(AwsRecommendation)
}

func (svc *AwsRecommendationsV1) Patch(recommendations RecommendationStateV1) (response RecommendationResponseV1, err error) {
	return svc.client.Recommendations.patch(AwsRecommendation, recommendations)
}

// GetReport This is an experimental feature. Returned RecommendationID's are not guaranteed to be correct. Scoped to Lacework Account/Subaccount
func (svc *AwsRecommendationsV1) GetReport(reportType string) (response []RecommendationV1, err error) {
	awsCfg, err := svc.client.Integrations.ListAwsCfg()
	if err != nil {
		return []RecommendationV1{}, err
	}
	if len(awsCfg.Data) == 0 {
		return []RecommendationV1{}, errors.Wrap(err, "unable to find an AWS cloud account integration")
	}
	
	accountID := awsCfg.Data[0].Data.GetAccountID()

	cfg := ComplianceAwsReportConfig{AccountID: accountID, Type: reportType}

	var res complianceAwsReportResponse
	res, err = svc.client.Compliance.GetAwsReport(cfg)
	if err != nil {
		return []RecommendationV1{}, err
	}

	var recommendationIDs []string

	for _, rec := range res.Data[0].Recommendations {
		recommendationIDs = append(recommendationIDs, rec.RecID)
	}

	schema := ReportSchema{res.Data[0].ReportType, recommendationIDs}

	// fetch all aws recommendations
	allRecommendations, err := svc.client.Recommendations.Aws.List()
	filteredRecommendations := filterRecommendations(allRecommendations, schema)
	return filteredRecommendations, nil
}
