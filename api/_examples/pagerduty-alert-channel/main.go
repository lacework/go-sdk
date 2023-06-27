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
	)
	if err != nil {
		log.Fatal(err)
	}

	alert := api.NewAlertChannel("pagerduty-alert-from-golang",
		api.PagerDutyApiAlertChannelType,
		api.PagerDutyApiDataV2{
			IntegrationKey: "1234abc8901abc567abc123abc78e012",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(alert)
	if err != nil {
		log.Fatal(err)
	}

	// Output: PagerDuty alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("PagerDuty alert channel created: %s", response.Data.IntgGuid)
}
