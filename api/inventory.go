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

type InventoryType string

// Deprecated: use `InventoryType` instead
type inventoryDataset string

const (
	AwsInventoryType   InventoryType = "AWS"
	AzureInventoryType InventoryType = "Azure"
	GcpInventoryType   InventoryType = "GCP"

	// Deprecated
	AwsInventoryDataset inventoryDataset = "AwsCompliance"
)

// Search retrieves information about resources in your cloud integrations, such as virtual machines,
// S3 buckets, security groups, and more. This function expects the response and the search filters.
//
// e.g.
//
//	 var (
//		  awsInventorySearchResponse api.InventoryAwsResponse
//		  filter = api.InventorySearch{
//			  SearchFilter: api.SearchFilter{
//				  Filters: []api.Filter{{
//					  Expression: "eq",
//					  Field:      "urn",
//					  Value:      "arn:aws:s3:::my-bucket",
//				  }},
//			  },
//			  Csp: api.AwsInventoryType,
//		  }
//	 )
//	 lacework.V2.Inventory.Search(&awsInventorySearchResponse, filters)
func (svc *InventoryService) Search(response interface{}, filters SearchableFilter) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2InventorySearch, filters, response)
}

// Scan triggers a resource inventory scan
func (svc *InventoryService) Scan(cloud InventoryType) (response InventoryScanResponse, err error) {
	url := fmt.Sprintf(apiV2InventoryScanCsp, cloud)
	err = svc.client.RequestEncoderDecoder("POST", url, nil, &response)
	return
}

type InventorySearch struct {
	SearchFilter
	Csp InventoryType `json:"csp"`

	// Deprecated: use `Csp` instead
	Dataset inventoryDataset `json:"dataset,omitempty"`
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

type InventoryGcpResponse struct {
	Data   []InventoryGcp `json:"data"`
	Paging V2Pagination   `json:"paging"`

	v2PageMetadata `json:"-"`
}

func (r InventoryGcpResponse) GetDataLength() int {
	return len(r.Data)
}

func (r InventoryGcpResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *InventoryGcpResponse) ResetPaging() {
	r.Paging = V2Pagination{}
	r.Data = nil
}

type InventoryGcp struct {
	InventoryCommon

	CloudDetails struct {
		FolderIds        []string      `json:"folderIds"`
		FolderNames      []interface{} `json:"folderNames"`
		ParentResourceID string        `json:"parentResourceId"`
		ProjectID        string        `json:"projectId"`
		ProjectName      string        `json:"projectName"`
		ProjectNumber    string        `json:"projectNumber"`
	} `json:"cloudDetails"`
	ResourceConfig struct {
		CreationTimestamp string `json:"creationTimestamp"`
		Description       string `json:"description"`
		DestRange         string `json:"destRange"`
		ID                string `json:"id"`
		Name              string `json:"name"`
		Network           string `json:"network"`
		NextHopNetwork    string `json:"nextHopNetwork"`
		Priority          int    `json:"priority"`
		SelfLink          string `json:"selfLink"`
	} `json:"resourceConfig"`
}

type InventoryCommon struct {
	ApiKey         string    `json:"apiKey"`
	Csp            string    `json:"csp"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	ResourceID     string    `json:"resourceId"`
	ResourceRegion string    `json:"resourceRegion"`
	ResourceType   string    `json:"resourceType"`
	Service        string    `json:"service"`
	Urn            string    `json:"urn"`

	Status struct {
		FormatVersion int    `json:"formatVersion"`
		Props         any    `json:"props"`
		Status        string `json:"status"`
		ErrorType     string `json:"errorType,omitempty"`
		ErrorMessage  string `json:"errorMessage,omitempty"`
	} `json:"status"`
}

type InventoryRawResponse struct {
	Data   []InventoryCommon `json:"data"`
	Paging V2Pagination      `json:"paging"`

	v2PageMetadata `json:"-"`
}

func (r InventoryRawResponse) GetDataLength() int {
	return len(r.Data)
}

func (r InventoryRawResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *InventoryRawResponse) ResetPaging() {
	r.Paging = V2Pagination{}
	r.Data = nil
}
