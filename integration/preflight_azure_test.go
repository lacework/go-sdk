//go:build preflight

package integration

import (
	"os"
	"testing"

	"github.com/lacework/go-sdk/v2/lwpreflight/azure"
	"github.com/stretchr/testify/assert"
)

func TestPreflightAzure(t *testing.T) {
	const (
		SubscriptionID = "0c21d8f2-6a26-47b1-9d9c-ae935147b344"
		TenantID       = "3a376159-a4aa-4cab-8840-e9b286374a30"
	)

	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	preflight, err := azure.New(azure.Params{
		Agentless:      true,
		Config:         true,
		ActivityLog:    true,
		SubscriptionID: SubscriptionID,
		TenantID:       TenantID,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
	})

	assert.NoError(t, err)

	result, err := preflight.Run()

	assert.NoError(t, err)
	assert.NotEmpty(t, result.Caller.ObjectID)
	assert.False(t, result.Caller.IsAdmin)
	assert.NotEmpty(t, result.Caller.TenantID)
	assert.NotEmpty(t, result.Details.Regions)
	assert.Contains(t, result.Errors["azure_agentless"], "Required permission missing: Microsoft.Compute/virtualMachineScaleSets/read")
	assert.Contains(t, result.Errors["azure_activity_log"], "Required permission missing: Microsoft.Insights/eventtypes/values/read")
	assert.Contains(t, result.Errors["azure_config"], "Required permission missing: Microsoft.Resources/deployments/write")
}
