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

	res, err := lacework.V2.ResourceGroups.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range res.Data {
		var resourceGuid string
		resourceType := account.Type

		if account.Props != nil {
			resourceGuid = account.ResourceGuid
		} else {
			resourceGuid = account.ResourceGroupGuid
		}

		support := "Unsupported"
		switch resourceType {
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
		fmt.Printf("%s:%s:%s\n", resourceGuid, resourceType, support)
	}

	rgExpression := api.RGExpression{
		Operator: "AND",
		Children: []*api.RGChild{
			{
				FilterName: "filter1",
			},
		},
	}

	rgFilter := api.RGFilter{
		Field:     "Region",
		Operation: "STARTS_WITH",
		Values:    []string{"us-"},
	}

	rgQuery := api.RGQuery{
		Expression: &rgExpression,
		Filters: map[string]*api.RGFilter{
			"filter1": &rgFilter,
		},
	}

	myResourceGroupWithQuery := api.NewResourceGroupWithQuery(
		"resource-group-with-query-from-golang",
		api.AwsResourceGroup,
		"Resource groups in `us` regions",
		&rgQuery,
	)

	println("Creating a resource group v2")
	rgV2Resp, err := lacework.V2.ResourceGroups.Create(myResourceGroupWithQuery)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Succesfully created resource group \n %+v\n", rgV2Resp)

	println("Updating v2 resource group name")
	rgV2Resp.Data.NameV2 = "resource-group-with-query-from-golang-updated"

	updatedResponse, err := lacework.V2.ResourceGroups.UpdateAws(&rgV2Resp.Data)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Succesfully updated resource group \n %+v\n", updatedResponse)

	fmt.Printf("Deleting resource group %s\n", updatedResponse.Data.ResourceGroupGuid)
	err = lacework.V2.ResourceGroups.Delete(updatedResponse.Data.ResourceGroupGuid)
	if err != nil {
		log.Fatal(err)
	}
	println("Successfully deleted resource group")

	props := api.LwAccountResourceGroupProps{
		Description: "All Lacework accounts",
		LwAccounts:  []string{"tech-ally"},
	}

	myResourceGroup := api.NewResourceGroup(
		"resource-group-from-golang",
		api.LwAccountResourceGroup,
		props,
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
	fmt.Printf("Resource Group created: %s", response.Data.ResourceGuid)
}
