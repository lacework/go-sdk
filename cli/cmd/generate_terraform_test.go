package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/stretchr/testify/assert"
)

// Test various scenarios for detection and/or installation of Terraform binary
func TestGenerationTfInstallation(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Run("forced installation should installed expected version", func(t *testing.T) {
			dir, err := ioutil.TempDir("", "tf-install-test")
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(dir)

			// Force install
			tf, err := LocateOrInstallTerraform(true, dir)
			installedTempPath := filepath.Dir(tf.ExecPath())
			defer os.RemoveAll(installedTempPath)

			// Validate install works via version command and that installed version is correct
			version, _, _ := tf.Version(context.Background(), true)
			assert.Equal(t, installTerraformVersion, version.String())
			assert.NoError(t, err)
		})
		t.Run("existing terraform version should be located and used, if version is new enough", func(t *testing.T) {
			dir, err := ioutil.TempDir("", "tf-install-test")
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(dir)

			// Install some version that should be detected
			installer := &releases.ExactVersion{
				Product: product.Terraform,
				Version: version.Must(version.NewVersion("1.0.10")),
			}

			installPath, err := installer.Install(context.Background())
			if err != nil {
				panic(err)
			}
			installedTempPath := filepath.Dir(installPath)
			defer os.RemoveAll(installedTempPath)

			// Update path to find the new version
			oldPath := os.Getenv("PATH")
			os.Setenv("PATH", fmt.Sprintf("%s:%s", installedTempPath, oldPath))
			defer os.Setenv("PATH", oldPath)

			// Locate or install should find version 1.0.10
			tf, err := LocateOrInstallTerraform(false, dir)
			version, _, _ := tf.Version(context.Background(), true)
			assert.Equal(t, "1.0.10", version.String())
			assert.Equal(t, installedTempPath, filepath.Dir(tf.ExecPath()))
			assert.NoError(t, err)
		})
		t.Run("installed version of terraform that is too old gets ephemeral installation", func(t *testing.T) {
			dir, err := ioutil.TempDir("", "tf-install-test")
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(dir)

			// Install some version that should be detected
			installer := &releases.ExactVersion{
				Product: product.Terraform,
				Version: version.Must(version.NewVersion("0.13.7")),
			}

			installPath, err := installer.Install(context.Background())
			if err != nil {
				panic(err)
			}
			installedTempPath := filepath.Dir(installPath)
			defer os.RemoveAll(installedTempPath)

			// Update path to find the new version
			oldPath := os.Getenv("PATH")
			os.Setenv("PATH", fmt.Sprintf("%s:%s", installedTempPath, oldPath))
			defer os.Setenv("PATH", oldPath)

			// Locate or install should find version previously installed version, determine it's too old, install newer version
			tf, err := LocateOrInstallTerraform(false, dir)
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(filepath.Dir(tf.ExecPath()))
			version, _, _ := tf.Version(context.Background(), true)

			assert.NoError(t, err)
			assert.Equal(t, installTerraformVersion, version.String())
			assert.NotEqual(t, installedTempPath, filepath.Dir(tf.ExecPath())) // We shouldn't be using the old version we installed
		})
		t.Run("installed version of terraform that is too old to support version checking gets ephemeral installation", func(t *testing.T) {
			dir, err := ioutil.TempDir("", "tf-install-test")
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(dir)

			// Install some version that should be detected
			installer := &releases.ExactVersion{
				Product: product.Terraform,
				Version: version.Must(version.NewVersion("0.12.31")),
			}

			installPath, err := installer.Install(context.Background())
			if err != nil {
				panic(err)
			}
			installedTempPath := filepath.Dir(installPath)
			defer os.RemoveAll(installedTempPath)

			// Update path to find the new version
			oldPath := os.Getenv("PATH")
			os.Setenv("PATH", fmt.Sprintf("%s:%s", installedTempPath, oldPath))
			defer os.Setenv("PATH", oldPath)

			// Locate or install should find version 0.12.31, determine it's too old, install newer version
			tf, err := LocateOrInstallTerraform(false, dir)
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(filepath.Dir(tf.ExecPath()))
			version, _, _ := tf.Version(context.Background(), true)

			assert.NoError(t, err)
			assert.Equal(t, installTerraformVersion, version.String())
			assert.NotEqual(t, installedTempPath, filepath.Dir(tf.ExecPath())) // We shouldn't be using the old version we installed
		})
	}
}
