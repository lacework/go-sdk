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
		api.WithApiV2(),
	)
	if err != nil {
		log.Fatal(err)
	}

	reportTypes, err := lacework.V2.ReportDefinitions.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, report := range reportTypes.Data {
		fmt.Println(report.ReportName)
	}
}
