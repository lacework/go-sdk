package cmd

import (
	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUniqueHostEvalGUID(t *testing.T) {
	expectedEvalGUID := "12345"
	actualEvalGUID := getUniqueHostEvalGUID(mockVulnerabilitiesHostResponse(expectedEvalGUID))

	assert.Equal(t, expectedEvalGUID, actualEvalGUID)
}

func mockVulnerabilitiesHostResponse(evalGUID string) api.VulnerabilitiesHostResponse {
	return api.VulnerabilitiesHostResponse{
		Data: []api.VulnerabilityHost{
			{
				EvalGUID:  "54321",
				StartTime: time.Now().AddDate(0, 0, -2),
			},
			{
				EvalGUID:  "54321",
				StartTime: time.Now().AddDate(0, 0, -2),
			},
			{
				EvalGUID:  evalGUID,
				StartTime: time.Now(),
			},
			{
				EvalGUID:  "98765",
				StartTime: time.Now().AddDate(0, 0, -1),
			},
			{
				EvalGUID:  evalGUID,
				StartTime: time.Now(),
			},
		},
	}
}
