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
	"encoding/base64"
	"fmt"
)

const (
	JiraCloudAlertType  = "JIRA_CLOUD"
	JiraServerAlertType = "JIRA_SERVER"
)

// NewJiraAlertChannel returns an instance of JiraAlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new JiraAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   jiraAlert := api.NewJiraAlertChannel("foo",
//     api.JiraAlertChannelData{
//       JiraType:         api.JiraCloudAlertType,
//       JiraUrl:          "mycompany.atlassian.net",
//       IssueType:        "Bug",
//       ProjectID:        "EXAMPLE",
//       Username:         "me",
//       ApiToken:         "my-api-token",
//       IssueGrouping:    "Resources",
//     },
//   )
//
//   client.Integrations.CreateJiraAlertChannel(jiraAlert)
//
func NewJiraAlertChannel(name string, data JiraAlertChannelData) JiraAlertChannel {
	return JiraAlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    JiraIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// NewJiraCloudAlertChannel returns a JiraAlertChannel instance preconfigured as a JIRA_CLOUD type
func NewJiraCloudAlertChannel(name string, data JiraAlertChannelData) JiraAlertChannel {
	data.JiraType = JiraCloudAlertType
	return NewJiraAlertChannel(name, data)
}

// NewJiraServerAlertChannel returns a JiraAlertChannel instance preconfigured as a JIRA_SERVER type
func NewJiraServerAlertChannel(name string, data JiraAlertChannelData) JiraAlertChannel {
	data.JiraType = JiraServerAlertType
	return NewJiraAlertChannel(name, data)
}

// CreateJiraAlertChannel creates a jira alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateJiraAlertChannel(integration JiraAlertChannel) (
	response JiraAlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetJiraAlertChannel gets a jira alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetJiraAlertChannel(guid string) (
	response JiraAlertChannelResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateJiraAlertChannel updates a single jira alert channel integration
func (svc *IntegrationsService) UpdateJiraAlertChannel(data JiraAlertChannel) (
	response JiraAlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListJiraAlertChannel lists the JIRA external integrations available on the Lacework Server
func (svc *IntegrationsService) ListJiraAlertChannel() (response JiraAlertChannelResponse, err error) {
	err = svc.listByType(JiraIntegration, &response)
	return
}

type JiraAlertChannelResponse struct {
	Data    []JiraAlertChannel `json:"data"`
	Ok      bool               `json:"ok"`
	Message string             `json:"message"`
}

type JiraAlertChannel struct {
	commonIntegrationData
	Data JiraAlertChannelData `json:"DATA"`
}

type JiraAlertChannelData struct {
	JiraType      string `json:"JIRA_TYPE" mapstructure:"JIRA_TYPE"`
	JiraUrl       string `json:"JIRA_URL" mapstructure:"JIRA_URL"`
	IssueType     string `json:"ISSUE_TYPE" mapstructure:"ISSUE_TYPE"`
	ProjectID     string `json:"PROJECT_ID" mapstructure:"PROJECT_ID"`
	Username      string `json:"USERNAME" mapstructure:"USERNAME"`
	ApiToken      string `json:"API_TOKEN,omitempty" mapstructure:"API_TOKEN"` // Jira Cloud
	Password      string `json:"PASSWORD,omitempty" mapstructure:"PASSWORD"`   // Jira Server
	IssueGrouping string `json:"ISSUE_GROUPING,omitempty" mapstructure:"ISSUE_GROUPING"`

	// This field must be a base64 encode with the following format:
	//
	// "data:application/json;name=i.json;base64,[ENCODING]"
	//
	// [ENCODING] is the the base64 encode, use EncodeCustomTemplateFile() to encode a JSON template
	CustomTemplateFile string `json:"CUSTOM_TEMPLATE_FILE,omitempty"`
}

func (jira *JiraAlertChannelData) EncodeCustomTemplateFile(template string) {
	encodedTemplate := base64.StdEncoding.EncodeToString([]byte(template))
	jira.CustomTemplateFile = fmt.Sprintf("data:application/json;name=i.json;base64,%s", encodedTemplate)
}
