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

	myServiceNowChannel := api.NewAlertChannel("service-now-alert-from-golang",
		api.ServiceNowRestAlertChannelType,
		api.ServiceNowRestDataV2{
			InstanceURL:   "https://dev123.service-now.com",
			Username:      "snow-user",
			Password:      "snow-password",
			IssueGrouping: "Events",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myServiceNowChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Service Now alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Service Now alert channel created: %s", response.Data.IntgGuid)
}
