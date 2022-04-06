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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lacework/go-sdk/internal/databox"
)

// AzureRecommendationsV1 is a service that interacts with the V1 Recommendations
// endpoints from the Lacework Server
type AzureRecommendationsV1 struct {
	client *Client
}

const azureCIS = "Azure_CIS_(?!131)\\w+"
const azureCIS131 = "Azure_CIS_131\\w+"

func (svc *AzureRecommendationsV1) List() ([]RecommendationV1, error) {
	return svc.client.Recommendations.list(AzureRecommendation)
}

func (svc *AzureRecommendationsV1) Patch(recommendations RecommendationStateV1) (RecommendationResponseV1, error) {
	return svc.client.Recommendations.patch(AzureRecommendation, recommendations)
}

// GetReport This is an experimental feature. Returned RecommendationID's are not guaranteed to be correct. Scoped to Lacework Account/Subaccount
func (svc *AzureRecommendationsV1) GetReport(reportType string) ([]RecommendationV1, error) {
	var (
		schemaBytes []byte
		ok          bool
	)
	report := struct {
		Ids []string `json:"recommendation_ids"`
	}{}

	switch reportType {
	case "AZURE_CIS":
		schemaBytes, ok = databox.Get("/reports/azure/cis.json")
		if !ok {
			return []RecommendationV1{}, errors.New(
				"compliance report schema not found",
			)
		}
	case "AZURE_CIS_131":
		schemaBytes, ok = databox.Get("/reports/azure/cis_131.json")
		if !ok {
			return []RecommendationV1{}, errors.New(
				"compliance report schema not found",
			)
		}
	default:
		return nil, errors.New(fmt.Sprintf("unable to find recommendations for report type %s", reportType))
	}

	err := json.Unmarshal(schemaBytes, &report)
	if err != nil {
		return []RecommendationV1{}, err
	}

	schema := ReportSchema{reportType, report.Ids}

	// fetch all azure recommendations
	allRecommendations, err := svc.client.Recommendations.Azure.List()
	if err != nil {
		return []RecommendationV1{}, err
	}
	filteredRecommendations := filterRecommendations(allRecommendations, schema)
	return filteredRecommendations, nil
}
