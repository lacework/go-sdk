package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
	)
	if err != nil {
		log.Fatal(err)
	}

	// List all Aws recommendations
	fmt.Println("\n AWS Recommendations\n ------------------")
	res, err := lacework.Recommendations.AwsList()
	if err != nil {
		log.Fatal(err)
	}

	for _, rec := range res {
		fmt.Printf("%s: %t \n", rec.ID, rec.State)
	}

	// List all Gcp recommendations
	fmt.Println("\n GCP Recommendations\n ------------------")
	res, err = lacework.Recommendations.GcpList()
	if err != nil {
		log.Fatal(err)
	}

	for _, rec := range res {
		fmt.Printf("%s: %t \n", rec.ID, rec.State)
	}

	// List all Azure recommendations
	fmt.Println("\n Azure Recommendations\n ------------------")
	res, err = lacework.Recommendations.AzureList()
	if err != nil {
		log.Fatal(err)
	}

	for _, rec := range res {
		fmt.Printf("%s: %t \n", rec.ID, rec.State)
	}

	// Enable/Disable a single RecommendationV1
	response, err := lacework.Recommendations.PatchAws(api.RecommendationStateV1{"LW_S3_1": "disable"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Recommendations Patched: %v", response.Data[0])

	response, err = lacework.Recommendations.PatchAws(api.RecommendationStateV1{"LW_S3_1": "enable"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Recommendations Patched: %v", response.Data)
}
