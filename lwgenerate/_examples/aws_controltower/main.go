package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwgenerate/aws_controltower"
)

func controlTowerWithCloudtrail() {
	hcl, err := aws_controltower.NewTerraform("us-east-1", "s3Arn", "Arn").Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(hcl)
}

func main() {
	controlTowerWithCloudtrail()
}
