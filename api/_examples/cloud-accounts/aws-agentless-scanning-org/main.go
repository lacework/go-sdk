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

	res, err := lacework.V2.CloudAccounts.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range res.Data {
		support := "Unsupported"
		switch account.Type {
		case api.AwsSidekickOrgCloudAccount.String():
			support = "Supported"
		}

		// Output: INTEGRATION-GUID:INTEGRATION-TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", account.IntgGuid, account.Type, support)
	}

	awsSidekickOrgData := api.AwsSidekickOrgData{
		ScanFrequency:           24,
		ScanContainers:          true,
		ScanHostVulnerabilities: true,
		MonitoredAccounts:       "r-1234",
		ManagementAccount:       "000123456789",
		ScanningAccount:         "123456789000",
	}

	awsSidekickOrgAccount := api.NewCloudAccount(
		fmt.Sprintf("%s-from-golang", api.AwsSidekickOrgCloudAccount.String()),
		api.AwsSidekickOrgCloudAccount,
		awsSidekickOrgData,
	)

	awsSidekickOrgResponse, err := lacework.V2.CloudAccounts.CreateAwsSidekickOrg(awsSidekickOrgAccount)
	if err != nil {
		log.Fatal(err)
	}

	// Output: AwsSidekickOrg Cloud Account created: THE-INTEGRATION-GUID
	fmt.Printf("Cloud Account created: %s", awsSidekickOrgResponse.Data.IntgGuid)
}
