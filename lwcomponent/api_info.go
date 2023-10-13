package lwcomponent

import (
	"github.com/Masterminds/semver"
)

type ApiInfo interface {
	Id() int32

	LatestVersion() *semver.Version

	AllVersions() []semver.Version
}

type apiInfo struct {
	id          int32
	name        string
	version     semver.Version
	allVersions []semver.Version
	desc        string
	sizeKB      int64
}

func NewAPIInfo(id int32, name string, version *semver.Version, allVersions []semver.Version, desc string, size int64) ApiInfo {
	return &apiInfo{
		id:          id,
		name:        name,
		version:     *version,
		allVersions: allVersions,
		desc:        desc,
		sizeKB:      size,
	}
}

func (a *apiInfo) Id() int32 {
	return a.id
}

func (a *apiInfo) LatestVersion() *semver.Version {
	return &a.version
}

func (a *apiInfo) AllVersions() []semver.Version {
	return a.allVersions
}
