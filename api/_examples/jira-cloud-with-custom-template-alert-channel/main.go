package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
		api.WithApiV2())
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
	jiraAlert := api.NewAlertChannel("integration_name",
		api.JiraAlertChannelType,
		jira)

	response, err := lacework.V2.AlertChannels.Create(jiraAlert)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("New Jira Cloud integration created: %s", response.Data.IntgGuid)
}
