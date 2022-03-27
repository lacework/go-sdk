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

	myGcpPubSubChannel := api.NewAlertChannel("gcp-pub-sub-alert-from-golang",
		api.GcpPubSubAlertChannelType,
		api.GcpPubSubDataV2{
			ProjectID: "my-sample-project-191923",
			TopicID:   "mytopic",
			Credentials: api.GcpPubSubCredentials{
				ClientID:     "client_id",
				ClientEmail:  "foo@example.iam.gserviceaccount.com",
				PrivateKey:   "priv_key",
				PrivateKeyID: "p_key_id",
			},
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myGcpPubSubChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Gcp Pub Sub alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Gcp Pub Sub alert channel created: %s", response.Data.IntgGuid)
}
