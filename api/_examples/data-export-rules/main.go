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

	rule := api.DataExportRule{
		Filter: api.DataExportRuleFilter{
			Name:            "export-rule-created-by-golang",
			ProfileVersions: []string{"V1"},
		},
		Type: "Dataexport",
		IDs:  []string{"TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A"},
	}

	createResponse, err := lacework.V2.DataExportRules.Create(rule)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created new rule: %s\n", createResponse.Data.ID)

	getResponse, _ := lacework.V2.DataExportRules.Get(createResponse.Data.ID)
	fmt.Printf("Retrieved rule: %s\n", getResponse.Data.ID)

	createResponse.Data.Filter.Name = "export-rule-created-by-golang-updated"

	updateResponse, err := lacework.V2.DataExportRules.Update(createResponse.Data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated rule: %s\n", updateResponse.Data.ID)

	searchResponse, err := lacework.V2.DataExportRules.Search(api.SearchFilter{
		Filters: []api.Filter{{
			Expression: "like",
			Field:      "filters.name",
			Value:      "export-rule-",
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, rule := range searchResponse.Data {
		fmt.Printf("Found rule: %s\n", rule.ID)
	}

	listResponse, _ := lacework.V2.DataExportRules.List()
	for _, rule := range listResponse.Data {
		fmt.Printf("Listing rule: %s\n", rule.ID)
	}

	err = lacework.V2.DataExportRules.Delete(createResponse.Data.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted rule: %s\n", createResponse.Data.ID)
}
