package lwcomponent

import (
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
	DevelopmentFile = ".dev"
)

type HostInfo interface {
	Delete() error

	Development() bool

	Dir() string

	Name() string

	Signature() (sig []byte, err error)

	Validate() error

	Version() (*semver.Version, error)
}

type hostInfo struct {
	dir string
}

func NewHostInfo(dir string) HostInfo {
	return &hostInfo{dir}
}

func (h *hostInfo) Delete() error {
	return os.RemoveAll(h.dir)
}

func (h *hostInfo) Development() bool {
	return file.FileExists(filepath.Join(h.dir, DevelopmentFile))
}

func (h *hostInfo) Dir() string {
	return h.dir
}

// Returns the Component name
//
// The Component name is the same as the name of the base directory
func (h *hostInfo) Name() string {
	return filepath.Base(h.dir)
}

func (h *hostInfo) Signature() (sig []byte, err error) {
	_, err = os.Stat(h.dir)
	if os.IsNotExist(err) {
		return
	}

	path := filepath.Join(h.dir, SignatureFile)
	if !file.FileExists(path) {
		return
	}

	sig, err = os.ReadFile(path)
	if err != nil {
		return
	}

	return
}

func (h *hostInfo) Version() (version *semver.Version, err error) {
	_, err = os.Stat(h.dir)
	if os.IsNotExist(err) {
		return
	}

	path := filepath.Join(h.dir, VersionFile)
	if !file.FileExists(path) {
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	return semver.NewVersion(strings.TrimSpace(string(data)))
}

func (h *hostInfo) Validate() (err error) {
	data, err := os.ReadFile(filepath.Join(h.dir, VersionFile))
	if err != nil {
		return
	}

	version := string(data)

	_, err = semver.NewVersion(strings.TrimSpace(version))
	if err != nil {
		return
	}

	componentName := h.Name()

	if !file.FileExists(filepath.Join(h.dir, SignatureFile)) {
		return errors.New(fmt.Sprintf("missing file '%s'", componentName))
	}

	if !file.FileExists(filepath.Join(h.dir, componentName)) {
		return errors.New(fmt.Sprintf("missing file '%s'", componentName))
	}

	return
}
