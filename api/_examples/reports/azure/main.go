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

	azureReportCfg := api.AzureReportConfig{TenantID: "example-tenant", SubscriptionID: "example-subscription", Type: api.AZURE_CIS}
	azureReport, err := lacework.V2.Reports.Azure.Get(azureReportCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Azure Report: ...
	fmt.Printf("Azure Report: %v", azureReport.Data)
}
