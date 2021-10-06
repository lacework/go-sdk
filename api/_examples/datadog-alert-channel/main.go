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

	myDatadogChannel := api.NewAlertChannel("datadog-alert-from-golang",
		api.DatadogAlertChannelType,
		api.DatadogDataV2{
			DatadogSite: api.DatadogSiteEu,
			DatadogType: api.DatadogServiceEventsSummary,
			ApiKey:      "datadog-key",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myDatadogChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Datadog alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Datadog alert channel created: %s", response.Data.IntgGuid)
}
