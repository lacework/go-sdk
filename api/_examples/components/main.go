package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	var (
		OS   string = "linux"
		arch string = "amd64"
	)

	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := lacework.V2.Components.ListComponents(OS, arch)
	if err != nil {
		log.Fatal(err)
	}

	components := resp.Data[0].Components

	for _, component := range components {
		fmt.Println(component)

		resp, err := lacework.V2.Components.ListComponentVersions(1, OS, arch)
		if err != nil {
			log.Fatal(err)
		}

		versions := resp.Data[0].Versions

		fmt.Println(versions)

		resArtifact, err := lacework.V2.Components.FetchComponentArtifact(1, OS, arch, versions[len(versions)-1])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resArtifact.Data[0])
	}

}
