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
		Signature:   "f35b88ae47f9778061543a33ad6799a0d16adf9af02f7527bdf053d81f9a0607",
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
		Name:        "lacework-content-library",
		Description: "Lacework Content Library",
		Version:     "0.1.0",
		Signature:   "99032d1fb22b1ea6119f5f728cadae6feddfa62a45a66d52f77558cffd80b7f2",
		CLICommand:  false,
		CommandName: "",
		Binary:      true,
		Library:     false,
		Standalone:  false,
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
		Signature:   "0df2d5957dd7583361dcc3a888b2ad9e3fa29a413bbf711a572f65348227d898",
		CLICommand:  false,
		CommandName: "",
		Binary:      true,
		Library:     false,
		Standalone:  false,
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
	}
	mockLCL, _ = cmd.LoadLCL(mockLWComponentState)
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

func TestListQueries(t *testing.T) {
	assert.Equal(t, len(mockLCL.Queries), len(mockLCL.ListQueries().Data))
}

type getNewQueryTest struct {
	Name    string
	Library cmd.LaceworkContentLibrary
	QueryID string
	Error   error
}

var getNewQueryTests = []getNewQueryTest{
	getNewQueryTest{
		Name:  "NoQueryID",
		Error: errors.New("query ID must be provided"),
	},
	getNewQueryTest{
		Name:    "MalformedQuery",
		Library: malformedLCL,
		QueryID: "my_query",
		Error:   errors.New("query exists but is malformed"),
	},
	getNewQueryTest{
		Name:    "QueryOK",
		Library: *mockLCL,
		QueryID: "LW_Custom_AWS_CTA_AuroraPasswordChange",
		Error:   nil,
	},
}

func TestGetNewQuery(t *testing.T) {
	ept, err := ensureMockLCL(getMockLCLBinaryName())
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	for _, gnqt := range getNewQueryTests {
		t.Run(gnqt.Name, func(t *testing.T) {
			actualNewQuery, actualError := gnqt.Library.GetNewQuery(gnqt.QueryID)

			if gnqt.Error != nil {
				assert.Equal(t, gnqt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
				assert.Equal(t, gnqt.QueryID, actualNewQuery.QueryID)
			}
		})
	}
}
