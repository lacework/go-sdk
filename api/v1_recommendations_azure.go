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

// AzureRecommendationsV1 is a service that interacts with the V1 Recommendations
// endpoints from the Lacework Server
type AzureRecommendationsV1 struct {
	client *Client
}

func (svc *AzureRecommendationsV1) List() ([]RecommendationV1, error) {
	return svc.client.Recommendations.list(AzureRecommendation)
}

func (svc *AzureRecommendationsV1) Patch(recommendations RecommendationStateV1) (response RecommendationResponseV1, err error) {
	return svc.client.Recommendations.patch(AzureRecommendation, recommendations)
}

func (svc *AzureRecommendationsV1) GetReport(reportType string) (response []RecommendationV1, err error) {
	return []RecommendationV1{}, nil
}
