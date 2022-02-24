package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -1) // last day
	)

	filters := api.SearchFilter{
		TimeFilter: &api.TimeFilter{
			StartTime: &before,
			EndTime:   &now,
		},
	}

	response, err := lacework.V2.Vulnerabilities.Containers.Search(filters)
	if err != nil {
		log.Fatal(err)
	}

	if response.Paging.Urls.NextPage != "" {
		// Output: Reading 3 pages
		fmt.Printf("Reading %d pages\n", response.Paging.TotalRows/response.Paging.Rows)
	}

	for {
		// Output: Rows read: 5000
		// Output: Rows read: 5000
		// Output: Rows read: 5000
		// Output: Rows read: 3388
		fmt.Printf("Rows read: %d\n", response.Paging.Rows)

		pageOk, err := lacework.NextPage(&response)
		if err == nil && pageOk {
			continue
		}
		break
	}
}
