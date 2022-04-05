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

// RecommendationsServiceV1 is a service that interacts with the V1 Recommendations
// endpoints from the Lacework Server
type RecommendationsServiceV1 struct {
	client *Client
	Aws    recommendationServiceV1
	Azure  recommendationServiceV1
	Gcp    recommendationServiceV1
}

type recommendationServiceV1 interface {
	List() ([]RecommendationV1, error)
	Patch(recommendations RecommendationStateV1) (response RecommendationResponseV1, err error)
}

type RecommendationTypeV1 string

const (
	AwsRecommendation   RecommendationTypeV1 = "aws"
	AzureRecommendation RecommendationTypeV1 = "azure"
	GcpRecommendation   RecommendationTypeV1 = "gcp"
)

func (svc *RecommendationsServiceV1) list(cloudType RecommendationTypeV1) ([]RecommendationV1, error) {
	var response RecommendationResponseV1
	err := svc.client.RequestDecoder("GET", fmt.Sprintf(apiRecommendations, cloudType), nil, &response)
	return response.RecommendationList(), err
}

func (svc *RecommendationsServiceV1) patch(cloudType RecommendationTypeV1, recommendations RecommendationStateV1) (
	response RecommendationResponseV1,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("PATCH", fmt.Sprintf(apiRecommendations, cloudType), recommendations, &response)
	return
}

type RecommendationStateV1 map[string]string

type RecommendationDataV1 map[string]RecommendationEnabledV1

type RecommendationV1 struct {
	ID    string
	State bool
}

type RecommendationEnabledV1 struct {
	Enabled bool `json:"enabled"`
}

type RecommendationResponseV1 struct {
	Data    []RecommendationDataV1 `json:"data"`
	Ok      bool                   `json:"ok"`
	Message string                 `json:"message"`
}

func (res *RecommendationResponseV1) RecommendationList() (recommendations []RecommendationV1) {
	if len(res.Data) > 0 {
		for k, v := range res.Data[0] {
			recommendations = append(recommendations, RecommendationV1{k, v.Enabled})
		}
	}
	return
}
