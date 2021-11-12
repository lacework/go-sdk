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

	myWebhookChannel := api.NewAlertChannel("webhook-alert-from-golang",
		api.WebhookAlertChannelType,
		api.WebhookDataV2{
			WebhookUrl: "https://mywebhook.com/?api-token=123",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myWebhookChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Webhook alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Webhook alert channel created: %s", response.Data.IntgGuid)
}
