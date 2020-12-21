package main

import (
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient("account", api.WithApiKeys("KEY", "SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	myAwsS3Channel := api.NewAwsS3AlertChannel("s3-alert-from-golang",
		api.AwsS3ChannelData{
			Credentials: api.AwsS3Creds{
				RoleArn:    "arn:aws:iam::1234567890:role/lacework_iam_example_role",
				BucketArn:  "arn:aws:s3:::bucket_name/key_name",
				ExternalID: "0123456789",
			}
		},
	)

	response, err := lacework.Integrations.CreateAwsS3AlertChannel(myAwsS3Channel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Aws S3 alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Aws S3 alert channel created: %s", response.Data[0].IntgGuid)
}
