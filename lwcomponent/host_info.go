package lwcomponent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/lacework/go-sdk/internal/file"
	"github.com/pkg/errors"
)

var (
	VersionFile     = ".version"
	SignatureFile   = ".signature"
	InfoFile        = ".info"
	DevelopmentFile = ".dev"
)

type HostInfo struct {
	ComponentType Type   `json:"type"`
	Description   string `json:"description"`
	Dir           string `json:"-"`
}

func NewHostInfo(dir string) (*HostInfo, error) {
	path := filepath.Join(dir, InfoFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf("unable to read %s file", path)
	}

	info := HostInfo{}

	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, errors.Errorf("unable to unmarshal %s file", path)
	}

	info.Dir = dir

	return &info, nil
}

func CreateHostInfo(dir string, desc string, componentType Type) (*HostInfo, error) {
	path := filepath.Join(dir, InfoFile)

	if !file.FileExists(path) {
		info := &HostInfo{
			Dir:           dir,
			ComponentType: componentType,
			Description:   desc,
		}

		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(info); err != nil {
			return nil, err
		}

		return info, os.WriteFile(path, buf.Bytes(), 0644)
	}

	return NewHostInfo(dir)
}

func (h *HostInfo) Delete() error {
	return os.RemoveAll(h.Dir)
}

func (h *HostInfo) Development() bool {
	return file.FileExists(filepath.Join(h.Dir, DevelopmentFile))
}

// Returns the Component name
//
// The Component name is the same as the name of the base directory
func (h *HostInfo) Name() string {
	return filepath.Base(h.Dir)
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

	componentName := h.Name()

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
