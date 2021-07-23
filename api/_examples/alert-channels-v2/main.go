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

	alertChannel, err := lacework.V2.AlertChannels.GetEmailUser("CUSTOMER_8EB5E8092016A0B8CBD8CB591362344E3A87761B997ABA0")
	if err != nil {
		log.Fatal(err)
	}

	// Output: Alert channel: THE-INTEGRATION-GUID
	fmt.Printf("Alert channel: %s", alertChannel.Data.Data.ChannelProps.Recipients[0])
}
