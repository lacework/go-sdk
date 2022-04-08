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
	res, err := lacework.Recommendations.Aws.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, rec := range res {
		fmt.Printf("%s: %t \n", rec.ID, rec.State)
	}

	// List all Gcp recommendations
	fmt.Println("\n GCP Recommendations\n ------------------")
	res, err = lacework.Recommendations.Gcp.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, rec := range res {
		fmt.Printf("%s: %t \n", rec.ID, rec.State)
	}

	// List all Azure recommendations
	fmt.Println("\n Azure Recommendations\n ------------------")
	res, err = lacework.Recommendations.Azure.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, rec := range res {
		fmt.Printf("%s: %t \n", rec.ID, rec.State)
	}

	// Enable/Disable a single RecommendationV1
	response, err := lacework.Recommendations.Aws.Patch(api.RecommendationStateV1{"LW_S3_1": "disable"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Recommendations Patched: %v", response.Data[0])

	response, err = lacework.Recommendations.Aws.Patch(api.RecommendationStateV1{"LW_S3_1": "enable"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Recommendations Patched: %v", response.Data)

	// List all Recommendation IDs for a given report type
	reportSchema, err := lacework.Recommendations.Aws.GetReport("AWS_CIS_S3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s Report Schema \n ---------------------------- \n%v\n", "AWS_CIS_S3", reportSchema)

	// Enable/Disable all recommendations of a given report type
	enableAll := api.NewRecommendationV1State(reportSchema, true)

	response, err = lacework.Recommendations.Aws.Patch(enableAll)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nEnabled Recommendations \n ---------------------------- \n")
	for k, v := range enableAll {
		fmt.Printf("%s:%s \n", k, v)
	}
}
