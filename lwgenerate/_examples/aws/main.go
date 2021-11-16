package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwgenerate"
)

func main() {
	hcl, err := lwgenerate.NewAwsTerraform("us-east-1", true, true).Generate()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(hcl)
}
