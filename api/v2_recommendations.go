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
	"fmt"
)

// RecommendationsServiceV2 is a service that interacts with the V2 Recommendations
// endpoints from the Lacework Server
type RecommendationsServiceV2 struct {
	client *Client
	Aws    recommendationServiceV2
	Azure  recommendationServiceV2
	Gcp    recommendationServiceV2
}

type recommendationServiceV2 interface {
	List() ([]RecV2, error)
	Patch(recommendations RecommendationStateV2) (RecommendationResponseV2, error)
	GetReport(reportType string) ([]RecV2, error)
}

type RecommendationTypeV2 string

const (
	AwsRecommendation   RecommendationTypeV2 = "aws"
	AzureRecommendation RecommendationTypeV2 = "azure"
	GcpRecommendation   RecommendationTypeV2 = "gcp"
)

func (svc *RecommendationsServiceV2) list(cloudType RecommendationTypeV2) ([]RecV2, error) {
	var response RecommendationResponseV2
	err := svc.client.RequestDecoder("GET", fmt.Sprintf(apiRecommendations, cloudType), nil, &response)
	return response.RecommendationList(), err
}

func (svc *RecommendationsServiceV2) patch(cloudType RecommendationTypeV2, recommendations RecommendationStateV2) (
	response RecommendationResponseV2,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("PATCH", fmt.Sprintf(apiRecommendations, cloudType), recommendations, &response)
	return
}

type RecommendationStateV2 map[string]string

type RecommendationDataV2 map[string]RecommendationEnabledV2

type RecV2 struct {
	ID    string
	State bool
}

type RecommendationEnabledV2 struct {
	Enabled bool `json:"enabled"`
}

type RecommendationResponseV2 struct {
	Data    []RecommendationDataV2 `json:"data"`
	Ok      bool                   `json:"ok"`
	Message string                 `json:"message"`
}

func (res *RecommendationResponseV2) RecommendationList() (recommendations []RecV2) {
	if len(res.Data) > 0 {
		for k, v := range res.Data[0] {
			recommendations = append(recommendations, RecV2{k, v.Enabled})
		}
	}
	return
}

type ReportSchema struct {
	Name              string            `json:"name"`
	RecommendationIDs map[string]string `json:"recommendationIDs"`
}

func NewRecommendationV2State(recommendations []RecV2, state bool) RecommendationStateV2 {
	request := make(map[string]string)
	for _, rec := range recommendations {
		if state {
			request[rec.ID] = "enable"

		} else {
			request[rec.ID] = "disable"
		}
	}
	return request
}

func NewRecommendationV2(recommendations []RecV2) RecommendationStateV2 {
	request := make(map[string]string)
	for _, rec := range recommendations {
		if rec.State {
			request[rec.ID] = "enable"

		} else {
			request[rec.ID] = "disable"
		}
	}
	return request
}

// ReportStatus This is an experimental feature. Returned RecommendationID's are not guaranteed to be correct.
func (res *RecommendationResponseV2) ReportStatus() map[string]bool {
	var recommendations = make(map[string]bool)

	for _, rec := range res.RecommendationList() {
		recommendations[rec.ID] = rec.State
	}

	return recommendations
}

// filterRecommendations This is an experimental feature. Returned RecommendationID's are not guaranteed to be correct.
func filterRecommendations(allRecommendations []RecV2, schema ReportSchema) []RecV2 {
	var recommendations []RecV2

	for _, rec := range allRecommendations {
		_, ok := schema.RecommendationIDs[rec.ID]
		if ok {
			recommendations = append(recommendations, rec)
		}
	}
	return recommendations
}
