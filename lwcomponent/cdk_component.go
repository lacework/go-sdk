package lwcomponent

import (
	"github.com/Masterminds/semver"
	"github.com/fatih/color"
)

const (
	DevelopmentEnv = "LW_CDK_DEV_COMPONENT"
)

type CDKComponent struct {
	Name        string
	Description string
	SizeKB      int64
	Type        Type
	Status      Status

	Exec Executer

	apiInfo  ApiInfo
	hostInfo HostInfo
	stage    Stager
}

func NewCDKComponent(name string, componentType Type, apiInfo ApiInfo, hostInfo HostInfo) CDKComponent {
	var (
		exec Executer = &nonExecutable{}
	)

	status := status(apiInfo, hostInfo)

	switch status {
	case Installed, UpdateAvailable, Deprecated, Development:
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
		Name:     name,
		Type:     componentType,
		Status:   status,
		Exec:     exec,
		apiInfo:  apiInfo,
		hostInfo: hostInfo,
	}
}

func (c *CDKComponent) InstalledVersion() (version *semver.Version) {
	var err error

	if c.hostInfo != nil {
		version, err = c.hostInfo.Version()
		if err == nil {
			return
		}
	}

	return
}

func (c *CDKComponent) LatestVersion() (version *semver.Version) {
	if c.apiInfo != nil {
		version = c.apiInfo.LatestVersion()
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
	case Installed, Deprecated, Development, UpdateAvailable:
		version, err = c.hostInfo.Version()
		if err != nil {
			panic(err)
		}
	case NotInstalled:
		version = c.apiInfo.LatestVersion()
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
				return Deprecated
			}
		}
	}

	if apiInfo != nil && hostInfo == nil {
		return NotInstalled
	}

	return status
}
