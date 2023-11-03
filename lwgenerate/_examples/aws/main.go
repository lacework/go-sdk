package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwgenerate/aws"
)

func basic() {
	hcl, err := aws.NewTerraform("us-east-1", false, true, true, true).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(hcl)
}

func customAwsProfile() {
	hcl, err := aws.NewTerraform(
		"us-east-1",
		false,
		true,
		true,
		true,
		aws.WithAwsProfile("mycorp-profile")).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func consolidatedTrailWithSubAccounts() {
	hcl, err := aws.NewTerraform(
		"us-east-1",
		false,
		true,
		true,
		true,
		aws.WithSubaccounts(
			aws.NewAwsSubAccount("subaccount1", "us-east-2"),
			aws.NewAwsSubAccount("subaccount2", "us-east-1"),
		),
		aws.UseConsolidatedCloudtrail(),
		aws.WithAwsProfile("mycorp-profile")).Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n-----\n%s", hcl)
}

func main() {
	basic()
	customAwsProfile()
	consolidatedTrailWithSubAccounts()
}
