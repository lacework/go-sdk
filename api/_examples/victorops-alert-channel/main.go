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

	myVictorOpsChannel := api.NewAlertChannel("victorops-alert-from-golang",
		api.VictorOpsAlertChannelType,
		api.VictorOpsDataV2{
			Url: "https://alert.victorops.com/integrations/generic/20131114/alert/31e945ee-5cad-44e7-afb0-97c20ea80dd8/database",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myVictorOpsChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: VictorOps alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("VictorOps alert channel created: %s", response.Data.IntgGuid)
}
