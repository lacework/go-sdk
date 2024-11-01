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
	"strings"
	"time"

	"github.com/pkg/errors"
)

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
	MachineResourceGroupQueryTemplate string
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
	description string, query *RGQuery) ResourceGroupData {
	return ResourceGroupData{
		Name:        name,
		Type:        iType.String(),
		Enabled:     1,
		Query:       query,
		Description: description,
	}
}

func (svc *ResourceGroupsService) List() (response ResourceGroupsResponse, err error) {
	var rawResponse ResourceGroupsResponse
	err = svc.client.RequestDecoder("GET", apiV2ResourceGroups, nil, &rawResponse)
	if err != nil {
		return rawResponse, err
	}

	err = sanitizeFieldsInRawResponseList(&rawResponse, &response)
	if err != nil {
		return rawResponse, err
	}

	return rawResponse, nil
}

func sanitizeFieldsInRawResponse(rawResponse *ResourceGroupResponse, response interface{}) error {
	// update filters keys to match the query template
	updateFiltersKeys(&rawResponse.Data)

	j, err := json.Marshal(rawResponse)
	if err != nil {
		return err
	}

	return json.Unmarshal(j, &response)
}

func sanitizeFieldsInRawResponseList(rawResponse *ResourceGroupsResponse, response interface{}) error {
	for i := range rawResponse.Data {
		// update filters keys to match the query template
		updateFiltersKeys(&rawResponse.Data[i])
	}

	j, err := json.Marshal(rawResponse)
	if err != nil {
		return err
	}

	return json.Unmarshal(j, &response)
}

func (svc *ResourceGroupsService) Create(group ResourceGroupData) (
	response ResourceGroupResponse,
	err error,
) {
	var rawResponse ResourceGroupResponse
	err = svc.create(group, &rawResponse)
	if err != nil {
		return
	}

	err = sanitizeFieldsInRawResponse(&rawResponse, &response)

	return
}

func (svc *ResourceGroupsService) Update(data *ResourceGroupData) (
	response ResourceGroupResponse,
	err error,
) {
	if data == nil {
		err = errors.New("resource group must not be empty")
		return
	}
	guid := data.ID()
	data.ResetResourceGUID()

	var rawResponse ResourceGroupResponse
	err = svc.update(guid, data, &rawResponse)

	if err != nil {
		return
	}

	err = sanitizeFieldsInRawResponse(&rawResponse, &response)

	return
}

func collectFilterNames(children []*RGChild, filterNames map[string]string) {
	for _, child := range children {
		if child.FilterName != "" {
			normalizedKey := strings.ReplaceAll(strings.ToLower(child.FilterName), "_", "")
			filterNames[normalizedKey] = child.FilterName
		}
		if len(child.Children) > 0 {
			collectFilterNames(child.Children, filterNames)
		}
	}
}

/*
updateFiltersKeys updates the keys in the Filters map of ResourceGroupData to ensure they match the filter names
defined in the nested children of the query expression. This is necessary because JSON decoding/encoding can
convert keys to camel case, causing mismatches. The function normalizes the keys by removing underscores and
converting them to lower case, then compares them with the filter names. If a mismatch is found, the key is
updated to the value in RGExpression.Children
*/
func updateFiltersKeys(data *ResourceGroupData) {
	filterNames := make(map[string]string)
	collectFilterNames(data.Query.Expression.Children, filterNames)

	updatedFilters := make(map[string]*RGFilter)
	for key, value := range data.Query.Filters {
		normalizedKey := strings.ReplaceAll(strings.ToLower(key), "_", "")
		if _, exists := filterNames[normalizedKey]; exists {
			updatedFilters[filterNames[normalizedKey]] = value
		} else {
			updatedFilters[key] = value
		}
	}

	data.Query.Filters = updatedFilters
}

func (group *ResourceGroupData) ResetResourceGUID() {
	group.ResourceGroupGuid = ""
	group.UpdatedBy = ""
	group.UpdatedTime = nil
	group.CreatedBy = ""
	group.CreatedTime = nil
	group.IsDefaultBoolean = nil
}

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

func (svc *ResourceGroupsService) Get(guid string, response interface{}) error {
	var rawResponse ResourceGroupResponse

	err := svc.get(guid, &rawResponse)
	if err != nil {
		return err
	}

	err = sanitizeFieldsInRawResponse(&rawResponse, response)
	if err != nil {
		return err
	}

	return nil
}

func (svc *ResourceGroupsService) create(data interface{}, response interface{}) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2ResourceGroups, data, response)
}

func (svc *ResourceGroupsService) get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify an resource group guid")
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

type ResourceGroupsService struct {
	client *Client
}

type RGExpression struct {
	Operator string     `json:"operator"`
	Children []*RGChild `json:"children"`
}

type RGChild struct {
	Operator   string     `json:"operator,omitempty"`
	FilterName string     `json:"filterName,omitempty"`
	Children   []*RGChild `json:"children,omitempty"`
}

type RGFilter struct {
	Field     string   `json:"field"`
	Operation string   `json:"operation"`
	Values    []string `json:"values"`
	Key       string   `json:"key,omitempty"`
}

type RGQuery struct {
	Filters    map[string]*RGFilter `json:"filters"`
	Expression *RGExpression        `json:"expression"`
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

func (group *ResourceGroupData) ID() string {
	return group.ResourceGroupGuid
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
	CreatedTime       *time.Time `json:"createdTime,omitempty"`
	CreatedBy         string     `json:"createdBy,omitempty"`
	UpdatedTime       *time.Time `json:"updatedTime,omitempty"`
	UpdatedBy         string     `json:"updatedBy,omitempty"`
	IsDefaultBoolean  *bool      `json:"isDefaultBoolean,omitempty"`
	Type              string     `json:"resourceType"`
	Enabled           int        `json:"enabled"`
}
