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
		case api.GcpGkeAuditCloudAccount.String():
			support = "Supported"
		}

		// Output: INTEGRATION-GUID:INTEGRATION-TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", account.IntgGuid, account.Type, support)
	}

	gcpGkeAuditData := api.GcpGkeAuditData{
		Credentials: api.GcpGkeAuditCredentials{
			ClientEmail:  "ross.moles@lacework.net",
			ClientId:     "0123456789",
			PrivateKey:   "",
			PrivateKeyId: "",
		},
		IntegrationType:  "Project",
		OrganizationId:   "OrgId",
		ProjectId:        "ProjectMcProjectFace",
		SubscriptionName: "projects/ProjectMcProjectFace/subscriptions/SubscribeyMcSubscribeFace",
	}

	gcpGkeAuditCloudAccount := api.NewCloudAccount(
		"cloud-from-golang",
		api.GcpGkeAuditCloudAccount,
		gcpGkeAuditData,
	)

	gcpGkeAuditResponse, err := lacework.V2.CloudAccounts.Create(gcpGkeAuditCloudAccount)
	if err != nil {
		log.Fatal(err)
	}

	// Output: GcpGkeAudit Cloud Account created: THE-INTEGRATION-GUID
	fmt.Printf("Cloud Account created: %s", gcpGkeAuditResponse.Data.IntgGuid)
}
