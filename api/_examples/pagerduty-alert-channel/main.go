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

	alert := api.NewPagerDutyAlertChannel("pagerduty-alert-from-golang",
		api.PagerDutyData{
			IntegrationKey:   "1234abc8901abc567abc123abc78e012",
			MinAlertSeverity: 5,
		},
	)

	response, err := lacework.Integrations.CreatePagerDutyAlertChannel(alert)
	if err != nil {
		log.Fatal(err)
	}

	// Output: PagerDuty alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("PagerDuty alert channel created: %s", response.Data[0].IntgGuid)
}
