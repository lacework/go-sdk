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
)

// IntegrationsService is a service that interacts with the integrations
// endpoints from the Lacework Server
type IntegrationsService struct {
	client *Client
}

type integrationType int

const (
	// type that defines a non-existing integration
	NoneIntegration integrationType = iota

	// AWS Config integration type
	AwsCfgIntegration

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

// IntegrationTypes is the list of available integration types
var IntegrationTypes = map[integrationType]string{
	NoneIntegration:             "NONE",
	AwsCfgIntegration:           "AWS_CFG",
	AwsCloudTrailIntegration:    "AWS_CT_SQS",
	GcpCfgIntegration:           "GCP_CFG",
	GcpAuditLogIntegration:      "GCP_AT_SES",
	AzureCfgIntegration:         "AZURE_CFG",
	AzureActivityLogIntegration: "AZURE_AL_SEQ",
}

// String returns the string representation of an integration type
func (i integrationType) String() string {
	return IntegrationTypes[i]
}

// FindIntegrationType looks up inside the list of available integration types
// the matching type from the provided string, if none, returns NoneIntegration
func FindIntegrationType(t string) (integrationType, bool) {
	for iType, str := range IntegrationTypes {
		if str == t {
			return iType, true
		}
	}
	return NoneIntegration, false
}

// Get gets a single integration matching the integration guid on the Lacework Server,
// the returned integration contains the 'Data' field raw (map of interfaces)
func (svc *IntegrationsService) Get(guid string) (
	response RawIntegrationsResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// Delete deletes a single integration matching the integration guid on the Lacework Server
// the returned integration contains the 'Data' field raw (map of interfaces)
func (svc *IntegrationsService) Delete(guid string) (
	response RawIntegrationsResponse,
	err error,
) {
	err = svc.delete(guid, &response)
	return
}

// List lists the external integrations available on the Lacework Server
func (svc *IntegrationsService) List() (response RawIntegrationsResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiIntegrations, nil, &response)
	return
}

// ListByType lists the external integrations from the provided type that are available
// on the Lacework Server
func (svc *IntegrationsService) ListByType(iType integrationType) (response RawIntegrationsResponse, err error) {
	err = svc.listByType(iType, &response)
	return
}

// GetSchema get the integration schema for the provided integration type
func (svc *IntegrationsService) GetSchema(iType integrationType) (
	response map[string]interface{},
	err error,
) {
	apiPath := fmt.Sprintf(apiIntegrationSchema, iType.String())
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *IntegrationsService) get(guid string, response interface{}) error {
	apiPath := fmt.Sprintf(apiIntegrationByGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, response)
}

func (svc *IntegrationsService) create(data interface{}, response interface{}) error {
	body, err := jsonReader(data)
	if err != nil {
		return err
	}

	err = svc.client.RequestDecoder("POST", apiIntegrations, body, response)
	return err
}

func (svc *IntegrationsService) update(guid string, data interface{}, response interface{}) error {
	var (
		apiPath   = fmt.Sprintf(apiIntegrationByGUID, guid)
		body, err = jsonReader(data)
	)
	if err != nil {
		return err
	}

	return svc.client.RequestDecoder("PATCH", apiPath, body, response)
}

func (svc *IntegrationsService) delete(guid string, response interface{}) error {
	apiPath := fmt.Sprintf(apiIntegrationByGUID, guid)
	return svc.client.RequestDecoder("DELETE", apiPath, nil, response)
}

func (svc *IntegrationsService) listByType(iType integrationType, response interface{}) error {
	apiPath := fmt.Sprintf(apiIntegrationsByType, iType.String())
	return svc.client.RequestDecoder("GET", apiPath, nil, &response)
}

type commonIntegrationData struct {
	IntgGuid             string            `json:"INTG_GUID,omitempty"`
	Name                 string            `json:"NAME"`
	CreatedOrUpdatedTime string            `json:"CREATED_OR_UPDATED_TIME,omitempty"`
	CreatedOrUpdatedBy   string            `json:"CREATED_OR_UPDATED_BY,omitempty"`
	Type                 string            `json:"TYPE"`
	Enabled              int               `json:"ENABLED"`
	State                *IntegrationState `json:"STATE,omitempty"`
	IsOrg                int               `json:"IS_ORG,omitempty"`
	TypeName             string            `json:"TYPE_NAME,omitempty"`
}

func (c commonIntegrationData) Status() string {
	if c.Enabled == 1 {
		return "Enabled"
	}
	return "Disabled"
}

func (c commonIntegrationData) StateString() string {
	if c.State != nil && c.State.Ok {
		return "Ok"
	}
	return "Check"
}

type IntegrationState struct {
	Ok                 bool   `json:"ok"`
	LastUpdatedTime    string `json:"lastUpdatedTime"`
	LastSuccessfulTime string `json:"lastSuccessfulTime"`
}

type RawIntegration struct {
	commonIntegrationData
	Data map[string]interface{} `json:"DATA"`
}

type RawIntegrationsResponse struct {
	Data    []RawIntegration `json:"data"`
	Ok      bool             `json:"ok"`
	Message string           `json:"message"`
}
