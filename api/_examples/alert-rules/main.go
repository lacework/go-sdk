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

	res, err := lacework.V2.AlertRules.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, rule := range res.Data {
		support := "Unsupported"
		switch rule.Type {
		case api.AlertRuleEventType:
			support = "Supported"
		}

		// Output: GUID:RULE_TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", rule.Guid, rule.Type, support)
	}

	rule := api.AlertRuleConfig{
		Channels:        []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
		Description:     "This is a test alert rule",
		Severities:      api.AlertRuleSeverities{api.AlertRuleSeverityHigh},
		ResourceGroups:  []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
		EventCategories: []string{"Compliance"},
	}

	myAlertRule := api.NewAlertRule("MyTestAlertRule",
		rule,
	)

	response, err := lacework.V2.AlertRules.Create(myAlertRule)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Alert Rule created: GUID
	fmt.Printf("Alert Rule Group created: %s", response.Data.Guid)
}
