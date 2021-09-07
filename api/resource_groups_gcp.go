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

// Wildcard values for selecting all gcp projects
var (
	GcpResourceGroupAllProjects = []string{"*"}
)

// GetGcp gets a single Gcp ResourceGroup matching the
// provided resource guid
func (svc *ResourceGroupsService) GetGcp(guid string) (
	response GcpResourceGroupResponse,
	err error,
) {
	var rawResponse resourceGroupWorkaroundResponse
	err = svc.get(guid, &rawResponse)
	if err != nil {
		return
	}

	return setGcpResponse(rawResponse)
}

// UpdateGcp updates a single Gcp ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) UpdateGcp(data ResourceGroup) (
	response GcpResourceGroupResponse,
	err error,
) {
	if data == nil {
		return GcpResourceGroupResponse{}, errors.New("resource group must not be empty")
	}
	guid := data.ID()
	data.ResetResourceGUID()

	err = svc.update(guid, data, &response)
	return
}

// CreateGcp creates a single Gcp ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) CreateGcp(data ResourceGroup) (
	response GcpResourceGroupResponse,
	err error,
) {
	err = svc.create(data, &response)
	if err != nil {
		return
	}

	return
}

func setGcpResponse(response resourceGroupWorkaroundResponse) (gcp GcpResourceGroupResponse, err error) {
	var props GcpResourceGroupJsonStringProps

	isDefault, err := strconv.Atoi(response.Data.IsDefault)
	if err != nil {
		return
	}

	gcp = GcpResourceGroupResponse{
		Data: GcpResourceGroupData{
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
	gcp.Data.Props = GcpResourceGroupProps(props)
	return
}

type GcpResourceGroupResponse struct {
	Data GcpResourceGroupData `json:"data"`
}

type GcpResourceGroupData struct {
	Guid         string                `json:"guid,omitempty"`
	IsDefault    int                   `json:"isDefault,omitempty"`
	ResourceGuid string                `json:"resourceGuid,omitempty"`
	Name         string                `json:"resourceName"`
	Type         string                `json:"resourceType"`
	Enabled      int                   `json:"enabled,omitempty"`
	Props        GcpResourceGroupProps `json:"props"`
}

type GcpResourceGroupProps struct {
	Description  string   `json:"description,omitempty"`
	Organization string   `json:"organization"`
	Projects     []string `json:"projects"`
	UpdatedBy    string   `json:"updatedBy,omitempty"`
	LastUpdated  int      `json:"lastUpdated,omitempty"`
}

// Workaround for props being returned as a json string
type GcpResourceGroupJsonStringProps struct {
	Description  string   `json:"DESCRIPTION,omitempty"`
	Organization string   `json:"ORGANIZATION"`
	Projects     []string `json:"PROJECTS"`
	UpdatedBy    string   `json:"UPDATED_BY,omitempty"`
	LastUpdated  int      `json:"LAST_UPDATED,omitempty"`
}
