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

	alert := api.NewAwsCloudWatchAlertChannel("aws-cloudwatch-alert-from-golang",
		api.AwsCloudWatchData{
			EventBusArn:      "arn:aws:events:us-west-2:1234567890:event-bus/default",
			MinAlertSeverity: 1,
		},
	)

	response, err := lacework.Integrations.CreateAwsCloudWatchAlertChannel(alert)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Aws CloudWatch alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Aws CloudWatch alert channel created: %s", response.Data[0].IntgGuid)
}
