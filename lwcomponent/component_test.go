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

package lwcomponent_test

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	capturer "github.com/lacework/go-sdk/internal/capturer"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/lacework/go-sdk/lwcomponent"
)

var (
	mockVersion, _ = semver.NewVersion("1.0.0")
	mockComponent  = lwcomponent.Component{
		Name:          "lacework-mock-component",
		Description:   "This is a mock component",
		LatestVersion: *mockVersion,
		Type:          "BINARY",
	}
	mockComponent2 = lwcomponent.Component{
		Name: "lacework-mock-component2",
	}
	//go:embed test_resources/hello-world2.sh
	helloWorld []byte
	//go:embed test_resources/hello-world2.sig
	helloWorldSig string
)

func ensureMockComponent(version, signature string) (string, error) {
	cPath, err := mockComponent.Path()
	if err.Error() != "component not found on disk" {
		return "", err
	}

	if version != "" {
		ver := mockComponent.LatestVersion
		newv, _ := semver.NewVersion(version)
		mockComponent.LatestVersion = *newv
		defer func() {
			mockComponent.LatestVersion = ver
		}()
	}
	if err := mockComponent.WriteVersion(); err != nil {
		return "", err
	}

	if signature != "" {
		if err := mockComponent.WriteSignature([]byte(signature)); err != nil {
			return "", err
		}
	}

	return cPath, os.WriteFile(cPath, helloWorld, 0744)
}

func TestGetComponent(t *testing.T) {
	componentDir, err := ensureMockComponent("", "")
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
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Components",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "{\"components\": [],\"version\": \"0.1.0\"}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	state, err := lwcomponent.LoadState(c)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(state.Components))
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
		Component: mockComponent2,
		Error:     errors.New("component version file does not exist"),
	},
	// currentVersionTest{
	// Name:      "bad",
	// Component: mockComponent2,
	// Version:   "not a semver",
	// Error:     errors.New("unable to parse component version"),
	// },
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
			componentDir, err := ensureMockComponent(cvt.Version, "")
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

type currentSignatureTest struct {
	Name      string
	Component lwcomponent.Component
	Signature string
	Expected  []byte
	Error     error
}

var currentSignatureTests = []currentSignatureTest{
	currentSignatureTest{
		Name:      "notfound",
		Component: mockComponent2,
		Error:     errors.New("component signature file does not exist"),
	},
	currentSignatureTest{
		Name:      "bad",
		Component: mockComponent,
		Signature: "-",
		Expected:  []byte{},
		Error:     errors.New("unable to decode component signature"),
	},
	currentSignatureTest{
		Name:      "ok",
		Component: mockComponent,
		Signature: base64.StdEncoding.EncodeToString([]byte("mysig")),
		Expected:  []byte("mysig"),
		Error:     nil,
	},
}

func TestCurrentSignature(t *testing.T) {
	for _, cst := range currentSignatureTests {
		t.Run(cst.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent("", cst.Signature)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

			actualCV, actualError := cst.Component.SignatureFromDisk()
			assert.Equal(t, cst.Expected, actualCV)
			if cst.Error != nil {
				assert.Equal(t, cst.Error.Error(), actualError.Error())
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
			componentDir, err := ensureMockComponent(uat.Version, "")
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

			actual, err := uat.Component.UpdateAvailable()
			assert.NoError(t, err)
			assert.Equal(t, uat.Expected, actual)
		})
	}
}

func TestComponentStatus(t *testing.T) {
	componentDir, err := ensureMockComponent("", "")
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
	Signature      string
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
		Signature:      helloWorldSig,
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
			componentDir, err := ensureMockComponent(rart.Version, rart.Signature)
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
	Signature string
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
		Signature: helloWorldSig,
		Args:      []string{"World"},
		Stdin:     strings.NewReader("Error"),
		Expected:  "Hello World!\nHello !\n",
		Error:     nil,
	},
}

func TestRunAndOutput(t *testing.T) {
	for _, raot := range RunAndOutputTests {
		t.Run(raot.Name, func(t *testing.T) {
			componentDir, err := ensureMockComponent(raot.Version, raot.Signature)
			if err != nil {
				assert.FailNow(t, "Unable to ensureMockComponent")
			}
			defer os.RemoveAll(componentDir)

			actual := capturer.CaptureOutput(func() {
				// @afiune fake STDIN
				actualError := raot.Component.RunAndOutput(raot.Args)

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
