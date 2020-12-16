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

	myWebhookChannel := api.NewWebhookAlertChannel("webhook-alert-from-golang",
		api.WebhookChannelData{
			WebhookUrl: "https://mywebhook.com/?api-token=123",
		},
	)

	response, err := lacework.Integrations.CreateWebhookAlertChannel(myWebhookChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Webhook alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Webhook alert channel created: %s", response.Data[0].IntgGuid)
}
