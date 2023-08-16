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

	awsReportCfg := api.AwsReportConfig{AccountID: "123456789", Type: api.AWS_CMMC_1_02}
	awsReport, err := lacework.V2.Reports.Aws.Get(awsReportCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Aws Report: ...
	fmt.Printf("Aws Report: %v", awsReport.Data)
}
