package cmd

import (
	"regexp"
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

func TestArnRegex(t *testing.T) {
	ok, _ := regexp.MatchString(AwsArnRegex, "invalidarnstring")
	assert.False(t, ok, "aws arn cannot be an arbitrary string")

	ok, _ = regexp.MatchString(AwsArnRegex, "$#@###$#%^")
	assert.False(t, ok, "aws arn cannot be an arbitrary string")

	ok, _ = regexp.MatchString(AwsArnRegex, "arn:partition:service:region:account-id:resource-id")
	assert.True(t, ok, "aws arn do not have to have resource-type supplied")

	ok, _ = regexp.MatchString(AwsArnRegex, "arn:partition:service:region:account-id:resource-type/resource-id")
	assert.True(t, ok, "aws arn is valid with both /resource-id or :resource-id")

	ok, _ = regexp.MatchString(AwsArnRegex, "arn:partition:service:region:account-id:resource-type:resource-id")
	assert.True(t, ok, "aws arn is valid with both /resource-id or :resource-id")

	ok, _ = regexp.MatchString(AwsArnRegex, "arn:partition:service:region:account-id:resource-type/*")
	assert.True(t, ok, "aws arn supports wildcard syntax")
}

func TestAwsRegionRegex(t *testing.T) {
	ok, _ := regexp.MatchString(AwsRegionRegex, "invalidarnstring")
	assert.False(t, ok, "aws region cannot be an arbitrary string")

	ok, _ = regexp.MatchString(AwsRegionRegex, "us-gov-east-1")
	assert.False(t, ok, "aws gov cloud regions not currently supported")

	ok, _ = regexp.MatchString(AwsRegionRegex, "us-east-1")
	assert.True(t, ok, "aws region us-east-1 is valid")

	ok, _ = regexp.MatchString(AwsRegionRegex, "ap-northeast-1")
	assert.True(t, ok, "aws region ap-norteast-1 is valid")
}

func TestAwsSubAccountValidation(t *testing.T) {
	ret := validateAwsSubAccounts([]string{"invalid"})
	assert.Error(t, ret, "subaccount can't be arbitrary string")

	ret = validateAwsSubAccounts([]string{"profilename:us-east-1"})
	assert.Nil(t, ret, "should be valid")

	ret = validateAwsSubAccounts([]string{"profilename:us-east-1", "profile2:ap-southeast-2"})
	assert.Nil(t, ret, "should be valid")

	ret = validateAwsSubAccounts([]string{"profilename:us-east-1", "invalid"})
	assert.Error(t, ret, "should can't be arbitrary string")
}
