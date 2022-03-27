package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"), api.WithApiV2(),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")))
	if err != nil {
		log.Fatal(err)
	}

	mySlackChannel := api.NewAlertChannel("slack-alert-from-golang",
		api.SlackChannelAlertChannelType,
		api.SlackChannelDataV2{
			SlackUrl: "https://hooks.slack.com/services/ABCD/12345/abcd1234",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(mySlackChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Slack alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Slack alert channel created: %s", response.Data.IntgGuid)
}
