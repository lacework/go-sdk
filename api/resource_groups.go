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
	"fmt"
	"strconv"

	"github.com/lacework/go-sdk/lwtime"
	"github.com/pkg/errors"
)

// ResourceGroupsService is the service that interacts with
// the ResourceGroups schema from the Lacework APIv2 Server
type ResourceGroupsService struct {
	client *Client
}

type ResourceGroupProps interface {
	GetBaseProps() ResourceGroupPropsBase
}

type ResourceGroupPropsBase struct {
	Description string       `json:"description"`
	UpdatedBy   string       `json:"updatedBy,omitempty"`
	LastUpdated *lwtime.Epoch `json:"lastUpdated,omitempty"`
}

type ResourceGroup interface {
	ID() string
	ResourceGroupType() ResourceGroupType
	ResetResourceGUID()
}

type ResourceGroupType int

const (
	// type that defines a non-existing Resource Group
	NoneResourceGroup ResourceGroupType = iota
	AwsResourceGroup
	AzureResourceGroup
	ContainerResourceGroup
	GcpResourceGroup
	MachineResourceGroup

	// requires Org Access account client.WithOrgAccess()
	LwAccountResourceGroup
)

// ResourceGroupTypes is the list of available Resource Group types
var ResourceGroupTypes = map[ResourceGroupType]string{
	NoneResourceGroup:      "None",
	AwsResourceGroup:       "AWS",
	AzureResourceGroup:     "AZURE",
	ContainerResourceGroup: "CONTAINER",
	GcpResourceGroup:       "GCP",
	LwAccountResourceGroup: "LW_ACCOUNT",
	MachineResourceGroup:   "MACHINE",
}

// String returns the string representation of a Resource Group type
func (i ResourceGroupType) String() string {
	return ResourceGroupTypes[i]
}

// NewResourceGroup returns an instance of the ResourceGroupData struct with the
// provided ResourceGroup type, name and the props field as an interface{}.
//
// NOTE: This function must be used by any ResourceGroup type.
//
// Basic usage: Initialize a new ContainerResourceGroup struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   group := api.NewResourceGroup("container resource group",
//     api.ContainerResourceGroup,
//     api.ContainerResourceGroupData{
//       Props: api.ContainerResourceGroupProps{
//			Description:     "all containers,
//			ContainerLabels: ContainerResourceGroupAllLabels,
//			ContainerTags:   ContainerResourceGroupAllTags,
//		},
//     },
//   )
//
//   client.V2.ResourceGroups.Create(group)
//
func NewResourceGroup(name string, iType ResourceGroupType, props interface{}) ResourceGroupData {
	return ResourceGroupData{
		Name:    name,
		Type:    iType.String(),
		Enabled: 1,
		Props:   props,
	}
}

// FindResourceGroupType looks up inside the list of available resource group types
// the matching type from the provided string, if none, returns NoneResourceGroup
func FindResourceGroupType(resourceGroup string) (ResourceGroupType, bool) {
	for resType, resStr := range ResourceGroupTypes {
		if resStr == resourceGroup {
			return resType, true
		}
	}
	return NoneResourceGroup, false
}

// List returns a list of Resource Groups
func (svc *ResourceGroupsService) List() (response ResourceGroupsResponse, err error) {
	var rawResponse resourceGroupsWorkaroundResponse
	err = svc.client.RequestDecoder("GET", apiV2ResourceGroups, nil, &rawResponse)
	if err != nil {
		return
	}
	response, err = setResourceGroupsResponse(rawResponse)
	if err != nil {
		return
	}

	return
}

// Create creates a single Resource Group
func (svc *ResourceGroupsService) Create(group ResourceGroupData) (
	response ResourceGroupResponse,
	err error,
) {
	err = svc.create(group, &response)
	return
}

// Update updates a single ResourceGroup on the Lacework Server
func (svc *ResourceGroupsService) Update(data ResourceGroup) (
	response ResourceGroupResponse,
	err error,
) {
	if data == nil {
		err = errors.New("resource group must not be empty")
		return
	}
	guid := data.ID()
	data.ResetResourceGUID()

	err = svc.update(guid, data, &response)
	if err != nil {
		return
	}

	return
}

func castResourceGroupResponse(data resourceGroupWorkaroundData, response interface{}) error {
	isDefault, err := strconv.Atoi(data.IsDefault)
	if err != nil {
		return err
	}
	group := ResourceGroupResponse{
		Data: ResourceGroupData{
			Guid:         data.Guid,
			IsDefault:    isDefault,
			ResourceGuid: data.ResourceGuid,
			Name:         data.Name,
			Type:         data.Type,
			Enabled:      data.Enabled,
			Props:        data.Props,
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

func setResourceGroupResponse(response resourceGroupWorkaroundData) (ResourceGroupResponse, error) {
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
}

func setResourceGroupsResponse(workaround resourceGroupsWorkaroundResponse) (ResourceGroupsResponse, error) {
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

// Delete deletes a Resource Group that matches the provided resource guid
func (svc *ResourceGroupsService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify a resourceGuid")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2ResourceGroupsFromGUID, guid),
		nil,
		nil,
	)
}

// Get returns a raw response of the Resource Group with the matching resource guid.
//
// To return a more specific Go struct of a Resource Group, use the proper
// method such as GetContainerResourceGroup() where the function name is composed by:
//
//  Get<Type>(guid)
//
//    Where <Type> is the Resource Group type.
func (svc *ResourceGroupsService) Get(guid string, response interface{}) error {
	var rawResponse resourceGroupWorkaroundResponse
	err := svc.get(guid, &rawResponse)
	if err != nil {
		return err
	}

	return castResourceGroupResponse(rawResponse.Data, &response)
}

func (svc *ResourceGroupsService) create(data interface{}, response interface{}) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2ResourceGroups, data, response)
}

func (svc *ResourceGroupsService) get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify an resourceGuid")
	}
	apiPath := fmt.Sprintf(apiV2ResourceGroupsFromGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, response)
}

func (svc *ResourceGroupsService) update(guid string, data interface{}, response interface{}) error {
	if guid == "" {
		return errors.New("specify a resource group guid")
	}

	apiPath := fmt.Sprintf(apiV2ResourceGroupsFromGUID, guid)
	return svc.client.RequestEncoderDecoder("PATCH", apiPath, data, response)
}

func (group ResourceGroupData) ResourceGroupType() ResourceGroupType {
	t, _ := FindResourceGroupType(group.Type)
	return t
}

func (group ResourceGroupData) ID() string {
	return group.ResourceGuid
}

func (group *ResourceGroupData) ResetResourceGUID() {
	group.ResourceGuid = ""
}

func (group ResourceGroupData) Status() string {
	if group.Enabled == 1 {
		return "Enabled"
	}
	return "Disabled"
}

type ResourceGroupResponse struct {
	Data ResourceGroupData `json:"data"`
}

type ResourceGroupsResponse struct {
	Data []ResourceGroupData `json:"data"`
}

type ResourceGroupData struct {
	Guid         string      `json:"guid,omitempty"`
	IsDefault    int         `json:"isDefault,omitempty"`
	ResourceGuid string      `json:"resourceGuid,omitempty"`
	Name         string      `json:"resourceName"`
	Type         string      `json:"resourceType"`
	Enabled      int         `json:"enabled,omitempty"`
	Props        interface{} `json:"props"`
}

// RAIN-21510 workaround
type resourceGroupWorkaroundResponse struct {
	Data resourceGroupWorkaroundData `json:"data"`
}

type resourceGroupsWorkaroundResponse struct {
	Data []resourceGroupWorkaroundData `json:"data"`
}

type resourceGroupWorkaroundData struct {
	Guid         string      `json:"guid,omitempty"`
	IsDefault    string      `json:"isDefault,omitempty"`
	ResourceGuid string      `json:"resourceGuid,omitempty"`
	Name         string      `json:"resourceName"`
	Type         string      `json:"resourceType"`
	Enabled      int         `json:"enabled,omitempty"`
	Props        interface{} `json:"props"`
}
