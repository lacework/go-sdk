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

import "github.com/pkg/errors"

// GcpRecommendationsV1 is a service that interacts with the V1 Recommendations
// endpoints from the Lacework Server
type GcpRecommendationsV1 struct {
	client *Client
}

func (svc *GcpRecommendationsV1) List() ([]RecommendationV1, error) {
	return svc.client.Recommendations.list(GcpRecommendation)
}

func (svc *GcpRecommendationsV1) Patch(recommendations RecommendationStateV1) (response RecommendationResponseV1, err error) {
	return svc.client.Recommendations.patch(GcpRecommendation, recommendations)
}

// GetReport This is an experimental feature. Returned RecommendationID's are not guaranteed to be correct. Scoped to Lacework Account/Subaccount
func (svc *GcpRecommendationsV1) GetReport(reportType string) (response []RecommendationV1, err error) {
	gcpCfg, err := svc.client.Integrations.ListGcpCfg()
	if err != nil {
		return []RecommendationV1{}, err
	}

	// cis1.3 is off by default, so we won't be able to fetch previous reports.
	
	if len(gcpCfg.Data) == 0 {
		return []RecommendationV1{}, errors.Wrap(err, "unable to find a GCP cloud account integration")
	}

	var projectID = gcpCfg.Data[0].Data.ID
	var orgID = "n/a"

	// TODO get all gcp projects for org
	// cli/cmd/compliance_gcp.go:496
	//if gcpCfg.Data[0].Data.IDType == "ORGANIZATION" {
	//
	//}

	cfg := ComplianceGcpReportConfig{
		OrganizationID: orgID,
		ProjectID:      projectID,
		Type:           reportType,
	}

	var res complianceGcpReportResponse
	res, err = svc.client.Compliance.GetGcpReport(cfg)
	if err != nil {
		return []RecommendationV1{}, err
	}

	var recommendationIDs []string

	for _, rec := range res.Data[0].Recommendations {
		recommendationIDs = append(recommendationIDs, rec.RecID)
	}

	schema := ReportSchema{res.Data[0].ReportType, recommendationIDs}

	// fetch all gcp recommendations
	allRecommendations, err := svc.client.Recommendations.Gcp.List()
	if err != nil {
		return []RecommendationV1{}, err
	}
	filteredRecommendations := filterRecommendations(allRecommendations, schema)
	return filteredRecommendations, nil
}
