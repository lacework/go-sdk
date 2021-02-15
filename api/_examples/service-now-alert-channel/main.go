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

	myServiceNowChannel := api.NewServiceNowAlertChannel("service-now-alert-from-golang",
		api.ServiceNowChannelData{
			InstanceURL:   "snow-lacework.com",
			Username:      "snow-user",
			Password:      "snow-password",
			IssueGrouping: "Events",
		},
	)

	response, err := lacework.Integrations.CreateServiceNowAlertChannel(myServiceNowChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Service Now alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Service Now alert channel created: %s", response.Data[0].IntgGuid)
}
