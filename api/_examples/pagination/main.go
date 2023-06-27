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

	response, err := lacework.V2.Entities.ListMachineDetails()
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Output: Rows read: 5000
		// Output: Rows read: 5000
		// Output: Rows read: 21
		fmt.Printf("Rows read: %d\n", response.Paging.Rows)

		pageOk, err := lacework.NextPage(&response)
		if err == nil && pageOk {
			continue
		}
		break
	}
}
