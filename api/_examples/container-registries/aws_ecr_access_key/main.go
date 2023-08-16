package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
	)
	if err != nil {
		log.Fatal(err)
	}

	ecrData := api.AwsEcrAccessKeyData{
		RegistryDomain:   "example.domain",
		LimitNumImg:      5,
		NonOSPackageEval: true,
		AccessKeyCredentials: api.AwsEcrAccessKeyCredentials{
			AccessKeyID:     "example-key-id",
			SecretAccessKey: "ABC123",
		},
	}

	awsEcrContainerRegistry := api.NewContainerRegistry(
		"aws-ecr-access-key-from-golang",
		api.AwsEcrContainerRegistry,
		ecrData,
	)

	ecrResponse, err := lacework.V2.ContainerRegistries.Create(awsEcrContainerRegistry)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Aws Ecr Container Registry created: THE-INTEGRATION-GUID
	fmt.Printf("Container Registry created: %s", ecrResponse.Data.IntgGuid)
}
