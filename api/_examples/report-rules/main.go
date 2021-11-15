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

	res, err := lacework.V2.ReportRules.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, rule := range res.Data {
		support := "Unsupported"
		switch rule.Type {
		case api.ReportRuleEventType:
			support = "Supported"
		}

		// Output: GUID:RULE_TYPE:[Supported|Unsupported]
		fmt.Printf("%s:%s:%s\n", rule.Guid, rule.Type, support)
	}

	rule := api.ReportRuleConfig{
		EmailAlertChannels: []string{"TECHALLY_2F0C086E17AB64BEC84F4A5FF8A3F068CF2CE15847BCBCA"},
		Description:        "This is a test report rule",
		Severities:         api.ReportRuleSeverities{api.ReportRuleSeverityCritical},
		ResourceGroups:     []string{"TECHALLY_69B96544DE073FA270C1AFE4CA7FE74DA784C3A1B8EC77B"},
		NotificationTypes:  api.ReportRuleNotifications{api.AllGcpReportRuleNotifications, api.AllAwsReportRuleNotifications},
	}

	myReportRule, err := api.NewReportRule("MyTestReportRule",
		rule,
	)

	response, err := lacework.V2.ReportRules.Create(myReportRule)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Report Rule created: GUID
	fmt.Printf("Report Rule Group created: %s \n", response.Data.Guid)

	err = lacework.V2.ReportRules.Delete(response.Data.Guid)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Report Rule deleted: GUID
	fmt.Printf("Report Rule Group deleted: %s", response.Data.Guid)
}
