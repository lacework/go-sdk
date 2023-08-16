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

package cmd

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

	"github.com/lacework/go-sdk/lwcomponent"
)

var (
	mockPolicyReferences []LCLReference = []LCLReference{
		{
			ID:   "my_query",
			Type: "query",
			Path: "queries/my_query",
		},
		{
			ID:   "my_policy",
			Type: "policy",
			Path: "policies/my_policy",
		},
	}
	malformedLCL LaceworkContentLibrary = LaceworkContentLibrary{
		Queries: map[string]LCLQuery{
			"my_query": {References: []LCLReference{}},
		},
		Policies: map[string]LCLPolicy{
			"my_policy": {References: []LCLReference{
				LCLReference{},
			}},
		},
	}
	mockLCL LaceworkContentLibrary = LaceworkContentLibrary{
		Queries: map[string]LCLQuery{
			"my_query": {
				References: []LCLReference{
					{
						ID:   "my_query",
						Type: "query",
						Path: "queries/my_query",
					},
				},
			},
		},
		Policies: map[string]LCLPolicy{
			"my_policy": {
				References: mockPolicyReferences,
			},
		},
	}

	mockVersion, _                         = semver.NewVersion("0.1.0")
	mockLCLComponent lwcomponent.Component = lwcomponent.Component{
		Name:          "content-library",
		Description:   "Lacework Content Library",
		LatestVersion: *mockVersion,
		// @dhazekamp this only works for darwin-amd64 because we don't have per-package sigs
		Type: lwcomponent.BinaryType,
	}
	mockLWComponentState lwcomponent.State = lwcomponent.State{
		Components: []lwcomponent.Component{
			mockLCLComponent,
		},
	}
	nonZeroLCLComponent lwcomponent.Component = lwcomponent.Component{
		Name:          "content-library",
		Description:   "Lacework Content Library",
		LatestVersion: *mockVersion,
		Type:          lwcomponent.BinaryType,
	}
	nonZeroLWComponentState lwcomponent.State = lwcomponent.State{
		Components: []lwcomponent.Component{
			nonZeroLCLComponent,
		},
	}
	noParseLCLComponent lwcomponent.Component = lwcomponent.Component{
		Name:          "content-library",
		Description:   "Lacework Content Library",
		LatestVersion: *mockVersion,
		Type:          lwcomponent.BinaryType,
	}
	noParseLWComponentState lwcomponent.State = lwcomponent.State{
		Components: []lwcomponent.Component{
			noParseLCLComponent,
		},
	}
	//go:embed test_resources/content-library/*
	mockLCLBinaries embed.FS
)

func TestGetPolicyReference(t *testing.T) {
	ref, err := getPolicyReference([]LCLReference{})
	assert.NotNil(t, err)
	ref, _ = getPolicyReference(mockPolicyReferences)
	assert.Equal(t, mockPolicyReferences[1], ref)
}

type getQueryRefTest struct {
	Name      string
	Library   LaceworkContentLibrary
	QueryID   string
	Reference LCLReference
	Error     error
}

var getQueryRefTests = []getQueryRefTest{
	{
		Name:  "NoQueryID",
		Error: errors.New("query ID must be provided"),
	},
	{
		Name:    "QueryNotFound",
		QueryID: "my_query",
		Error:   errors.New("query does not exist in library"),
	},
	{
		Name:    "QueryMalformed",
		Library: malformedLCL,
		QueryID: "my_query",
		Error:   errors.New("query exists but is malformed"),
	},
	{
		Name:    "QueryOK",
		Library: mockLCL,
		QueryID: "my_query",
		Error:   nil,
	},
}

func TestGetQueryRef(t *testing.T) {
	for _, gqrt := range getQueryRefTests {
		t.Run(gqrt.Name, func(t *testing.T) {
			actualRef, actualError := gqrt.Library.getReferenceForQuery(gqrt.QueryID)

			if gqrt.Error != nil {
				assert.Equal(t, gqrt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
				assert.Equal(t, gqrt.QueryID, actualRef.ID)
			}
		})
	}
}

type getPolicyRefsTest struct {
	Name       string
	Library    LaceworkContentLibrary
	PolicyID   string
	References []LCLReference
	Error      error
}

var getPolicyRefsTests = []getPolicyRefsTest{
	{
		Name:  "NoPolicyID",
		Error: errors.New("policy ID must be provided"),
	},
	{
		Name:     "PolicyNotFound",
		PolicyID: "my_policy",
		Error:    errors.New("policy does not exist in library"),
	},
	{
		Name:     "PolicyMalformed",
		Library:  malformedLCL,
		PolicyID: "my_policy",
		Error:    errors.New("policy exists but is malformed"),
	},
	{
		Name:     "PolicyOK",
		Library:  mockLCL,
		PolicyID: "my_policy",
		Error:    nil,
	},
}

func TestGetPolicyRefs(t *testing.T) {
	for _, gprt := range getPolicyRefsTests {
		t.Run(gprt.Name, func(t *testing.T) {
			actualRefs, actualError := gprt.Library.getReferencesForPolicy(gprt.PolicyID)

			if gprt.Error != nil {
				assert.Equal(t, gprt.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
				assert.Equal(t, mockPolicyReferences, actualRefs)
			}
		})
	}
}

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

	cmpntBinaryPath := fmt.Sprintf("test_resources/content-library/%s", b)
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

func TestLoadLCLNilPtr(t *testing.T) {
	cli := new(cliState)

	_, err := cli.LoadLCL()
	assert.Equal(
		t,
		"unable to load Lacework Content Library",
		err.Error(),
	)
}

func TestLoadLCLNotFound(t *testing.T) {
	cli := cliState{LwComponents: new(lwcomponent.State)}

	// IsLCLInstalled
	assert.Equal(t, false, cli.isLCLInstalled())

	_, err := cli.LoadLCL()
	assert.Equal(
		t,
		"unable to load Lacework Content Library: component not installed",
		err.Error(),
	)
}

func TestLoadLCLNonZero(t *testing.T) {
	cli := cliState{LwComponents: &nonZeroLWComponentState}

	ept, err := ensureMockLCL("content-library-nonzero.sh")
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	_, err = cli.LoadLCL()
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
	cli := cliState{LwComponents: &noParseLWComponentState}

	ept, err := ensureMockLCL("content-library-noparse.sh")
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	_, err = cli.LoadLCL()
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

func _TestLoadLCLOK(t *testing.T) {
	cli := cliState{LwComponents: &mockLWComponentState}

	ept, err := ensureMockLCL(getMockLCLBinaryName())
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	// IsLCLInstalled
	assert.Equal(t, true, cli.isLCLInstalled())

	lcl, err := cli.LoadLCL()
	assert.Nil(t, err)
	_, ok := lcl.Queries["LW_Custom_Host_Activity_PotentialReverseShell"]
	assert.True(t, ok)
	_, ok = lcl.Policies["lwcustom-27"]
	assert.True(t, ok)
	assert.NotEqual(t, "", lcl.Policies["lwcustom-27"].PolicyID)
	assert.NotEqual(t, "", lcl.Policies["lwcustom-27"].QueryID)
}

func TestGetQueryNoID(t *testing.T) {
	lcl := LaceworkContentLibrary{}
	_, actualError := lcl.GetQuery("")
	assert.Equal(t, "query ID must be provided", actualError.Error())
}

func TestGetQueryMalformed(t *testing.T) {
	malformedLCL := LaceworkContentLibrary{
		Queries: map[string]LCLQuery{
			"my_query": {References: []LCLReference{
				LCLReference{},
			}},
		},
	}

	_, actualError := malformedLCL.GetQuery("my_query")
	assert.Equal(t, "query exists but is malformed", actualError.Error())
}

func _TestGetQueryOK(t *testing.T) {
	cli := cliState{LwComponents: &mockLWComponentState}

	queryID := "LW_Custom_AWS_CTA_AuroraPasswordChange"
	ept, err := ensureMockLCL(getMockLCLBinaryName())
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	lcl, err := cli.LoadLCL()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	actualQuery, actualError := lcl.GetQuery(queryID)
	assert.Nil(t, actualError)
	assert.Contains(t, actualQuery, queryID)
}

func TestGetPolicyNoID(t *testing.T) {
	lcl := LaceworkContentLibrary{}
	_, actualError := lcl.GetPolicy("")
	assert.Equal(t, "policy ID must be provided", actualError.Error())
}

func TestGetPolicyMalformed(t *testing.T) {
	_, actualError := malformedLCL.GetPolicy("my_policy")
	assert.Equal(t, "policy exists but is malformed", actualError.Error())
}

func _TestGetPolicyOK(t *testing.T) {
	cli := cliState{LwComponents: &mockLWComponentState}

	policyID := "lwcustom-28"
	ept, err := ensureMockLCL(getMockLCLBinaryName())
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	lcl, err := cli.LoadLCL()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	actualString, actualError := lcl.GetPolicy(policyID)
	assert.Nil(t, actualError)
	assert.Contains(t, actualString, fmt.Sprintf("$account-%s", policyID))
}

func _TestGetPoliciesByTagMalformed(t *testing.T) {
	policies := malformedLCL.GetPoliciesByTag("lwredteam")
	assert.Equal(t, 0, len(policies))
}

func _TestGetPoliciesByTag(t *testing.T) {
	cli := cliState{LwComponents: &mockLWComponentState}

	ept, err := ensureMockLCL(getMockLCLBinaryName())
	defer removeMockLCL(ept)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	lcl, err := cli.LoadLCL()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	// no such tag
	policies := lcl.GetPoliciesByTag("nosuchtag")
	assert.Equal(t, 0, len(policies))

	// ok tag
	policies = lcl.GetPoliciesByTag("lwredteam")
	assert.Equal(t, 4, len(policies))
}
