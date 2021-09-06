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
	"strconv"

	"github.com/pkg/errors"
)

// Wildcard values for selecting all azure subscriptions
var (
	AzureResourceGroupAllSubscriptions = []string{"*"}
)

// GetAzureResourceGroup gets a single Azure ResourceGroup matching the
// provided resource guid
func (svc *ResourceGroupsService) GetAzureResourceGroup(guid string) (
	response AzureResourceGroupResponse,
	err error,
) {
	var rawResponse resourceGroupWorkaroundResponse
	err = svc.get(guid, &rawResponse)
	if err != nil {
		return
	}

	return setAzureResponse(rawResponse)
}

// UpdateAzureResourceGroup updates a single Azure ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) UpdateAzureResourceGroup(data ResourceGroup) (
	response AzureResourceGroupResponse,
	err error,
) {
	var rawResponse resourceGroupWorkaroundResponse
	err = svc.update(data.ID(), data, &rawResponse)
	if err != nil {
		return
	}

	return setAzureResponse(rawResponse)
}

// CreateAzureResourceGroup creates a single Azure ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) CreateAzureResourceGroup(data ResourceGroup) (
	response AzureResourceGroupResponse,
	err error,
) {
	var rawResponse ResourceGroupResponse
	err = svc.create(data, &rawResponse)
	if err != nil {
		return
	}

	return setAzureResourceGroupCreateResponse(rawResponse)
}

func setAzureResourceGroupCreateResponse(response ResourceGroupResponse) (azure AzureResourceGroupResponse, err error) {
	var props AzureResourceGroupProps

	azure = AzureResourceGroupResponse{
		Data: AzureResourceGroupData{
			Guid:         response.Data.Guid,
			IsDefault:    response.Data.IsDefault,
			ResourceGuid: response.Data.ResourceGuid,
			Name:         response.Data.Name,
			Type:         response.Data.Type,
			Enabled:      response.Data.Enabled,
		},
	}

	propsString, ok := response.Data.Props.(string)
	if !ok {
		err = errors.New("unable to cast props field from API response")
		return
	}

	err = json.Unmarshal([]byte(propsString), &props)
	if err != nil {
		return
	}
	azure.Data.Props = props
	return
}

func setAzureResponse(response resourceGroupWorkaroundResponse) (az AzureResourceGroupResponse, err error) {
	var props AzureResourceGroupProps

	isDefault, err := strconv.Atoi(response.Data.IsDefault)
	if err != nil {
		return
	}

	az = AzureResourceGroupResponse{
		Data: AzureResourceGroupData{
			Guid:         response.Data.Guid,
			IsDefault:    isDefault,
			ResourceGuid: response.Data.ResourceGuid,
			Name:         response.Data.Name,
			Type:         response.Data.Type,
			Enabled:      response.Data.Enabled,
		},
	}

	propsString, ok := response.Data.Props.(string)
	if !ok {
		err = errors.New("unable to cast props field from API response")
		return
	}

	err = json.Unmarshal([]byte(propsString), &props)
	if err != nil {
		return
	}
	az.Data.Props = props
	return
}

type AzureResourceGroupResponse struct {
	Data AzureResourceGroupData `json:"data"`
}

type AzureResourceGroupData struct {
	Guid         string                  `json:"guid,omitempty"`
	IsDefault    int                     `json:"isDefault,omitempty"`
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
