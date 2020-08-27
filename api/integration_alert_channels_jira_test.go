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

package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestIntegrationsJiraAlertChannelTypes(t *testing.T) {
	assert.Equal(t, "JIRA_CLOUD", api.JiraCloudAlertType)
	assert.Equal(t, "JIRA_SERVER", api.JiraServerAlertType)
}

func TestIntegrationsNewJiraAlertChannel(t *testing.T) {
	subject := api.NewJiraAlertChannel("integration_name",
		api.JiraAlertChannelData{
			JiraType:      api.JiraCloudAlertType,
			JiraUrl:       "mycompany.atlassian.net",
			IssueType:     "Bug",
			ProjectID:     "TEST",
			Username:      "my@username.com",
			ApiToken:      "my-api-token",
			IssueGrouping: "Resources",
		},
	)
	assert.Equal(t, api.JiraIntegration.String(), subject.Type)
	assert.Equal(t, api.JiraCloudAlertType, subject.Data.JiraType)
}

func TestIntegrationsNewJiraAlertChannelWithCustomTemplateFile(t *testing.T) {
	templateJSON := `{
      "fields": {
          "labels": [
              "myLabel"
          ],
          "priority": 
          {
              "id": "1"
          }
      }
  }`
	jira := api.JiraAlertChannelData{
		JiraType:      api.JiraCloudAlertType,
		JiraUrl:       "mycompany.atlassian.net",
		IssueType:     "Bug",
		ProjectID:     "TEST",
		Username:      "my@username.com",
		ApiToken:      "my-api-token",
		IssueGrouping: "Resources",
	}
	jira.EncodeCustomTemplateFile(templateJSON)

	subject := api.NewJiraAlertChannel("integration_name", jira)
	assert.Equal(t, api.JiraIntegration.String(), subject.Type)
	assert.Equal(t, api.JiraCloudAlertType, subject.Data.JiraType)
	assert.Contains(t,
		subject.Data.CustomTemplateFile,
		"data:application/json;name=i.json;base64,",
		"check the custom_template_file encoder",
	)
}

func TestIntegrationsNewJiraCloudAlertChannel(t *testing.T) {
	subject := api.NewJiraCloudAlertChannel("integration_name",
		api.JiraAlertChannelData{
			JiraUrl:       "mycompany.atlassian.net",
			IssueType:     "Bug",
			ProjectID:     "TEST",
			Username:      "my@username.com",
			ApiToken:      "my-api-token",
			IssueGrouping: "Resources",
		},
	)
	assert.Equal(t, api.JiraIntegration.String(), subject.Type)
	assert.Equal(t, api.JiraCloudAlertType, subject.Data.JiraType)
}

func TestIntegrationsNewJiraServerAlertChannel(t *testing.T) {
	subject := api.NewJiraServerAlertChannel("integration_name",
		api.JiraAlertChannelData{
			JiraUrl:       "mycompany.atlassian.net",
			IssueType:     "Bug",
			ProjectID:     "TEST",
			Username:      "my@username.com",
			Password:      "my-password",
			IssueGrouping: "Resources",
		},
	)
	assert.Equal(t, api.JiraIntegration.String(), subject.Type)
	assert.Equal(t, api.JiraServerAlertType, subject.Data.JiraType)
}

func TestIntegrationsCreateJiraAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateJiraAlertChannel should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "JIRA", "wrong integration type")
			assert.Contains(t, body, "JIRA_CLOUD", "wrong jira type")
			assert.Contains(t, body, "mycompany.atlassian.net", "wrong jira url")
			assert.Contains(t, body, "Bug", "wrong issue type")
			assert.Contains(t, body, "TEST", "wrong project_id")
			assert.Contains(t, body, "my@username.com", "wrong username")
			assert.Contains(t, body, "my-api-token", "wrong api token")
			assert.Contains(t, body, "Resources", "wrong issue grouping")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, jiraIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewJiraAlertChannel("integration_name",
		api.JiraAlertChannelData{
			JiraType:      api.JiraCloudAlertType,
			JiraUrl:       "mycompany.atlassian.net",
			IssueType:     "Bug",
			ProjectID:     "TEST",
			Username:      "my@username.com",
			ApiToken:      "my-api-token",
			IssueGrouping: "Resources",
		},
	)
	assert.Equal(t, "integration_name", data.Name, "JIRA integration name mismatch")
	assert.Equal(t, "JIRA", data.Type, "a new JIRA integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new JIRA integration should be enabled")

	response, err := c.Integrations.CreateJiraAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "mycompany.atlassian.net", resData.Data.JiraUrl)
		assert.Equal(t, "JIRA_CLOUD", resData.Data.JiraType)
		assert.Equal(t, "Bug", resData.Data.IssueType)
		assert.Equal(t, "TEST", resData.Data.ProjectID)
		assert.Equal(t, "my@username.com", resData.Data.Username)
		assert.Equal(t, "Resources", resData.Data.IssueGrouping)
	}
}

func TestIntegrationsGetJiraAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetJiraAlertChannel should be a GET method")
		fmt.Fprintf(w, jiraIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetJiraAlertChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "mycompany.atlassian.net", resData.Data.JiraUrl)
		assert.Equal(t, "JIRA_CLOUD", resData.Data.JiraType)
		assert.Equal(t, "Bug", resData.Data.IssueType)
		assert.Equal(t, "TEST", resData.Data.ProjectID)
		assert.Equal(t, "my@username.com", resData.Data.Username)
		assert.Equal(t, "Resources", resData.Data.IssueGrouping)
	}
}

func TestIntegrationsUpdateJiraAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateJiraAlertChannel should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "JIRA", "wrong integration type")
			assert.Contains(t, body, "JIRA_CLOUD", "wrong jira type")
			assert.Contains(t, body, "mycompany.atlassian.net", "wrong jira url")
			assert.Contains(t, body, "Bug", "wrong issue type")
			assert.Contains(t, body, "TEST", "wrong project_id")
			assert.Contains(t, body, "my@username.com", "wrong username")
			assert.Contains(t, body, "my-api-token", "wrong api token")
			assert.Contains(t, body, "Resources", "wrong issue grouping")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, jiraIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewJiraCloudAlertChannel("integration_name",
		api.JiraAlertChannelData{
			JiraUrl:       "mycompany.atlassian.net",
			IssueType:     "Bug",
			ProjectID:     "TEST",
			Username:      "my@username.com",
			ApiToken:      "my-api-token",
			IssueGrouping: "Resources",
		},
	)
	assert.Equal(t, "integration_name", data.Name, "JIRA integration name mismatch")
	assert.Equal(t, "JIRA", data.Type, "a new JIRA integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new JIRA integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateJiraAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListJiraAlertChannel(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/JIRA",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListJiraAlertChannel should be a GET method")
			fmt.Fprintf(w, jiraMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListJiraAlertChannel()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func jiraIntegrationJsonResponse(intgGUID string) string {
	return `
{
  "data": [` + singleJiraIntegration(intgGUID) + `],
  "ok": true,
  "message": "SUCCESS"
}
`
}

func jiraMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleJiraIntegration(guid))
	}
	return `
{
"data": [` + strings.Join(integrations, ", ") + `],
"ok": true,
"message": "SUCCESS"
}
`
}

// @afiune heads-up: MIN_ALERT_SEVERITY is a legacy field
func singleJiraIntegration(id string) string {
	return `
{
  "INTG_GUID": "` + id + `",
  "CREATED_OR_UPDATED_BY": "user@email.com",
  "CREATED_OR_UPDATED_TIME": "2020-Jul-16 19:59:22 UTC",
  "DATA": {
    "ISSUE_GROUPING": "Resources",
    "ISSUE_TYPE": "Bug",
    "JIRA_TYPE": "JIRA_CLOUD",
    "JIRA_URL": "mycompany.atlassian.net",
    "MIN_ALERT_SEVERITY": 1,
    "PROJECT_ID": "TEST",
    "USERNAME": "my@username.com"
  },
  "ENABLED": 1,
  "IS_ORG": 0,
  "NAME": "integration_name",
  "STATE": {
    "lastSuccessfulTime": "2020-Jul-16 18:26:54 UTC",
    "lastUpdatedTime": "2020-Jul-16 18:26:54 UTC",
    "ok": true
  },
  "TYPE": "JIRA",
  "TYPE_NAME": "JIRA"
}
`
}
