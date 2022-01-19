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
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

type State struct {
	Version    string      `json:"version"`
	Components []Component `json:"components"`
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

// @dhazekamp need to avoid cache poisoning with respect retrieving component signature
func LoadState() (*State, error) {
	state := new(State)

	cacheDir, err := cacheDir()
	if err != nil {
		return state, err
	}

	componentsFile := path.Join(cacheDir, "components")
	if fileExists(componentsFile) {
		componentState, err := ioutil.ReadFile(componentsFile)
		if err != nil {
			return state, err
		}

		err = json.Unmarshal(componentState, state)
		if err != nil {
			return state, err
		}
	}

	return state, nil
}

var baseRunErr string = "unable to run component"

type Component struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	Signature   string `json:"signature"`
	//Size ?

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
}

// @dhazekamp validate component name
func (cmpnt Component) Path() (string, error) {
	cacheDir, err := cacheDir()
	if err != nil {
		return "", err
	}
	cmpntPath := path.Join(cacheDir, cmpnt.Name, cmpnt.Name)
	if !fileExists(cmpntPath) {
		return cmpntPath, errors.New("component does not exist")
	}
	return cmpntPath, nil
}

// @dhazekamp replace sha256 validation with minisign
func (cmpnt Component) isVerified() (bool, error) {
	baseErr := "unable to verify component"

	// ensure we have a component signature
	if cmpnt.Signature == "" {
		return false, errors.Wrap(errors.New("component has no signature"), baseErr)
	}
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
	if cmpnt.Signature != hex.EncodeToString(h.Sum(nil)) {
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
