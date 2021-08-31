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

// GetAzureResourceGroup gets a single Azure ResourceGroup matching the
// provided resource guid
func (svc *ResourceGroupsService) GetAzureResourceGroup(guid string) (
	response AzureResourceGroupResponse,
	err error,
) {
	var rawResponse ResourceGroupResponse
	err = svc.get(guid, &rawResponse)
	if err != nil {
		return AzureResourceGroupResponse{}, err
	}

	return convertAzureResponse(rawResponse)
}

// UpdateAzureResourceGroup updates a single Azure ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) UpdateAzureResourceGroup(data ResourceGroup) (
	response AzureResourceGroupResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

func convertAzureResponse(rawResponse ResourceGroupResponse) (azureResponse AzureResourceGroupResponse, err error) {
	var props AzureResourceGroupProps
	err = json.Unmarshal([]byte(rawResponse.Data.Props.(string)), &props)
	if err != nil {
		return AzureResourceGroupResponse{}, err
	}

	azureResponse, err = castAzureResponse(rawResponse)
	if err != nil {
		return AzureResourceGroupResponse{}, err
	}

	azureResponse.Data.Props = props
	return azureResponse, nil
}

func castAzureResponse(res interface{}) (r AzureResourceGroupResponse, err error) {
	var j []byte
	j, err = json.Marshal(res)
	if err != nil {
		return
	}
	err = json.Unmarshal(j, &r)
	return
}

type AzureResourceGroupResponse struct {
	Data AzureResourceGroupData `json:"data"`
}

type AzureResourceGroupData struct {
	Guid         string                  `json:"guid,omitempty"`
	IsDefault    string                  `json:"isDefault,omitempty"`
	ResourceGuid string                  `json:"resourceGuid,omitempty"`
	Name         string                  `json:"resourceName"`
	Type         string                  `json:"resourceType"`
	Enabled      int                     `json:"enabled,omitempty"`
	Props        AzureResourceGroupProps `json:"props"`
}

type AzureResourceGroupProps struct {
	Description   string   `json:"DESCRIPTION,omitempty"`
	Tenant        string   `json:"TENANT"`
	Subscriptions []string `json:"SUBSCRIPTIONS"`
	UpdatedBy     string   `json:"UPDATED_BY,omitempty"`
	LastUpdated   int      `json:"LAST_UPDATED,omitempty"`
}
