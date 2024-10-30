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
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// ResourceGroupsService is the service that interacts with
// the ResourceGroups schema from the Lacework APIv2 Server
type ResourceGroupsService struct {
	client *Client
}

type ResourceGroup interface {
	ID() string
	ResourceGroupType() resourceGroupType
	ResetResourceGUID()
	ResetRGV2Fields()
	IsV2Group() bool
}

type resourceGroupType int

const (
	// type that defines a non-existing Resource Group
	NoneResourceGroup resourceGroupType = iota
	AwsResourceGroup
	AzureResourceGroup
	ContainerResourceGroup
	GcpResourceGroup
	MachineResourceGroup
	OciResourceGroup
	KubernetesResourceGroup
)

// query templates
var (
	NoneResourceGroupQueryTemplate string = ""
	//go:embed _templates/resource_groups/aws.json
	AwsResourceGroupQueryTemplate string
	//go:embed _templates/resource_groups/azure.json
	AzureResourceGroupQueryTemplate string
	//go:embed _templates/resource_groups/container.json
	ContainerResourceGroupQueryTemplate string
	//go:embed _templates/resource_groups/gcp.json
	GcpResourceGroupQueryTemplate string
	//go:embed _templates/resource_groups/machine.json
	MachineResourceGroupQueryTemplate   string
	LwAccountResourceGroupQueryTemplate string = ""
	//go:embed _templates/resource_groups/oci.json
	OciResourceGroupQueryTemplate string
	//go:embed _templates/resource_groups/kubernetes.json
	KubernetesResourceGroupQueryTemplate string
)

type resourceGroupContext struct {
	resourceGroupType string
	queryTemplate     string
}

// ResourceGroupTypes is the list of available Resource Group types
var ResourceGroupTypes = map[resourceGroupType]resourceGroupContext{
	NoneResourceGroup:       {resourceGroupType: "None", queryTemplate: NoneResourceGroupQueryTemplate},
	AwsResourceGroup:        {resourceGroupType: "AWS", queryTemplate: AwsResourceGroupQueryTemplate},
	AzureResourceGroup:      {resourceGroupType: "AZURE", queryTemplate: AzureResourceGroupQueryTemplate},
	ContainerResourceGroup:  {resourceGroupType: "CONTAINER", queryTemplate: ContainerResourceGroupQueryTemplate},
	GcpResourceGroup:        {resourceGroupType: "GCP", queryTemplate: GcpResourceGroupQueryTemplate},
	MachineResourceGroup:    {resourceGroupType: "MACHINE", queryTemplate: MachineResourceGroupQueryTemplate},
	OciResourceGroup:        {resourceGroupType: "OCI", queryTemplate: OciResourceGroupQueryTemplate},
	KubernetesResourceGroup: {resourceGroupType: "KUBERNETES", queryTemplate: KubernetesResourceGroupQueryTemplate},
}

func NewResourceGroup(name string, iType resourceGroupType,
	description string, query *RGQuery) ResourceGroupDataWithQuery {
	return ResourceGroupDataWithQuery{
		Name:        name,
		Type:        iType.String(),
		Enabled:     1,
		Query:       query,
		Description: description,
	}
}

// String returns the string representation of a Resource Group type
func (i resourceGroupType) String() string {
	return ResourceGroupTypes[i].resourceGroupType
}

// QueryTemplate returns the resource group type's query template
func (i resourceGroupType) QueryTemplate() string {
	return ResourceGroupTypes[i].queryTemplate
}

// FindResourceGroupType looks up inside the list of available resource group types
// the matching type from the provided string, if none, returns NoneResourceGroup
func FindResourceGroupType(typ string) (resourceGroupType, bool) {
	for i, ctx := range ResourceGroupTypes {
		if typ == ctx.resourceGroupType {
			return i, true
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
	group := ResourceGroupResponse{
		Data: ResourceGroupData{
			IsDefaultBoolean:  data.IsDefaultBoolean,
			ResourceGroupGuid: data.ResourceGroupGuid,
			Name:              data.Name,
			Type:              data.Type,
			Enabled:           data.Enabled,
			Query:             data.Query,
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

func setResourceGroupResponse(response resourceGroupWorkaroundData) (ResourceGroupResponse,
	error) {
	return ResourceGroupResponse{
		Data: ResourceGroupData(response),
	}, nil
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
//	Get<Type>(guid)
//
//	  Where <Type> is the Resource Group type.
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

func (group ResourceGroupData) ResourceGroupType() resourceGroupType {
	t, _ := FindResourceGroupType(group.Type)
	return t
}

func (group ResourceGroupData) ID() string {
	return group.ResourceGroupGuid
}

func (group *ResourceGroupData) ResetResourceGUID() {
	group.ResourceGroupGuid = ""
	group.UpdatedBy = ""
	group.UpdatedTime = nil
	group.CreatedBy = ""
	group.CreatedTime = nil
	group.IsDefaultBoolean = nil
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
	Name              string     `json:"name,omitempty"`
	Query             *RGQuery   `json:"query,omitempty"`
	Description       string     `json:"description,omitempty"`
	ResourceGroupGuid string     `json:"resourceGroupGuid,omitempty"`
	CreatedTime       *time.Time `json:"lastUpdated,omitempty"`
	CreatedBy         string     `json:"createdBy,omitempty"`
	UpdatedTime       *time.Time `json:"updatedTime,omitempty"`
	UpdatedBy         string     `json:"updatedBy,omitempty"`
	IsDefaultBoolean  *bool      `json:"isDefaultBoolean,omitempty"`
	Type              string     `json:"resourceType"`
	Enabled           int        `json:"enabled"`
}

// RAIN-21510 workaround
type resourceGroupWorkaroundResponse struct {
	Data resourceGroupWorkaroundData `json:"data"`
}

type resourceGroupsWorkaroundResponse struct {
	Data []resourceGroupWorkaroundData `json:"data"`
}

type resourceGroupWorkaroundData struct {
	Name              string     `json:"name,omitempty"`
	Query             *RGQuery   `json:"query,omitempty"`
	Description       string     `json:"description,omitempty"`
	ResourceGroupGuid string     `json:"resourceGroupGuid,omitempty"`
	CreatedTime       *time.Time `json:"lastUpdated,omitempty"`
	CreatedBy         string     `json:"createdBy,omitempty"`
	UpdatedTime       *time.Time `json:"updatedTime,omitempty"`
	UpdatedBy         string     `json:"updatedBy,omitempty"`
	IsDefaultBoolean  *bool      `json:"isDefaultBoolean,omitempty"`
	Type              string     `json:"resourceType"`
	Enabled           int        `json:"enabled,omitempty"`
}
