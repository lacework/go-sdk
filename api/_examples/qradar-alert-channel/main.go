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

	myQRadarChannel := api.NewQRadarAlertChannel("qradar-alert-from-golang",
		api.QRadarChannelData{
			HostURL:           "https://qradar-lacework.com",
			HostPort:          8080,
			CommunicationType: api.QRadarCommHttps,
		},
	)

	response, err := lacework.Integrations.CreateQRadarAlertChannel(myQRadarChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: QRadar alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("QRadar alert channel created: %s", response.Data[0].IntgGuid)
}
