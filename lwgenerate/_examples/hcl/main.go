package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwgenerate"
)

func exampleModule() {
	data, err := lwgenerate.NewModule("servers", "./app-clusters",
		lwgenerate.HclModuleWithAttributes(map[string]interface{}{
			"servers": 5,
		}),
	).ToBlock()

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Print(lwgenerate.CreateHclStringOutput(lwgenerate.CombineHclBlocks(data)))
}

func exampleProvider() {
	attrs := map[string]interface{}{"region": "us-east-2"}
	data, err := lwgenerate.NewProvider("aws", lwgenerate.HclProviderWithAttributes(attrs)).ToBlock()

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Print(lwgenerate.CreateHclStringOutput(lwgenerate.CombineHclBlocks(data)))
}

func exampleRequiredProviders() {
	data, err := lwgenerate.CreateRequiredProviders(
		lwgenerate.NewRequiredProvider("aws",
			lwgenerate.HclRequiredProviderWithSource("hashicorp/aws"),
			lwgenerate.HclRequiredProviderWithVersion("~> 0.1")),
		lwgenerate.NewRequiredProvider("kubernetes",
			lwgenerate.HclRequiredProviderWithSource("hashicorp/kubernetes"),
			lwgenerate.HclRequiredProviderWithVersion(">= 2.0.0")))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Print(lwgenerate.CreateHclStringOutput(lwgenerate.CombineHclBlocks(data)))
}

func exampleSimpleTraversal() {
	data, err := lwgenerate.NewModule("mymodule", "./mymodule",
		lwgenerate.HclModuleWithAttributes(map[string]interface{}{
			"someparameter": lwgenerate.CreateSimpleTraversal([]string{"module", "myothermodule", "propertyvalue"}),
		}),
	).ToBlock()

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Print(lwgenerate.CreateHclStringOutput(lwgenerate.CombineHclBlocks(data)))
}

func exampleGenericHclBlock() {
	data, err := lwgenerate.HclCreateGenericBlock(
		"mycustomresource",
		[]string{"label1", "label2"},
		map[string]interface{}{
			"property1": 5,
			"property2": "test",
			"property3": map[string]interface{}{
				"one": 1,
				"two": "a",
			},
			"property4": []string{"one", "two", "three"},
			"property5": []interface{}{"one", 2, true},
		},
	)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Print(lwgenerate.CreateHclStringOutput(lwgenerate.CombineHclBlocks(data)))
}

func main() {
	exampleRequiredProviders()
	exampleProvider()
	exampleModule()
	exampleSimpleTraversal()
	exampleGenericHclBlock()
}
