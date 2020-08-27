package main

import (
	"log"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient("account", api.WithApiKeys("KEY", "SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	jiraTemplateJSON := `{
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
	jira.EncodeCustomTemplateFile(jiraTemplateJSON)
	jiraAlert := api.NewJiraAlertChannel("integration_name", jira)
	client.Integrations.CreateJiraAlertChannel(jiraAlert)
}
