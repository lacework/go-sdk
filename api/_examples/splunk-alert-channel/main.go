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

	mySplunkChannel := api.NewSplunkAlertChannel("splunk-alert-from-golang",
		api.SplunkChannelData{
			Channel:  "channel-name",
			HecToken: "AA111111-11AA-1AA1-11AA-11111AA1111A",
			Host:     "localhost",
			Port:     80,
			Ssl:      false,
			EventData: api.SplunkEventData{
				Index:  "test-index",
				Source: "test-source",
			},
		},
	)

	response, err := lacework.Integrations.CreateSplunkAlertChannel(mySplunkChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Splunk alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Splunk alert channel created: %s", response.Data[0].IntgGuid)
}
