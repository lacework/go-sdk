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

	// setup inline scanner data
	identifierTags := make([]map[string]string, 0)

	inlineScannerData := api.InlineScannerData{
		RegistryType:  api.InlineScannerContainerRegistry.String(),
		LimitNumScan:  "60",
		IdentifierTag: identifierTags,
	}

	inlineScanner := api.NewContainerRegistry(
		"test",
		api.InlineScannerContainerRegistry,
		inlineScannerData,
	)

	// create inline scanner
	inlineScannerCreateResponse, errCreate := lacework.V2.ContainerRegistries.Create(
		inlineScanner,
	)
	if errCreate != nil {
		log.Fatal(errCreate)
	}
	fmt.Printf("Created inline scanner guid: %s\n", inlineScannerCreateResponse.Data.IntgGuid)
	fmt.Printf("Created inline scanner server token: %s\n", inlineScannerCreateResponse.Data.ServerToken.ServerToken)

	// get created inline scanner
	inlineScannerGetResponse, errGet := lacework.V2.ContainerRegistries.GetInlineScanner(
		inlineScannerCreateResponse.Data.IntgGuid,
	)
	if errGet != nil {
		log.Fatal(errGet)
	}
	fmt.Printf("Found inline scanner guid: %s\n", inlineScannerGetResponse.Data.IntgGuid)
	fmt.Printf("Found inline scanner token: %s\n", inlineScannerGetResponse.Data.ServerToken.ServerToken)

	// delete created inline scanner
	errDelete := lacework.V2.ContainerRegistries.Delete(
		inlineScannerCreateResponse.Data.IntgGuid,
	)
	if errDelete != nil {
		log.Fatal(errDelete)
	}
	fmt.Printf("Deleted inline scanner guid: %s\n", inlineScannerCreateResponse.Data.IntgGuid)
}
