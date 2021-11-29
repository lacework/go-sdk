package cmd

import (
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/stretchr/testify/assert"
)

func toggleNonInteractive() {
	cli.nonInteractive = !cli.nonInteractive
}

func TestGenerateMostBasicArgs(t *testing.T) {
	toggleNonInteractive()
	defer toggleNonInteractive()

	data := aws.GenerateAwsTfConfigurationArgs{}
	data.Cloudtrail = true
	data.Config = true
	data.AwsRegion = "us-east-2"
	output := "/tmp"
	err := promptAwsGenerate(&data, &aws.ExistingIamRoleDetails{}, &output)

	assert.Nil(t, err)
}

func TestMissingValidEntityToConfigure(t *testing.T) {
	toggleNonInteractive()
	defer toggleNonInteractive()

	data := aws.GenerateAwsTfConfigurationArgs{}
	output := "/tmp"
	err := promptAwsGenerate(&data, &aws.ExistingIamRoleDetails{}, &output)
	assert.Error(t, err)
	assert.Equal(t, "Must enable cloudtrail or config!", err.Error())
}
