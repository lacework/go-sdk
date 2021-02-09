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

	myNewRelicChannel := api.NewNewRelicAlertChannel("new-relic-alert-from-golang",
		api.NewRelicChannelData{
			AccountID: 2338053,
			InsertKey: "x-xx-xxxxxxxxxxxxxxxxxx",
		},
	)

	response, err := lacework.Integrations.CreateNewRelicAlertChannel(myNewRelicChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: New Relic alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("New Relic alert channel created: %s", response.Data[0].IntgGuid)
}
