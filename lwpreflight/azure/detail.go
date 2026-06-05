package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

type Details struct {
	Regions []string // Available regions for the subscription
}

func FetchDetails(p *Preflight) error {
	client, err := armsubscriptions.NewClient(p.azureConfig.cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create subscriptions client: %v", err)
	}

	p.verboseWriter.Write("Discovering available regions")

	// Get available locations using the pager
	pager := client.NewListLocationsPager(p.azureConfig.subscriptionID, nil)
	regions := make([]string, 0)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return fmt.Errorf("failed to list locations: %v", err)
		}

		for _, location := range page.Value {
			if location.Name == nil {
				continue
			}

			// filter out edge zones
			if location.Type != nil && *location.Type != armsubscriptions.LocationTypeRegion {
				continue
			}

			// filter out non-physical regions
			if location.Metadata != nil && location.Metadata.RegionType != nil {
				if *location.Metadata.RegionType != armsubscriptions.RegionTypePhysical {
					continue
				}
			}

			// filter out Microsoft-internal staging/canary regions (e.g. eastusstg, centraluseuap).
			// Note: we intentionally do NOT filter on RegionCategoryRecommended — Azure marks
			// fully functional regions like westus, eastus2 and westus3 as "Other" (CAD-2132).
			if strings.HasSuffix(*location.Name, "stg") || strings.HasSuffix(*location.Name, "euap") {
				continue
			}

			// add to available regions
			regions = append(regions, *location.Name)
		}
	}

	p.details = Details{
		Regions: regions,
	}

	return nil
}
