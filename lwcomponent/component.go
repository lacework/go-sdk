//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// The Lacework component package facilitates loading and executing components
package lwcomponent

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"aead.dev/minisign"
	"github.com/Masterminds/semver"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/cache"
	"github.com/lacework/go-sdk/internal/file"
)

// State holds the components specification
//
// You can load the state from the Lacework API server by passing an `api.Client`.
//
// client, err := api.NewClient(account, opts...)
// cState, err := lwcomponent.LoadState(client)
//
// Or, you can load the state from the local storage.
//
// cState, err := lwcomponent.LocalState()
//
type State struct {
	Version    string      `json:"version"`
	Components []Component `json:"components"`
}

// LoadState loads the state from the Lacework API server
func LoadState(client *api.Client) (*State, error) {
	if client != nil {
		s := new(State)
		err := client.RequestDecoder("GET", "v2/Components", nil, s)
		if err != nil {
			return s, err
		}

		s.loadDevComponent()

		return s, s.WriteState()
	}
	return nil, errors.New("invalid api client")
}

// loadDevComponent will load a component that is under development,
// developers need to export the environment variable 'LW_CDK_DEV_COMPONENT'
func (s *State) loadDevComponent() {
	if devComponent := os.Getenv("LW_CDK_DEV_COMPONENT"); devComponent != "" {
		for i := range s.Components {
			if s.Components[i].Name == devComponent {
				// existing component being developed
				if err := s.Components[i].loadDevSpecs(); err != nil {
					s.Components[i].Description = err.Error()
				}
				return
			}
		}

		// component is not yet defined, add it to the state
		dev := Component{Name: devComponent}
		if err := dev.loadDevSpecs(); err != nil {
			dev.Description = err.Error()
		}
		s.Components = append(s.Components, dev)
	}
}

// LocalState loads the state from the local storage ("Dir()/state")
func LocalState() (*State, error) {
	state := new(State)
	componentsFile, err := Dir()
	if err != nil {
		return state, err
	}

	stateFile := filepath.Join(componentsFile, "state")

	stateBytes, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return state, err
	}
	err = json.Unmarshal(stateBytes, state)
	return state, err
}

// GetComponent returns the pointer of a component, if the component is not
// found, this function will return a `nil` pointer and `false`
//
// Usage:
// ```go
// component, found := s.GetComponent(name)
// if !found {
//   fmt.Println("Component %s not found", name)
// }
// ```
func (s State) GetComponent(name string) (*Component, bool) {
	for i := range s.Components {
		if s.Components[i].Name == name {
			return &s.Components[i], true
		}
	}
	return nil, false
}

// WriteState stores the components state to disk
func (s State) WriteState() error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	stateFile := filepath.Join(dir, "state")
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(s); err != nil {
		return err
	}

	if err := ioutil.WriteFile(stateFile, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

// Dir returns the directory where the components will be stored
func Dir() (string, error) {
	cacheDir, err := cache.CacheDir()
	if err != nil {
		return "", errors.Wrap(err, "unable to locate components directory")
	}
	return filepath.Join(cacheDir, "components"), nil
}

func (s State) Install(name string) error {
	component, found := s.GetComponent(name)
	if !found {
		return errors.New("component not found")
	}

	rPath, err := component.RootPath()
	if err != nil {
		return err
	}

	// verify development mode
	if component.underDevelopment() {
		p, _ := component.Path() // @afiune we don't care if the component exists or not
		msg := "components under development can't be installed.\n\n" +
			"Deploy the component manually at '" + p + "'"
		return errors.New(msg)
	}

	// @afiune verify if component is in latest

	// @afiune install
	if err := os.MkdirAll(rPath, os.ModePerm); err != nil {
		return err
	}

	artifact, found := component.ArtifactForRunningHost()
	if !found {
		return errors.New("unsupported platform")
	}

	path, err := component.Path()
	// it is ok if the component is not found yet
	if err != nil && !IsNotFound(err) {
		return err
	}

	err = downloadFile(path, artifact.URL)
	if err != nil {
		return errors.Wrap(err, "unable to download component artifact")
	}

	// @afiune check 1) cross-platform and 2) correct permissions
	// if the file has permissions already, can we avoid this?
	if component.IsExecutable() {
		if err := os.Chmod(path, 0744); err != nil {
			return errors.Wrap(err, "unable to make component executable")
		}
	}

	if err := component.WriteSignature([]byte(artifact.Signature)); err != nil {
		return err
	}

	if err := component.WriteVersion(); err != nil {
		return err
	}

	// verify component
	if err := component.isVerified(); err != nil {
		// @afiune notify and remove installed component
		defer os.RemoveAll(rPath)
		return err
	}

	return nil
}

var (
	baseRunErr string = "unable to run component"
)

type Artifact struct {
	OS        string `json:"os"`
	ARCH      string `json:"arch"`
	URL       string `json:"url"`
	Signature string `json:"signature"`
	//Size ?
}

// Components should leave a trail/crumb after installation or update,
// these messages will be shown by the Lacework CLI
type Breadcrumbs struct {
	InstallationMessage string `json:"installationMessage,omitempty"`
	UpdateMessage       string `json:"updateMessage,omitempty"`
}

type Component struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Type          Type           `json:"type"`
	LatestVersion semver.Version `json:"version"`
	Artifacts     []Artifact     `json:"artifacts"`
	Breadcrumbs   Breadcrumbs    `json:"breadcrumbs,omitempty"`

	// @dhazekamp command_name required when CLICommand is true?
	CommandName string `json:"command_name"`
}

// RootPath returns the component's root path ("Dir()/{name}")
func (c Component) RootPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, c.Name), nil
}

// Path returns the path to the component ("RootPath()/{name}")
func (c Component) Path() (string, error) {
	dir, err := c.RootPath()
	if err != nil {
		return "", err
	}

	// @afiune maybe component/version/bin
	// but why would we want older versions?
	cPath := filepath.Join(dir, c.Name)
	if runtime.GOOS == "windows" {
		cPath += ".exe"
	}

	if file.FileExists(cPath) {
		return cPath, nil
	}
	return cPath, ErrComponentNotFound
}

// CurrentVersion returns the current installed version of the component
func (c Component) CurrentVersion() (*semver.Version, error) {
	// development mode, avoid loading the current version,
	// return latest which is what's inside the '.dev' specs
	if c.underDevelopment() {
		return &c.LatestVersion, nil
	}

	dir, err := c.RootPath()
	if err != nil {
		return nil, err
	}

	cvPath := filepath.Join(dir, ".version")
	if !file.FileExists(cvPath) {
		// @afiune help the user fix this issue with a better error message
		return nil, errors.New("component version file does not exist")
	}

	dat, err := os.ReadFile(cvPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read component version file")
	}

	cv, err := semver.NewVersion(strings.TrimSpace(string(dat)))
	if err != nil {
		err = errors.New("unable to parse component version")
	}
	return cv, err
}

// SignatureFromDisk returns the component signature stored on disk ("RootPath()/.signature")
func (c Component) SignatureFromDisk() ([]byte, error) {
	var sig []byte

	dir, err := c.RootPath()
	if err != nil {
		return nil, err
	}

	csPath := filepath.Join(dir, ".signature")
	if !file.FileExists(csPath) {
		return sig, errors.New("component signature file does not exist")
	}

	dat, err := os.ReadFile(csPath)
	if err != nil {
		return sig, errors.Wrap(err, "unable to read component signature file")
	}

	// decode artifact signature
	sig, err = base64.StdEncoding.DecodeString(string(dat))
	if err != nil {
		return sig, errors.New("unable to decode component signature")
	}
	return sig, nil
}

// WriteSignature stores the component signature on disk
func (c Component) WriteSignature(signature []byte) error {
	dir, err := c.RootPath()
	if err != nil {
		return err
	}

	cvPath := filepath.Join(dir, ".signature")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return ioutil.WriteFile(cvPath, signature, 0644)
}

// WriteVersion stores the component version on disk
func (c Component) WriteVersion() error {
	dir, err := c.RootPath()
	if err != nil {
		return err
	}

	cvPath := filepath.Join(dir, ".version")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return ioutil.WriteFile(cvPath, []byte(c.LatestVersion.String()), 0644)
}

// UpdateAvailable returns true if there is a newer version of the component
func (c Component) UpdateAvailable() (bool, error) {
	cv, err := c.CurrentVersion()
	if err != nil {
		return false, err
	}

	return c.LatestVersion.GreaterThan(cv), nil
}

// Status returns the component status
func (c Component) Status() Status {
	_, err := c.Path()
	if err == nil {
		// check if the component has an update
		update, err := c.UpdateAvailable()
		if err == nil && update {
			return UpdateAvailable

		}
		return Installed
	}
	if IsNotFound(err) {
		return NotInstalled
	}
	return UnknownStatus
}

// ArtifactForRunningHost returns the right component artifact for the running host,
func (c Component) ArtifactForRunningHost() (*Artifact, bool) {
	for _, artifact := range c.Artifacts {
		if artifact.OS == runtime.GOOS && artifact.ARCH == runtime.GOARCH {
			return &artifact, true
		}
	}
	return nil, false
}

// loadDevSpecs will lookup for the '.dev' specs file under the
// component root path to load it into the component itself
func (c *Component) loadDevSpecs() error {
	dir, err := c.RootPath()
	if err != nil {
		return errors.New("unable to detect RootPath")
	}

	devSpecs := filepath.Join(dir, ".dev")
	if file.FileExists(devSpecs) {
		devSpecsBytes, err := ioutil.ReadFile(devSpecs)
		if err != nil {
			return errors.Errorf("unable to read %s file", devSpecs)
		}
		err = json.Unmarshal(devSpecsBytes, c)
		if err != nil {
			return errors.Errorf("unable to unmarshal %s file", devSpecs)
		}
	} else {
		return errors.Errorf("create dev specs file '%s'", devSpecs)
	}

	return nil
}

// underDevelopment returns true if the component is under development
// that is, if the component root path has the '.dev' specs file or, if
// the environment variable 'LW_CDK_DEV_COMPONENT' matches the component name
func (c Component) underDevelopment() bool {
	if os.Getenv("LW_CDK_DEV_COMPONENT") == c.Name {
		return true
	}

	dir, err := c.RootPath()
	if err != nil {
		return false
	}

	return file.FileExists(filepath.Join(dir, ".dev"))
}

// isVerified checks if the component has a valid signature
func (c Component) isVerified() error {
	// development mode, avoid verifying
	if c.underDevelopment() {
		return nil
	}

	// get component signature
	sig, err := c.SignatureFromDisk()
	if err != nil {
		return err
	}

	// get component path
	cPath, err := c.Path()
	if err != nil {
		return err
	}

	// open the component
	f, err := os.ReadFile(cPath)
	if err != nil {
		return errors.New("unable to read component file")
	}

	// load public key
	rootPublicKey := minisign.PublicKey{}
	if err := rootPublicKey.UnmarshalText([]byte(publicKey)); err != nil {
		return errors.Wrap(err, "unable to load root public key")
	}

	// validate the signature
	return verifySignature(rootPublicKey, f, sig)
}
