package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwgenerate/azure"
)

func ActivityLogWithoutConfig() {
	hcl, err := azure.NewTerraform(false, true, true).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(hcl)
}

func ActivityLogWithConfig() {
	hcl, err := azure.NewTerraform(true, true, true).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func ConfigWithoutActivityLog() {
	hcl, err := azure.NewTerraform(true, false, true).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func CustomActiveDirectory() {
	hcl, err := azure.NewTerraform(true, true, false,
		azure.WithConfigIntegrationName("Test Config Rename"),
		azure.WithActivityLogIntegrationName("Test Activity Log Rename"),
		azure.WithAdApplicationPassword("AD-Test-Password"),
		azure.WithAdServicePrincipalId("AD-Test-Principal-ID"),
		azure.WithAdApplicationId("AD-Test-Application-ID"),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func ActivityLogWithSubscriptionsList() {
	testIds := []string{"test-id-1", "test-id-2", "test-id-3"}
	hcl, err := azure.NewTerraform(false, true, true,
		azure.WithSubscriptionIds(testIds),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func ConfigWithSubscriptionsList() {
	testIds := []string{"test-id-1", "test-id-2", "test-id-3"}
	hcl, err := azure.NewTerraform(true, false, true,
		azure.WithSubscriptionIds(testIds),
	).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func main() {
	ActivityLogWithoutConfig()
	ActivityLogWithConfig()
	ConfigWithoutActivityLog()
	CustomActiveDirectory()
	ActivityLogWithSubscriptionsList()
	ConfigWithSubscriptionsList()
}
