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

import "fmt"

// RecommendationsService is a service that interacts with the Recommendations
// endpoints from the Lacework Server
type RecommendationsService struct {
	client *Client
}

type RecommendationType string

const (
	AwsRecommendation   RecommendationType = "aws"
	AzureRecommendation RecommendationType = "azure"
	GcpRecommendation   RecommendationType = "gcp"
)

func (svc *RecommendationsService) AwsList() ([]Recommendation, error) {
	return svc.list(AwsRecommendation)
}

func (svc *RecommendationsService) AzureList() ([]Recommendation, error) {
	return svc.list(AzureRecommendation)
}

func (svc *RecommendationsService) GcpList() ([]Recommendation, error) {
	return svc.list(GcpRecommendation)
}

func (svc *RecommendationsService) PatchAws(recommendations RecommendationState) (response RecommendationResponse, err error) {
	return svc.patch(AwsRecommendation, recommendations)
}

func (svc *RecommendationsService) PatchAzure(recommendations RecommendationState) (response RecommendationResponse, err error) {
	return svc.patch(AzureRecommendation, recommendations)
}

func (svc *RecommendationsService) PatchGcp(recommendations RecommendationState) (response RecommendationResponse, err error) {
	return svc.patch(GcpRecommendation, recommendations)
}

func (svc *RecommendationsService) list(cloudType RecommendationType) ([]Recommendation, error) {
	var response RecommendationResponse
	err := svc.client.RequestDecoder("GET", fmt.Sprintf(apiRecommendations, cloudType), nil, &response)
	return response.RecommendationList(), err
}

func (svc *RecommendationsService) patch(cloudType RecommendationType, recommendations RecommendationState) (
	response RecommendationResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("PATCH", fmt.Sprintf(apiRecommendations, cloudType), recommendations, &response)
	return
}

type RecommendationState map[string]string

type RecommendationData map[string]RecommendationEnabled

type Recommendation struct {
	ID    string
	State bool
}

type RecommendationEnabled struct {
	Enabled bool `json:"enabled"`
}

type RecommendationResponse struct {
	Data    []RecommendationData `json:"data"`
	Ok      bool                 `json:"ok"`
	Message string               `json:"message"`
}

func (res *RecommendationResponse) RecommendationList() (recommendations []Recommendation) {
	if len(res.Data) > 0 {
		for k, v := range res.Data[0] {
			recommendations = append(recommendations, Recommendation{k, v.Enabled})
		}
	}
	return
}
