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

package cmd_test

import (
	"embed"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwcomponent"
)

var (
	mockLCLComponent lwcomponent.Component = lwcomponent.Component{
		Name:        "lacework-content-library",
		Description: "Lacework Content Library",
		Version:     "0.1.0",
		// @dhazekamp this only works for darwin-amd64 because we don't have per-package sigs
		CLICommand:  false,
		CommandName: "",
		Binary:      true,
		Library:     false,
		Standalone:  false,
		Artifacts: []lwcomponent.Artifact{
			lwcomponent.Artifact{
				OS:        "darwin",
				ARCH:      "amd64",
				Signature: "f35b88ae47f9778061543a33ad6799a0d16adf9af02f7527bdf053d81f9a0607",
			},
			lwcomponent.Artifact{
				OS:        "darwin",
				ARCH:      "arm64",
				Signature: "1033f26e03deed726311383ea175ad51632af516def4968b4b7fc39ec9a7d815",
			},
			lwcomponent.Artifact{
				OS:        "linux",
				ARCH:      "amd64",
				Signature: "5d75c71af8068832cf079ba697df54f2bd1bfdaea20bbe8c022d71ba6e420e10",
			},
			lwcomponent.Artifact{
				OS:        "linux",
				ARCH:      "arm64",
				Signature: "854a5e4ed7dd5f4d7c019a849e28da5e7ad944785fbeb37f525f295cb169d971",
			},
		},
	}
	mockLWComponentState lwcomponent.State = lwcomponent.State{
		Components: []lwcomponent.Component{
			mockLCLComponent,
		},
	}
	nonZeroLCLComponent lwcomponent.Component = lwcomponent.Component{
		Name:        "lacework-content-library",
		Description: "Lacework Content Library",
		Version:     "0.1.0",
		CLICommand:  false,
		CommandName: "",
		Binary:      true,
		Library:     false,
		Standalone:  false,
		Artifacts: []lwcomponent.Artifact{
			lwcomponent.Artifact{
				OS:        "darwin",
				ARCH:      "amd64",
				Signature: "99032d1fb22b1ea6119f5f728cadae6feddfa62a45a66d52f77558cffd80b7f2",
			},
			lwcomponent.Artifact{
				OS:        "darwin",
				ARCH:      "arm64",
				Signature: "99032d1fb22b1ea6119f5f728cadae6feddfa62a45a66d52f77558cffd80b7f2",
			},
			lwcomponent.Artifact{
				OS:        "linux",
				ARCH:      "amd64",
				Signature: "99032d1fb22b1ea6119f5f728cadae6feddfa62a45a66d52f77558cffd80b7f2",
			},
			lwcomponent.Artifact{
				OS:        "linux",
				ARCH:      "arm64",
				Signature: "99032d1fb22b1ea6119f5f728cadae6feddfa62a45a66d52f77558cffd80b7f2",
			},
		},
	}
	nonZeroLWComponentState lwcomponent.State = lwcomponent.State{
		Components: []lwcomponent.Component{
			nonZeroLCLComponent,
		},
	}
	noParseLCLComponent lwcomponent.Component = lwcomponent.Component{
		Name:        "lacework-content-library",
		Description: "Lacework Content Library",
		Version:     "0.1.0",
		CLICommand:  false,
		CommandName: "",
		Binary:      true,
		Library:     false,
		Standalone:  false,
		Artifacts: []lwcomponent.Artifact{
			lwcomponent.Artifact{
				OS:        "darwin",
				ARCH:      "amd64",
				Signature: "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
			},
			lwcomponent.Artifact{
				OS:        "darwin",
				ARCH:      "arm64",
				Signature: "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
			},
			lwcomponent.Artifact{
				OS:        "linux",
				ARCH:      "amd64",
				Signature: "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
			},
			lwcomponent.Artifact{
				OS:        "linux",
				ARCH:      "arm64",
				Signature: "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
			},
		},
	}
	noParseLWComponentState lwcomponent.State = lwcomponent.State{
		Components: []lwcomponent.Component{
			noParseLCLComponent,
		},
	}
	//go:embed test_resources/lacework-content-library/*
	mockLCLBinaries embed.FS

	malformedLCL cmd.LaceworkContentLibrary = cmd.LaceworkContentLibrary{
		Queries: map[string]cmd.LCLQuery{
			"my_query": cmd.LCLQuery{References: []cmd.LCLReference{
				cmd.LCLReference{},
			}},
		},
		Policies: map[string]cmd.LCLPolicy{
			"my_policy": cmd.LCLPolicy{References: []cmd.LCLReference{
				cmd.LCLReference{},
				cmd.LCLReference{},
			}},
		},
	}
)

type mockLCLPlacementType int64

const (
	mockLCLNoop mockLCLPlacementType = iota
	mockLCLReplaced
	mockLCLPlaced
)

func getMockLCLBinaryName() string {
	osName, arch := runtime.GOOS, runtime.GOARCH
	ext := ""
	if osName == "windows" {
		ext = ".exe"
	}
	return fmt.Sprintf("%s-%s-%s%s", mockLCLComponent.Name, osName, arch, ext)
}

func ensureMockLCL(b string) (mockLCLPlacementType, error) {
	placementType := mockLCLNoop

	cmpntPath, _ := mockLCLComponent.Path()
	if cmpntPath == "" {
		return placementType, errors.New("unable to ensure mock LCL component")
	}

	if mockLCLComponent.Status() == lwcomponent.Installed {
		os.Rename(cmpntPath, cmpntPath+".bak")
		placementType = mockLCLReplaced
	} else {
		placementType = mockLCLPlaced
		dir, _ := path.Split(cmpntPath)
		os.MkdirAll(dir, os.ModePerm)
	}

	cmpntBinaryPath := fmt.Sprintf("test_resources/lacework-content-library/%s", b)
	cmpntBytes, err := mockLCLBinaries.ReadFile(cmpntBinaryPath)
	if err != nil {
		return placementType, err
	}

	return placementType, os.WriteFile(cmpntPath, cmpntBytes, 0777)
}

func removeMockLCL(ept mockLCLPlacementType) {
	if ept == mockLCLNoop {
		return
	}

	cmpntPath, _ := mockLCLComponent.Path()
	if ept == mockLCLReplaced {
		os.Rename(cmpntPath+".bak", cmpntPath)
		return
	}
	os.Remove(cmpntPath)
}

func TestLoadLCLNotFound(t *testing.T) {
	state := *new(lwcomponent.State)

	// IsLCLInstalled
	assert.Equal(t, false, cmd.IsLCLInstalled(state))

	_, err := cmd.LoadLCL(state)
	if err == nil {
		assert.NotNil(t, err)
	} else {
		assert.Equal(
			t,
			"unable to load Lacework Content Library: Lacework Content Library is not installed",
			err.Error(),
		)
	}
}

func TestLoadLCLNonZero(t *testing.T) {
	ept, err := ensureMockLCL("lacework-content-library-nonzero.sh")
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	_, err = cmd.LoadLCL(nonZeroLWComponentState)
	if err == nil {
		assert.NotNil(t, err)
	} else {
		assert.Equal(
			t,
			"unable to load Lacework Content Library: unable to run component: exit status 1",
			err.Error(),
		)
	}
}

func TestLoadLCLNoParse(t *testing.T) {
	ept, err := ensureMockLCL("lacework-content-library-noparse.sh")
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	_, err = cmd.LoadLCL(noParseLWComponentState)
	if err == nil {
		assert.NotNil(t, err)
	} else {
		assert.Equal(
			t,
			"unable to load Lacework Content Library: invalid character 'H' looking for beginning of value",
			err.Error(),
		)
	}
}

func TestLoadLCLOK(t *testing.T) {
	ept, err := ensureMockLCL(getMockLCLBinaryName())
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	// IsLCLInstalled
	assert.Equal(t, true, cmd.IsLCLInstalled(mockLWComponentState))

	lcl, err := cmd.LoadLCL(mockLWComponentState)
	assert.Nil(t, err)
	_, ok := lcl.Queries["LW_Custom_Host_Activity_PotentialReverseShell"]
	assert.True(t, ok)
	_, ok = lcl.Policies["lwcustom-27"]
	assert.True(t, ok)
}

func TestGetQueryNoID(t *testing.T) {
	lcl := cmd.LaceworkContentLibrary{}
	_, actualError := lcl.GetQuery("")
	assert.Equal(t, "query ID must be provided", actualError.Error())
}

func TestGetQueryMalformed(t *testing.T) {
	malformedLCL := cmd.LaceworkContentLibrary{
		Queries: map[string]cmd.LCLQuery{
			"my_query": cmd.LCLQuery{References: []cmd.LCLReference{
				cmd.LCLReference{},
			}},
		},
	}

	_, actualError := malformedLCL.GetQuery("my_query")
	assert.Equal(t, "query exists but is malformed", actualError.Error())
}

func TestGetQueryOK(t *testing.T) {
	queryID := "LW_Custom_AWS_CTA_AuroraPasswordChange"
	ept, err := ensureMockLCL(getMockLCLBinaryName())
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	lcl, err := cmd.LoadLCL(mockLWComponentState)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	actualQuery, actualError := lcl.GetQuery(queryID)
	assert.Nil(t, actualError)
	assert.Contains(t, actualQuery, queryID)
}

func TestListPolicies(t *testing.T) {
	ept, err := ensureMockLCL(getMockLCLBinaryName())
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	lcl, err := cmd.LoadLCL(mockLWComponentState)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	policiesResponse, err := lcl.ListPolicies()
	assert.Nil(t, err)
	assert.Equal(t, len(lcl.Policies), len(policiesResponse.Data))
}

func TestGetPolicyNoID(t *testing.T) {
	lcl := cmd.LaceworkContentLibrary{}
	_, actualError := lcl.GetPolicy("")
	assert.Equal(t, "policy ID must be provided", actualError.Error())
}

func TestGetPolicyMalformed(t *testing.T) {
	_, actualError := malformedLCL.GetPolicy("my_policy")
	assert.Equal(t, "policy exists but is malformed", actualError.Error())
}

func TestGetPolicyOK(t *testing.T) {
	policyID := "lwcustom-28"
	ept, err := ensureMockLCL(getMockLCLBinaryName())
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	lcl, err := cmd.LoadLCL(mockLWComponentState)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	actualString, actualError := lcl.GetPolicy(policyID)
	assert.Nil(t, actualError)
	assert.Contains(t, actualString, fmt.Sprintf("$account-%s", policyID))
}
