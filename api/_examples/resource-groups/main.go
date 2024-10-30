package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/v2/api"
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
		resourceGuid = account.ResourceGroupGuid

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
		case api.MachineResourceGroup.String():
			support = "Supported"
		case api.OciResourceGroup.String():
			support = "Supported"
		case api.KubernetesResourceGroup.String():
			support = "Supported"
		}

		// Output: RESOURCE_GROUP_GUID:RESOURCE_TYPE:[Supported|Unsupported]
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

	myResourceGroupWithQuery := api.NewResourceGroup(
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
	rgV2Resp.Data.Name = "resource-group-with-query-from-golang-updated"

	updatedResponse, err := lacework.V2.ResourceGroups.Update(&rgV2Resp.Data)

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
}
