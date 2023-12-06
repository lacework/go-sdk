package lwcomponent

import (
	"github.com/Masterminds/semver"
	"github.com/fatih/color"
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
	case Installed, InstalledDeprecated, NotInstalledDeprecated, Development, UpdateAvailable, Tainted:
		version, err = c.HostInfo.Version()
		if err != nil {
			panic(err)
		}
	case NotInstalled:
		version = c.ApiInfo.LatestVersion()
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
			if hostInfo.Development() {
				return Development
			} else {
				return InstalledDeprecated
			}
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
