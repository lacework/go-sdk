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

	myGcpPubSubChannel := api.NewGcpPubSubAlertChannel("gcp-pub-sub-alert-from-golang",
		api.GcpPubSubChannelData{
			ProjectID: "my-sample-project-191923",
			TopicID:   "mytopic",
			Credentials: api.GcpCredentials{
				ClientID:     "client_id",
				ClientEmail:  "foo@example.iam.gserviceaccount.com",
				PrivateKey:   "priv_key",
				PrivateKeyID: "p_key_id",
			},
		},
	)

	response, err := lacework.Integrations.CreateGcpPubSubAlertChannel(myGcpPubSubChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Gcp Pub Sub alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Gcp Pub Sub alert channel created: %s", response.Data[0].IntgGuid)
}
