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
		SubscriptionID = "1fe75302-1906-45bc-bdc1-79b76799dd74"
	)

	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	tenantID := os.Getenv("AZURE_TENANT_ID")

	preflight, err := azure.New(azure.Params{
		Agentless:      true,
		Config:         true,
		ActivityLog:    true,
		SubscriptionID: SubscriptionID,
		TenantID:       tenantID,
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
	assert.Contains(t, result.Errors["azure_activity_log"], "Required permission missing: Microsoft.Insights/diagnosticSettings/delete")
	assert.Contains(t, result.Errors["azure_config"], "Required permission missing: Microsoft.Authorization/roleAssignments/write")
}
