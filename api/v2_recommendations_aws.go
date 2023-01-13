//
// Author:: Ross Moles (<ross.moles@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	"encoding/json"
	"errors"

	"github.com/lacework/go-sdk/internal/databox"
)

// AwsRecommendationsV2 is a service that interacts with the V2 Recommendations
// endpoints from the Lacework Server
type AwsRecommendationsV2 struct {
	client *Client
}

func (svc *AwsRecommendationsV2) List() ([]RecV2, error) {
	return svc.client.V2.Recommendations.list(AwsRecommendation)
}

func (svc *AwsRecommendationsV2) Patch(recommendations RecommendationStateV2) (RecommendationResponseV2, error) {
	return svc.client.V2.Recommendations.patch(AwsRecommendation, recommendations)
}

// GetReport This is an experimental feature. Returned RecommendationID's are not guaranteed to be correct. Scoped to Lacework Account/Subaccount
func (svc *AwsRecommendationsV2) GetReport(reportType string) ([]RecV2, error) {
	report := struct {
		Ids map[string]string `json:"recommendation_ids"`
	}{}

	schemaBytes, ok := databox.Get("/reports/aws/cis.json")
	if !ok {
		return []RecV2{}, errors.New(
			"compliance report schema not found",
		)
	}

	err := json.Unmarshal(schemaBytes, &report)
	if err != nil {
		return []RecV2{}, err
	}

	schema := ReportSchema{reportType, report.Ids}

	// fetch all aws recommendations
	allRecommendations, err := svc.client.V2.Recommendations.Aws.List()
	if err != nil {
		return []RecV2{}, err
	}
	filteredRecommendations := filterRecommendations(allRecommendations, schema)
	return filteredRecommendations, nil
}
