package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/azure"
	"github.com/stretchr/testify/assert"
)

func toggleAzureNonInteractive() {
	cli.noCache = !cli.noCache
	cli.nonInteractive = !cli.nonInteractive
}

func TestGenerationActivityLogWithConfig(t *testing.T) {
	toggleAzureNonInteractive()
	defer toggleAzureNonInteractive()

	data := azure.GenerateAzureTfConfigurationArgs{}
	data.Config = true
	data.ActivityLog = true
	data.CreateAdIntegration = true
	err := promptAzureGenerate(&data, &AzureGenerateCommandExtraState{Output: "/tmp"})
	assert.Nil(t, err)

}

func TestMissingValidEntity(t *testing.T) {
	toggleNonInteractive()
	defer toggleNonInteractive()

	data := azure.GenerateAzureTfConfigurationArgs{}
	data.Config = false
	data.ActivityLog = false

	err := promptAzureGenerate(&data, &AzureGenerateCommandExtraState{Output: "/tmp"})
	assert.Error(t, err)
	assert.Equal(t, "must enable activity log or config", err.Error())
}

func TestValidStorageLocations(t *testing.T) {
	err := validateStorageLocation("East US")
	assert.Nil(t, err)
	err = validateStorageLocation("Brazil Southeast")
	assert.Nil(t, err)

}

func TestInvalidStorageLocations(t *testing.T) {
	err := validateStorageLocation("Mars")
	assert.Error(t, err)
	assert.Equal(t, "invalid storage location supplied", err.Error())
	err = validateStorageLocation("Jupiter")
	assert.Error(t, err)
	assert.Equal(t, "invalid storage location supplied", err.Error())
}

func TestAzureGenerationCache(t *testing.T) {
	t.Run("extra state shouldn't be written if empty", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		extraState := &AzureGenerateCommandExtraState{}
		extraState.writeCache()
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAzureAssetExtraState)))
	})
	t.Run("extra state should be written if not empty", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		extraState := AzureGenerateCommandExtraState{Output: "/tmp"}
		extraState.writeCache()
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAzureAssetExtraState)))
	})
	t.Run("iac params should not be cached when empty", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		args := azure.GenerateAzureTfConfigurationArgs{}
		writeAzureGenerationArgsCache(&args)
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAzureAssetIacParams)))
	})
	t.Run("iac params should be cached when not empty", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "lacework-cli-cache")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		args := azure.GenerateAzureTfConfigurationArgs{ActivityLog: true}
		writeAzureGenerationArgsCache(&args)
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedAzureAssetIacParams)))
	})
}
