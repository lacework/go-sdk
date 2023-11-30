package lwcomponent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

type DevInfo struct {
	ComponentType Type
	Desc          string
	Name          string
	Version       string
}

func NewDevInfo(dir string) (*DevInfo, error) {
	path := filepath.Join(dir, DevelopmentFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf("unable to read %s file", path)
	}

	info := DevInfo{}

	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, errors.Errorf("unable to unmarshal %s file", path)
	}

	_, err = semver.NewVersion(info.Version)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("development component '%s' version '%s'", info.Name, info.Version))
	}

	return &info, nil
}
