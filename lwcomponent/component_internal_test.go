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
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	mockComponent = Component{
		Name:        "lacework-mock-component",
		Description: "This is a mock mock component",
		Version:     "0.1.0",
		Status:      "installed",
		Signature:   "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
		CLICommand:  false,
		CommandName: "",
		Binary:      true,
		Library:     false,
		Standalone:  false,
	}
)

func ensureMockComponent() (string, error) {
	componentData := `#!/bin/bash
echo "Hello World!"
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

type PathTest struct {
	Name      string
	Component Component
	Error     error
}

var PathTests = []PathTest{
	PathTest{
		"NotExists",
		Component{
			Name:       "no-such-component",
			Binary:     true,
			Standalone: false,
		},
		errors.New("component does not exist"),
	},
	PathTest{
		"Exists",
		mockComponent,
		nil,
	},
}

func TestPath(t *testing.T) {
	cacheDir, err := cacheDir()
	if err != nil {
		assert.FailNow(t, "Unable to determine cacheDir")
	}

	componentPath, err := ensureMockComponent()
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.Remove(componentPath)

	for _, lpt := range PathTests {
		t.Run(lpt.Name, func(t *testing.T) {
			expectedLoc := path.Join(cacheDir, lpt.Component.Name, lpt.Component.Name)
			actualLoc, actualError := lpt.Component.Path()

			assert.Equal(t, actualLoc, expectedLoc)
			if lpt.Error != nil {
				assert.Equal(t, lpt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
			}
		})
	}
}

type IsVerifiedTest struct {
	Name      string
	Component Component
	Expected  bool
	Error     error
}

var IsVerifiedTests = []IsVerifiedTest{
	IsVerifiedTest{
		"NoSignature",
		Component{Name: "lacework-mock-component", Signature: ""},
		false,
		errors.New("unable to verify component: component has no signature"),
	},
	IsVerifiedTest{
		"NoPath",
		Component{Name: "lacework-notexists-component", Signature: "foo", Standalone: true},
		false,
		errors.New("unable to verify component: component does not exist"),
	},
	IsVerifiedTest{
		"Mismatch",
		Component{Name: "lacework-mock-component", Signature: "foo"},
		false,
		errors.New("unable to verify component: signature mismatch"),
	},
	IsVerifiedTest{
		"Verified",
		mockComponent,
		true,
		nil,
	},
}

func TestIsVerified(t *testing.T) {
	componentPath, err := ensureMockComponent()
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.Remove(componentPath)

	for _, ivt := range IsVerifiedTests {
		t.Run(ivt.Name, func(t *testing.T) {
			actualVerified, actualError := ivt.Component.isVerified()
			assert.Equal(t, ivt.Expected, actualVerified)
			if ivt.Error != nil {
				assert.Equal(t, ivt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
			}
		})
	}
}

type RunTest struct {
	Name      string
	Component Component
	Cmd       *exec.Cmd
	Error     error
}

var RunTests = []RunTest{
	RunTest{
		Name:      "IsNotBinary",
		Component: Component{},
		Error:     errors.New("unable to run component: component is not a binary"),
	},
	RunTest{
		Name:      "IsNotVerified",
		Component: Component{Binary: true},
		Error:     errors.New("unable to run component: unable to verify component: component has no signature"),
	},
	RunTest{
		Name:      "OK",
		Component: mockComponent,
		Error:     nil,
	},
	RunTest{
		Name:      "Error",
		Component: mockComponent,
		Error:     errors.New(`unable to run component: exec: "nosuchcmd": executable file not found in $PATH`),
	},
}

func TestRun(t *testing.T) {
	componentPath, err := ensureMockComponent()
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.Remove(componentPath)

	for _, rt := range RunTests {
		t.Run(rt.Name, func(t *testing.T) {
			if rt.Name == "Error" {
				rt.Cmd = exec.Command("nosuchcmd")
			} else {
				loc, err := mockComponent.Path()
				assert.Nil(t, err)
				rt.Cmd = exec.Command(loc)
			}

			actualError := rt.Component.run(rt.Cmd)
			if rt.Error != nil {
				assert.Equal(t, rt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
			}
		})
	}
}
