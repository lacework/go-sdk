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

	myQRadarChannel := api.NewAlertChannel("qradar-alert-from-golang",
		api.IbmQRadarAlertChannelType,
		api.IbmQRadarDataV2{
			HostURL:        "https://qradar-lacework.com",
			HostPort:       8080,
			QRadarCommType: api.QRadarCommHttps,
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myQRadarChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: QRadar alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("QRadar alert channel created: %s", response.Data.IntgGuid)
}
