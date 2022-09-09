package main

import (
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient("",
		api.WithApiKeys("", ""),
		api.WithApiV2(),
	)
	if err != nil {
		log.Fatal(err)
	}

	filter := api.SearchFilter{
		Filters: []api.Filter{{
			Expression: "eq",
			Field:      "eventType",
			Value:      "SuspiciousFile",
		}},
	}

	res, err := lacework.V2.Events.Search(filter)
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range res.Data {
		fmt.Printf("%d: %s: %s\n", event.Id, event.EventType, event.Severity())
	}
}
