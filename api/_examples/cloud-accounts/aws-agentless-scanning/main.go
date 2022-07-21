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

	res, err := lacework.V2.CloudAccounts.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range res.Data {
		support := "Unsupported"
		switch account.Type {
		case api.awsSidekickCloudAccount.String():
			support = "Supported"
		}

		// Output: INTEGRATION-GUID:INTEGRATION-TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", account.IntgGuid, account.Type, support)
	}

	awsAgentlessScanningData := api.AwsAgentlessScanningData{
		ScanFrequency:           24,
		ScanContainers:          true,
		ScanHostVulnerabilities: true,
	}

	awsAgentlessScanningAccount := api.NewCloudAccount(
		fmt.Sprintf("%s-from-golang", api.awsSidekickCloudAccount.String()),
		api.awsSidekickCloudAccount,
		awsAgentlessScanningData,
	)

	awsAgentlessScanningResponse, err := lacework.V2.CloudAccounts.Create(awsAgentlessScanningAccount)
	if err != nil {
		log.Fatal(err)
	}

	// Output: AwsAgentlessScanning Cloud Account created: THE-INTEGRATION-GUID
	fmt.Printf("Cloud Account created: %s", awsAgentlessScanningResponse.Data.IntgGuid)
}
