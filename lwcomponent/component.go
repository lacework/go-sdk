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

// A development kit for the cloud based of modular components.
package lwcomponent

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"aead.dev/minisign"
	"github.com/Masterminds/semver"
	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/cache"
	"github.com/lacework/go-sdk/internal/file"
	"github.com/lacework/go-sdk/internal/archive"
)

// State holds the components specification
//
// You can load the state from the Lacework API server by passing an `api.Client`.
//
//	client, err := api.NewClient(account, opts...)
//	cState, err := lwcomponent.LoadState(client)
//
// Or, you can load the state from the local storage.
//
//	cState, err := lwcomponent.LocalState()
type State struct {
	Version    string      `json:"version"`
	Components []Component `json:"components"`
}

// LoadState loads the state from the Lacework API server
func LoadState(client *api.Client) (*State, error) {
	if client != nil {
		s := new(State)

		// load remote components - this involves a network call that may fail due
		// to network issues or rate limits, so we retry this if it fails
		err := backoff.Retry(func() error {
			return client.RequestDecoder("GET", "v2/Components", nil, s)
		}, backoffStrategy())
		if err != nil {
			return s, err
		}

		// load local components
		s.loadComponentsFromDisk()

		// load dev components
		s.loadDevComponents()

		return s, s.WriteState()
	}
	return nil, errors.New("invalid api client")
}

func backoffStrategy() *backoff.ExponentialBackOff {
	strategy := backoff.NewExponentialBackOff()
	strategy.InitialInterval = 2 * time.Second
	strategy.MaxElapsedTime = 1 * time.Minute
	return strategy
}

// loadComponentsFromDisk will load all component from disk (local)
func (s *State) loadComponentsFromDisk() {
	if dir, err := Dir(); err == nil {
		components, err := os.ReadDir(dir)
		if err != nil {
			return
		}

		// traverse components dir
		for _, c := range components {
			if !c.IsDir() {
				continue
			}

			// load components that are not already registered
			if _, found := s.GetComponent(c.Name()); !found {
				component := Component{Name: c.Name()}

				// verify that the directory is a component, that means that the
				// directory contains either a '.dev' file or, both '.version'
				// and '.signature' files
				//
				// TODO @afiune maybe we should deploy a .specs file?
				err := component.isVerified()

				if component.UnderDevelopment() || err == nil {
					s.Components = append(s.Components, component)
				}
			}
		}
	}
}

// loadDevComponents will load all components that are under development
func (s *State) loadDevComponents() {
	for i := range s.Components {
		if s.Components[i].UnderDevelopment() {
			// existing component being developed
			if err := s.Components[i].loadDevSpecs(); err != nil {
				s.Components[i].Description = err.Error()
			}
		}
	}

	if devComponent := os.Getenv("LW_CDK_DEV_COMPONENT"); devComponent != "" {
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

	stateBytes, err := os.ReadFile(stateFile)
	if err != nil {
		return state, err
	}
	if err := json.Unmarshal(stateBytes, state); err != nil {
		return state, err
	}

	// load local components
	state.loadComponentsFromDisk()

	// load dev components
	state.loadDevComponents()

	return state, nil
}

// GetComponent returns the pointer of a component, if the component is not
// found, this function will return a `nil` pointer and `false`
//
// Usage:
//
//	component, found := state.GetComponent(name)
//
//	if !found {
//		fmt.Println("Component %s not found", name)
//	}
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

	if err := os.WriteFile(stateFile, buf.Bytes(), 0644); err != nil {
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

func (s State) Install(component *Component, version string) error {
	rPath, err := component.RootPath()
	if err != nil {
		return err
	}

	// verify development mode
	if component.UnderDevelopment() {
		p, _ := component.Path() // @afiune we don't care if the component exists or not
		msg := "components under development can't be installed.\n\n" +
			"Deploy the component manually at '" + p + "'"
		return errors.New(msg)
	}

	// @afiune verify if component is in latest
	artifact, found := component.ArtifactForRunningHost(version)
	if !found {
		return errors.Errorf(
			"could not find an artifact for version %s on the current platform (%s/%s)",
			version, runtime.GOOS, runtime.GOARCH,
		)
	}

	path, err := component.Path()
	// it is ok if the component is not found yet
	if err != nil && !IsNotFound(err) {
		return err
	}

	// download to temp dir, this must be different from staging dir
	// because the archive may have the same name as the extracted folder
	downloadDir, err := os.MkdirTemp("", "cdk-component-stage-download")
	if err != nil {
		return err
	}
	defer os.RemoveAll(downloadDir)
	downloadPath := filepath.Join(downloadDir, component.Name)

	// Stage to temp dir before installing
	stagingDir, err := os.MkdirTemp("", "cdk-component-stage-extract")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stagingDir)
	
	stagingPath := filepath.Join(stagingDir, component.Name)

	// Slow S3 downloads
	downloadTimeout := 0 * time.Minute
	switch component.Name {
	case "sast":
		downloadTimeout = 5 * time.Minute
	}

	err = DownloadFile(downloadPath, artifact.URL, time.Duration(downloadTimeout))
	if err != nil {
		return errors.Wrap(err, "unable to download component artifact")
	}

	
	// if the component is a tgz archive unpack it, otherwise leave it alone
	if err = archive.DetectTGZAndUnpack(downloadPath, stagingDir); err != nil {
		return err
	}

	//if the component was not an archive then nothing was created in the staging dir
	//we must move it over
	if _, err := os.Stat(stagingPath); errors.Is(err, os.ErrNotExist) {
		err = os.Rename(downloadPath, stagingPath)
		if err != nil {
			return err
		}
	}
	
	// if the component is not an archive make a dir for it to live in
	f, err := os.Stat(stagingPath)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		if err := os.MkdirAll(rPath, os.ModePerm); err != nil {
			return err
		}
		//move the component from the staging dir to it's path
		if err = os.Rename(stagingPath, path); err != nil {
			return err
		}
	}	else {
			//move the component from the staging dir to it's root path
			if err = os.Rename(stagingPath, rPath); err != nil {
				return err
			}
	}
	
	if err := component.WriteVersion(artifact.Version); err != nil {
		return err
	}

	// @afiune check 1) cross-platform and 2) correct permissions
	// if the file has permissions already, can we avoid this?
	if component.IsExecutable() {
		if err := os.Chmod(path, 0744); err != nil {
			return errors.Wrap(err, "unable to make component executable")
		}
	}

	return nil
}

func (s State) Verify(component *Component, version string) error {
	artifact, found := component.ArtifactForRunningHost(version)
	if !found {
		return errors.Errorf(
			"could not find an artifact for version %s on the current platform (%s/%s)",
			version, runtime.GOOS, runtime.GOARCH,
		)
	}

	if err := component.WriteSignature([]byte(artifact.Signature)); err != nil {
		return err
	}

	rPath, err := component.RootPath()
	if err != nil {
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
	OS            string `json:"os"`
	ARCH          string `json:"arch"`
	URL           string `json:"url,omitempty"`
	Signature     string `json:"signature"`
	Version       string `json:"version"`
	UpdateMessage string `json:"updateMessage"`
	//Size ?
}

// Components should leave a trail/crumb after installation or update,
// these messages will be shown by the Lacework CLI
type Breadcrumbs struct {
	InstallationMessage string `json:"installationMessage,omitempty"`
	UpdateMessage       string `json:"updateMessage,omitempty"`
}

// Component can be a command-line tool, a new command that extends the Lacework CLI, or a library that
// contains files used by another Lacework component.
type Component struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Type          Type           `json:"type"`
	LatestVersion semver.Version `json:"-"`
	Artifacts     []Artifact     `json:"artifacts"`
	Breadcrumbs   Breadcrumbs    `json:"breadcrumbs,omitempty"`

	// @dhazekamp command_name required when CLICommand is true?
	CommandName string `json:"command_name,omitempty"`
}

func (c *Component) UnmarshalJSON(data []byte) error {
	type ComponentAlias Component
	type T struct {
		*ComponentAlias     `json:",inline"`
		LatestVersionString string `json:"version"`
	}

	temp := &T{ComponentAlias: (*ComponentAlias)(c)}
	err := json.Unmarshal(data, temp)
	if err != nil {
		return err
	}

	latestVersion, err := semver.NewVersion(temp.LatestVersionString)
	if err != nil {
		return err
	}
	c.LatestVersion = *latestVersion

	return nil
}

func (c Component) MarshalJSON() ([]byte, error) {
	type ComponentAlias Component
	type T struct {
		ComponentAlias      `json:",inline"`
		LatestVersionString string `json:"version"`
	}

	obj := &T{ComponentAlias: (ComponentAlias)(c), LatestVersionString: c.LatestVersion.String()}
	return json.Marshal(obj)
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
	if c.UnderDevelopment() {
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

	return os.WriteFile(cvPath, signature, 0644)
}

// WriteVersion stores the component version on disk
func (c Component) WriteVersion(installed string) error {
	dir, err := c.RootPath()
	if err != nil {
		return err
	}

	cvPath := filepath.Join(dir, ".version")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if installed == "" {
		installed = c.LatestVersion.String()
	}
	return os.WriteFile(cvPath, []byte(installed), 0644)
}

// UpdateAvailable returns true if there is a newer version of the component
func (c Component)  UpdateAvailable() (bool, error) {
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
func (c Component) ArtifactForRunningHost(version string) (*Artifact, bool) {
	for _, artifact := range c.Artifacts {
		if artifact.OS == runtime.GOOS && artifact.ARCH == runtime.GOARCH && artifact.Version == version {
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
		devSpecsBytes, err := os.ReadFile(devSpecs)
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

// UnderDevelopment returns true if the component is under development
// that is, if the component root path has the '.dev' specs file or, if
// the environment variable 'LW_CDK_DEV_COMPONENT' matches the component name
func (c Component) UnderDevelopment() bool {
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
	if c.UnderDevelopment() {
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

func (c Component) EnterDevelopmentMode() error {
	if c.UnderDevelopment() {
		return errors.New("component already under development.")
	}

	dir, err := c.RootPath()
	if err != nil {
		return errors.New("unable to detect RootPath")
	}

	devSpecs := filepath.Join(dir, ".dev")
	if !file.FileExists(devSpecs) {
		// remove prod artifacts
		c.Artifacts = make([]Artifact, 0)

		// configure dev version
		cv, _ := semver.NewVersion("0.0.0-dev")
		c.LatestVersion = *cv

		// update description
		c.Description = fmt.Sprintf("(dev-mode) %s", c.Description)

		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(c); err != nil {
			return err
		}

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		return os.WriteFile(devSpecs, buf.Bytes(), 0644)
	}

	return nil
}

func (c Component) getVersionsAndBreadcrumbs() ([]*semver.Version, map[semver.Version]string) {
	versionToBreadcrumb := make(map[semver.Version]string)
	allVersions := make([]*semver.Version, 0)
	versionToBreadcrumb[c.LatestVersion] = c.Breadcrumbs.UpdateMessage
	allVersions = append(allVersions, &c.LatestVersion)
	for _, artifact := range c.Artifacts {
		parsedVersion, err := semver.NewVersion(artifact.Version)
		if err == nil {
			// We don't expect invalid versions from the server, but if we do get them
			// let's just recover and ignore that version rather than crashing.
			_, alreadySeen := versionToBreadcrumb[*parsedVersion]
			if !alreadySeen {
				versionToBreadcrumb[*parsedVersion] = artifact.UpdateMessage
				allVersions = append(allVersions, parsedVersion)
			}
		}
	}
	sort.Sort(semver.Collection(allVersions))
	return allVersions, versionToBreadcrumb
}

func (c Component) MakeUpdateMessage(from, to semver.Version) string {
	if from.LessThan(&to) {
		versions, breadcrumbs := c.getVersionsAndBreadcrumbs()
		updateMessage := ""
		for _, version := range versions {
			if version.LessThan(&from) || version.Equal(&from) {
				// We're before the breadcrumbs we care about, we don't include this one but keep going
				continue
			}
			if version.GreaterThan(&to) {
				// We're past the breadcrumbs we care about, we can stop iterating
				break
			}
			if breadcrumbs[*version] != "" {
				// We've found a breadcrumb in the range (from, to] which is one we care about
				updateMessage += "\n"
				updateMessage += breadcrumbs[*version]
			}
		}
		return fmt.Sprintf("Successfully upgraded component from %s to %s%s", from.String(), to.String(), updateMessage)
	}
	return fmt.Sprintf("Successfully downgraded component from %s to %s", from.String(), to.String())
}

func (c Component) ListVersions(installed *semver.Version) string {
	versions, _ := c.getVersionsAndBreadcrumbs()
	result := "The following versions of this component are available to install:"
	foundInstalled := false
	for _, version := range versions {
		result += "\n"
		result += " - " + version.String()
		if installed != nil && version.Equal(installed) {
			result += " (installed)"
			foundInstalled = true
		}
	}
	if installed != nil && !foundInstalled {
		result += fmt.Sprintf(
			"\n\nThe currently installed version %s is no longer available to install.",
			installed.String(),
		)
	}
	return result
}
