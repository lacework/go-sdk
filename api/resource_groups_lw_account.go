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

// GetContainerResourceGroup gets a single LwAccount ResourceGroup matching the
// provided resource guid
func (svc *ResourceGroupsService) GetLwAccountResourceGroup(guid string) (
	response LwAccountResourceGroupResponse,
	err error,
) {
	var rawResponse ResourceGroupResponse
	err = svc.get(guid, &rawResponse)
	if err != nil {
		return LwAccountResourceGroupResponse{}, err
	}

	return convertLwAccountResponse(rawResponse)
}

// UpdateLwAccountResourceGroup updates a single LwAccount ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) UpdateLwAccountResourceGroup(data ResourceGroup) (
	response LwAccountResourceGroupResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

func convertLwAccountResponse(rawResponse ResourceGroupResponse) (lwAccountResponse LwAccountResourceGroupResponse, err error) {
	var props LwAccountResourceGroupProps
	err = json.Unmarshal([]byte(rawResponse.Data.Props.(string)), &props)
	if err != nil {
		return LwAccountResourceGroupResponse{}, err
	}

	lwAccountResponse, err = castLwAccountResponse(rawResponse)
	if err != nil {
		return LwAccountResourceGroupResponse{}, err
	}

	lwAccountResponse.Data.Props = props
	return lwAccountResponse, nil
}

func castLwAccountResponse(res interface{}) (r LwAccountResourceGroupResponse, err error) {
	var j []byte
	j, err = json.Marshal(res)
	if err != nil {
		return
	}
	err = json.Unmarshal(j, &r)
	return
}

type LwAccountResourceGroupResponse struct {
	Data LwAccountResourceGroupData `json:"data"`
}

type LwAccountResourceGroupData struct {
	Guid         string                      `json:"guid,omitempty"`
	IsDefault    string                      `json:"isDefault,omitempty"`
	ResourceGuid string                      `json:"resourceGuid,omitempty"`
	Name         string                      `json:"resourceName"`
	Type         string                      `json:"resourceType"`
	Enabled      int                         `json:"enabled,omitempty"`
	Props        LwAccountResourceGroupProps `json:"props"`
}

type LwAccountResourceGroupProps struct {
	Description string   `json:"DESCRIPTION,omitempty"`
	LwAccounts  []string `json:"LW_ACCOUNTS"`
	UpdatedBy   string   `json:"UPDATED_BY,omitempty"`
	LastUpdated int      `json:"LAST_UPDATED,omitempty"`
}
