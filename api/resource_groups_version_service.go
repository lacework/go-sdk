//
// Author:: Zeki Sherif(<zeki.sherif@lacework.net>)
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
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type ResourceGroupsVersionService struct {
	client                 *Client
	v1ResourceGroupService *ResourceGroupsService
	v2ResourceGroupService *ResourceGroupsV2Service
	featureFlagService     *FeatureFlagsService
}

type ResourceGroupsInterfaceData interface {
	GetProps() interface{}
	GetQuery() *RGQuery
}

func (group ResourceGroupData) GetProps() interface{} {
	return group.Props
}

func (group ResourceGroupData) GetQuery() *RGQuery {
	return nil
}

func NewResourceGroupsVersionService(c *Client) *ResourceGroupsVersionService {
	return &ResourceGroupsVersionService{
		c,
		&ResourceGroupsService{c},
		&ResourceGroupsV2Service{c},
		&FeatureFlagsService{c},
	}
}

// NewResourceGroup returns an instance of the ResourceGroupData struct with the
// provided ResourceGroup type, name and the props field as an interface{}.
//
// NOTE: This function must be used by any ResourceGroup type.
//
// Basic usage: Initialize a new ContainerResourceGroup struct, then
//
//	             use the new instance to do CRUD operations
//
//	  client, err := api.NewClient("account")
//	  if err != nil {
//	    return err
//	  }
//
//	  group := api.NewResourceGroup("container resource group",
//	    api.ContainerResourceGroup,
//	    api.ContainerResourceGroupData{
//	      Props: api.ContainerResourceGroupProps{
//				Description:     "all containers,
//				ContainerLabels: ContainerResourceGroupAllLabels,
//				ContainerTags:   ContainerResourceGroupAllTags,
//			},
//	    },
//	  )
//
//	  client.V2.ResourceGroups.Create(group)
func NewResourceGroup(name string, iType ResourceGroupType, props interface{}) ResourceGroupData {
	return ResourceGroupData{
		Name:    name,
		Type:    iType.String(),
		Enabled: 1,
		Props:   props,
	}
}

// NewResourceGroupWithQuery Only available with RGv2 beta
func NewResourceGroupWithQuery(name string, iType ResourceGroupType,
	description string, query *RGQuery) ResourceGroupDataWithQuery {
	return ResourceGroupDataWithQuery{
		Name:        name,
		Type:        iType.String(),
		Enabled:     1,
		Query:       query,
		Description: description,
	}
}

func isRGV2FlagEnabled(featureFlagService *FeatureFlagsService) bool {
	response, err := featureFlagService.GetFeatureFlagsMatchingPrefix(ApiV2CliFeatureFlag)

	if err != nil {
		return false
	}

	return len(response.Data.Flags) >= 1
}

func (svc *ResourceGroupsVersionService) Get(guid string, response interface{}) error {
	var rawResponse resourceGroupWorkaroundResponse
	err := svc.get(guid, &rawResponse)
	if err != nil {
		return err
	}

	if rawResponse.Data.Query != nil {
		return castRGV2WorkAroundResponse(rawResponse, response)
	} else {
		return castRGV1WorkAroundResponse(rawResponse, response)
	}
}

func (svc *ResourceGroupsVersionService) Create(group ResourceGroupsInterfaceData) (
	response ResourceGroupResponse,
	err error,
) {
	isV2FlagEnabled := isRGV2FlagEnabled(svc.featureFlagService)

	if group.GetProps() == nil && !isV2FlagEnabled && group.GetQuery() == nil {
		if isV2FlagEnabled {
			err = errors.New("Invalid request. Missing `query` field.")
		} else {
			err = errors.New("Invalid request. Missing `props` field.")
		}

		return
	}

	if group.GetProps() != nil {
		response, err = svc.v1ResourceGroupService.Create(group.(ResourceGroupData))
		return
	}

	createResponse, createErr := svc.v2ResourceGroupService.Create(group.(ResourceGroupDataWithQuery))
	if createErr != nil {
		err = createErr
		return
	}

	err = castResourceGroupV2Response(createResponse, &response)
	return
}

func (svc *ResourceGroupsVersionService) Update(group ResourceGroupsInterfaceData) (
	response ResourceGroupResponse,
	err error,
) {
	if group.GetProps() != nil {
		response, err = svc.v1ResourceGroupService.Update(group.(ResourceGroup))
		return
	}

	isV2FlagEnabled := isRGV2FlagEnabled(svc.featureFlagService)

	if isV2FlagEnabled {
		updateResponse, updateErr := svc.v2ResourceGroupService.Update(group.(ResourceGroup))
		if updateErr != nil {
			err = updateErr
			return
		}

		err = castResourceGroupV2Response(updateResponse, &response)
		return
	}

	err = errors.New("Unable to update resource group")
	return
}

func (svc *ResourceGroupsVersionService) Delete(guid string) error {
	// It doesn't matcher which version of service we use as api-server handles
	// delete for both v1 and v2 resource groups
	err := svc.v1ResourceGroupService.Delete(guid)

	if err != nil {
		return err
	}

	return nil
}

func (svc *ResourceGroupsVersionService) List() (response ResourceGroupsResponse, err error) {
	var rawResponse resourceGroupsWorkaroundResponse
	err = svc.client.RequestDecoder("GET", apiV2ResourceGroups, nil, &rawResponse)

	if err != nil {
		return
	}

	return setResourceGroupsVersionUnawareResponse(rawResponse)
}

func castRGV1WorkAroundResponse(data resourceGroupWorkaroundResponse, response interface{}) error {
	isDefault, err := strconv.Atoi(data.Data.IsDefault)
	if err != nil {
		return err
	}
	group := ResourceGroupResponse{
		Data: ResourceGroupData{
			Guid:         data.Data.Guid,
			IsDefault:    isDefault,
			ResourceGuid: data.Data.ResourceGuid,
			Name:         data.Data.Name,
			Type:         data.Data.Type,
			Enabled:      data.Data.Enabled,
			Props:        data.Data.Props,
		},
	}

	j, err := json.Marshal(group)
	if err != nil {
		return err
	}

	err = json.Unmarshal(j, &response)
	if err != nil {
		return err
	}

	return nil
}

func castRGV2WorkAroundResponse(data resourceGroupWorkaroundResponse, response interface{}) error {
	group := ResourceGroupResponse{
		Data: ResourceGroupData{
			Type:              data.Data.Type,
			Enabled:           data.Data.Enabled,
			NameV2:            data.Data.NameV2,
			Query:             data.Data.Query,
			Description:       data.Data.Description,
			ResourceGroupGuid: data.Data.ResourceGroupGuid,
			CreatedTime:       data.Data.CreatedTime,
			CreatedBy:         data.Data.CreatedBy,
			UpdatedTime:       data.Data.UpdatedTime,
			UpdatedBy:         data.Data.UpdatedBy,
			IsDefaultBoolean:  data.Data.IsDefaultBoolean,
			IsOrg:             data.Data.IsOrg,
		},
	}

	j, err := json.Marshal(group)
	if err != nil {
		return err
	}

	err = json.Unmarshal(j, &response)
	if err != nil {
		return err
	}

	return nil
}

func castResourceGroupV2Response(data ResourceGroupV2Response, response interface{}) error {
	group := ResourceGroupResponse{
		Data: ResourceGroupData{
			Type:              data.Data.Type,
			Enabled:           data.Data.Enabled,
			NameV2:            data.Data.Name,
			Query:             data.Data.Query,
			Description:       data.Data.Description,
			ResourceGroupGuid: data.Data.ResourceGroupGuid,
			CreatedTime:       data.Data.CreatedTime,
			CreatedBy:         data.Data.CreatedBy,
			UpdatedTime:       data.Data.UpdatedTime,
			UpdatedBy:         data.Data.UpdatedBy,
			IsDefaultBoolean:  data.Data.IsDefaultBoolean,
			IsOrg:             data.Data.IsOrg,
		},
	}

	j, err := json.Marshal(group)
	if err != nil {
		return err
	}

	err = json.Unmarshal(j, &response)
	if err != nil {
		return err
	}

	return nil
}

func (svc *ResourceGroupsVersionService) get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify an resourceGuid")
	}
	apiPath := fmt.Sprintf(apiV2ResourceGroupsFromGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, response)
}

func (svc *ResourceGroupsVersionService) create(data interface{}, response interface{}) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2ResourceGroups, data, response)
}

func (svc *ResourceGroupsVersionService) update(guid string, data interface{}, response interface{}) error {
	if guid == "" {
		return errors.New("specify a resource group guid")
	}

	apiPath := fmt.Sprintf(apiV2ResourceGroupsFromGUID, guid)
	return svc.client.RequestEncoderDecoder("PATCH", apiPath, data, response)
}

func setResourceGroupResponse(response resourceGroupWorkaroundData) (ResourceGroupResponse,
	error) {

	if response.Props != nil {
		isDefault, err := strconv.Atoi(response.IsDefault)
		if err != nil {
			return ResourceGroupResponse{}, err
		}
		return ResourceGroupResponse{
			Data: ResourceGroupData{
				Guid:         response.Guid,
				IsDefault:    isDefault,
				ResourceGuid: response.ResourceGuid,
				Name:         response.Name,
				Type:         response.Type,
				Enabled:      response.Enabled,
				Props:        response.Props,
			},
		}, nil
	} else {
		return ResourceGroupResponse{
			Data: ResourceGroupData{
				Type:              response.Type,
				Enabled:           response.Enabled,
				NameV2:            response.NameV2,
				Query:             response.Query,
				Description:       response.Description,
				ResourceGroupGuid: response.ResourceGroupGuid,
				CreatedTime:       response.CreatedTime,
				CreatedBy:         response.CreatedBy,
				UpdatedTime:       response.UpdatedTime,
				UpdatedBy:         response.UpdatedBy,
				IsDefaultBoolean:  response.IsDefaultBoolean,
				IsOrg:             response.IsOrg,
			},
		}, nil
	}
}

func setResourceGroupsVersionUnawareResponse(workaround resourceGroupsWorkaroundResponse) (
	ResourceGroupsResponse, error) {
	var data []ResourceGroupData
	for _, r := range workaround.Data {
		group, err := setResourceGroupResponse(r)
		if err != nil {
			return ResourceGroupsResponse{}, err
		}
		data = append(data, group.Data)
	}

	return ResourceGroupsResponse{Data: data}, nil
}
