package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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
	err := promptAwsGenerate(&data, &aws.AwsGenerateCommandExtraState{Output: "/tmp"})

	assert.Nil(t, err)
}

func TestMissingValidEntityToConfigure(t *testing.T) {
	toggleNonInteractive()
	defer toggleNonInteractive()

	data := aws.GenerateAwsTfConfigurationArgs{}
	err := promptAwsGenerate(&data, &aws.AwsGenerateCommandExtraState{Output: "/tmp"})
	assert.Nil(t, err)
	err = data.Validate()
	assert.Error(t, err)
	assert.Equal(t, "Agentless, CloudTrail or Config integration must be enabled", err.Error())
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

func TestGenerationCache(t *testing.T) {
	t.Run("extra state shouldn't be written if empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		extraState := &aws.AwsGenerateCommandExtraState{}
		writeExtraStateCache(extraState)
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAwsExtraStateKey)))
	})
	t.Run("extra state should be written if not empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		extraState := &aws.AwsGenerateCommandExtraState{Output: "/tmp"}
		writeExtraStateCache(extraState)
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAwsExtraStateKey)))
	})
	t.Run("iac params should not be cached when empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		args := aws.GenerateAwsTfConfigurationArgs{}
		writeArgsCache(&args)
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAwsArgsKey)))
	})
	t.Run("iac params should be cached when not empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		args := aws.GenerateAwsTfConfigurationArgs{AwsRegion: "us-east-2", Agentless: true}
		writeArgsCache(&args)
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAwsArgsKey)))
	})
}
