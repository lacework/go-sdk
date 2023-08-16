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

	gcpReportCfg := api.GcpReportConfig{OrganizationID: "123456789", ProjectID: "example-project", Type: api.GCP_CIS13}
	gcpReport, err := lacework.V2.Reports.Gcp.Get(gcpReportCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Gcp Report: ...
	fmt.Printf("Gcp Report: %v", gcpReport.Data)
}
