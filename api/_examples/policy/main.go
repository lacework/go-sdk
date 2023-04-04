package policy

import (
	"fmt"
	"log"
	"os"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwseverity"
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

	state := true
	policyList := api.BulkUpdatePolicies{{
		PolicyID: "lacework-global-39",
		Enabled:  &state,
		Severity: lwseverity.High.String(),
	}}
	res, err := lacework.V2.Policy.UpdateMany(policyList)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
