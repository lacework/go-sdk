package lwcomponent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/lacework/go-sdk/v2/internal/file"
	"github.com/pkg/errors"
)

var (
	VersionFile     = ".version"
	SignatureFile   = ".signature"
	InfoFile        = ".info"
	DevelopmentFile = ".dev"
)

type HostInfo struct {
	Name          string `json:"name"`
	ComponentType Type   `json:"type"`
	Desc          string `json:"description"`
	Dir           string `json:"-"`
}

func LoadHostInfo(dir string) (*HostInfo, error) {
	path := filepath.Join(dir, InfoFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf("unable to read %s file", path)
	}

	hostInfo := HostInfo{}

	err = json.Unmarshal(data, &hostInfo)
	if err != nil {
		return nil, errors.Errorf("unable to unmarshal %s file", path)
	}

	hostInfo.Dir = dir

	if hostInfo.Name == "" {
		hostInfo.Name = filepath.Base(dir)
	}

	return &hostInfo, nil
}

func NewHostInfo(dir string, desc string, componentType Type) (*HostInfo, error) {
	path := filepath.Join(dir, InfoFile)

	if !file.FileExists(path) {
		info := &HostInfo{
			Name:          filepath.Base(dir),
			Dir:           dir,
			ComponentType: componentType,
			Desc:          desc,
		}

		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(info); err != nil {
			return nil, err
		}

		return info, os.WriteFile(path, buf.Bytes(), 0644)
	}

	return LoadHostInfo(dir)
}

func (h *HostInfo) Delete() error {
	return os.RemoveAll(h.Dir)
}

func (h *HostInfo) Development() bool {
	return file.FileExists(filepath.Join(h.Dir, DevelopmentFile))
}

func (h *HostInfo) Signature() (sig []byte, err error) {
	_, err = os.Stat(h.Dir)
	if os.IsNotExist(err) {
		return
	}

	path := filepath.Join(h.Dir, SignatureFile)
	if !file.FileExists(path) {
		return
	}

	sig, err = os.ReadFile(path)
	if err != nil {
		return
	}

	return
}

func (h *HostInfo) Version() (version *semver.Version, err error) {
	_, err = os.Stat(h.Dir)
	if os.IsNotExist(err) {
		return
	}

	path := filepath.Join(h.Dir, VersionFile)
	if !file.FileExists(path) {
		return nil, errors.New("missing .version file")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	return semver.NewVersion(strings.TrimSpace(string(data)))
}

func (h *HostInfo) Validate() (err error) {
	data, err := os.ReadFile(filepath.Join(h.Dir, VersionFile))
	if err != nil {
		return
	}

	version := string(data)

	_, err = semver.NewVersion(strings.TrimSpace(version))
	if err != nil {
		return
	}

	componentName := h.Name

	if !file.FileExists(filepath.Join(h.Dir, SignatureFile)) {
		return errors.New(fmt.Sprintf("missing file '%s'", componentName))
	}

	path := filepath.Join(h.Dir, componentName)

	if operatingSystem == "windows" {
		path = fmt.Sprintf("%s.exe", path)
	}

	if !file.FileExists(path) {
		return errors.New(fmt.Sprintf("missing file '%s'", componentName))
	}

	return
}
