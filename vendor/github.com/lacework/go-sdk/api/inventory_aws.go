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

type InventoryAwsResponse struct {
	Data   []InventoryAws `json:"data"`
	Paging V2Pagination   `json:"paging"`
}

func (r InventoryAwsResponse) GetDataLength() int {
	return len(r.Data)
}

func (r InventoryAwsResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *InventoryAwsResponse) ResetPaging() {
	r.Paging = V2Pagination{}
	r.Data = nil
}

type InventoryAws struct {
	ApiKey         string `json:"apiKey"`
	Csp            string `json:"csp"`
	EndTime        string `json:"endTime"`
	StartTime      string `json:"startTime"`
	ResourceId     string `json:"resourceId"`
	ResourceRegion string `json:"resourceRegion"`
	ResourceTags   any    `json:"resourceTags"`
	ResourceType   string `json:"resourceType"`
	Service        string `json:"service"`
	Urn            string `json:"urn"`
	CloudDetails   struct {
		AccountAlias string `json:"accountAlias"`
		AccountID    string `json:"accountID"`
	} `json:"cloudDetails"`
	Status struct {
		FormatVersion int    `json:"formatVersion"`
		Props         any    `json:"props"`
		Status        string `json:"status"`
		// Error status
		ErrorMessage string `json:"errorMessage,omitempty"`
		ErrorType    string `json:"errorType,omitempty"`
	} `json:"status"`
	ResourceConfig any `json:"resourceConfig"`
}
