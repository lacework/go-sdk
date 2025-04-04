package lwcomponent_test

import (
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/lacework/go-sdk/v2/lwcomponent"
	"github.com/stretchr/testify/assert"
)

func TestNewCDKComponent(t *testing.T) {
	_, home := FakeHome()
	defer ResetHome(home)

	t.Run("UnknownStatus", func(t *testing.T) {
		component := lwcomponent.NewCDKComponent(nil, &lwcomponent.HostInfo{Name: "unknown"})

		assert.Equal(t, lwcomponent.UnknownStatus, component.Status)
	})

	t.Run("NotInstalled", func(t *testing.T) {
		var name string = "not-installed"

		version, _ := semver.NewVersion("1.1.1")
		allVersions := []*semver.Version{version}

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, false, lwcomponent.BinaryType)

		component := lwcomponent.NewCDKComponent(apiInfo, nil)

		assert.Equal(t, lwcomponent.NotInstalled, component.Status)
	})

	t.Run("Installed", func(t *testing.T) {
		var (
			name string = "installed"
			ver  string = "1.1.1"
		)

		version, _ := semver.NewVersion(ver)
		allVersions := []*semver.Version{version}

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, false, lwcomponent.BinaryType)

		CreateLocalComponent(name, ver, false, lwcomponent.BinaryType)

		dir, _ := lwcomponent.CatalogCacheDir()

		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(apiInfo, hostInfo)

		assert.Equal(t, lwcomponent.Installed, component.Status)
	})

	t.Run("UpdateAvailable", func(t *testing.T) {
		var (
			name       string = "update-available"
			apiVer     string = "1.1.1"
			installVer string = "0.0.1"
		)

		version, _ := semver.NewVersion(apiVer)
		installedVersion, _ := semver.NewVersion(installVer)
		allVersions := []*semver.Version{version, installedVersion}

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, false, lwcomponent.BinaryType)

		CreateLocalComponent(name, installVer, false, lwcomponent.BinaryType)

		dir, _ := lwcomponent.CatalogCacheDir()

		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(apiInfo, hostInfo)

		assert.Equal(t, lwcomponent.UpdateAvailable, component.Status)
	})

	t.Run("Installed No API", func(t *testing.T) {
		var (
			name string = "installed-no-api"
			ver  string = "1.1.1"
		)

		CreateLocalComponent(name, ver, false, lwcomponent.BinaryType)

		dir, _ := lwcomponent.CatalogCacheDir()

		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(nil, hostInfo)

		assert.Equal(t, lwcomponent.Installed, component.Status)
	})

	t.Run("Installed Deprecated API", func(t *testing.T) {
		var (
			name string = "installed-deprecated-api"
			ver  string = "1.1.1"
		)

		version, _ := semver.NewVersion(ver)
		allVersions := []*semver.Version{version}

		CreateLocalComponent(name, ver, false, lwcomponent.BinaryType)

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, true, lwcomponent.BinaryType)

		dir, _ := lwcomponent.CatalogCacheDir()
		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(apiInfo, hostInfo)

		assert.Equal(t, lwcomponent.InstalledDeprecated, component.Status)
	})

	t.Run("Not Installed Deprecated", func(t *testing.T) {
		var (
			name string = "not-installed-deprecated"
		)

		version, _ := semver.NewVersion("1.1.1")
		allVersions := []*semver.Version{version}

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, true, lwcomponent.BinaryType)

		component := lwcomponent.NewCDKComponent(apiInfo, nil)

		assert.Equal(t, lwcomponent.NotInstalledDeprecated, component.Status)
	})

	t.Run("Tainted", func(t *testing.T) {
		var (
			name       string = "tainted"
			apiVer     string = "1.1.1"
			installVer string = "0.0.1"
		)

		version, _ := semver.NewVersion(apiVer)
		allVersions := []*semver.Version{version}

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, false, lwcomponent.BinaryType)

		CreateLocalComponent(name, installVer, false, lwcomponent.BinaryType)

		dir, _ := lwcomponent.CatalogCacheDir()

		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(apiInfo, hostInfo)

		assert.Equal(t, lwcomponent.Tainted, component.Status)
	})

	t.Run("Development", func(t *testing.T) {
		var (
			name string = "development"
			ver  string = "1.1.1"
		)

		CreateLocalComponent(name, ver, true, lwcomponent.BinaryType)

		dir, _ := lwcomponent.CatalogCacheDir()

		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(nil, hostInfo)

		assert.Equal(t, lwcomponent.Development, component.Status)
	})

	t.Run("binary-exec", func(t *testing.T) {
		var (
			name string = "binary-exec"
			ver  string = "1.1.1"
		)

		version, _ := semver.NewVersion(ver)
		allVersions := []*semver.Version{version}

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, false, lwcomponent.BinaryType)

		CreateLocalComponent(name, ver, false, lwcomponent.BinaryType)

		dir, _ := lwcomponent.CatalogCacheDir()

		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(apiInfo, hostInfo)

		_, _, err := component.Exec.Execute([]string{}, "")
		assert.Nil(t, err)
	})

	t.Run("command-exec", func(t *testing.T) {
		var (
			name string = "command-exec"
			ver  string = "1.1.1"
		)

		version, _ := semver.NewVersion(ver)
		allVersions := []*semver.Version{version}

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, false, lwcomponent.BinaryType)

		CreateLocalComponent(name, ver, false, lwcomponent.BinaryType)

		dir, _ := lwcomponent.CatalogCacheDir()

		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(apiInfo, hostInfo)

		_, _, err := component.Exec.Execute([]string{}, "")
		assert.Nil(t, err)
	})

	t.Run("library-non-exec", func(t *testing.T) {
		var (
			name string = "lib-non-exec"
			ver  string = "1.1.1"
		)

		version, _ := semver.NewVersion(ver)
		allVersions := []*semver.Version{version}

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, false, lwcomponent.LibraryType)

		CreateLocalComponent(name, ver, false, lwcomponent.LibraryType)

		dir, _ := lwcomponent.CatalogCacheDir()

		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(apiInfo, hostInfo)

		_, _, err := component.Exec.Execute([]string{}, "")
		assert.Equal(t, lwcomponent.ErrNonExecutable, err)
	})

	t.Run("stand-alone-non-exec", func(t *testing.T) {
		var (
			name string = "stand-alone-non-exec"
			ver  string = "1.1.1"
		)

		version, _ := semver.NewVersion(ver)
		allVersions := []*semver.Version{version}

		apiInfo := lwcomponent.NewAPIInfo(1, name, version, allVersions, "", 1, false, lwcomponent.StandaloneType)

		CreateLocalComponent(name, ver, false, lwcomponent.StandaloneType)

		dir, _ := lwcomponent.CatalogCacheDir()

		hostInfo, _ := lwcomponent.LoadHostInfo(filepath.Join(dir, name))

		component := lwcomponent.NewCDKComponent(apiInfo, hostInfo)

		_, _, err := component.Exec.Execute([]string{}, "")
		assert.Equal(t, lwcomponent.ErrNonExecutable, err)
	})
}
