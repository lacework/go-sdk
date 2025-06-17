package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

type Details struct {
	Regions []string // Available regions for the subscription
}

func FetchDetails(p *Preflight) error {
	client, err := armsubscriptions.NewClient(p.cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create subscriptions client: %v", err)
	}

	// Get available locations using the pager
	pager := client.NewListLocationsPager(p.subscriptionID, nil)
	regions := make([]string, 0)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return fmt.Errorf("failed to list locations: %v", err)
		}
		for _, location := range page.Value {
			if location.Name != nil {
				regions = append(regions, *location.Name)
			}
		}
	}

	p.details = Details{
		Regions: regions,
	}

	return nil
}
