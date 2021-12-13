//go:build generation || (!windows && generation)

package integration

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
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/stretchr/testify/assert"
)

// Test new tf installation
func TestGenerationTfInstallation(t *testing.T) {
	dir, err := ioutil.TempDir("", "tf-install-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Force install
	tf, err := cmd.LocateOrInstallTerraform(true, dir)
	installedTempPath := filepath.Dir(tf.ExecPath())
	defer os.RemoveAll(installedTempPath)

	// Validate install works via version command and that installed version is correct
	version, _, _ := tf.Version(context.Background(), true)
	assert.Equal(t, cmd.InstallTerraformVersion, version.String())
	assert.NoError(t, err)
}

// Test existing TF installations are found
func TestGenerationUsingExistingTfVersion(t *testing.T) {
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
	tf, err := cmd.LocateOrInstallTerraform(false, dir)
	version, _, _ := tf.Version(context.Background(), true)
	assert.Equal(t, "1.0.10", version.String())
	assert.Equal(t, installedTempPath, filepath.Dir(tf.ExecPath()))
	assert.NoError(t, err)
}

// Test too old of installation gets an ephemeral version installed
func TestGenerationOutdatedTfVersionGetsNewInstall(t *testing.T) {
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
	tf, err := cmd.LocateOrInstallTerraform(false, dir)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(filepath.Dir(tf.ExecPath()))
	version, _, _ := tf.Version(context.Background(), true)

	assert.NoError(t, err)
	assert.Equal(t, cmd.InstallTerraformVersion, version.String())
	assert.NotEqual(t, installedTempPath, filepath.Dir(tf.ExecPath())) // We shouldn't be using the old version we installed
}

// Test we handle prior to 0.13.x correctly; which does NOT support the --json argument of `terraform version`
func TestGenerationPre13xTfVersionGetsNewInstall(t *testing.T) {
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
	tf, err := cmd.LocateOrInstallTerraform(false, dir)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(filepath.Dir(tf.ExecPath()))
	version, _, _ := tf.Version(context.Background(), true)

	assert.NoError(t, err)
	assert.Equal(t, cmd.InstallTerraformVersion, version.String())
	assert.NotEqual(t, installedTempPath, filepath.Dir(tf.ExecPath())) // We shouldn't be using the old version we installed
}
