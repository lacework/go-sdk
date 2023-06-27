package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")))
	if err != nil {
		log.Fatal(err)
	}

	alert := api.NewAlertChannel("aws-cloudwatch-alert-from-golang",
		api.CloudwatchEbAlertChannelType,
		api.CloudwatchEbDataV2{
			EventBusArn: "arn:aws:events:us-west-2:1234567890:event-bus/default",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(alert)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Aws CloudWatch alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Aws CloudWatch alert channel created: %s", response.Data.IntgGuid)
}
