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
		case api.AzureAdAlCloudAccount.String():
			support = "Supported"
		}

		// Output: INTEGRATION-GUID:INTEGRATION-TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", account.IntgGuid, account.Type, support)
	}

	azureAdAlData := api.AzureAdAlData{
		Credentials: api.AzureAdAlCredentials{
			ClientID:     "client-id",
			ClientSecret: "some-secret",
		},
		TenantID:          "tenant-id",
		EventHubNamespace: "EventHubNamespace",
		EventHubName:      "EventHubName",
	}

	azureAdAlCloudAccount := api.NewCloudAccount(
		"cloud-from-golang",
		api.AzureAdAlCloudAccount,
		azureAdAlData,
	)

	azureAdAlIntegrationResponse, err := lacework.V2.CloudAccounts.Create(azureAdAlCloudAccount)
	if err != nil {
		log.Fatal(err)
	}

	// Output: AzureAdAl Cloud Account created: THE-INTEGRATION-GUID
	fmt.Printf("Cloud Account created: %s", azureAdAlIntegrationResponse.Data.IntgGuid)
}
