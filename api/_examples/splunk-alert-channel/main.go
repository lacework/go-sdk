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
		api.WithApiV2(),
	)
	if err != nil {
		log.Fatal(err)
	}

	mySplunkChannel := api.NewAlertChannel("splunk-alert-from-golang",
		api.SplunkHecAlertChannelType,
		api.SplunkHecDataV2{
			Channel:  "channel-name",
			HecToken: "AA111111-11AA-1AA1-11AA-11111AA1111A",
			Host:     "localhost",
			Port:     80,
			Ssl:      false,
			EventData: api.SplunkHecEventDataV2{
				Index:  "test-index",
				Source: "test-source",
			},
		},
	)

	response, err := lacework.V2.AlertChannels.Create(mySplunkChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Splunk alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Splunk alert channel created: %s", response.Data.IntgGuid)
}
