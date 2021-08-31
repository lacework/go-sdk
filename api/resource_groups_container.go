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

// GetContainerResourceGroup gets a single Container ResourceGroup matching the
// provided resource guid
func (svc *ResourceGroupsService) GetContainerResourceGroup(guid string) (
	response ContainerResourceGroupResponse,
	err error,
) {
	var rawResponse ResourceGroupResponse
	err = svc.get(guid, &rawResponse)
	if err != nil {
		return ContainerResourceGroupResponse{}, err
	}

	return convertContainerResponse(rawResponse)
}

// UpdateContainerResourceGroup updates a single Container ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) UpdateContainerResourceGroup(data ResourceGroup) (
	response ContainerResourceGroupResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

func convertContainerResponse(rawResponse ResourceGroupResponse) (containerResponse ContainerResourceGroupResponse, err error) {
	var props ContainerResourceGroupProps
	err = json.Unmarshal([]byte(rawResponse.Data.Props.(string)), &props)
	if err != nil {
		return ContainerResourceGroupResponse{}, err
	}

	containerResponse, err = castContainerResponse(rawResponse)
	if err != nil {
		return ContainerResourceGroupResponse{}, err
	}

	containerResponse.Data.Props = props
	return containerResponse, nil
}

func castContainerResponse(res interface{}) (r ContainerResourceGroupResponse, err error) {
	var j []byte
	j, err = json.Marshal(res)
	if err != nil {
		return
	}
	err = json.Unmarshal(j, &r)
	return
}

type ContainerResourceGroupResponse struct {
	Data ContainerResourceGroupData `json:"data"`
}

type ContainerResourceGroupData struct {
	Guid         string                      `json:"guid,omitempty"`
	IsDefault    string                      `json:"isDefault,omitempty"`
	ResourceGuid string                      `json:"resourceGuid,omitempty"`
	Name         string                      `json:"resourceName"`
	Type         string                      `json:"resourceType"`
	Enabled      int                         `json:"enabled,omitempty"`
	Props        ContainerResourceGroupProps `json:"props"`
}

type ContainerResourceGroupProps struct {
	Description     string              `json:"DESCRIPTION,omitempty"`
	ContainerLabels []map[string]string `json:"CONTAINER_LABELS"`
	ContainerTags   []string            `json:"CONTAINER_TAGS"`
	UpdatedBy       string              `json:"UPDATED_BY,omitempty"`
	LastUpdated     int                 `json:"LAST_UPDATED,omitempty"`
}
