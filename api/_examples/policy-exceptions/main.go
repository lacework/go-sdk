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

	policyID := "lacework-global-39"
	res, err := lacework.V2.PolicyExceptions.List(policyID)
	if err != nil {
		log.Fatal(err)
	}

	for _, exception := range res.Data {
		switch exception.ExceptionID {
		}

		// Output: ExceptionID: [ID]
		fmt.Printf("ExceptionID:%s\n", exception.ExceptionID)
	}

	myPolicyException := api.PolicyException{
		Description: "Exception created by the go-sdk",
		Constraints: []api.PolicyExceptionConstraint{{FieldKey: "accountIds", FieldValues: []string{"*"}}},
	}

	response, err := lacework.V2.PolicyExceptions.Create(policyID, myPolicyException)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Policy Exception created: ID
	fmt.Printf("Policy Exception created: %s", response.Data)

	err = lacework.V2.PolicyExceptions.Delete(policyID, response.Data.ExceptionID)
	if err != nil {
		log.Fatal(err)
	}

	// Output: Policy Exception deleted: ID
	fmt.Printf("Policy Exception deleted: %s", response.Data.ExceptionID)
}
