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

import "go.uber.org/zap"

// NewAzureIntegration returns an instance of AzureIntegration with the provided
// integration type, name and data. The type can only be AzureCfgIntegration or
// AzureActivityLogIntegration
//
// Basic usage: Initialize a new AzureIntegration struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   azure := api.NewAzureIntegration("bar",
//     api.AzureActivityLogIntegration,
//     api.AzureIntegrationData{
//       TenantID: "tenant_id",
//       QueueUrl: "https://abc.queue.core.windows.net/123",
//       Credentials: api.AzureIntegrationCreds{
//         ClientID: "client_id",
//         ClientSecret: "secret",
//       },
//     },
//   )
//   if err != nil {
//     return err
//   }
//
//   client.Integrations.CreateAzure(azure)
//
func NewAzureIntegration(name string, iType integrationType, data AzureIntegrationData) AzureIntegration {
	return AzureIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    iType.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// NewAzureCfgIntegration returns an instance of AzureIntegration of type AZURE_CFG
func NewAzureCfgIntegration(name string, data AzureIntegrationData) AzureIntegration {
	return NewAzureIntegration(name, AzureCfgIntegration, data)
}

// NewAzureActivityLogIntegration returns an instance of AzureIntegration of type AZURE_AL_SEQ
func NewAzureActivityLogIntegration(name string, data AzureIntegrationData) AzureIntegration {
	return NewAzureIntegration(name, AzureActivityLogIntegration, data)
}

// CreateAzure creates a single Azure integration on the Lacework Server
func (svc *IntegrationsService) CreateAzure(integration AzureIntegration) (
	response AzureIntegrationsResponse,
	err error,
) {
	err = svc.create(integration, &response)
	if err != nil {
		return
	}

	// WORKAROUND (@afiune) The backend is currently not triggering an initial
	// report automatically after creation of Cloud Account (CFG) Integrations,
	// we are implementing this trigger here until we implement it in the backend
	// with RAIN-13422
	if len(response.Data) == 0 {
		return
	}
	if integration.Type != AzureCfgIntegration.String() {
		return
	}

	intgGuid := response.Data[0].IntgGuid
	svc.client.log.Info("triggering compliance report",
		zap.String("cloud_integration", integration.Type),
		zap.String("int_guid", intgGuid),
	)
	_, errComplianceReport := svc.client.Compliance.RunIntegrationReport(intgGuid)
	if errComplianceReport != nil {
		svc.client.log.Warn("unable to trigger compliance report",
			zap.String("cloud_integration", integration.Type),
			zap.String("int_guid", intgGuid),
			zap.String("error", errComplianceReport.Error()),
		)
	}

	return
}

// GetAzure gets a single Azure integration matching the integration guid on
// the Lacework Server
func (svc *IntegrationsService) GetAzure(guid string) (
	response AzureIntegrationsResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAzure updates a single Azure integration on the Lacework Server
func (svc *IntegrationsService) UpdateAzure(data AzureIntegration) (
	response AzureIntegrationsResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// DeleteAzure deletes a single Azure integration matching the integration on
// the Lacework Server
func (svc *IntegrationsService) DeleteAzure(guid string) (
	response AzureIntegrationsResponse,
	err error,
) {
	err = svc.delete(guid, &response)
	return

}

// ListAzureCfg lists the AZURE_CFG external integrations available on the Lacework Server
func (svc *IntegrationsService) ListAzureCfg() (
	response AzureIntegrationsResponse, err error,
) {
	err = svc.listByType(AzureCfgIntegration, &response)
	return
}

// ListAzureActivityLog lists the AZURE_AL_SEQ external integrations available
// on the Lacework Server
func (svc *IntegrationsService) ListAzureActivityLog() (
	response AzureIntegrationsResponse, err error,
) {
	err = svc.listByType(AzureActivityLogIntegration, &response)
	return
}

type AzureIntegrationsResponse struct {
	Data    []AzureIntegration `json:"data"`
	Ok      bool               `json:"ok"`
	Message string             `json:"message"`
}

type AzureIntegration struct {
	commonIntegrationData
	Data AzureIntegrationData `json:"DATA"`
}

type AzureIntegrationData struct {
	Credentials AzureIntegrationCreds `json:"CREDENTIALS" mapstructure:"CREDENTIALS"`
	TenantID    string                `json:"TENANT_ID" mapstructure:"TENANT_ID"`

	// QueueUrl is a field that exists and is required for the AWS_CT_SQS integration,
	// though, it doesn't exist for AZURE_CFG integrations, that's why we omit it if empty
	QueueUrl string `json:"QUEUE_URL,omitempty" mapstructure:"QUEUE_URL"`
}

type AzureIntegrationCreds struct {
	ClientID     string `json:"CLIENT_ID" mapstructure:"CLIENT_ID"`
	ClientSecret string `json:"CLIENT_SECRET" mapstructure:"CLIENT_SECRET"`
}
