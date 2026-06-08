//go:build preflight

package integration

import (
	"os"
	"testing"

	"github.com/lacework/go-sdk/v2/lwpreflight/aws"
	"github.com/stretchr/testify/assert"
)

func TestPreflightAWS(t *testing.T) {
	const (
		accountID = "159822160952"
		arn       = "arn:aws:iam::159822160952:user/preflight-test"
		userID    = "AKIASKNREUQ4P2U6GSKM"
		name      = "preflight-test"
	)

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	preflight, err := aws.New(aws.Params{
		Agentless:       false,
		Config:          true,
		CloudTrail:      false,
		Region:          "us-west-2",
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	})

	assert.NoError(t, err)

	result, err := preflight.Run()

	assert.NoError(t, err)
	assert.Equal(t, accountID, result.Caller.AccountID)
	assert.Equal(t, arn, result.Caller.ARN)
	assert.Equal(t, userID, result.Caller.UserID)
	assert.Equal(t, name, result.Caller.Name)
	assert.Equal(t, 18, len(result.Details.Regions))
	assert.Contains(t, result.Errors["aws_agentless"], "Required permission missing: ec2:AssociateRouteTable")
	assert.Contains(t, result.Errors["aws_cloudtrail"], "Required permission missing: cloudtrail:AddTags")
}
