package integration

import (
	"os"
	"testing"

	"github.com/lacework/go-sdk/v2/lwpreflight/aws"
	"github.com/stretchr/testify/assert"
)

const (
	accountID = "441565585060"
	arn       = "arn:aws:iam::441565585060:user/preflight-test"
	userID    = "AIDAWNT24D2SO25OXHPQE"
	name      = "preflight-test"
)

func TestPreflightAWS(t *testing.T) {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	preflight, err := aws.New(aws.Params{
		Agentless:       true,
		Config:          true,
		CloudTrail:      true,
		Region:          "us-west-2",
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	})

	assert.NoError(t, err)

	result, err := preflight.Run()

	assert.NoError(t, err)
	assert.Equal(t, result.Caller.AccountID, accountID)
	assert.Equal(t, result.Caller.ARN, arn)
	assert.Equal(t, result.Caller.UserID, userID)
	assert.Equal(t, result.Caller.Name, name)
	assert.Equal(t, len(result.Details.Regions), 18)
	assert.Contains(t, result.Errors["aws_agentless"], "Required permission missing: ec2:AssociateRouteTable")
	assert.Contains(t, result.Errors["aws_config"], "Required permission missing: cloudformation:CreateStack")
	assert.Contains(t, result.Errors["aws_cloudtrail"], "Required permission missing: cloudtrail:AddTags")
}
