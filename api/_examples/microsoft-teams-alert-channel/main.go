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

	myTeamsChannel := api.NewAlertChannel("microsoft-teams-alert-from-golang",
		api.MicrosoftTeamsAlertChannelType,
		api.MicrosoftTeamsData{
			TeamsURL: "https://outlook.office.com/webhook/api-token",
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myTeamsChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Microsoft Teams alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Microsoft Teams alert channel created: %s", response.Data.IntgGuid)
}
