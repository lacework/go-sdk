package lwcomponent

import (
	"github.com/Masterminds/semver"
)

type ApiInfo struct {
	Id            int32             `json:"id"`
	Name          string            `json:"name"`
	Version       *semver.Version   `json:"version"`
	AllVersions   []*semver.Version `json:"allVersions"`
	Desc          string            `json:"desc"`
	SizeKB        int64             `json:"sizeKB"`
	Deprecated    bool              `json:"deprecated"`
	ComponentType Type              `json:"componentType"`
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
) *ApiInfo {
	return &ApiInfo{
		Id:            id,
		Name:          name,
		Version:       version,
		AllVersions:   allVersions,
		Desc:          desc,
		SizeKB:        size,
		Deprecated:    deprecated,
		ComponentType: componentType,
	}
}
