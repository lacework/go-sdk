package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

// CheckVNetQuota checks if the current VNet count is approaching or exceeding the quota limit for the given region
// For Azure Agentless only
func CheckVNetQuota(p *Preflight) error {
	region := p.azureConfig.region
	if region == "" {
		region = "eastus" // fallback default
	}
	p.verboseWriter.Write(fmt.Sprintf("Checking VNet quota for region %s\n", region))
	usageClient, err := armnetwork.NewUsagesClient(p.azureConfig.subscriptionID, p.azureConfig.cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create usage client: %w", err)
	}

	pager := usageClient.NewListPager(region, nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get usage page: %w", err)
		}
		for _, usage := range page.Value {
			if usage.Name != nil && usage.Name.Value != nil && *usage.Name.Value == "VirtualNetworks" {
				current := int32(0)
				limit := int32(0)
				if usage.CurrentValue != nil {
					current = int32(*usage.CurrentValue)
				}
				if usage.Limit != nil {
					limit = int32(*usage.Limit)
				}
				if current >= limit {
					p.errors[Agentless] = append(
						p.errors[Agentless],
						fmt.Sprintf("VNet quota limit exceeded in region %s (Current: %d, Limit: %d)", region, current, limit),
					)
				}
				p.verboseWriter.Write(fmt.Sprintf("VNet usage in region %s: %d/%d\n", region, current, limit))
				return nil // Only one VNet quota per region
			}
		}
	}
	return nil
}
