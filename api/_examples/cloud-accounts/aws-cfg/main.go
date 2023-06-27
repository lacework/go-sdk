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
		case api.AwsCfgCloudAccount.String():
			support = "Supported"
		}

		// Output: INTEGRATION-GUID:INTEGRATION-TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", account.IntgGuid, account.Type, support)
	}

	awsCfgData := api.AwsCfgData{
		Credentials: api.AwsCfgCredentials{
			RoleArn:    "arn:aws:iam::1234567890:role/lacework_iam_example_role",
			ExternalID: "abc123",
		},
	}

	awsCtSqsCloudAccount := api.NewCloudAccount(
		"aws-cfg-from-golang",
		api.AwsCfgCloudAccount,
		awsCfgData,
	)

	awsCfgResponse, err := lacework.V2.CloudAccounts.Create(awsCtSqsCloudAccount)
	if err != nil {
		log.Fatal(err)
	}

	// Output: AwsCfg Cloud Account created: THE-INTEGRATION-GUID
	fmt.Printf("Cloud Account created: %s", awsCfgResponse.Data.IntgGuid)
}
