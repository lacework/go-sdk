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
		case api.GcpAlPubSubCloudAccount.String():
			support = "Supported"
		}

		// Output: INTEGRATION-GUID:INTEGRATION-TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", account.IntgGuid, account.Type, support)
	}

	gcpAlPubSubData := api.GcpAlPubSubSesData{
		Credentials: api.GcpAlPubSubCredentials{
			ClientEmail:  "test@project.iam.gserviceaccount.com",
			ClientId:     "0123456789",
			PrivateKey:   "",
			PrivateKeyId: "",
		},
		IntegrationType:  "PROJECT",
		OrganizationId:   "TestOrgId",
		ProjectId:        "TestProjectId",
		SubscriptionName: "projects/test-project/subscriptions/test-subscription",
		TopicId:          "projects/test-project/subscriptions/test-topic",
	}

	gcpAlPubSubCloudAccount := api.NewCloudAccount(
		"cloud-from-golang",
		api.GcpAlPubSubCloudAccount,
		gcpAlPubSubData,
	)

	gcpAlPubSubResponse, err := lacework.V2.CloudAccounts.Create(gcpAlPubSubCloudAccount)
	if err != nil {
		log.Fatal(err)
	}

	// Output: GcpAlPubSub Cloud Account created: THE-INTEGRATION-GUID
	fmt.Printf("Cloud Account created: %s", gcpAlPubSubResponse.Data.IntgGuid)
}
