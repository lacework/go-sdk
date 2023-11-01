package lwcomponent

import (
	"github.com/Masterminds/semver"
)

type ApiInfo interface {
	Id() int32

	LatestVersion() *semver.Version

	AllVersions() []*semver.Version

	Deprecated() bool
}

type apiInfo struct {
	id          int32
	name        string
	version     semver.Version
	allVersions []*semver.Version
	desc        string
	sizeKB      int64
	deprecated  bool
}

func NewAPIInfo(
	id int32,
	name string,
	version *semver.Version,
	allVersions []*semver.Version,
	desc string,
	size int64,
	deprecated bool,
) ApiInfo {
	return &apiInfo{
		id:          id,
		name:        name,
		version:     *version,
		allVersions: allVersions,
		desc:        desc,
		sizeKB:      size,
		deprecated:  deprecated,
	}
}

func (a *apiInfo) Id() int32 {
	return a.id
}

// AllVersions implements ApiInfo.
func (a *apiInfo) AllVersions() []*semver.Version {
	return a.allVersions
}

// LatestVersion implements ApiInfo.
func (a *apiInfo) LatestVersion() *semver.Version {
	return &a.version
}

// Deprecated implements ApiInfo.
func (a *apiInfo) Deprecated() bool {
	return a.deprecated
}
