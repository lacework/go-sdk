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

	myCiscoWebexChannel := api.NewCiscoWebexAlertChannel("cisco-webex-alert-from-golang",
		api.CiscoWebexChannelData{
			WebhookURL: "https://webexapis.com/v1/webhooks/incoming/api-token",
		},
	)

	response, err := lacework.Integrations.CreateCiscoWebexAlertChannel(myCiscoWebexChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Cisco Webex alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Cisco Webex alert channel created: %s", response.Data[0].IntgGuid)
}
