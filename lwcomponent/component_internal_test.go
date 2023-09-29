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
	"strings"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/internal/cache"
)

var (
	mockVersion, _ = semver.NewVersion("0.1.0")
	mockComponent  = Component{
		Name:          "lacework-mock-component",
		Description:   "This is a mock mock component",
		LatestVersion: *mockVersion,
		Breadcrumbs: Breadcrumbs{
			InstallationMessage: "Successfully installed!",
			UpdateMessage:       "Now supports feature 3",
		},
		Type: "BINARY",
		Artifacts: []Artifact{{
			OS:            "darwin",
			ARCH:          "arm64",
			URL:           "https://someurl.com/",
			Signature:     "abcdef123456789",
			Version:       "0.0.1",
			UpdateMessage: "Now supports feature 1",
		}, {
			OS:            "darwin",
			ARCH:          "arm64",
			URL:           "https://someurl.com/",
			Signature:     "abcdef123456789",
			Version:       "0.0.2",
			UpdateMessage: "Now supports feature 2",
		}, {
			OS:            "windows",
			ARCH:          "arm64",
			URL:           "https://someurl.com/",
			Signature:     "abcdef123456789",
			Version:       "0.0.2",
			UpdateMessage: "Now supports feature 2",
		}, {
			OS:            "darwin",
			ARCH:          "arm64",
			URL:           "https://someurl.com/",
			Signature:     "abcdef123456789",
			Version:       "0.1.0",
			UpdateMessage: "Now supports feature 3",
		}, {
			OS:            "darwin",
			ARCH:          "arm64",
			URL:           "https://someurl.com/",
			Signature:     "abcdef123456789",
			Version:       "0.1.1",
			UpdateMessage: "Now supports feature 4",
		}},
	}
	//go:embed test_resources/hello-world.sh
	helloWorld []byte
	//go:embed test_resources/hello-world.sig
	helloWorldSig           string
	helloWorldSigDecoded, _ = base64.StdEncoding.DecodeString(helloWorldSig)
)

func ensureMockComponent(version, signature string) (string, error) {
	cmpntPath, err := mockComponent.Path()
	if err.Error() != "component not found on disk" {
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
	{
		"NotExists",
		Component{
			Name: "no-such-component",
			Type: "STANDALONE",
		},
		errors.New("component not found on disk"),
	},
	{
		"Exists",
		mockComponent,
		nil,
	},
}

func TestPath(t *testing.T) {
	cacheDir, err := cache.CacheDir()
	if err != nil {
		assert.FailNow(t, "Unable to determine cacheDir")
	}

	componentDir, err := ensureMockComponent("", "")
	if err != nil {
		assert.FailNowf(t, "Unable to ensureMockComponent.", "Error: %s", err.Error())
	}
	defer os.RemoveAll(componentDir)

	for _, lpt := range pathTests {
		t.Run(lpt.Name, func(t *testing.T) {
			expectedLoc := path.Join(cacheDir, "components", lpt.Component.Name, lpt.Component.Name)
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
	{
		Name: "NoSignature",
		Component: Component{
			Name: "lacework-mock-component",
		},
		Version: "0.1.0",
		Error:   errors.New("component signature file does not exist"),
	},
	{
		Name: "Mismatch",
		Component: Component{
			Name: "lacework-mock-component",
		},
		Version:   "0.1.0",
		Signature: base64.StdEncoding.EncodeToString([]byte("blah blah blah")),
		Error:     errors.New("unable to parse signature"),
	},
	{
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
	{
		Name:      "IsNotBinary",
		Component: Component{Name: "IsNotBinary"},
		Version:   "0.1.0",
		Error:     errors.New("unable to run component: component IsNotBinary is not a binary"),
	},
	{
		Name:      "IsNotVerified",
		Component: Component{Name: "IsNotVerified", Type: "BINARY"},
		Version:   "0.1.0",
		Error:     errors.New("unable to run component: component signature file does not exist"),
	},
	{
		Name:      "OK",
		Component: mockComponent,
		Version:   "0.1.0",
		Signature: helloWorldSig,
		Error:     nil,
	},
	{
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

func TestMakeUpdateMessage(t *testing.T) {
	from, _ := semver.NewVersion("0.0.1")
	to, _ := semver.NewVersion("0.1.0")
	message := mockComponent.MakeUpdateMessage(*from, *to)
	assert.Contains(t, message, "from 0.0.1 to 0.1.0")
	assert.NotContains(t, message, "Now supports feature 1")
	assert.Contains(t, message, "Now supports feature 2")
	assert.Contains(t, message, "Now supports feature 3")
	assert.NotContains(t, message, "Now supports feature 4")
}

func TestListVersionsWithInstalled(t *testing.T) {
	installed, _ := semver.NewVersion("0.1.0")
	versions := mockComponent.ListVersions(installed)
	assert.Equal(t, 1, strings.Count(versions, "0.0.1"))
	assert.Equal(t, 1, strings.Count(versions, "0.0.2"))
	assert.Equal(t, 1, strings.Count(versions, "0.1.0 (installed)"))
	assert.Equal(t, 1, strings.Count(versions, "0.1.1"))
}

func TestListVersionsWithMissingInstalled(t *testing.T) {
	installed, _ := semver.NewVersion("0.0.0")
	versions := mockComponent.ListVersions(installed)
	assert.Equal(t, 1, strings.Count(versions, "0.0.1"))
	assert.Equal(t, 1, strings.Count(versions, "0.0.2"))
	assert.Equal(t, 1, strings.Count(versions, "0.1.0"))
	assert.Equal(t, 1, strings.Count(versions, "0.1.1"))
	assert.Contains(t, versions, "currently installed version 0.0.0 is no longer available")
}

func TestListVersionsWithoutInstalled(t *testing.T) {
	versions := mockComponent.ListVersions(nil)
	assert.Equal(t, 1, strings.Count(versions, "0.0.1"))
	assert.Equal(t, 1, strings.Count(versions, "0.0.2"))
	assert.Equal(t, 1, strings.Count(versions, "0.1.0"))
	assert.Equal(t, 1, strings.Count(versions, "0.1.1"))
	assert.NotContains(t, versions, "installed")
}
