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
		Artifacts: []Artifact{
			Artifact{
				OS:        "darwin",
				ARCH:      "amd64",
				Signature: "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
				Version:   *mockVersion,
			},
			Artifact{
				OS:        "darwin",
				ARCH:      "arm64",
				Signature: "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
				Version:   *mockVersion,
			},
			Artifact{
				OS:        "linux",
				ARCH:      "amd64",
				Signature: "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
				Version:   *mockVersion,
			},
			Artifact{
				OS:        "linux",
				ARCH:      "arm64",
				Signature: "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
				Version:   *mockVersion,
			},
		},
	}
)

func ensureMockComponent(version string) (string, error) {
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

	componentDir, err := ensureMockComponent("")
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

type getArtifactTest struct {
	Name      string
	Component Component
	Version   string
	Expected  Artifact
	Error     error
}

var getArtifactTests = []getArtifactTest{
	getArtifactTest{
		Name:      "NoCurrentVersion",
		Component: Component{Name: "lacework-mock-component2"},
		Error:     errors.New("component does not exist"),
	},
	getArtifactTest{
		Name:      "UnsupportedPlatform",
		Component: Component{Name: "lacework-mock-component"},
		Version:   "1.0.0",
		Error:     errors.New("artifact not found"),
	},
	getArtifactTest{
		Name:      "UnsupportedPlatform",
		Component: mockComponent,
		Version:   "0.1.0",
		Error:     nil,
	},
}

func TestGetArtifact(t *testing.T) {
	for _, gat := range getArtifactTests {
		t.Run(gat.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(gat.Version)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

			actualArtifact, actualError := gat.Component.getArtifact()
			if gat.Error != nil {
				assert.Equal(t, gat.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
				assert.NotEqual(t, "", actualArtifact.Signature)
			}
		})
	}
}

type isVerifiedTest struct {
	Name      string
	Component Component
	Version   string
	Expected  bool
	Error     error
}

var isVerifiedTests = []isVerifiedTest{
	isVerifiedTest{
		Name:      "ArtifactNotFound",
		Component: Component{Name: "lacework-mock-component"},
		Version:   "0.1.0",
		Error: errors.New(
			"unable to verify component: artifact not found"),
	},
	isVerifiedTest{
		Name: "NoSignature",
		Component: Component{
			Name: "lacework-mock-component",
			Artifacts: []Artifact{
				Artifact{
					OS:      "darwin",
					ARCH:    "amd64",
					Version: *mockVersion,
				},
				Artifact{
					OS:      "darwin",
					ARCH:    "arm64",
					Version: *mockVersion,
				},
				Artifact{
					OS:      "linux",
					ARCH:    "amd64",
					Version: *mockVersion,
				},
				Artifact{
					OS:      "linux",
					ARCH:    "arm64",
					Version: *mockVersion,
				},
			},
		},
		Version: "0.1.0",
		Error:   errors.New("component has no signature"),
	},
	isVerifiedTest{
		Name: "Mismatch",
		Component: Component{
			Name: "lacework-mock-component",
			Artifacts: []Artifact{
				Artifact{
					OS:        "darwin",
					ARCH:      "amd64",
					Signature: "foo",
					Version:   *mockVersion,
				},
				Artifact{
					OS:        "darwin",
					ARCH:      "arm64",
					Signature: "foo",
					Version:   *mockVersion,
				},
				Artifact{
					OS:        "linux",
					ARCH:      "amd64",
					Signature: "foo",
					Version:   *mockVersion,
				},
				Artifact{
					OS:        "linux",
					ARCH:      "arm64",
					Signature: "foo",
					Version:   *mockVersion,
				},
			},
		},
		Version: "0.1.0",
		Error:   errors.New("signature mismatch"),
	},
	isVerifiedTest{
		Name:      "Verified",
		Component: mockComponent,
		Version:   "0.1.0\n",
		Expected:  true,
		Error:     nil,
	},
}

func TestIsVerified(t *testing.T) {
	for _, ivt := range isVerifiedTests {
		t.Run(ivt.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(ivt.Version)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

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

type runTest struct {
	Name      string
	Component Component
	Version   string
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
		Error: errors.New(
			"unable to run component: unable to verify component: component does not exist"),
	},
	runTest{
		Name:      "OK",
		Component: mockComponent,
		Version:   "0.1.0",
		Error:     nil,
	},
	runTest{
		Name:      "Error",
		Component: mockComponent,
		Version:   "0.1.0",
		Error:     errors.New(`unable to run component: exec: "nosuchcmd": executable file not found in $PATH`),
	},
}

func TestRun(t *testing.T) {
	for _, rt := range runTests {
		t.Run(rt.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(rt.Version)
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
