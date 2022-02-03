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

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwcomponent"
)

var (
	mockVersion, _                         = semver.NewVersion("0.1.0")
	mockLCLComponent lwcomponent.Component = lwcomponent.Component{
		Name:          "lacework-content-library",
		Description:   "Lacework Content Library",
		LatestVersion: *mockVersion,
		// @dhazekamp this only works for darwin-amd64 because we don't have per-package sigs
		CLICommand:  false,
		CommandName: "",
		Binary:      true,
		Library:     false,
		Standalone:  false,
	}
	mockLWComponentState lwcomponent.State = lwcomponent.State{
		Components: []lwcomponent.Component{
			mockLCLComponent,
		},
	}
	nonZeroLCLComponent lwcomponent.Component = lwcomponent.Component{
		Name:          "lacework-content-library",
		Description:   "Lacework Content Library",
		LatestVersion: *mockVersion,
		CLICommand:    false,
		CommandName:   "",
		Binary:        true,
		Library:       false,
		Standalone:    false,
	}
	nonZeroLWComponentState lwcomponent.State = lwcomponent.State{
		Components: []lwcomponent.Component{
			nonZeroLCLComponent,
		},
	}
	noParseLCLComponent lwcomponent.Component = lwcomponent.Component{
		Name:          "lacework-content-library",
		Description:   "Lacework Content Library",
		LatestVersion: *mockVersion,
		CLICommand:    false,
		CommandName:   "",
		Binary:        true,
		Library:       false,
		Standalone:    false,
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

	dir, _ := path.Split(cmpntPath)
	cmpntVersionPath := path.Join(dir, ".version")
	cmpntSignaturePath := path.Join(dir, ".signature")

	if mockLCLComponent.Status() == lwcomponent.Installed {
		os.Rename(cmpntPath, cmpntPath+".bak")
		os.Rename(cmpntVersionPath, cmpntVersionPath+".bak")
		os.Rename(cmpntSignaturePath, cmpntSignaturePath+".bak")
		placementType = mockLCLReplaced
	} else {
		placementType = mockLCLPlaced
		os.MkdirAll(dir, os.ModePerm)
	}

	cmpntBinaryPath := fmt.Sprintf("test_resources/lacework-content-library/%s", b)
	cmpntBytes, err := mockLCLBinaries.ReadFile(cmpntBinaryPath)
	if err != nil {
		return placementType, err
	}
	os.WriteFile(cmpntPath, cmpntBytes, 0777)

	canonSignaturePath := fmt.Sprintf("%s.sig", cmpntBinaryPath)
	cmpntSignatureBytes, err := mockLCLBinaries.ReadFile(canonSignaturePath)
	if err != nil {
		return placementType, err
	}
	os.WriteFile(cmpntSignaturePath, cmpntSignatureBytes, 0666)

	return placementType, os.WriteFile(
		cmpntVersionPath, []byte(mockLCLComponent.LatestVersion.String()), 0666)
}

func removeMockLCL(ept mockLCLPlacementType) {
	if ept == mockLCLNoop {
		return
	}

	cmpntPath, _ := mockLCLComponent.Path()
	dir, _ := path.Split(cmpntPath)
	cmpntVersionPath := path.Join(dir, ".version")
	cmpntSignaturePath := path.Join(dir, ".signature")

	if ept == mockLCLReplaced {
		os.Rename(cmpntPath+".bak", cmpntPath)
		os.Rename(cmpntVersionPath+".bak", cmpntVersionPath)
		os.Rename(cmpntSignaturePath+".bak", cmpntSignaturePath)
		return
	}
	os.RemoveAll(dir)
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
	assert.NotEqual(t, "", lcl.Policies["lwcustom-27"].PolicyID)
	assert.NotEqual(t, "", lcl.Policies["lwcustom-27"].QueryID)
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
