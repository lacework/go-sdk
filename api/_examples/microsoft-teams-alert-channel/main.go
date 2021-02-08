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

	myTeamsChannel := api.NewMicrosoftTeamsAlertChannel("microsoft-teams-alert-from-golang",
		api.MicrosoftTeamsChannelData{
			TeamsURL: "https://outlook.office.com/webhook/api-token",
		},
	)

	response, err := lacework.Integrations.CreateMicrosoftTeamsAlertChannel(myTeamsChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Microsoft Teams alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Microsoft Teams alert channel created: %s", response.Data[0].IntgGuid)
}
