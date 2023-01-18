//
// Author:: Ross Moles (<ross.moles@lacework.net>)
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
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// SuppressionsServiceV2 is a service that interacts with the V2 Suppressions
// endpoints from the Lacework Server
type SuppressionsServiceV2 struct {
	client *Client
	Aws    suppressionServiceV2
	Azure  suppressionServiceV2
	Gcp    suppressionServiceV2
}

type suppressionServiceV2 interface {
	List() (map[string]SuppressionV2, error)
}

type SuppressionTypeV2 string

const (
	AwsSuppression   SuppressionTypeV2 = "aws"
	AzureSuppression SuppressionTypeV2 = "azure"
	GcpSuppression   SuppressionTypeV2 = "gcp"
)

func (svc *SuppressionsServiceV2) list(cloudType SuppressionTypeV2) (map[string]SuppressionV2,
	error) {
	var response SuppressionResponseV2
	err := svc.client.RequestDecoder("GET", fmt.Sprintf(apiSuppressions, cloudType), nil, &response)
	return response.SuppressionList(), err
}

type SuppressionResponseV2 struct {
	Data    []SuppressionDataV2 `json:"data"`
	Ok      bool                `json:"ok"`
	Message string              `json:"message"`
}

type SuppressionDataV2 struct {
	RecommendationSuppressions map[string]map[string]interface{} `json:"recommendationExceptions"`
}

type SuppressionV2 struct {
	Enabled               bool                    `json:"enabled"`
	SuppressionConditions []SuppressionConditions `json:"suppressionConditions"`
}

type SuppressionConditions struct {
	AccountIds         []string            `json:"accountIds,omitempty"`
	OrganizationIds    []string            `json:"organizationIds,omitempty"`
	ProjectIds         []string            `json:"projectIds,omitempty"`
	RegionNames        []string            `json:"regionNames,omitempty"`
	ResourceLabels     []map[string]string `json:"resourceLabels,omitempty"`
	ResourceGroupNames []string            `json:"resourceGroupNames,omitempty"`
	ResourceNames      []string            `json:"resourceNames,omitempty"`
	ResourceTags       []map[string]string `json:"resourceTags,omitempty"`
	SubscriptionIds    []string            `json:"subscriptionIds,omitempty"`
	TenantIds          []string            `json:"tenantIds,omitempty"`
	Comment            string              `json:"comments,omitempty"`
}

func (res *SuppressionResponseV2) SuppressionList() (suppressions map[string]SuppressionV2) {
	if len(res.Data) > 0 {
		suppressions = make(map[string]SuppressionV2)
		for _, v := range res.Data {
			for key, val := range v.RecommendationSuppressions {
				var sup SuppressionV2

				err := mapstructure.Decode(val, &sup)
				if err != nil {
					return nil
				}

				suppressions[key] = sup
			}
		}
	}
	return
}
