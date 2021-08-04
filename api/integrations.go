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

	// AWS Gov Cloud Config integration type
	AwsGovCloudCfgIntegration

	// AWS GovCloud CloudTrail integration type
	AwsGovCloudCTIntegration

	// AWS S3 channel integration type
	AwsS3ChannelIntegration

	// Datadog channel integration type
	DatadogChannelIntegration

	// GCP Config integration type
	GcpCfgIntegration

	// GCP Audit Log integration type
	GcpAuditLogIntegration

	// GCP Pub Sub alert channel integration type
	GcpPubSubChannelIntegration

	// New Relic Insights alert channel integration type
	NewRelicChannelIntegration

	// Azure Config integration type
	AzureCfgIntegration

	// Azure Activity Log integration type
	AzureActivityLogIntegration

	// Cisco Webex integration type
	CiscoWebexChannelIntegration

	// Container registry integration type
	ContainerRegistryIntegration

	// Microsoft Teams channel integration type
	MicrosoftTeamsChannelIntegration

	// QRadar channel integration type
	QRadarChannelIntegration

	// Slack channel integration type
	SlackChannelIntegration

	// Sevice Now alert channel integration type
	ServiceNowChannelIntegration

	// Splunk channel integration type
	SplunkIntegration

	// AWS CloudWatch integration type
	AwsCloudWatchIntegration

	// Pager Duty integration type
	PagerDutyIntegration

	// Jira integration type
	JiraIntegration

	// Email integration type
	EmailIntegration

	// VictorOps channel integration type
	VictorOpsChannelIntegration

	// Webhook channel integration type
	WebhookIntegration
)

type  integration struct {
	name string
	schema integrationSchema
}

// IntegrationTypes is the list of available integration types
var IntegrationTypes = map[integrationType]integration {
	NoneIntegration:                  {"NONE", None},
	AwsCfgIntegration:                {"AWS_CFG", CloudAccounts},
	AwsCloudTrailIntegration:         {"AWS_CT_SQS",CloudAccounts},
	AwsGovCloudCfgIntegration:        {"AWS_US_GOV_CFG",CloudAccounts},
	AwsGovCloudCTIntegration:         {"AWS_US_GOV_CT_SQS",CloudAccounts},
	AwsS3ChannelIntegration:          {"AWS_S3",AlertChannels},
	CiscoWebexChannelIntegration:     {"CISCO_SPARK_WEBHOOK",AlertChannels},
	DatadogChannelIntegration:        {"DATADOG",AlertChannels},
	GcpCfgIntegration:                {"GCP_CFG",CloudAccounts},
	GcpAuditLogIntegration:           {"GCP_AT_SES",CloudAccounts},
	GcpPubSubChannelIntegration:      {"GCP_PUBSUB",AlertChannels},
	NewRelicChannelIntegration:       {"NEW_RELIC_INSIGHTS",AlertChannels},
	AzureCfgIntegration:              {"AZURE_CFG",CloudAccounts},
	AzureActivityLogIntegration:      {"AZURE_AL_SEQ",CloudAccounts},
	ContainerRegistryIntegration:     {"CONT_VULN_CFG",ContainerRegistries},
	QRadarChannelIntegration:         {"IBM_QRADAR",AlertChannels},
	MicrosoftTeamsChannelIntegration: {"MICROSOFT_TEAMS",AlertChannels},
	SlackChannelIntegration:          {"SLACK_CHANNEL",AlertChannels},
	SplunkIntegration:                {"SPLUNK_HEC",AlertChannels},
	ServiceNowChannelIntegration:     {"SERVICE_NOW_REST",AlertChannels},
	AwsCloudWatchIntegration:         {"CLOUDWATCH_EB",AlertChannels},
	PagerDutyIntegration:             {"PAGER_DUTY_API",AlertChannels},
	JiraIntegration:                  {"JIRA",AlertChannels},
	EmailIntegration:                 {"EMAIL_USER",AlertChannels},
	VictorOpsChannelIntegration:      {"VICTOR_OPS",AlertChannels},
	WebhookIntegration:               {"WEBHOOK",AlertChannels},
}

// String returns the string representation of an integration type
func (i integrationType) String() string {
	return IntegrationTypes[i].name
}

// Schema returns the integration type
func (i integrationType) Schema() integrationSchema {
	return IntegrationTypes[i].schema
}

// FindIntegrationType looks up inside the list of available integration types
// the matching type from the provided string, if none, returns NoneIntegration
func FindIntegrationType(t string) (integrationType, bool) {
	for iType, str := range IntegrationTypes {
		if str.name == t {
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
	apiPath := fmt.Sprintf(apiIntegrationFromGUID, guid)
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
		apiPath   = fmt.Sprintf(apiIntegrationFromGUID, guid)
		body, err = jsonReader(data)
	)
	if err != nil {
		return err
	}

	return svc.client.RequestDecoder("PATCH", apiPath, body, response)
}

func (svc *IntegrationsService) delete(guid string, response interface{}) error {
	apiPath := fmt.Sprintf(apiIntegrationFromGUID, guid)
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
	Details 		   map[string]interface{} `json:"details,omitempty"`
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
