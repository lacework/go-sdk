package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")))
	if err != nil {
		log.Fatal(err)
	}

	myAwsS3Channel := api.NewAlertChannel("s3-alert-from-golang",
		api.AwsS3AlertChannelType,
		api.AwsS3DataV2{
			Credentials: api.AwsS3Credentials{
				RoleArn:    "arn:aws:iam::1234567890:role/lacework_iam_example_role",
				BucketArn:  "arn:aws:s3:::bucket_name/key_name",
				ExternalID: "0123456789",
			},
		},
	)

	response, err := lacework.V2.AlertChannels.Create(myAwsS3Channel)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Aws S3 alert channel created: THE-INTEGRATION-GUID
	fmt.Printf("Aws S3 alert channel created: %s", response.Data.IntgGuid)
}
