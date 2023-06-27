package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")))
	if err != nil {
		log.Fatal(err)
	}

	myCiscoSparkWebhookChannel := api.NewAlertChannel("cisco-webex-alert-from-golang",
		api.CiscoSparkWebhookAlertChannelType,
		api.CiscoSparkWebhookDataV2{
			Webhook: "https://webexapis.com/v1/webhooks/incoming/api-token",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myCiscoSparkWebhookChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Cisco Webex alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Cisco Webex alert channel created: %s", response.Data.IntgGuid)
}
