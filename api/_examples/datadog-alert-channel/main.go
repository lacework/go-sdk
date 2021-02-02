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

	myDatadogChannel := api.NewDatadogAlertChannel("datadog-alert-from-golang",
		api.DatadogChannelData{
			DatadogSite:    "eu",
			DatadogService: "Events Summary",
			ApiKey:         "datadog-key",
		},
	)

	response, err := lacework.Integrations.CreateDatadogAlertChannel(myDatadogChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Datadog alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Datadog alert channel created: %s", response.Data[0].IntgGuid)
}
