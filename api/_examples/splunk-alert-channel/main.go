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

	mySplunkChannel := api.NewAlertChannel("splunk-alert-from-golang",
		api.SplunkAlertChannelType,
		api.SplunkDataV2{
			Channel:  "channel-name",
			HecToken: "AA111111-11AA-1AA1-11AA-11111AA1111A",
			Host:     "localhost",
			Port:     80,
			Ssl:      false,
			EventData: api.SplunkEventDataV2{
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
