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
		fmt.Println(err.Error())
		os.Exit(1)
	}

	blocks, err := lwgenerate.CombineHclBlocks(data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(lwgenerate.CreateHclStringOutput(blocks))
}

func exampleProvider() {
	attrs := map[string]interface{}{"region": "us-east-2"}
	data, err := lwgenerate.NewProvider("aws", lwgenerate.HclProviderWithAttributes(attrs)).ToBlock()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	blocks, err := lwgenerate.CombineHclBlocks(data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(lwgenerate.CreateHclStringOutput(blocks))
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
		fmt.Println(err.Error())
		os.Exit(1)
	}

	blocks, err := lwgenerate.CombineHclBlocks(data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(lwgenerate.CreateHclStringOutput(blocks))
}

func exampleSimpleTraversal() {
	data, err := lwgenerate.NewModule("mymodule", "./mymodule",
		lwgenerate.HclModuleWithAttributes(map[string]interface{}{
			"someparameter": lwgenerate.CreateSimpleTraversal([]string{"module", "myothermodule", "propertyvalue"}),
		}),
	).ToBlock()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	blocks, err := lwgenerate.CombineHclBlocks(data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(lwgenerate.CreateHclStringOutput(blocks))
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
		},
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	blocks, err := lwgenerate.CombineHclBlocks(data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(lwgenerate.CreateHclStringOutput(blocks))
}

func main() {
	exampleRequiredProviders()
	exampleProvider()
	exampleModule()
	exampleSimpleTraversal()
	exampleGenericHclBlock()
}
