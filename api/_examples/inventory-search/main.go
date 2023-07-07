package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lacework/go-sdk/api"
)

// Execute: When having the Lacework CLI installed and configured, it is possible
// to execute the following example with the command:
//
//	LW_API_KEY=$(lacework configure show api_key) \
//		LW_API_SECRET=$(lacework configure show api_secret) \
//		LW_ACCOUNT=$(lacework configure show account) \
//		LW_SUBACCOUNT=$(lacework configure show subaccount) \
//		go run api/_examples/inventory-search/main.go
func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
	)
	if err != nil {
		log.Fatal(err)
	}

	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -1) // last day
	)

	var (
		awsInventorySearchResponse api.InventoryAwsResponse
		filters                    = api.InventorySearch{
			SearchFilter: api.SearchFilter{
				TimeFilter: &api.TimeFilter{
					StartTime: &before,
					EndTime:   &now,
				},
			},
			Csp: api.AwsInventoryType,
		}
	)

	err = lacework.V2.Inventory.Search(&awsInventorySearchResponse, filters)
	if err != nil {
		log.Fatal(err)
	}

	if awsInventorySearchResponse.Paging.Urls.NextPage != "" {
		fmt.Printf(
			">> Reading %d pages\n",
			awsInventorySearchResponse.Paging.TotalRows/awsInventorySearchResponse.Paging.Rows,
		)
	}

	for {
		fmt.Printf(">> Resources read: %d\n", awsInventorySearchResponse.Paging.Rows)

		for _, resource := range awsInventorySearchResponse.Data {
			fmt.Printf("%s - %s %s\n", resource.StartTime, resource.EndTime, resource.ResourceId)
		}

		pageOk, err := lacework.NextPage(&awsInventorySearchResponse)
		if err == nil && pageOk {
			continue
		}
		break
	}
}
