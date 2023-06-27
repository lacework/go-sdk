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

	// List all report distributions
	reportDistributions, err := lacework.V2.ReportDistributions.List()
	if err != nil {
		log.Fatal(err)
	}

	// Get the details of a single report distribution
	reportDistribution := reportDistributions.Data[0]
	report, err := lacework.V2.ReportDistributions.Get(reportDistribution.ReportDistributionGuid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Retrieved my report distribution: %q", report.Data.DistributionName)

	// Delete a report
	err = lacework.V2.ReportDefinitions.Delete(report.Data.ReportDefinitionGuid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted report distribution: %q", report.Data.ReportDistributionGuid)
}
