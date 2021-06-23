package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	// TODO @afiune maybe think about a way to inject CI credentials and
	// run these examples as part of our CI pipelines
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
		case api.AwsCtSqsCloudAccount.String():
			//case api.AwsCfgCloudAccount:
			//case api.GcpCfgCloudAccount:
			//case api.GcpAtSesCloudAccount:
			support = "Supported"
		}

		// Output: INTEGRATION-GUID:INTEGRATION-TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", account.IntgGuid, account.Type, support)
	}

	awsCtSqsData := api.AwsCtSqsData{
		Credentials: api.AwsCtSqsCredentials{
			RoleArn:    "arn:aws:iam::123456789000:role/lw-iam-b8c91298",
			ExternalID: "abc123",
		},
		QueueUrl: "https://sqs.us-west-2.amazonaws.com/123456789000/lacework-ct-sqs-914bc324",
	}

	awsCtSqsData.EncodeAccountMappingFile([]byte(`{
                              "defaultLaceworkAccountAws": "demo",
                              "integration_mappings": {
                                "business-unit": {
                                  "aws_accounts": [
                                    "123456789000"
                                  ]
                                },
                                "development": {
                                  "aws_accounts": [
                                    "123456789001"
                                  ]
                                }
                              }
                            }`))

	myCloudAccount := api.NewCloudAccount(
		"cloud-from-golang",
		api.AwsCtSqsCloudAccount,
		awsCtSqsData,
	)

	// Account Mappings are only allowed at Organization level,
	// copy the client to make it and org client
	orgLwClient, err := api.CopyClient(lacework,
		api.WithOrgAccess(),
	)
	if err != nil {
		log.Fatal(err)
	}

	response, err := orgLwClient.V2.CloudAccounts.Create(myCloudAccount)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Cloud Account created: THE-INTEGRATION-GUID
	fmt.Printf("Cloud Account created: %s", response.Data.IntgGuid)
}
