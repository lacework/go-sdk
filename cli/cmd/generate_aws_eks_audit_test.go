package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/aws_eks_audit"

	"github.com/stretchr/testify/assert"
)

func toggleEksNonInteractive() {
	cli.noCache = !cli.noCache
	cli.nonInteractive = !cli.nonInteractive
}

func TestGenerateMostBasicEksArgs(t *testing.T) {
	toggleEksNonInteractive()
	defer toggleEksNonInteractive()

	data := aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs{}
	err := promptAwsEksAuditGenerate(
		&data,
		&aws_eks_audit.ExistingCrossAccountIamRoleDetails{},
		&AwsEksAuditGenerateCommandExtraState{Output: "/tmp"},
	)

	assert.Nil(t, err)
}

func TestAwsEksAuditRegionRegex(t *testing.T) {
	ok, _ := regexp.MatchString(AwsEksAuditRegionRegex, "invalidarnstring")
	assert.False(t, ok, "aws region cannot be an arbitrary string")

	ok, _ = regexp.MatchString(AwsEksAuditRegionRegex, "us-gov-east-1")
	assert.False(t, ok, "aws gov cloud regions not currently supported")

	ok, _ = regexp.MatchString(AwsEksAuditRegionRegex, "us-east-1")
	assert.True(t, ok, "aws region us-east-1 is valid")

	ok, _ = regexp.MatchString(AwsEksAuditRegionRegex, "ap-northeast-1")
	assert.True(t, ok, "aws region ap-norteast-1 is valid")
}

func TestEksGenerationCache(t *testing.T) {
	t.Run("extra state shouldn't be written if empty", func(t *testing.T) {

		dir, err := os.MkdirTemp("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {

			}
		}(dir)
		cli.InitCache(dir)

		extraState := &AwsEksAuditGenerateCommandExtraState{}
		extraState.writeCache()
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir,
			CachedAssetAwsEksAuditExtraState)))
	})
	t.Run("extra state should be written if not empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {

			}
		}(dir)
		cli.InitCache(dir)

		extraState := AwsEksAuditGenerateCommandExtraState{Output: "/tmp"}
		extraState.writeCache()
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAssetAwsEksAuditExtraState)))
	})
	t.Run("iac params should not be cached when empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {

			}
		}(dir)
		cli.InitCache(dir)

		args := aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs{}
		writeAwsEksAuditGenerationArgsCache(&args)
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir,
			CachedAssetAwsEksAuditIacParams)))
	})
	t.Run("iac params should be cached when not empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {

			}
		}(dir)
		cli.InitCache(dir)

		args := aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs{LaceworkProfile: "default"}
		writeAwsEksAuditGenerationArgsCache(&args)
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAssetAwsEksAuditIacParams)))
	})
}
