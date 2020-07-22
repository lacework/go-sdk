package main

import (
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient("account", api.WithApiKeys("KEY", "SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	mySlackChannel := api.NewSlackAlertChannel("slack-alert-from-golang",
		api.SlackChannelData{
			SlackUrl:         "https://hooks.slack.com/services/ABCD/12345/abcd1234",
			MinAlertSeverity: 3,
		},
	)

	response, err := lacework.Integrations.CreateSlackAlertChannel(mySlackChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Slack alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Slack alert channel created: %s", response.Data[0].IntgGuid)
}
