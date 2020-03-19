//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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
	"strings"
)

type integrationType int

const (
	// AWS Config integration type
	AwsCfgIntegration integrationType = iota

	// AWS CloudTrail integration type
	AwsCloudTrailIntegration

	// GCP Config integration type
	GcpCfgIntegration

	// GCP Audit Log integration type
	GcpAuditLogIntegration

	// Azure Config integration type
	AzureCfgIntegration

	// Azure Activity Log integration type
	AzureActivityLogIntegration
)

var integrationTypes = map[integrationType]string{
	AwsCfgIntegration:           "AWS_CFG",
	AwsCloudTrailIntegration:    "AWS_CT_SQS",
	GcpCfgIntegration:           "GCP_CFG",
	GcpAuditLogIntegration:      "GCP_AT_SES",
	AzureCfgIntegration:         "AZURE_CFG",
	AzureActivityLogIntegration: "AZURE_AL_SEQ",
}

func (i integrationType) String() string {
	return integrationTypes[i]
}

// GetIntegrations lists the external integrations available on the server
func (c *Client) GetIntegrations() (response integrationsResponse, err error) {
	err = c.RequestDecoder("GET", apiIntegrations, nil, &response)
	return
}

// GetIntegrationSchema get the integration schema for the provided integration type
func (c *Client) GetIntegrationSchema(iType integrationType) (response map[string]interface{}, err error) {
	apiPath := fmt.Sprintf(apiIntegrationSchema, iType.String())
	err = c.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (c *Client) createIntegration(data interface{}, response interface{}) error {
	body, err := jsonReader(data)
	if err != nil {
		return err
	}

	err = c.RequestDecoder("POST", apiIntegrations, body, response)
	return err
}

func (c *Client) getIntegration(intgGuid string, response interface{}) error {
	apiPath := fmt.Sprintf(apiIntegrationByGUID, intgGuid)
	return c.RequestDecoder("GET", apiPath, nil, response)
}

func (c *Client) updateIntegration(intgGuid string, data interface{}, response interface{}) error {
	body, err := jsonReader(data)
	if err != nil {
		return err
	}

	apiPath := fmt.Sprintf(apiIntegrationByGUID, intgGuid)
	err = c.RequestDecoder("PATCH", apiPath, body, response)
	return err
}

func (c *Client) deleteIntegration(intgGuid string, response interface{}) error {
	apiPath := fmt.Sprintf(apiIntegrationByGUID, intgGuid)
	return c.RequestDecoder("DELETE", apiPath, nil, response)
}

type commonIntegrationData struct {
	IntgGuid             string `json:"INTG_GUID,omitempty"`
	Name                 string `json:"NAME"`
	CreatedOrUpdatedTime string `json:"CREATED_OR_UPDATED_TIME,omitempty"`
	CreatedOrUpdatedBy   string `json:"CREATED_OR_UPDATED_BY,omitempty"`
	Type                 string `json:"TYPE"`
	Enabled              int    `json:"ENABLED"`
	State                state  `json:"STATE,omitempty"`
	IsOrg                int    `json:"IS_ORG,omitempty"`
	TypeName             string `json:"TYPE_NAME,omitempty"`
}

type state struct {
	Ok                 bool   `json:"ok"`
	LastUpdatedTime    string `json:"lastUpdatedTime"`
	LastSuccessfulTime string `json:"lastSuccessfulTime"`
}

type integrationsResponse struct {
	Data    []commonIntegrationData `json:"data"`
	Ok      bool                    `json:"ok"`
	Message string                  `json:"message"`
}

func (integrations *integrationsResponse) List() string {
	out := []string{}
	for _, integration := range integrations.Data {
		out = append(out, fmt.Sprintf("%s %s", integration.IntgGuid, integration.Type))
	}
	return strings.Join(out, "\n")
}
