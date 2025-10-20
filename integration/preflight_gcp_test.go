//go:build preflight

package integration

import (
	"os"
	"testing"

	"github.com/lacework/go-sdk/v2/lwpreflight/gcp"
	"github.com/stretchr/testify/assert"
)

func TestPreflightGCP(t *testing.T) {
	const (
		email  = "preflight-test@abc-demo-project-123.iam.gserviceaccount.com"
		userID = "110811735245298771692"
	)

	credentialsJSON := os.Getenv("GOOGLE_CREDENTIALS")

	preflight, err := gcp.New(gcp.Params{
		Agentless:       true,
		Config:          true,
		AuditLog:        true,
		Region:          "us-west2",
		ProjectID:       "abc-demo-project-123",
		CredentialsJSON: credentialsJSON,
	})

	assert.NoError(t, err)

	result, err := preflight.Run()

	assert.NoError(t, err)
	assert.Equal(t, email, result.Caller.Email)
	assert.Equal(t, userID, result.Caller.UserID)
	assert.GreaterOrEqual(t, len(result.Details.SchedulerRegions), 23)
	assert.Contains(t, result.Errors["gcp_agentless"], "Required permission missing: cloudscheduler.jobs.create")
	assert.Contains(t, result.Errors["gcp_audit_log"], "Required permission missing: compute.projects.get")
	assert.Contains(t, result.Errors["gcp_config"], "Required permission missing: iam.serviceAccountKeys.create")
}
