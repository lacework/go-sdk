package resource_groups

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

	res, err := lacework.V2.ResourceGroups.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range res.Data {
		support := "Unsupported"
		switch account.Type {
		case api.AwsResourceGroup.String():
			support = "Supported"
		case api.AzureResourceGroup.String():
			support = "Supported"
		case api.ContainerResourceGroup.String():
			support = "Supported"
		case api.GcpResourceGroup.String():
			support = "Supported"
		case api.LwAccountResourceGroup.String():
			support = "Supported"
		case api.MachineResourceGroup.String():
			support = "Supported"
		}

		// Output: RESOURCE_GUID:RESOURCE_TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", account.ResourceGuid, account.Type, support)
	}

	lwAccResourceGroupData := api.LwAccountResourceGroupData{
		Props: api.LwAccountResourceGroupProps{
			Description: "All AWS accounts",
			LwAccounts:  []string{"*"},
		},
	}

	myResourceGroup := api.NewResourceGroup(
		"resource-group-from-golang",
		api.LwAccountResourceGroup,
		lwAccResourceGroupData,
	)

	// LW_ACCOUNT resource groups are only allowed at Organization level,
	// copy the client to make it an org client
	orgLwClient, err := api.CopyClient(lacework,
		api.WithOrgAccess(),
	)
	if err != nil {
		log.Fatal(err)
	}

	response, err := orgLwClient.V2.ResourceGroups.Create(myResourceGroup)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Resource Group created: RESOURCE_GUID
	fmt.Printf("Resource Group created: %s", response.Data[0].ResourceGuid)
}
