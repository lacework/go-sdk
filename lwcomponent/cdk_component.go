package lwcomponent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/lacework/go-sdk/internal/file"
	"github.com/pkg/errors"
)

const (
	DevelopmentEnv = "LW_CDK_DEV_COMPONENT"
)

type CDKComponent struct {
	Name           string
	Description    string
	Type           Type
	Status         Status
	InstallMessage string
	UpdateMessage  string

	Exec Executer

	ApiInfo  ApiInfo
	HostInfo HostInfo
	stage    Stager
}

func NewCDKComponent(name string, desc string, componentType Type, apiInfo ApiInfo, hostInfo HostInfo) CDKComponent {
	var (
		exec Executer = &nonExecutable{}
	)

	status := status(apiInfo, hostInfo)

	switch status {
	case Installed, UpdateAvailable, InstalledDeprecated, Development:
		{
			dir := hostInfo.Dir()

			if componentType == BinaryType || componentType == CommandType {
				exec = NewExecuable(name, dir)
			}
		}
	default:
		{

		}
	}

	return CDKComponent{
		Name:        name,
		Description: desc,
		Type:        componentType,
		Status:      status,
		Exec:        exec,
		ApiInfo:     apiInfo,
		HostInfo:    hostInfo,
	}
}

func (c *CDKComponent) Dir() (string, error) {
	dir, err := CatalogCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, c.Name), nil
}

func (c *CDKComponent) EnterDevMode() error {
	if c.HostInfo != nil && c.HostInfo.Development() {
		return errors.New("component already under development.")
	}

	dir, err := c.Dir()
	if err != nil {
		return errors.New("unable to detect RootPath")
	}

	devFile := filepath.Join(dir, DevelopmentFile)
	if !file.FileExists(devFile) {
		devInfo := &DevInfo{
			ComponentType: c.Type,
			Desc:          fmt.Sprintf("(dev-mode) %s", c.Description),
			Name:          c.Name,
			Version:       "0.0.0-dev",
		}

		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(devInfo); err != nil {
			return err
		}

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		return os.WriteFile(devFile, buf.Bytes(), 0644)
	}

	return nil
}

func (c *CDKComponent) InstalledVersion() (version *semver.Version) {
	var err error

	if c.HostInfo != nil {
		version, err = c.HostInfo.Version()
		if err == nil {
			return
		}
	}

	return
}

func (c *CDKComponent) LatestVersion() (version *semver.Version) {
	if c.ApiInfo != nil {
		version = c.ApiInfo.LatestVersion()
	}

	return
}

func (c *CDKComponent) PrintSummary() []string {
	var (
		colorize *color.Color
		version  *semver.Version
		err      error
	)

	switch c.Status {
	case Installed, InstalledDeprecated, NotInstalledDeprecated, UpdateAvailable, Tainted:
		version, err = c.HostInfo.Version()
		if err != nil {
			panic(err)
		}
	case NotInstalled:
		version = c.ApiInfo.LatestVersion()
	case Development:
		devInfo, err := NewDevInfo(c.HostInfo.Dir())
		if err != nil {
			panic(err)
		}
		version, err = semver.NewVersion(devInfo.Version)
		if err != nil {
			panic(err)
		}
	default:
		version = &semver.Version{}
	}

	colorize = c.Status.Color()

	return []string{
		colorize.Sprintf(c.Status.String()),
		c.Name,
		version.String(),
		c.Description,
	}
}

func status(apiInfo ApiInfo, hostInfo HostInfo) Status {
	status := UnknownStatus

	if hostInfo != nil {
		if hostInfo.Development() {
			return Development
		}

		if err := hostInfo.Validate(); err != nil {
			return UnknownStatus
		}

		if apiInfo != nil {
			installedVer, err := hostInfo.Version()
			if err != nil {
				return UnknownStatus
			}

			if isTainted(apiInfo, installedVer) {
				return Tainted
			}

			if apiInfo.Deprecated() {
				return InstalledDeprecated
			}

			latestVer := apiInfo.LatestVersion()
			if latestVer.GreaterThan(installedVer) {
				return UpdateAvailable
			} else {
				return Installed
			}
		} else {
			return InstalledDeprecated
		}
	}

	if apiInfo != nil && hostInfo == nil {
		if apiInfo.Deprecated() {
			return NotInstalledDeprecated
		}

		return NotInstalled
	}

	return status
}

func isTainted(apiInfo ApiInfo, installedVer *semver.Version) bool {
	for _, ver := range apiInfo.AllVersions() {
		if ver.Equal(installedVer) {
			return false
		}
	}
	return true
}
