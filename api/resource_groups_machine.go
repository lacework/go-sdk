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

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// GetMachineResourceGroup gets a single Machine ResourceGroup matching the
// provided resource guid
func (svc *ResourceGroupsService) GetMachineResourceGroup(guid string) (
	response MachineResourceGroupResponse,
	err error,
) {
	var rawResponse ResourceGroupResponse
	err = svc.get(guid, &rawResponse)
	if err != nil {
		return
	}

	return setMachineAccountResponse(rawResponse)
}

// UpdateMachineResourceGroup updates a single Machine ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) UpdateMachineResourceGroup(data ResourceGroup) (
	response MachineResourceGroupResponse,
	err error,
) {
	var rawResponse ResourceGroupResponse
	err = svc.update(data.ID(), data, &rawResponse)
	if err != nil {
		return
	}

	return setMachineAccountResponse(rawResponse)
}

func setMachineAccountResponse(response ResourceGroupResponse) (machine MachineResourceGroupResponse, err error) {
	var props MachineResourceGroupProps
	machine = MachineResourceGroupResponse{
		Data: MachineResourceGroupData{
			Guid:         response.Data.Guid,
			IsDefault:    response.Data.IsDefault,
			ResourceGuid: response.Data.ResourceGuid,
			Name:         response.Data.Name,
			Type:         response.Data.Type,
			Enabled:      response.Data.Enabled,
		},
	}

	_, ok := response.Data.Props.(string)
	if !ok {
		err = errors.New("unable to cast props field from API response")
		return
	}

	err = json.Unmarshal([]byte(response.Data.Props.(string)), &props)
	if err != nil {
		return
	}
	machine.Data.Props = props
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
