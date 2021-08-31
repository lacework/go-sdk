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
	"fmt"

	"github.com/pkg/errors"
)

// ResourceGroupsService is the service that interacts with
// the ResourceGroups schema from the Lacework APIv2 Server
type ResourceGroupsService struct {
	client *Client
}

type ResourceGroup interface {
	ID() string
	ResourceGroupType() ResourceGroupType
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

// NewResourceGroup returns an instance of the ResourceGroupRaw struct with the
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
//			ContainerLabels: []map[string]string{{"*": "*"}},
//			ContainerTags:   []string{"*"},
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
	err = svc.client.RequestDecoder("GET", apiV2ResourceGroups, nil, &response)
	return
}

// Create creates a single Resource Group
func (svc *ResourceGroupsService) Create(integration ResourceGroupData) (
	response ResourceGroupsResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
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
	return svc.get(guid, &response)
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
		return errors.New("specify an intgGuid")
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

type ResourceGroupResponse struct {
	Data ResourceGroupData `json:"data"`
}

type ResourceGroupsResponse struct {
	Data []ResourceGroupData `json:"data"`
}

type ResourceGroupData struct {
	Guid         string      `json:"guid,omitempty"`
	IsDefault    string      `json:"isDefault,omitempty"`
	ResourceGuid string      `json:"resourceGuid,omitempty"`
	Name         string      `json:"resourceName"`
	Type         string      `json:"resourceType"`
	Enabled      int         `json:"enabled,omitempty"`
	Props        interface{} `json:"props,omitempty"`
}
