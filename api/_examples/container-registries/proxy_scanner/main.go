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
		api.WithApiV2(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// setup proxy scanner data
	repositoriesLimit := make([]string, 0)
	tagsLimit := make([]string, 0)
	labelLimit := make([]map[string]string, 0)
	limitByImg := 5

	proxyScannerData := api.ProxyScannerData{
		RegistryType: api.ProxyScannerContainerRegistry.String(),
		LimitNumImg:  limitByImg,
		LimitByRep:   repositoriesLimit,
		LimitByTag:   tagsLimit,
		LimitByLabel: labelLimit,
	}

	proxyScanner := api.NewContainerRegistry(
		"test",
		api.ProxyScannerContainerRegistry,
		proxyScannerData,
	)

	// create proxy scanner
	proxyScannerCreateResponse, errCreate := lacework.V2.ContainerRegistries.Create(
		proxyScanner,
	)
	if errCreate != nil {
		log.Fatal(errCreate)
	}
	fmt.Printf("Created proxy scanner guid: %s", proxyScannerCreateResponse.Data.IntgGuid)
	fmt.Printf("Created inline scanner server token: %s", proxyScannerCreateResponse.Data.ServerToken.ServerToken)

	// get created proxy scanner
	proxyScannerGetResponse, errGet := lacework.V2.ContainerRegistries.GetProxyScanner(
		proxyScannerCreateResponse.Data.IntgGuid,
	)
	if errGet != nil {
		log.Fatal(errGet)
	}
	fmt.Printf("Found proxy scanner guid: %s", proxyScannerGetResponse.Data.IntgGuid)

	// delete created proxy scanner
	errDelete := lacework.V2.ContainerRegistries.Delete(
		proxyScannerCreateResponse.Data.IntgGuid,
	)
	if errDelete != nil {
		log.Fatal(errDelete)
	}
	fmt.Printf("Deleted proxy scanner guid: %s", proxyScannerCreateResponse.Data.IntgGuid)
}
