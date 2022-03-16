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

func TestValidateSaCredFile(t *testing.T) {
	t.Run("JSON credentials file with client_email and private_key is valid", func(t *testing.T) {
		err := validateServiceAccountCredentialsFile("test_resources/generate_gcp_test_data/private_key_client_email_valid.json")
		assert.Equal(t, err, nil)
	})

	t.Run("JSON credentials file without client_email is not  valid", func(t *testing.T) {
		err := validateServiceAccountCredentialsFile("test_resources/generate_gcp_test_data/creds_no_client_email.json")
		assert.EqualError(t, err, "invalid GCP Service Account credentials file. The private_key and client_email fields MUST be present.")
	})

	t.Run("JSON credentials file without private_key is not valid", func(t *testing.T) {
		err := validateServiceAccountCredentialsFile("test_resources/generate_gcp_test_data/creds_no_private_key.json")
		assert.EqualError(t, err, "invalid GCP Service Account credentials file. The private_key and client_email fields MUST be present.")
	})

	t.Run("invalid JSON file", func(t *testing.T) {
		err := validateServiceAccountCredentialsFile("test_resources/generate_gcp_test_data/invalid_json.json")
		assert.EqualError(t, err, "unable to parse credentials file: invalid character '}' looking for beginning of object key string")
	})

	t.Run("non existent JSON file", func(t *testing.T) {
		err := validateServiceAccountCredentialsFile("test_resources/generate_gcp_test_data/foo.json")
		assert.EqualError(t, err, "provided GCP credentials file does not exist")
	})
}
