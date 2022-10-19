package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/gcp"

	"github.com/stretchr/testify/assert"
)

func toggleGcpNonInteractive() {
	cli.noCache = !cli.noCache
	cli.nonInteractive = !cli.nonInteractive
}

func TestGenerateMostBasicGcpArgs(t *testing.T) {
	toggleGcpNonInteractive()
	defer toggleGcpNonInteractive()

	data := gcp.GenerateGcpTfConfigurationArgs{}
	data.AuditLog = true
	data.Configuration = true
	err := promptGcpGenerate(&data, &gcp.ExistingServiceAccountDetails{}, &GcpGenerateCommandExtraState{Output: "/tmp"})

	assert.Nil(t, err)
}

func TestMissingValidGcpEntityToConfigure(t *testing.T) {
	toggleGcpNonInteractive()
	defer toggleGcpNonInteractive()

	data := gcp.GenerateGcpTfConfigurationArgs{}
	err := promptGcpGenerate(&data, &gcp.ExistingServiceAccountDetails{}, &GcpGenerateCommandExtraState{Output: "/tmp"})
	assert.Error(t, err)
	assert.Equal(t, "must enable audit log or configuration", err.Error())
}

func TestGcpBucketRegionRegex(t *testing.T) {
	ok, _ := regexp.MatchString(GcpRegionRegex, "invalidstring")
	assert.False(t, ok, "gcp region cannot be an arbitrary string")

	ok, _ = regexp.MatchString(GcpRegionRegex, "us-east1")
	assert.True(t, ok, "gcp region us-east1 is valid")

	ok, _ = regexp.MatchString(GcpRegionRegex, "asia-northeast1")
	assert.True(t, ok, "gcp region asia-northeast1 is valid")
}

func TestGcpGenerationCache(t *testing.T) {
	t.Run("extra state shouldn't be written if empty", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		extraState := &GcpGenerateCommandExtraState{}
		extraState.writeCache()
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAssetGcpExtraState)))
	})
	t.Run("extra state should be written if not empty", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		extraState := GcpGenerateCommandExtraState{Output: "/tmp"}
		extraState.writeCache()
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAssetGcpExtraState)))
	})
	t.Run("iac params should not be cached when empty", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		args := gcp.GenerateGcpTfConfigurationArgs{}
		writeGcpGenerationArgsCache(&args)
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedGcpAssetIacParams)))
	})
	t.Run("iac params should be cached when not empty", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		args := gcp.GenerateGcpTfConfigurationArgs{ServiceAccountCredentials: "default"}
		writeGcpGenerationArgsCache(&args)
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedGcpAssetIacParams)))
	})
}

func TestValidateGcpProjectId(t *testing.T) {
	tests := []struct {
		input string
		error bool
		desc  string
	}{
		{
			"valid123-1",
			false,
			"Valid project ID",
		},
		{
			"test1",
			true,
			"Too short",
		},
		{
			"1invalid-1",
			true,
			"Starts with digit",
		},
		{
			"-k1",
			true,
			"Starts with hyphen",
		},
		{
			"long-invalid-123456789-invalid1",
			true,
			"Too long",
		},
		{
			"invalid-123$",
			true,
			"Invalid character",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := validateGcpProjectId(tc.input)

			if tc.error {
				assert.Error(t, err, "Expected error")
			} else {
				assert.Nil(t, err, "Expected no error")
			}
		})
	}
}
