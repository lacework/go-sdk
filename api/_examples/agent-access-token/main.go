package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"), api.WithApiV2(),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
	)
	if err != nil {
		log.Fatal(err)
	}

	response, err := lacework.V2.AgentAccessTokens.SearchAlias("circle-ci-test-token")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Len: %d\n", len(response.Data))
}
