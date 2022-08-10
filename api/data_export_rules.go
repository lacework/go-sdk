//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

// DataExportRulesService is a service that interacts with the DataExportRule
// endpoints from the Lacework Server
type DataExportRulesService struct {
	client *Client
}

type DataExportRulesResponse struct {
	Data    []DataExportRule `json:"data"`
	Message string           `json:"message"`
}

type DataExportRuleResponse struct {
	Data    DataExportRule `json:"data"`
	Message string         `json:"message"`
}

type DataExportRule struct {
	ID     string               `json:"mcGuid,omitempty"`
	Filter DataExportRuleFilter `json:"filters"`
	Type   string               `json:"type"`
	IDs    []string             `json:"intgGuidList"`
}

type DataExportRuleFilter struct {
	Name            string   `json:"name"`
	CreatedBy       string   `json:"createdOrUpdatedBy,omitempty"`
	UpdatedTime     string   `json:"createdOrUpdatedTime,omitempty"`
	Enabled         int      `json:"enabled"`
	ProfileVersions []string `json:"profileVersions"`
}

// List returns a list of Data Export Rules
func (svc *DataExportRulesService) List() (
	response DataExportRulesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiV2DataExportRules, nil, &response)
	return
}

// Get returns a raw response of the Data Export Rule with the matching guid.
func (svc *DataExportRulesService) Get(id string) (
	response DataExportRuleResponse,
	err error,
) {
	if id == "" {
		err = errors.New("data export rule ID must be provided")
		return
	}
	apiPath := fmt.Sprintf(apiV2DataExportRulesFromGUID, id)

	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

// Create creates a single Data Export Rule
func (svc *DataExportRulesService) Create(rule DataExportRule) (response DataExportRuleResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2DataExportRules, rule, &response)
	return
}

// Delete deletes a Data Export Rule that matches the provided guid
func (svc *DataExportRulesService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2DataExportRulesFromGUID, guid),
		nil,
		nil,
	)
}

// Search returns a list of Data Export Rules
func (svc *DataExportRulesService) Search(filters SearchFilter) (
	response DataExportRulesResponse, err error,
) {
	err = svc.client.RequestEncoderDecoder(
		"POST", apiV2DataExportRulesSearch,
		filters, &response,
	)
	return
}
