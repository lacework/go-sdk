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
	// awsCFG - AWS Config integration type
	awsCFG integrationType = iota

	// awsCT - AWS CloudTrail integration type
	awsCT

	// gcpCFG - GCP Config integration type
	gcpCFG

	// gcpAT - GCP Audit Log integration type
	gcpAT

	// azureCFG - Azure Config integration type
	azureCFG

	// azureAL - Azure Activity Log integration type
	azureAL
)

var integrationTypes = map[integrationType]string{
	awsCFG:   "AWS_CFG",
	awsCT:    "AWS_CT_SQS",
	gcpCFG:   "GCP_CFG",
	gcpAT:    "GCP_AT_SES",
	azureCFG: "AZURE_CFG",
	azureAL:  "AZURE_AL_SEQ",
}

func (i integrationType) String() string {
	return integrationTypes[i]
}

// gcpResourceLevel determines Project or Organization level integration
type gcpResourceLevel int

const (
	// GcpProject level integration with GCP
	GcpProject gcpResourceLevel = iota

	// GcpOrganization level integration with GCP
	GcpOrganization
)

var gcpResourceLevels = map[gcpResourceLevel]string{
	GcpProject:      "PROJECT",
	GcpOrganization: "ORGANIZATION",
}

func (g gcpResourceLevel) String() string {
	return gcpResourceLevels[g]
}

// GetIntegrations lists the external integrations available on the server
func (c *Client) GetIntegrations() (response integrationsResponse, err error) {
	err = c.RequestDecoder("GET", apiIntegrations, nil, &response)
	return
}

func (c *Client) GetGCPIntegrations() (response gcpIntegrationsResponse, err error) {
	return
}
func (c *Client) GetAzureIntegrations() (response azureIntegrationsResponse, err error) {
	return
}
func (c *Client) GetAWSIntegrations() (response awsIntegrationsResponse, err error) {
	return
}

// NewGCPIntegrationData returns an instance of gcpIntegrationData
func NewGCPIntegrationData(name string, idType gcpResourceLevel) gcpIntegrationData {
	return gcpIntegrationData{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    gcpCFG.String(),
			Enabled: 1,
		},
		Data: gcpCfg{
			IdType: idType.String(),
		},
	}
}

// CreateGCPConfigIntegration creates a single integration on the server
func (c *Client) CreateGCPConfigIntegration(data gcpIntegrationData) (response gcpIntegrationsResponse, err error) {
	err = c.createIntegration(data, &response)
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

// GetGCPConfigIntegration gets a single integration matching the integration guid available on the server
func (c *Client) GetGCPConfigIntegration(intgGuid string) (response gcpIntegrationsResponse, err error) {
	err = c.getIntegration(intgGuid, &response)
	return
}

func (c *Client) getIntegration(intgGuid string, response interface{}) error {
	apiPath := fmt.Sprintf(apiIntegrationByGUID, intgGuid)
	return c.RequestDecoder("GET", apiPath, nil, response)
}

type commonIntegrationData struct {
	IntgGuid             string `json:"INTG_GUID"`
	Name                 string `json:"NAME"`
	CreatedOrUpdatedTime string `json:"CREATED_OR_UPDATED_TIME"`
	CreatedOrUpdatedBy   string `json:"CREATED_OR_UPDATED_BY"`
	Type                 string `json:"TYPE"`
	Enabled              int    `json:"ENABLED"`
	State                state  `json:"STATE"`
	IsOrg                int    `json:"IS_ORG"`
	TypeName             string `json:"TYPE_NAME"`
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

type state struct {
	Ok                 bool   `json:"ok"`
	LastUpdatedTime    string `json:"lastUpdatedTime"`
	LastSuccessfulTime string `json:"lastSuccessfulTime"`
}

type awsIntegrationsResponse struct {
	//Data    []gcpIntegrationData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}
type azureIntegrationsResponse struct {
	//Data    []gcpIntegrationData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type gcpIntegrationsResponse struct {
	Data    []gcpIntegrationData `json:"data"`
	Ok      bool                 `json:"ok"`
	Message string               `json:"message"`
}

type gcpIntegrationData struct {
	commonIntegrationData
	Data gcpCfg `json:"DATA"`
}

type gcpCfg struct {
	ID          string         `json:"ID"`
	IdType      string         `json:"ID_TYPE"`
	Credentials gcpCredentials `json:"CREDENTIALS"`
}

type gcpCredentials struct {
	ClientId     string `json:"CLIENT_ID"`
	ClientEmail  string `json:"CLIENT_EMAIL"`
	PrivateKeyId string `json:"PRIVATE_KEY_ID"`
	PrivateKey   string `json:"PRIVATE_KEY"`
}
