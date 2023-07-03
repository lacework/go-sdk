package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/gcp"
	"github.com/stretchr/testify/assert"
)

func toggleGkeNonInteractive() {
	cli.noCache = !cli.noCache
	cli.nonInteractive = !cli.nonInteractive
}

func TestGenerateMostBasicGkeArgs(t *testing.T) {
	toggleGkeNonInteractive()
	defer toggleGkeNonInteractive()

	data := gcp.GenerateGkeTfConfigurationArgs{}
	err := promptGkeGenerate(&data, &gcp.ServiceAccount{}, &GkeGenerateCommandExtraState{Output: "/tmp"})

	assert.Nil(t, err)
}

func TestGkeGenerationCache(t *testing.T) {
	t.Run("extra state shouldn't be written if empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "t")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		extraState := &GkeGenerateCommandExtraState{}
		extraState.writeCache()
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedGkeAssetExtraState)))
	})
	t.Run("extra state should be written if not empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "t")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		extraState := GkeGenerateCommandExtraState{Output: "/tmp"}
		extraState.writeCache()
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedGkeAssetExtraState)))
	})
	t.Run("iac params should not be cached when empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "t")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		args := gcp.GenerateGkeTfConfigurationArgs{}
		writeGkeGenerationArgsCache(&args)
		assert.NoFileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedGkeAssetIacParams)))
	})
	t.Run("iac params should be cached when not empty", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "t")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(dir)
		cli.InitCache(dir)

		args := gcp.GenerateGkeTfConfigurationArgs{ServiceAccountCredentials: "default"}
		writeGkeGenerationArgsCache(&args)
		assert.FileExists(t, filepath.FromSlash(fmt.Sprintf("%s/cache/standalone/%s", dir, CachedGkeAssetIacParams)))
	})
}
