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
	"strings"
)

// GcpRecommendationsV1 is a service that interacts with the V1 Recommendations
// endpoints from the Lacework Server
type GcpRecommendationsV1 struct {
	client *Client
}

const gcpCIS = "GCP_CIS_"
const gcpCIS12 = "GCP_CIS12_"

func (svc *GcpRecommendationsV1) List() ([]RecommendationV1, error) {
	return svc.client.Recommendations.list(GcpRecommendation)
}

func (svc *GcpRecommendationsV1) Patch(recommendations RecommendationStateV1) (response RecommendationResponseV1, err error) {
	return svc.client.Recommendations.patch(GcpRecommendation, recommendations)
}

// GetReport This is an experimental feature. Returned RecommendationID's are not guaranteed to be correct. Scoped to Lacework Account/Subaccount
func (svc *GcpRecommendationsV1) GetReport(reportType string) ([]RecommendationV1, error) {
	// fetch all gcp recommendations
	recommendations, err := svc.client.Recommendations.Gcp.List()
	if err != nil {
		return []RecommendationV1{}, err
	}

	switch reportType {
	case "GCP_CIS":
		return matchRecommendations(recommendations, gcpCIS), nil
	case "GCP_CIS12":
		return matchRecommendations(recommendations, gcpCIS12), nil
	default:
		return nil, errors.New(fmt.Sprintf("unable to find recommendations for report type %s", reportType))
	}
}

func matchRecommendations(allRecommendations []RecommendationV1, prefix string) []RecommendationV1 {
	var recommendations []RecommendationV1
	for _, rec := range allRecommendations {
		if strings.HasPrefix(rec.ID, prefix) {
			recommendations = append(recommendations, rec)
		}
	}
	return recommendations
}
