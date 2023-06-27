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

	// List all report definitions
	reportTypes, err := lacework.V2.ReportDefinitions.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, report := range reportTypes.Data {
		fmt.Println(report.ReportName)
	}

	// Get the details of a single report definition
	report, err := lacework.V2.ReportDefinitions.Get("107D750D552D999F6621084BB26D8070F49897D83EC39C205D2F6347")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Retrieved my report: %q", report.Data.ReportName)

	// Delete a report
	err = lacework.V2.ReportDefinitions.Delete(report.Data.ReportDefinitionGuid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted report: %q", report.Data.ReportName)
}
