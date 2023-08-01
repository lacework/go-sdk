package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwgenerate/aws_controltower"
)

func controlTowerWithCloudtrail() {
	hcl, err := aws_controltower.NewTerraform("arn:aws:s3:::aws-controltower-logs-0123456789-us-east-1",
		"arn:aws:sns:us-east-1:0123456789:aws-controltower-AllConfigNotifications").Generate()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(hcl)
}

func main() {
	controlTowerWithCloudtrail()
}
