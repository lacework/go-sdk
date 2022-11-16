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
	noPolicyInitGuid := "SNIFFTES_E7B0664F6EDC1A7BAABB001DF1024537CBBD580E5AC1995"

	// get created inline scanner with no policy assigned
	noPolicyGetResponse, errGet := lacework.V2.ContainerRegistries.GetInlineScanner(
		noPolicyInitGuid,
	)
	if errGet != nil {
		log.Fatal(errGet)
	}
	fmt.Printf("Found inline scanner guid: %s\n", noPolicyGetResponse.Data.IntgGuid)
	fmt.Printf("Found inline scanner token: %s\n", noPolicyGetResponse.Data.ServerToken.ServerToken)
	fmt.Printf("Found inline scanner props tags: %s\n", noPolicyGetResponse.Data.Props.Tags)
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
	policyInitGuid := "SNIFFTES_B8D069F9943F043D22582108DB4D0F8FED8024D035F94DD"

	// get created inline scanner with no policy assigned
	policyGetResponse, errGet := lacework.V2.ContainerRegistries.GetInlineScanner(
		policyInitGuid,
	)
	if errGet != nil {
		log.Fatal(errGet)
	}
	fmt.Printf("Found inline scanner guid: %s\n", policyGetResponse.Data.IntgGuid)
	fmt.Printf("Found inline scanner token: %s\n", policyGetResponse.Data.ServerToken.ServerToken)
	fmt.Printf("Found inline scanner props tags: %s\n", policyGetResponse.Data.Props.Tags)
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
