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

	"github.com/lacework/go-sdk/lwcomponent"
	"github.com/stretchr/testify/assert"
	capturer "github.com/zenizh/go-capturer"
)

var (
	mockComponent = lwcomponent.Component{
		Name:        "lacework-mock-component",
		Description: "This is a mock component",
		Version:     "0.1.0",
		Signature:   "d69669aadfa69e5a212c83d52d9e5ca257f6c8bfedf82f8e34eb9523e27e3a3f",
		CLICommand:  false,
		CommandName: "",
		Binary:      true,
		Library:     false,
		Standalone:  false,
	}
	mockComponent2 = lwcomponent.Component{
		Name: "lacework-mock-component2",
	}
)

func ensureMockComponent() (string, error) {
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
		return componentPath, err
	}
	return componentPath, os.WriteFile(componentPath, []byte(componentData), 0777)
}

func TestLoadState(t *testing.T) {
	_, err := lwcomponent.LoadState()
	assert.Nil(t, err)
}

func TestGetComponent(t *testing.T) {
	componentPath, err := ensureMockComponent()
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.Remove(componentPath)

	tempState := new(lwcomponent.State)
	tempState.Components = []lwcomponent.Component{
		mockComponent,
		mockComponent2,
	}

	assert.Equal(t, mockComponent, *tempState.GetComponent("lacework-mock-component"))
	assert.Nil(t, tempState.GetComponent("no-such-component"))
}

func TestComponentStatus(t *testing.T) {
	componentPath, err := ensureMockComponent()
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.Remove(componentPath)

	assert.Equal(t, "Installed", mockComponent.Status().String())
	assert.Equal(t, "Not Installed", mockComponent2.Status().String())
}

type RunAndReturnTest struct {
	Name           string
	Component      lwcomponent.Component
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
		Args:           []string{"World"},
		Stdin:          strings.NewReader("Error"),
		ExpectedStdout: "Hello World!\n",
		ExpectedStderr: "Hello Error!\n",
		Error:          nil,
	},
}

func TestRunAndReturn(t *testing.T) {
	componentPath, err := ensureMockComponent()
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.Remove(componentPath)

	for _, rart := range RunAndReturnTests {
		t.Run(rart.Name, func(t *testing.T) {
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
	Args      []string
	Stdin     io.Reader
	Expected  string
	Error     error
}

var RunAndOutputTests = []RunAndOutputTest{
	RunAndOutputTest{
		Name:      "OK",
		Component: mockComponent,
		Args:      []string{"World"},
		Stdin:     strings.NewReader("Error"),
		Expected:  "Hello World!\nHello Error!\n",
		Error:     nil,
	},
}

func TestRunAndOutput(t *testing.T) {
	componentPath, err := ensureMockComponent()
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.Remove(componentPath)

	for _, raot := range RunAndOutputTests {
		t.Run(raot.Name, func(t *testing.T) {
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
