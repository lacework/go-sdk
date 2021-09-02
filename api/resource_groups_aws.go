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

// Wildcard values for selecting all aws accounts
var (
	AwsResourceGroupAllAccounts = []string{"*"}
)

// GetAwsResourceGroup gets a single Aws ResourceGroup matching the
// provided resource guid
func (svc *ResourceGroupsService) GetAwsResourceGroup(guid string) (
	response AwsResourceGroupResponse,
	err error,
) {
	var rawResponse resourceGroupWorkaroundResponse
	err = svc.get(guid, &rawResponse)
	if err != nil {
		return
	}

	return setAwsResourceGroupResponse(rawResponse)
}

// UpdateAwsResourceGroup updates a single Aws ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) UpdateAwsResourceGroup(data ResourceGroup) (
	response AwsResourceGroupResponse,
	err error,
) {
	var rawResponse resourceGroupWorkaroundResponse
	err = svc.update(data.ID(), data, &rawResponse)
	if err != nil {
		return
	}

	return setAwsResourceGroupResponse(rawResponse)
}

func setAwsResourceGroupResponse(response resourceGroupWorkaroundResponse) (aws AwsResourceGroupResponse, err error) {
	var props AwsResourceGroupProps

	isDefault, err := strconv.Atoi(response.Data.IsDefault)
	if err != nil {
		return AwsResourceGroupResponse{}, err
	}

	aws = AwsResourceGroupResponse{
		Data: AwsResourceGroupData{
			Guid:         response.Data.Guid,
			IsDefault:    isDefault,
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
	aws.Data.Props = props
	return
}

type AwsResourceGroupResponse struct {
	Data AwsResourceGroupData `json:"data"`
}

type AwsResourceGroupData struct {
	Guid         string                `json:"guid,omitempty"`
	IsDefault    int                   `json:"isDefault,omitempty"`
	ResourceGuid string                `json:"resourceGuid,omitempty"`
	Name         string                `json:"resourceName"`
	Type         string                `json:"resourceType"`
	Enabled      int                   `json:"enabled,omitempty"`
	Props        AwsResourceGroupProps `json:"props"`
}

type AwsResourceGroupProps struct {
	Description string   `json:"DESCRIPTION,omitempty"`
	AccountIDs  []string `json:"ACCOUNT_IDS"`
	UpdatedBy   string   `json:"UPDATED_BY,omitempty"`
	LastUpdated int      `json:"LAST_UPDATED,omitempty"`
}
