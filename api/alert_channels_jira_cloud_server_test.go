//
// Author:: Vatasha White (<vatasha.white@lacework.net>)
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

package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestAlertChannelsService_GetJira(t *testing.T) {
	var (
		intgGUIDCloud  = intgguid.New()
		intgGUIDServer = intgguid.New()
		apiPathCloud   = fmt.Sprintf("AlertChannels/%s", intgGUIDCloud)
		apiPathServer  = fmt.Sprintf("AlertChannels/%s", intgGUIDServer)
		fakeServer     = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	type test struct {
		name     string
		intgGUID string
		apiPath  string
		jiraType string
		response string
	}

	tests := []test{
		{
			name:     "Get Jira Cloud",
			intgGUID: intgGUIDCloud,
			apiPath:  apiPathCloud,
			jiraType: api.JiraCloudAlertType,
			response: generateAlertChannelResponse(singleJiraCloudAlertChannel(intgGUIDCloud)),
		},
		{
			name:     "Get Jira Server",
			intgGUID: intgGUIDServer,
			apiPath:  apiPathServer,
			jiraType: api.JiraServerAlertType,
			response: generateAlertChannelResponse(singleJiraServerAlertChannel(intgGUIDServer)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeServer.MockAPI(tc.apiPath, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method, "GetJira() should be a GET method")
				fmt.Fprintf(w, tc.response)
			})

			c, err := api.NewClient("test",
				api.WithToken("TOKEN"),
				api.WithURL(fakeServer.URL()),
			)
			assert.Nil(t, err)

			response, err := c.V2.AlertChannels.GetJira(tc.intgGUID)
			assert.Nil(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, tc.intgGUID, response.Data.IntgGuid)
			assert.Equal(t, "integration_name", response.Data.Name)
			assert.Equal(t, "Jira", response.Data.Type)
			assert.True(t, response.Data.State.Ok)
			if tc.jiraType == api.JiraCloudAlertType {
				assert.Equal(t, "fake-api-token", response.Data.Data.ApiToken)
				assert.Equal(t, "JIRA_CLOUD", response.Data.Data.JiraType)
			} else if tc.jiraType == api.JiraServerAlertType {
				assert.Equal(t, "fake-password", response.Data.Data.Password)
				assert.Equal(t, "JIRA_SERVER", response.Data.Data.JiraType)
			}
			assert.Equal(t, "fake-custom-template-file", response.Data.Data.CustomTemplateFile)
			assert.Equal(t, "Events", response.Data.Data.IssueGrouping)
			assert.Equal(t, "fake-issue-type", response.Data.Data.IssueType)
			assert.Equal(t, "fake-jira-url", response.Data.Data.JiraUrl)
			assert.Equal(t, "fake-project-id", response.Data.Data.ProjectID)
			assert.Equal(t, "fake-username", response.Data.Data.Username)
		})
	}
}

func TestAlertChannelsService_UpdateJira(t *testing.T) {
	var (
		intgGUIDCloud  = intgguid.New()
		intgGUIDServer = intgguid.New()
		apiPathCloud   = fmt.Sprintf("AlertChannels/%s", intgGUIDCloud)
		apiPathServer  = fmt.Sprintf("AlertChannels/%s", intgGUIDServer)
		fakeServer     = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	type test struct {
		name     string
		intgGUID string
		apiPath  string
		jiraType string
		response string
		data     api.JiraDataV2
	}

	tests := []test{
		{
			name:     "Update Jira Cloud",
			intgGUID: intgGUIDCloud,
			apiPath:  apiPathCloud,
			jiraType: api.JiraCloudAlertType,
			response: generateAlertChannelResponse(singleJiraCloudAlertChannel(intgGUIDCloud)),
			data: api.JiraDataV2{
				ApiToken:           "fake-api-token",
				CustomTemplateFile: "fake-custom-template-file",
				IssueGrouping:      "Events",
				IssueType:          "fake-issue-type",
				JiraType:           api.JiraCloudAlertType,
				JiraUrl:            "fake-jira-url",
				ProjectID:          "fake-project-id",
				Username:           "fake-username",
			},
		},
		{
			name:     "Update Jira Server",
			intgGUID: intgGUIDServer,
			apiPath:  apiPathServer,
			jiraType: api.JiraServerAlertType,
			response: generateAlertChannelResponse(singleJiraServerAlertChannel(intgGUIDServer)),
			data: api.JiraDataV2{
				CustomTemplateFile: "fake-custom-template-file",
				IssueGrouping:      "Events",
				IssueType:          "fake-issue-type",
				JiraType:           api.JiraServerAlertType,
				JiraUrl:            "fake-jira-url",
				ProjectID:          "fake-project-id",
				Username:           "fake-username",
				Password:           "fake-password",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeServer.MockAPI(tc.apiPath, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "PATCH", r.Method, "UpdateJira() should be a PATCH method")

				if assert.NotNil(t, r.Body) {
					body := httpBodySniffer(r)
					assert.Contains(t, body, tc.intgGUID, "IntgGUID missing")
					assert.Contains(t, body, "name\":\"integration_name", "cloud account name is missing")
					assert.Contains(t, body, "type\":\"Jira", "wrong cloud account type")
					if tc.jiraType == api.JiraCloudAlertType {
						assert.Contains(t, body, "apiToken\":\"fake-api-token", "missing api token")
						assert.Contains(t, body, "jiraType\":\"JIRA_CLOUD", "missing jira type")
					} else if tc.jiraType == api.JiraServerAlertType {
						assert.Contains(t, body, "password\":\"fake-password", "missing password")
						assert.Contains(t, body, "jiraType\":\"JIRA_SERVER", "missing jira type")
					}
					assert.Contains(t, body, "customTemplateFile\":\"fake-custom-template-file", "missing custom template file")
					assert.Contains(t, body, "issueGrouping\":\"Events", "missing issue grouping")
					assert.Contains(t, body, "issueType\":\"fake-issue-type", "missing issue type")
					assert.Contains(t, body, "jiraUrl\":\"fake-jira-url", "missing jira url")
					assert.Contains(t, body, "projectId\":\"fake-project-id", "missing project id")
					assert.Contains(t, body, "username\":\"fake-username", "missing username")
					assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
				}

				fmt.Fprintf(w, tc.response)
			})

			c, err := api.NewClient("test",
				api.WithToken("TOKEN"),
				api.WithURL(fakeServer.URL()),
			)
			assert.Nil(t, err)

			jiraCloudAlertChan := api.NewAlertChannel("integration_name",
				api.JiraAlertChannelType,
				tc.data,
			)
			assert.Equal(t, "integration_name", jiraCloudAlertChan.Name)
			assert.Equal(t, "Jira", jiraCloudAlertChan.Type)
			assert.Equal(t, 1, jiraCloudAlertChan.Enabled)
			if tc.jiraType == api.JiraCloudAlertType {
				assert.Equal(t, "JIRA_CLOUD", jiraCloudAlertChan.Data.(api.JiraDataV2).JiraType)
			} else if tc.jiraType == api.JiraServerAlertType {
				assert.Equal(t, "JIRA_SERVER", jiraCloudAlertChan.Data.(api.JiraDataV2).JiraType)
			}
			jiraCloudAlertChan.IntgGuid = tc.intgGUID

			response, err := c.V2.AlertChannels.UpdateJira(jiraCloudAlertChan)
			if assert.NoError(t, err) {
				assert.NotNil(t, response)
				assert.Equal(t, tc.intgGUID, response.Data.IntgGuid)
				assert.True(t, response.Data.State.Ok)
				assert.Equal(t, "integration_name", response.Data.Name)
				if tc.jiraType == api.JiraCloudAlertType {
					assert.Equal(t, "fake-api-token", response.Data.Data.ApiToken)
					assert.Equal(t, "JIRA_CLOUD", response.Data.Data.JiraType)
				} else if tc.jiraType == api.JiraServerAlertType {
					assert.Equal(t, "fake-password", response.Data.Data.Password)
					assert.Equal(t, "JIRA_SERVER", response.Data.Data.JiraType)
				}
				assert.Equal(t, "fake-custom-template-file", response.Data.Data.CustomTemplateFile)
				assert.Equal(t, "Events", response.Data.Data.IssueGrouping)
				assert.Equal(t, "fake-issue-type", response.Data.Data.IssueType)
				assert.Equal(t, "fake-jira-url", response.Data.Data.JiraUrl)
				assert.Equal(t, "fake-project-id", response.Data.Data.ProjectID)
				assert.Equal(t, "fake-username", response.Data.Data.Username)
			}
		})
	}
}

func singleJiraCloudAlertChannel(id string) string {
	return fmt.Sprintf(`
	{
		"createdOrUpdatedBy": "vatasha.white@lacework.net",
		"createdOrUpdatedTime": "2021-09-29T117:55:47.277316",
		"data": {
			"apiToken": "fake-api-token",
			"customTemplateFile": "fake-custom-template-file",
			"issueGrouping": "Events",
			"issueType": "fake-issue-type",
			"jiraType": "JIRA_CLOUD",
			"jiraUrl": "fake-jira-url",
			"projectId": "fake-project-id",
			"username": "fake-username"
		},
		"enabled": 1,
		"intgGuid": %q,
		"isOrg": 0,
		"name": "integration_name",
		"state": {
		"details": {},
		"lastSuccessfulTime": 1632932665892,
			"lastUpdatedTime": 1632932665892,
			"ok": true
	},
		"type": "Jira"
	}
	`, id)
}

func singleJiraServerAlertChannel(id string) string {
	return fmt.Sprintf(`
	{
		"createdOrUpdatedBy": "vatasha.white@lacework.net",
		"createdOrUpdatedTime": "2021-09-29T117:55:47.277316",
		"data": {
			"password": "fake-password",
			"customTemplateFile": "fake-custom-template-file",
			"issueGrouping": "Events",
			"issueType": "fake-issue-type",
			"jiraType": "JIRA_SERVER",
			"jiraUrl": "fake-jira-url",
			"projectId": "fake-project-id",
			"username": "fake-username"
		},
		"enabled": 1,
		"intgGuid": %q,
		"isOrg": 0,
		"name": "integration_name",
		"state": {
		"details": {},
		"lastSuccessfulTime": 1632932665892,
			"lastUpdatedTime": 1632932665892,
			"ok": true
	},
		"type": "Jira"
	}
	`, id)
}

func TestJiraGroupings(t *testing.T) {
	surveyGroupings := []int{}
	for _, s := range api.JiraIssueGroupingsSurvey {
		surveyGroupings = append(surveyGroupings, int(s))
	}

	for i := range api.JiraIssueGroupings {
		assert.Contains(t, surveyGroupings, int(i))
	}
}
