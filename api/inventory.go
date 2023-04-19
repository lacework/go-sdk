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
	"fmt"
	"time"
)

type InventoryService struct {
	client *Client
}

type inventoryType string
type inventoryDataset string

const AwsInventoryType inventoryType = "AWS"
const AzureInventoryType inventoryType = "Azure"
const GcpInventoryType inventoryType = "GCP"
const AwsInventoryDataset inventoryDataset = "AwsCompliance"

// Search expects the response and the search filters
//
// e.g.
//
//  var (
//	  awsInventorySearchResponse api.InventoryAwsResponse
//	  filter = api.InventorySearch{
//		  SearchFilter: api.SearchFilter{
//			  Filters: []api.Filter{{
//				  Expression: "eq",
//				  Field:      "urn",
//				  Value:      arn:aws:s3:::my-bucket,
//			  }},
//		  },
//		  Dataset: api.AwsComplianceEvaluationDataset,
//	  }
//  )
//   lacework.V2.Inventory.Search(&awsInventorySearchResponse, filters)
//
func (svc *InventoryService) Search(response interface{}, filters SearchableFilter) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2InventorySearch, filters, response)
}

// Scan triggers a resource inventory scan
func (svc *InventoryService) Scan(cloud inventoryType) (response InventoryScanResponse, err error) {
	url := fmt.Sprintf(apiV2InventoryScanCsp, cloud)
	err = svc.client.RequestEncoderDecoder("POST", url, nil, &response)
	return
}

type InventorySearch struct {
	SearchFilter
	Csp     inventoryType    `json:"csp"`
	Dataset inventoryDataset `json:"dataset"`
}

func (i InventorySearch) GetTimeFilter() *TimeFilter {
	return i.TimeFilter
}

func (i InventorySearch) SetStartTime(time *time.Time) {
	i.TimeFilter.StartTime = time
}

func (i InventorySearch) SetEndTime(time *time.Time) {
	i.TimeFilter.EndTime = time
}

type InventoryScanResponse struct {
	Data struct {
		Status  string `json:"status"`
		Details string `json:"details"`
	} `json:"data"`
}
