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
	"fmt"

	"github.com/lacework/go-sdk/internal/databox"
)

// GcpRecommendationsV2 is a service that interacts with the V2 Recommendations
// endpoints from the Lacework Server
type GcpRecommendationsV2 struct {
	client *Client
}

func (svc *GcpRecommendationsV2) List() ([]RecV2, error) {
	return svc.client.V2.Recommendations.list(GcpRecommendation)
}

func (svc *GcpRecommendationsV2) Patch(recommendations RecommendationStateV2) (RecommendationResponseV2, error) {
	return svc.client.V2.Recommendations.patch(GcpRecommendation, recommendations)
}

// GetReport This is an experimental feature. Returned RecommendationID's are not guaranteed to be correct. Scoped to Lacework Account/Subaccount
func (svc *GcpRecommendationsV2) GetReport(reportType string) ([]RecV2, error) {
	var (
		schemaBytes []byte
		ok          bool
	)
	report := struct {
		Ids map[string]string `json:"recommendation_ids"`
	}{}

	switch reportType {
	case "CIS_1_0":
		schemaBytes, ok = databox.Get("/reports/gcp/cis.json")
		if !ok {
			return []RecV2{}, errors.New(
				"compliance report schema not found",
			)
		}
	case "CIS_1_2":
		schemaBytes, ok = databox.Get("/reports/gcp/cis_12.json")
		if !ok {
			return []RecV2{}, errors.New(
				"compliance report schema not found",
			)
		}
	default:
		return nil, fmt.Errorf("unable to find recommendations for report type %s", reportType)
	}

	err := json.Unmarshal(schemaBytes, &report)
	if err != nil {
		return []RecV2{}, err
	}

	schema := ReportSchema{reportType, report.Ids}

	// fetch all azure recommendations
	allRecommendations, err := svc.client.V2.Recommendations.Gcp.List()
	if err != nil {
		return []RecV2{}, err
	}
	filteredRecommendations := filterRecommendations(allRecommendations, schema)
	return filteredRecommendations, nil
}
