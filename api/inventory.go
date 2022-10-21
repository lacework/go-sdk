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

type InventoryService struct {
	client *Client
}

type inventoryType string
type inventoryDataset string

const AwsInventoryType inventoryType = "AWS"
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
func (svc *InventoryService) Search(response interface{}, filters InventorySearch) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2InventorySearch, filters, response)
}

type InventorySearch struct {
	SearchFilter
	Csp     inventoryType    `json:"csp"`
	Dataset inventoryDataset `json:"dataset"`
}

// WindowedSearch performs a new search of a specific time frame size,
// until response data is found or the max searchable days is reached
func (svc *InventoryService) WindowedSearch(response *InventoryAwsResponse, filter InventorySearch, size int, max int) error {
	for i := 0; i < max; i += size {
		err := svc.Search(&response, filter)
		if err != nil {
			return err
		}
		if len(response.Data) != 0 {
			return nil
		}

		//adjust window
		newStart := filter.StartTime.AddDate(0, 0, -size)
		newEnd := filter.EndTime.AddDate(0, 0, -size)

		// ensure we do not go over the max allowed searchable days
		rem := (i - max) % size
		if rem > 0 {
			newEnd = filter.EndTime.AddDate(0, 0, -rem)
		}
		filter.TimeFilter.StartTime = &newStart
		filter.TimeFilter.EndTime = &newEnd
	}
	return nil
}
