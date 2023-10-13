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
	"fmt"
	"time"

	"github.com/pkg/errors"
)

func (svc *ResourceGroupsV2Service) List() (response ResourceGroupsV2Response, err error) {
	var rawResponse ResourceGroupsV2Response
	err = svc.client.RequestDecoder("GET", apiV2ResourceGroups, nil, &rawResponse)
	if err != nil {
		return
	}

	return
}

func (svc *ResourceGroupsV2Service) Create(group ResourceGroupDataWithQuery) (
	response ResourceGroupV2Response,
	err error,
) {
	err = svc.create(group, &response)
	return
}

func (svc *ResourceGroupsV2Service) Update(data ResourceGroup) (
	response ResourceGroupV2Response,
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

func (svc *ResourceGroupsV2Service) Delete(guid string) error {
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

func (svc *ResourceGroupsV2Service) Get(guid string, response interface{}) error {
	var rawResponse resourceGroupWorkaroundResponse
	err := svc.get(guid, &rawResponse)
	if err != nil {
		return err
	}

	return castResourceGroupResponse(rawResponse.Data, &response)
}

func (svc *ResourceGroupsV2Service) create(data interface{}, response interface{}) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2ResourceGroups, data, response)
}

func (svc *ResourceGroupsV2Service) get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify an resourceGuid")
	}
	apiPath := fmt.Sprintf(apiV2ResourceGroupsFromGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, response)
}

func (svc *ResourceGroupsV2Service) update(guid string, data interface{}, response interface{}) error {
	if guid == "" {
		return errors.New("specify a resource group guid")
	}

	apiPath := fmt.Sprintf(apiV2ResourceGroupsFromGUID, guid)
	return svc.client.RequestEncoderDecoder("PATCH", apiPath, data, response)
}

type ResourceGroupV2Response struct {
	Data ResourceGroupDataWithQuery `json:"data"`
}

type ResourceGroupsV2Response struct {
	Data []ResourceGroupDataWithQuery `json:"data"`
}

type ResourceGroupsV2Service struct {
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
type ResourceGroupDataWithQuery struct {
	Name              string     `json:"name"`
	Type              string     `json:"resourceType"`
	Query             *RGQuery   `json:"query"`
	Description       string     `json:"description,omitempty"`
	ResourceGroupGuid string     `json:"resourceGroupGuid,omitempty"`
	CreatedTime       *time.Time `json:"lastUpdated,omitempty"`
	CreatedBy         string     `json:"createdBy,omitempty"`
	UpdatedTime       *time.Time `json:"updatedTime,omitempty"`
	UpdatedBy         string     `json:"updatedBy,omitempty"`
	Enabled           int        `json:"enabled,omitempty"`
	IsDefaultBoolean  *bool      `json:"isDefaultBoolean,omitempty"`
	IsOrg             *bool      `json:"isOrg,omitempty"`
}

func (group ResourceGroupDataWithQuery) GetProps() interface{} {
	return nil
}

func (group ResourceGroupDataWithQuery) GetQuery() *RGQuery {
	return group.Query
}

func (group ResourceGroupDataWithQuery) ResourceGroupType() resourceGroupType {
	t, _ := FindResourceGroupType(group.Type)
	return t
}

func (group ResourceGroupDataWithQuery) ID() string {
	return group.ResourceGroupGuid
}

func (group *ResourceGroupDataWithQuery) ResetRGV2Fields() {
	// no-op
}

func (group *ResourceGroupDataWithQuery) ResetResourceGUID() {
	group.ResourceGroupGuid = ""
}

func (group ResourceGroupDataWithQuery) IsV2Group() bool {
	return true
}
