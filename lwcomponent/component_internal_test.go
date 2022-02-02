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
	_ "embed"
	"encoding/base64"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	mockVersion, _ = semver.NewVersion("0.1.0")
	mockComponent  = Component{
		Name:          "lacework-mock-component",
		Description:   "This is a mock mock component",
		LatestVersion: *mockVersion,
		CLICommand:    false,
		CommandName:   "",
		Binary:        true,
		Library:       false,
		Standalone:    false,
	}
	//go:embed test_resources/hello-world.sh
	helloWorld []byte
	//go:embed test_resources/hello-world.sig
	helloWorldSig           string
	helloWorldSigDecoded, _ = base64.StdEncoding.DecodeString(helloWorldSig)
)

func ensureMockComponent(version, signature string) (string, error) {
	cmpntPath, err := mockComponent.Path()
	if err.Error() != "component does not exist" {
		return "", err
	}
	cmpntDir, _ := path.Split(cmpntPath)
	err = os.MkdirAll(cmpntDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	if version != "" {
		cmpntVersionPath := path.Join(cmpntDir, ".version")
		if err = os.WriteFile(cmpntVersionPath, []byte(version), 0666); err != nil {
			return "", err
		}
	}

	if signature != "" {
		cmpntSignaturePath := path.Join(cmpntDir, ".signature")
		if err = os.WriteFile(cmpntSignaturePath, []byte(signature), 0666); err != nil {
			return "", err
		}
	}

	return cmpntDir, os.WriteFile(cmpntPath, helloWorld, 0777)
}

type pathTest struct {
	Name      string
	Component Component
	Error     error
}

var pathTests = []pathTest{
	pathTest{
		"NotExists",
		Component{
			Name:       "no-such-component",
			Binary:     true,
			Standalone: false,
		},
		errors.New("component does not exist"),
	},
	pathTest{
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

	componentDir, err := ensureMockComponent("", "")
	if err != nil {
		assert.FailNow(t, "Unable to ensureMockComponent")
	}
	defer os.RemoveAll(componentDir)

	for _, lpt := range pathTests {
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

type isVerifiedTest struct {
	Name      string
	Component Component
	Version   string
	Signature string
	Error     error
}

var isVerifiedTests = []isVerifiedTest{
	isVerifiedTest{
		Name: "NoSignature",
		Component: Component{
			Name: "lacework-mock-component",
		},
		Version: "0.1.0",
		Error:   errors.New("component signature file does not exist"),
	},
	isVerifiedTest{
		Name: "Mismatch",
		Component: Component{
			Name: "lacework-mock-component",
		},
		Version:   "0.1.0",
		Signature: base64.StdEncoding.EncodeToString([]byte("blah blah blah")),
		Error:     errors.New("unable to verify component: unable to parse signature"),
	},
	isVerifiedTest{
		Name:      "Verified",
		Component: mockComponent,
		Version:   "0.1.0\n",
		Signature: helloWorldSig,
		Error:     nil,
	},
}

func TestIsVerified(t *testing.T) {
	for _, ivt := range isVerifiedTests {
		t.Run(ivt.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(ivt.Version, ivt.Signature)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

			actualError := ivt.Component.isVerified()
			if ivt.Error != nil {
				assert.Equal(t, ivt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
			}
		})
	}
}

type runTest struct {
	Name      string
	Component Component
	Version   string
	Signature string
	Cmd       *exec.Cmd
	Error     error
}

var runTests = []runTest{
	runTest{
		Name:      "IsNotBinary",
		Component: Component{},
		Version:   "0.1.0",
		Error:     errors.New("unable to run component: component is not a binary"),
	},
	runTest{
		Name:      "IsNotVerified",
		Component: Component{Binary: true},
		Version:   "0.1.0",
		Error:     errors.New("unable to run component: component does not exist"),
	},
	runTest{
		Name:      "OK",
		Component: mockComponent,
		Version:   "0.1.0",
		Signature: helloWorldSig,
		Error:     nil,
	},
	runTest{
		Name:      "Error",
		Component: mockComponent,
		Version:   "0.1.0",
		Signature: helloWorldSig,
		Error:     errors.New(`unable to run component: exec: "nosuchcmd": executable file not found in $PATH`),
	},
}

func TestRun(t *testing.T) {
	for _, rt := range runTests {
		t.Run(rt.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(rt.Version, rt.Signature)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

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
