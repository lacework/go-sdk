package cmd

import (
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/stretchr/testify/assert"
)

func toggleNonInteractive() {
	cli.noCache = !cli.noCache
	cli.nonInteractive = !cli.nonInteractive
}

func TestGenerateMostBasicArgs(t *testing.T) {
	toggleNonInteractive()
	defer toggleNonInteractive()

	data := aws.GenerateAwsTfConfigurationArgs{}
	data.Cloudtrail = true
	data.Config = true
	data.AwsRegion = "us-east-2"
	err := promptAwsGenerate(&data, &aws.ExistingIamRoleDetails{}, &AwsGenerateCommandExtraState{Output: "/tmp"})

	assert.Nil(t, err)
}

func TestMissingValidEntityToConfigure(t *testing.T) {
	toggleNonInteractive()
	defer toggleNonInteractive()

	data := aws.GenerateAwsTfConfigurationArgs{}
	err := promptAwsGenerate(&data, &aws.ExistingIamRoleDetails{}, &AwsGenerateCommandExtraState{Output: "/tmp"})
	assert.Error(t, err)
	assert.Equal(t, "must enable cloudtrail or config", err.Error())
}
