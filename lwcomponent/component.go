//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

// A Lacework component package to help facilitate the loading and execution of components
package lwcomponent

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

type State struct {
	Version    string      `json:"version"`
	Components []Component `json:"components"`
}

func (s State) GetComponent(name string) *Component {
	for i := range s.Components {
		if s.Components[i].Name == name {
			return &s.Components[i]
		}
	}
	return nil
}

func cacheDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".config", "lacework"), nil
}

// fileExists checks if a file exists and is not a directory
func fileExists(filename string) bool {
	f, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !f.IsDir()
}

//go:embed state
var componentState string

// @dhazekamp need to avoid cache poisoning with respect retrieving component signature
func LoadState() (*State, error) {
	state := new(State)

	if err := json.Unmarshal([]byte(componentState), state); err != nil {
		return state, err
	}
	return state, nil
}

type ComponentStatus int64

const (
	Unknown ComponentStatus = iota
	NotInstalled
	Installed
)

func (cs ComponentStatus) String() string {
	switch cs {
	case NotInstalled:
		return "Not Installed"
	case Installed:
		return "Installed"
	default:
		return "Unknown"
	}
}

var baseRunErr string = "unable to run component"
var cmpntNotFound string = "component does not exist"

type Artifact struct {
	OS        string `json:"os"`
	ARCH      string `json:"arch"`
	Signature string `json:"signature"`
	//Size ?
}

type Component struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	// will this component be accessible via the CLI
	CLICommand bool `json:"cli_command"`
	// @dhazekamp command_name required when CLICommand is true?
	CommandName string `json:"command_name"`

	// the component is a binary
	Binary bool `json:"binary"`

	// the component is a library, only provides content for the CLI or other components
	Library bool `json:"library"`

	// the component is standalone, should be available in $PATH
	Standalone bool `json:"standalone"`

	Artifacts []Artifact `json:"artifacts"`
}

// @dhazekamp validate component name
func (cmpnt Component) Path() (string, error) {
	cacheDir, err := cacheDir()
	if err != nil {
		return "", err
	}
	cmpntPath := path.Join(cacheDir, cmpnt.Name, cmpnt.Name)
	if !fileExists(cmpntPath) {
		return cmpntPath, errors.New(cmpntNotFound)
	}
	return cmpntPath, nil
}

func (cmpnt Component) Status() ComponentStatus {
	_, err := cmpnt.Path()
	if err == nil {
		return Installed
	}
	if err.Error() == cmpntNotFound {
		return NotInstalled
	}
	return Unknown
}

func (cmpnt Component) GetArtifact() (Artifact, error) {
	for _, a := range cmpnt.Artifacts {
		if a.OS == runtime.GOOS && a.ARCH == runtime.GOARCH {
			return a, nil
		}
	}
	return Artifact{}, errors.New("component not supported on this platform")
}

// @dhazekamp replace sha256 validation with minisign
func (cmpnt Component) isVerified() (bool, error) {
	var baseErr string = "unable to verify component"

	// get artifact
	a, err := cmpnt.GetArtifact()
	if err != nil {
		return false, errors.Wrap(err, baseErr)
	}
	// verify artifact has a signature
	if a.Signature == "" {
		return false, errors.Wrap(errors.New("component has no signature"), baseErr)
	}
	// get component path
	cmpntPath, err := cmpnt.Path()
	if err != nil {
		return false, errors.Wrap(err, baseErr)
	}
	// open the component
	f, err := os.Open(cmpntPath)
	if err != nil {
		return false, errors.Wrap(errors.Wrap(err, "unable to open component"), baseErr)
	}
	defer f.Close()
	// hash the component
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, errors.Wrap(errors.Wrap(err, "unable to hash component"), baseErr)
	}
	// validate the hash
	if a.Signature != hex.EncodeToString(h.Sum(nil)) {
		return false, errors.Wrap(errors.New("signature mismatch"), baseErr)
	}
	return true, nil
}

func (cmpnt Component) run(cmd *exec.Cmd) error {
	if !cmpnt.Binary {
		return errors.Wrap(errors.New("component is not a binary"), baseRunErr)
	}

	// verify component
	if isVerified, err := cmpnt.isVerified(); !isVerified {
		return errors.Wrap(err, baseRunErr)
	}

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, baseRunErr)
	}
	return nil
}

// RunAndOutput runs the command and outputs to os.Stdout and os.Stderr
func (cmpnt Component) RunAndOutput(args []string, stdin io.Reader) error {
	loc, err := cmpnt.Path()
	if err != nil {
		return errors.Wrap(err, baseRunErr)
	}

	cmd := exec.Command(loc, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmpnt.run(cmd)
}

// RunAndReturn runs the command and returns its standard output and standard error
func (cmpnt Component) RunAndReturn(args []string, stdin io.Reader) (
	stdout string,
	stderr string,
	err error,
) {
	var outBuff, errBuff bytes.Buffer

	loc, err := cmpnt.Path()
	if err != nil {
		err = errors.Wrap(err, baseRunErr)
		return
	}

	cmd := exec.Command(loc, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = stdin
	cmd.Stdout = &outBuff
	cmd.Stderr = &errBuff

	err = cmpnt.run(cmd)

	stdout, stderr = outBuff.String(), errBuff.String()
	return
}

// @hazekamp figure out LibraryComponent (if component is a library how do we interact with it)
