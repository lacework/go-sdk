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

	// List all report definitions
	reportTypes, err := lacework.V2.ReportDefinitions.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, report := range reportTypes.Data {
		fmt.Println(report.ReportName)
	}

	// Get the details of a single report definition
	report, err := lacework.V2.ReportDefinitions.Get("exampleID")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Retrieved my custom report: %q", report.Data.ReportName)

	// Create a new custom report definition
	myReport := api.ReportDefinitionConfig{
		ReportName:    "My Custom Report",
		ReportType:    "COMPLIANCE",
		SubReportType: "AWS",
		Sections: []api.ReportDefinitionSection{
			{Category: "Cust",
				Title:    "My Custom Category",
				Policies: []string{"AWS_CIS_2_6"},
			}},
		AlertChannels:    []string{"myAlertChannel"},
		DistributionType: "csv",
		Frequency:        "weekly",
	}

	newReport, err := lacework.V2.ReportDefinitions.Create(api.NewReportDefinition(myReport))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Create my new custom report: %q", newReport.Data.ReportName)
}
