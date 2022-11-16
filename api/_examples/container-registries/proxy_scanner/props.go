package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	// no policy initguid
	noPolicyInitGuid := "SNIFFTES_418DDCB85552330202BEDB52FCA60DBBD4C54BE3410974F"

	// get created proxy scanner with no policy assigned
	noPolicyGetResponse, errGet := lacework.V2.ContainerRegistries.GetProxyScanner(
		noPolicyInitGuid,
	)
	if errGet != nil {
		log.Fatal(errGet)
	}
	fmt.Printf("Found proxy scanner guid: %s\n", noPolicyGetResponse.Data.IntgGuid)
	fmt.Printf("Found proxy scanner token: %s\n", noPolicyGetResponse.Data.ServerToken.ServerToken)
	fmt.Printf("Found proxy scanner props tags: %s\n", noPolicyGetResponse.Data.Props.Tags)
	nop := noPolicyGetResponse.Data.Props.PolicyEvaluation
	if nop != nil {
		fmt.Printf("Found inline policy evaluation: %s\n", strconv.FormatBool(nop.Evaluate))
		if nop.Evaluate {
			for _, nog := range nop.PolicyGuids {
				fmt.Printf("Found inline policy guid: %s\n", nog)
			}
		}
	}

	// policy initguid
	policyInitGuid := "SNIFFTES_2FB4D9BB82B8A557E7B8EC791872AB850F0EB70E275FE6E"

	// get created proxy scanner with no policy assigned
	policyGetResponse, errGet := lacework.V2.ContainerRegistries.GetInlineScanner(
		policyInitGuid,
	)
	if errGet != nil {
		log.Fatal(errGet)
	}
	fmt.Printf("Found proxy scanner guid: %s\n", policyGetResponse.Data.IntgGuid)
	fmt.Printf("Found proxy scanner token: %s\n", policyGetResponse.Data.ServerToken.ServerToken)
	fmt.Printf("Found proxy scanner props tags: %s\n", policyGetResponse.Data.Props.Tags)
	p := policyGetResponse.Data.Props.PolicyEvaluation
	if p != nil {
		fmt.Printf("Found inline policy evaluation: %s\n", strconv.FormatBool(p.Evaluate))
		if p.Evaluate {
			for _, g := range p.PolicyGuids {
				fmt.Printf("Found inline policy guid: %s\n", g)
			}
		}
	}

}
