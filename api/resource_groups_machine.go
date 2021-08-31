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

// GetMachineResourceGroup gets a single Machine ResourceGroup matching the
// provided resource guid
func (svc *ResourceGroupsService) GetMachineResourceGroup(guid string) (
	response MachineResourceGroupResponse,
	err error,
) {
	var rawResponse ResourceGroupResponse
	err = svc.get(guid, &rawResponse)
	if err != nil {
		return MachineResourceGroupResponse{}, err
	}

	return convertMachineResponse(rawResponse)
}

// UpdateMachineResourceGroup updates a single Machine ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) UpdateMachineResourceGroup(data ResourceGroup) (
	response MachineResourceGroupResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

func convertMachineResponse(rawResponse ResourceGroupResponse) (machineResponse MachineResourceGroupResponse, err error) {
	var props MachineResourceGroupProps
	err = json.Unmarshal([]byte(rawResponse.Data.Props.(string)), &props)
	if err != nil {
		return MachineResourceGroupResponse{}, err
	}

	machineResponse, err = castMachineResponse(rawResponse)
	if err != nil {
		return MachineResourceGroupResponse{}, err
	}

	machineResponse.Data.Props = props
	return machineResponse, nil
}

func castMachineResponse(res interface{}) (r MachineResourceGroupResponse, err error) {
	var j []byte
	j, err = json.Marshal(res)
	if err != nil {
		return
	}
	err = json.Unmarshal(j, &r)
	return
}

type MachineResourceGroupResponse struct {
	Data MachineResourceGroupData `json:"data"`
}

type MachineResourceGroupData struct {
	Guid         string                    `json:"guid,omitempty"`
	IsDefault    string                    `json:"isDefault,omitempty"`
	ResourceGuid string                    `json:"resourceGuid,omitempty"`
	Name         string                    `json:"resourceName"`
	Type         string                    `json:"resourceType"`
	Enabled      int                       `json:"enabled,omitempty"`
	Props        MachineResourceGroupProps `json:"props"`
}

type MachineResourceGroupProps struct {
	Description string              `json:"DESCRIPTION,omitempty"`
	MachineTags []map[string]string `json:"MACHINE_TAGS"`
	UpdatedBy   string              `json:"UPDATED_BY,omitempty"`
	LastUpdated int                 `json:"LAST_UPDATED,omitempty"`
}
