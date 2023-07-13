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

	alertChannels, err := lacework.V2.AlertChannels.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, channel := range alertChannels.Data {
		fmt.Printf("Alert channel: %s\n", channel.Name)
	}
	// Output:
	//
	// Alert channel: DEFAULT EMAIL
}
