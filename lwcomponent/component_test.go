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

package lwcomponent_test

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	capturer "github.com/zenizh/go-capturer"

	"github.com/lacework/go-sdk/lwcomponent"
)

var (
	mockVersion, _ = semver.NewVersion("1.0.0")
	mockComponent  = lwcomponent.Component{
		Name:          "lacework-mock-component",
		Description:   "This is a mock component",
		LatestVersion: *mockVersion,
		CLICommand:    false,
		CommandName:   "",
		Binary:        true,
		Library:       false,
		Standalone:    false,
		Artifacts: []lwcomponent.Artifact{
			lwcomponent.Artifact{
				OS:        "darwin",
				ARCH:      "amd64",
				Signature: "d69669aadfa69e5a212c83d52d9e5ca257f6c8bfedf82f8e34eb9523e27e3a3f",
				Version:   *mockVersion,
			},
			lwcomponent.Artifact{
				OS:        "darwin",
				ARCH:      "arm64",
				Signature: "d69669aadfa69e5a212c83d52d9e5ca257f6c8bfedf82f8e34eb9523e27e3a3f",
				Version:   *mockVersion,
			},
			lwcomponent.Artifact{
				OS:        "linux",
				ARCH:      "amd64",
				Signature: "d69669aadfa69e5a212c83d52d9e5ca257f6c8bfedf82f8e34eb9523e27e3a3f",
				Version:   *mockVersion,
			},
			lwcomponent.Artifact{
				OS:        "linux",
				ARCH:      "arm64",
				Signature: "d69669aadfa69e5a212c83d52d9e5ca257f6c8bfedf82f8e34eb9523e27e3a3f",
				Version:   *mockVersion,
			},
		},
	}
	mockComponent2 = lwcomponent.Component{
		Name: "lacework-mock-component2",
	}
)

func ensureMockComponent(version string) (string, error) {
	componentData := `#!/bin/bash
echo "Hello $1!"
read line
echo "Hello $line!" >&2
`

	componentPath, err := mockComponent.Path()
	if err.Error() != "component does not exist" {
		return "", err
	}
	dir, _ := path.Split(componentPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	if version != "" {
		componentVersionPath := path.Join(dir, ".version")
		os.WriteFile(componentVersionPath, []byte(version), 0666)
		if err != nil {
			return "", err
		}
	}
	return dir, os.WriteFile(componentPath, []byte(componentData), 0777)
}

func TestGetComponent(t *testing.T) {
	componentDir, err := ensureMockComponent("")
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.RemoveAll(componentDir)

	tempState := new(lwcomponent.State)
	tempState.Components = []lwcomponent.Component{
		mockComponent,
		mockComponent2,
	}

	assert.Equal(t, mockComponent, *tempState.GetComponent("lacework-mock-component"))
	assert.Nil(t, tempState.GetComponent("no-such-component"))
}

func TestLoadState(t *testing.T) {
	_, err := lwcomponent.LoadState()
	assert.Nil(t, err)
}

type currentVersionTest struct {
	Name      string
	Component lwcomponent.Component
	Version   string
	Expected  *semver.Version
	Error     error
}

var currentVersionTests = []currentVersionTest{
	currentVersionTest{
		Name:      "notfound",
		Component: mockComponent,
		Error:     errors.New("component version file does not exist"),
	},
	currentVersionTest{
		Name:      "bad",
		Component: mockComponent,
		Version:   "not a semver",
		Error:     errors.New("unable to parse component version"),
	},
	currentVersionTest{
		Name:      "ok",
		Component: mockComponent,
		Version:   "1.0.0",
		Expected:  mockVersion,
		Error:     nil,
	},
}

func TestCurrentVersion(t *testing.T) {
	for _, cvt := range currentVersionTests {
		t.Run(cvt.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(cvt.Version)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

			actualCV, actualError := cvt.Component.CurrentVersion()
			assert.Equal(t, cvt.Expected, actualCV)
			if cvt.Error != nil {
				assert.Equal(t, cvt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
			}
		})
	}
}

type updateAvailableTest struct {
	Name      string
	Component lwcomponent.Component
	Version   string
	Expected  bool
}

var updateAvailableTests = []updateAvailableTest{
	updateAvailableTest{
		Name:      "error",
		Component: mockComponent,
	},
	updateAvailableTest{
		Name:      "yes",
		Component: mockComponent,
		Version:   "0.9.0",
		Expected:  true,
	},
	updateAvailableTest{
		Name:      "no",
		Component: mockComponent,
		Version:   "1.0.5",
	},
}

func TestUpdateAvailable(t *testing.T) {
	for _, uat := range updateAvailableTests {
		t.Run(uat.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(uat.Version)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

			actual := uat.Component.UpdateAvailable()
			assert.Equal(t, uat.Expected, actual)
		})
	}
}

func TestComponentStatus(t *testing.T) {
	componentDir, err := ensureMockComponent("")
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.RemoveAll(componentDir)

	assert.Equal(t, "Installed", mockComponent.Status().String())
	assert.Equal(t, "Not Installed", mockComponent2.Status().String())
}

type RunAndReturnTest struct {
	Name           string
	Component      lwcomponent.Component
	Version        string
	Args           []string
	Stdin          io.Reader
	ExpectedStdout string
	ExpectedStderr string
	Error          error
}

var RunAndReturnTests = []RunAndReturnTest{
	RunAndReturnTest{
		Name:           "OK",
		Component:      mockComponent,
		Version:        "1.0.0",
		Args:           []string{"World"},
		Stdin:          strings.NewReader("Error"),
		ExpectedStdout: "Hello World!\n",
		ExpectedStderr: "Hello Error!\n",
		Error:          nil,
	},
}

func TestRunAndReturn(t *testing.T) {
	for _, rart := range RunAndReturnTests {
		t.Run(rart.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(rart.Version)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

			actualStdout, actualStderr, actualError := rart.Component.RunAndReturn(rart.Args, rart.Stdin)

			if rart.Error != nil {
				assert.Equal(t, rart.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
			}
			assert.Equal(t, rart.ExpectedStdout, actualStdout)
			assert.Equal(t, rart.ExpectedStderr, actualStderr)
		})
	}
}

type RunAndOutputTest struct {
	Name      string
	Component lwcomponent.Component
	Version   string
	Args      []string
	Stdin     io.Reader
	Expected  string
	Error     error
}

var RunAndOutputTests = []RunAndOutputTest{
	RunAndOutputTest{
		Name:      "OK",
		Component: mockComponent,
		Version:   "1.0.0",
		Args:      []string{"World"},
		Stdin:     strings.NewReader("Error"),
		Expected:  "Hello World!\nHello Error!\n",
		Error:     nil,
	},
}

func TestRunAndOutput(t *testing.T) {
	for _, raot := range RunAndOutputTests {
		t.Run(raot.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(raot.Version)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

			actual := capturer.CaptureOutput(func() {
				actualError := raot.Component.RunAndOutput(raot.Args, raot.Stdin)

				if raot.Error != nil {
					assert.Equal(t, raot.Error.Error(), actualError.Error())
				} else {
					assert.Nil(t, actualError)
				}
			})
			assert.Equal(t, raot.Expected, actual)
		})
	}
}
