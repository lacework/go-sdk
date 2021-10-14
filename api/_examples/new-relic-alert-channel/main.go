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

	myNewRelicChannel := api.NewAlertChannel("new-relic-alert-from-golang",
		api.NewRelicAlertChannelType,
		api.NewRelicDataV2{
			AccountID: 2338053,
			InsertKey: "x-xx-xxxxxxxxxxxxxxxxxx",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myNewRelicChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: New Relic alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("New Relic alert channel created: %s", response.Data.IntgGuid)
}
