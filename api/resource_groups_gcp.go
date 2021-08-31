//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

import "encoding/json"

// GetGcpResourceGroup gets a single Gcp ResourceGroup matching the
// provided resource guid
func (svc *ResourceGroupsService) GetGcpResourceGroup(guid string) (
	response GcpResourceGroupResponse,
	err error,
) {
	var rawResponse ResourceGroupResponse
	err = svc.get(guid, &rawResponse)
	if err != nil {
		return GcpResourceGroupResponse{}, err
	}

	return convertGcpResponse(rawResponse)
}

// UpdateGcpResourceGroup updates a single Gcp ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) UpdateGcpResourceGroup(data ResourceGroup) (
	response GcpResourceGroupResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

func convertGcpResponse(rawResponse ResourceGroupResponse) (gcpResponse GcpResourceGroupResponse, err error) {
	var props GcpResourceGroupProps
	err = json.Unmarshal([]byte(rawResponse.Data.Props.(string)), &props)
	if err != nil {
		return GcpResourceGroupResponse{}, err
	}

	gcpResponse, err = castGcpResponse(rawResponse)
	if err != nil {
		return GcpResourceGroupResponse{}, err
	}

	gcpResponse.Data.Props = props
	return gcpResponse, nil
}

func castGcpResponse(res interface{}) (r GcpResourceGroupResponse, err error) {
	var j []byte
	j, err = json.Marshal(res)
	if err != nil {
		return
	}
	err = json.Unmarshal(j, &r)
	return
}

type GcpResourceGroupResponse struct {
	Data GcpResourceGroupData `json:"data"`
}

type GcpResourceGroupData struct {
	Guid         string                `json:"guid,omitempty"`
	IsDefault    string                `json:"isDefault,omitempty"`
	ResourceGuid string                `json:"resourceGuid,omitempty"`
	Name         string                `json:"resourceName"`
	Type         string                `json:"resourceType"`
	Enabled      int                   `json:"enabled,omitempty"`
	Props        GcpResourceGroupProps `json:"props"`
}

type GcpResourceGroupProps struct {
	Description  string   `json:"DESCRIPTION,omitempty"`
	Organization string   `json:"ORGANIZATION"`
	Projects     []string `json:"PROJECTS"`
	UpdatedBy    string   `json:"UPDATED_BY,omitempty"`
	LastUpdated  int      `json:"LAST_UPDATED,omitempty"`
}
