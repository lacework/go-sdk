package lwcomponent

import (
	"github.com/Masterminds/semver"
)

type ApiInfo interface {
	Id() int32

	LatestVersion() *semver.Version

	AllVersions() []*semver.Version

	Deprecated() bool

	Size() int64
}

type apiInfo struct {
	id            int32
	name          string
	version       semver.Version
	allVersions   []*semver.Version
	desc          string
	sizeKB        int64
	deprecated    bool
	componentType Type
}

func NewAPIInfo(
	id int32,
	name string,
	version *semver.Version,
	allVersions []*semver.Version,
	desc string,
	size int64,
	deprecated bool,
	componentType Type,
) ApiInfo {
	return &apiInfo{
		id:            id,
		name:          name,
		version:       *version,
		allVersions:   allVersions,
		desc:          desc,
		sizeKB:        size,
		deprecated:    deprecated,
		componentType: componentType,
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

// Size implements ApiInfo.
func (a *apiInfo) Size() int64 {
	return a.sizeKB
}
