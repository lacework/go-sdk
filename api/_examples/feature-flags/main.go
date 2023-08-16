package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
	)
	if err != nil {
		log.Fatal(err)
	}

	response, err := lacework.V2.FeatureFlags.GetFeatureFlagsMatchingPrefix("PUBLIC.sca")
	if err != nil {
		log.Fatal(err)
	}

	// Output: boolean flags evaluating with truthy values and starting with "PUBLIC.sca"
	for _, flag := range response.Data.Flags {
		fmt.Println(flag)
	}
}
