package client

import (
	"fmt"
	"strings"
)

// GetIntegrations lists the external integrations available on the server
func (c *client) GetIntegrations() (response integrationsResponse, err error) {
	err = c.RequestDecoder("GET", apiIntegrations, nil, &response)
	return
}

func (c *client) GetGCPIntegrations() (response gcpIntegrationsResponse, err error) {
	return
}
func (c *client) GetAzureIntegrations() (response azureIntegrationsResponse, err error) {
	return
}
func (c *client) GetAwsIntegrations() (response awsIntegrationsResponse, err error) {
	return
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
	//Data    []gcpIntegrationData `json:"data"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type gcpIntegrationData struct {
	commonIntegrationData
	Data gcpCfg `json:"DATA"`
}

type gcpCfg struct {
	ID            string         `json:"ID"`
	IdType        string         `json:"ID_TYPE"`
	IssueGrouping string         `json:"ISSUE_GROUPING"`
	Credentials   gcpCredentials `json:"CREDENTIALS"`
}

type gcpCredentials struct {
	ClientId     string `json:"CLIENT_ID"`
	ClientEmail  string `json:"CLIENT_EMAIL"`
	PrivateKeyId string `json:"PRIVATE_KEY_ID"`
	PrivateKey   string `json:"PRIVATE_KEY"`
}
