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
	"path"
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

		return s, s.WriteState()
	}
	return nil, errors.New("invalid api client")
}

// LocalState loads the state from the local storage ("Dir()/state")
func LocalState() (*State, error) {
	state := new(State)
	componentsFile, err := Dir()
	if err != nil {
		return state, err
	}

	stateFile := path.Join(componentsFile, "state")

	stateBytes, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return state, err
	}
	err = json.Unmarshal(stateBytes, state)
	return state, err
}

// GetComponent returns the pointer of a component, if the component does not
// exit, this function will return `nil`
func (s State) GetComponent(name string) *Component {
	for i := range s.Components {
		if s.Components[i].Name == name {
			return &s.Components[i]
		}
	}
	return nil
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

	stateFile := path.Join(dir, "state")
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
	return path.Join(cacheDir, "components"), nil
}

func (s State) Install(name string) error {
	component := s.GetComponent(name)
	if component == nil {
		return errors.New("component not found")
	}

	rPath, err := component.RootPath()
	if err != nil {
		return err
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

	err = downloadFile(path.Join(rPath, component.Name), artifact.URL)
	if err != nil {
		return errors.Wrap(err, "unable to download component artifact")
	}

	// @afiune check 1) cross-platform and 2) correct permissions
	// if the file has permissions already, can we avoid this?
	if component.IsExecutable() {
		if err := os.Chmod(path.Join(rPath, component.Name), 0744); err != nil {
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
	baseRunErr    string = "unable to run component"
	cmpntNotFound string = "component not found on disk"
)

type Artifact struct {
	OS        string `json:"os"`
	ARCH      string `json:"arch"`
	URL       string `json:"url"`
	Signature string `json:"signature"`
	//Size ?
}

type Component struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Type          Type           `json:"type"`
	LatestVersion semver.Version `json:"version"`
	Artifacts     []Artifact     `json:"artifacts"`

	// @dhazekamp command_name required when CLICommand is true?
	CommandName string `json:"command_name"`
}

// RootPath returns the component's root path ("Dir()/{name}")
func (c Component) RootPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, c.Name), nil
}

// Path returns the path to the component ("RootPath()/{name}")
func (c Component) Path() (string, error) {
	dir, err := c.RootPath()
	if err != nil {
		return "", err
	}

	// @afiune maybe component/version/bin
	// but why would we want older versions?
	cPath := path.Join(dir, c.Name)
	if file.FileExists(cPath) {
		return cPath, nil
	}
	return cPath, errors.New(cmpntNotFound)
}

// CurrentVersion returns the current installed version of the component
func (c Component) CurrentVersion() (*semver.Version, error) {
	dir, err := c.RootPath()
	if err != nil {
		return nil, err
	}

	cvPath := path.Join(dir, ".version")
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

	csPath := path.Join(dir, ".signature")
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

	cvPath := path.Join(dir, ".signature")
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

	cvPath := path.Join(dir, ".version")

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
	if err.Error() == cmpntNotFound {
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

func (c Component) isVerified() error {
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
