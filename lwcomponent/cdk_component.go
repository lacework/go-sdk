package lwcomponent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/lacework/go-sdk/v2/internal/file"
	"github.com/pkg/errors"
)

const (
	DevelopmentEnv = "LW_CDK_DEV_COMPONENT"
)

type CDKComponent struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Type           Type   `json:"type"`
	Status         Status `json:"-"`
	InstallMessage string `json:"-"`
	UpdateMessage  string `json:"-"`

	Exec Executer `json:"-"`

	ApiInfo  *ApiInfo  `json:"apiInfo,omitempty"`
	HostInfo *HostInfo `json:"-"`
	stage    Stager
}

func NewCDKComponent(apiInfo *ApiInfo, hostInfo *HostInfo) CDKComponent {
	var (
		exec Executer = &nonExecutable{}
	)

	status := status(apiInfo, hostInfo)

	switch status {
	case Installed, UpdateAvailable, InstalledDeprecated, Development:
		{
			dir := hostInfo.Dir

			if hostInfo.ComponentType == BinaryType || hostInfo.ComponentType == CommandType {
				exec = NewExecuable(hostInfo.Name, dir)
			}
		}
	default:
		{

		}
	}

	if apiInfo != nil {
		return CDKComponent{
			Name:        apiInfo.Name,
			Description: apiInfo.Desc,
			Type:        apiInfo.ComponentType,
			Status:      status,
			Exec:        exec,
			ApiInfo:     apiInfo,
			HostInfo:    hostInfo,
		}
	}

	if hostInfo != nil {
		return CDKComponent{
			Name:        hostInfo.Name,
			Description: hostInfo.Desc,
			Type:        hostInfo.ComponentType,
			Status:      status,
			Exec:        exec,
			ApiInfo:     apiInfo,
			HostInfo:    hostInfo,
		}
	}

	return CDKComponent{}
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

func (c *CDKComponent) InstalledVersion() *semver.Version {
	if c.HostInfo != nil {
		version, err := c.HostInfo.Version()
		if err == nil {
			return version
		}

		if componentDir, err := c.Dir(); err == nil {
			if devInfo, err := newDevInfo(componentDir); err == nil {
				version, err = semver.NewVersion(devInfo.Version)
				if err == nil {
					return version
				}
			}
		}
	}

	return nil
}

func (c *CDKComponent) LatestVersion() *semver.Version {
	if c.ApiInfo != nil {
		return c.ApiInfo.Version
	}

	return nil
}

func (c *CDKComponent) PrintSummary() []string {
	var (
		colorize *color.Color
		version  *semver.Version
	)

	switch c.Status {
	case Installed, InstalledDeprecated, UpdateAvailable, Development, Tainted:
		version = c.InstalledVersion()
	case NotInstalled, NotInstalledDeprecated:
		version = c.ApiInfo.Version
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

func status(apiInfo *ApiInfo, hostInfo *HostInfo) Status {
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

			if apiInfo.Deprecated {
				return InstalledDeprecated
			}

			latestVer := apiInfo.Version
			if latestVer.GreaterThan(installedVer) {
				return UpdateAvailable
			} else {
				return Installed
			}
		} else {
			return Installed
		}
	}

	if apiInfo != nil && hostInfo == nil {
		if apiInfo.Deprecated {
			return NotInstalledDeprecated
		}

		return NotInstalled
	}

	return status
}

func isTainted(apiInfo *ApiInfo, installedVer *semver.Version) bool {
	if len(apiInfo.AllVersions) == 0 {
		return false
	}

	for _, ver := range apiInfo.AllVersions {
		if ver.Equal(installedVer) {
			return false
		}
	}
	return true
}
