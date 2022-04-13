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
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestRemoveInactivePackagesFromManifest(t *testing.T) {
	manifest := new(api.PackageManifest)
	subject := cli.removeInactivePackagesFromManifest(manifest, "rpm")
	assert.Equal(t, manifest, subject)
}

func TestRemoveInactivePackagesFromManifestRemoveKernelRPM(t *testing.T) {
	manifest := &api.PackageManifest{
		OsPkgInfoList: []api.OsPkgInfo{
			api.OsPkgInfo{
				Os: "amzn", OsVer: "2",
				Pkg: "kernel", PkgVer: "0:4.14.209-160.339.amzn2", // with EPOCH
			},
			api.OsPkgInfo{
				Os: "amzn", OsVer: "2",
				Pkg: "kernel", PkgVer: "4.14.203-156.331.amzn2", // without EPOCH
			},
		},
	}
	subject := cli.removeInactivePackagesFromManifest(manifest, "rpm")
	assert.Empty(t, subject)
}

func TestRemoveInactivePackagesFromManifestRemoveKernelDPKG(t *testing.T) {
	manifest := &api.PackageManifest{
		OsPkgInfoList: []api.OsPkgInfo{
			api.OsPkgInfo{
				Os: "ubuntu", OsVer: "18.04",
				Pkg: "linux-image-5.3.0-1035-aws", PkgVer: "5.3.0-1035.37",
			},
			api.OsPkgInfo{
				Os: "ubuntu", OsVer: "18.04",
				Pkg: "sudo", PkgVer: "1.8.21p2-3ubuntu1.2", // not a kernel pkg
			},
		},
	}
	subject := cli.removeInactivePackagesFromManifest(manifest, "dpkg-query")
	assert.NotEmpty(t, subject)
	assert.Equal(t, &api.PackageManifest{
		OsPkgInfoList: []api.OsPkgInfo{
			api.OsPkgInfo{
				Os: "ubuntu", OsVer: "18.04",
				Pkg: "sudo", PkgVer: "1.8.21p2-3ubuntu1.2", // this pkg should persist
			},
		},
	}, subject)
}

func TestRemoveInactivePackagesFromManifestUnknownManager(t *testing.T) {
	manifest := new(api.PackageManifest)
	subject := cli.removeInactivePackagesFromManifest(manifest, "apk")
	assert.Equal(t, manifest, subject)
}

func TestRemoveEpochFromPkgVersion(t *testing.T) {
	assert.Equal(t,
		"4.14.209-160.339.amzn2",
		removeEpochFromPkgVersion("4.14.209-160.339.amzn2"))
	assert.Equal(t,
		"4.14.209-160.339.amzn2",
		removeEpochFromPkgVersion("0:4.14.209-160.339.amzn2"))
	assert.Equal(t,
		"version",
		removeEpochFromPkgVersion("epoch:version"))
}

func TestSplitPackageManifest(t *testing.T) {
	cases := []struct {
		chunks       int
		size         int
		expectedSize int
	}{
		{expectedSize: 100,
			size:   500,
			chunks: 5},
		{expectedSize: 45,
			size:   45000,
			chunks: 1000},
		{expectedSize: 50,
			size:   100,
			chunks: 2},
		{expectedSize: 2,
			size:   1001,
			chunks: 1000},
		{expectedSize: 28,
			size:   55000,
			chunks: 2000},
		{expectedSize: 1,
			size:   123,
			chunks: 1000},
	}
	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			manifest := &api.PackageManifest{
				OsPkgInfoList: make([]api.OsPkgInfo, kase.size),
			}
			subject := splitPackageManifest(manifest, kase.chunks)
			assert.Equal(t, kase.expectedSize, len(subject))
		})
	}
}

func TestFanOutHostScans(t *testing.T) {
	// mock the api client
	client, err := api.NewClient("test", api.WithToken("mock"))
	assert.Nil(t, err)
	client.Vulnerabilities = api.NewVulnerabilityService(client)
	cli.LwApi = client
	defer func() {
		cli.LwApi = nil
	}()

	subject, err := fanOutHostScans()
	assert.Nil(t, err)
	assert.Equal(t, api.HostVulnScanPkgManifestResponse{}, subject)

	subject, err = fanOutHostScans(nil)
	assert.Nil(t, err)
	assert.Equal(t, api.HostVulnScanPkgManifestResponse{}, subject)

	// more than 10 morkers should return an error
	multiManifests := make([]*api.PackageManifest, 11)
	subject, err = fanOutHostScans(multiManifests...)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(),
			"limit of packages exceeded",
		)
	}
	assert.Equal(t, api.HostVulnScanPkgManifestResponse{}, subject)

	subject, err = fanOutHostScans(&api.PackageManifest{})
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(),
			"[403] Unauthorized Access", // intentional error since we are mocking the api token
		)
	}
	assert.Equal(t, api.HostVulnScanPkgManifestResponse{}, subject)
}

func TestMergeHostVulnScanPkgManifestResponses(t *testing.T) {
	cases := []struct {
		expected api.HostVulnScanPkgManifestResponse
		from     api.HostVulnScanPkgManifestResponse
		to       api.HostVulnScanPkgManifestResponse
	}{
		// empty responses
		{expected: api.HostVulnScanPkgManifestResponse{},
			from: api.HostVulnScanPkgManifestResponse{},
			to:   api.HostVulnScanPkgManifestResponse{}},
		// responses should return an Ok status
		{expected: api.HostVulnScanPkgManifestResponse{
			Ok: true},
			from: api.HostVulnScanPkgManifestResponse{
				Ok: true},
			to: api.HostVulnScanPkgManifestResponse{
				Ok: false}},
		// messages should change only if the previous one is empty or different
		{expected: api.HostVulnScanPkgManifestResponse{
			Message: "SUCCESS"},
			from: api.HostVulnScanPkgManifestResponse{
				Message: "SUCCESS"},
			to: api.HostVulnScanPkgManifestResponse{
				Message: ""}},
		{expected: api.HostVulnScanPkgManifestResponse{
			Message: "YES"},
			from: api.HostVulnScanPkgManifestResponse{
				Message: ""},
			to: api.HostVulnScanPkgManifestResponse{
				Message: "YES"}},
		{expected: api.HostVulnScanPkgManifestResponse{
			Message: "OLD,NEW"},
			from: api.HostVulnScanPkgManifestResponse{
				Message: "NEW"},
			to: api.HostVulnScanPkgManifestResponse{
				Message: "OLD"}},
		// merge two responses into one single response 1 + 1 = 2
		{
			expected: api.HostVulnScanPkgManifestResponse{
				Vulns: []api.HostScanPackageVulnDetails{
					api.HostScanPackageVulnDetails{}, api.HostScanPackageVulnDetails{},
				},
			},
			from: api.HostVulnScanPkgManifestResponse{
				Vulns: []api.HostScanPackageVulnDetails{api.HostScanPackageVulnDetails{}}},
			to: api.HostVulnScanPkgManifestResponse{
				Vulns: []api.HostScanPackageVulnDetails{api.HostScanPackageVulnDetails{}}},
		},
	}
	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			mergeHostVulnScanPkgManifestResponses(&kase.to, &kase.from)
			assert.Equal(t, kase.expected, kase.to)
		})
	}
}

func TestIsEsmEnabled(t *testing.T) {
	assert.False(t, cli.IsEsmEnabled())
}

func TestParseOsRelease(t *testing.T) {
	file, err := ioutil.TempFile("", "os-release")
	assert.Nil(t, err)
	_, err = file.WriteString(mockUbuntuOSReleaseFile)
	assert.Nil(t, err)

	defer os.Remove(file.Name())

	os, err := openOsReleaseFile(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, mockUbuntu.Name, os.Name)
	assert.Equal(t, mockUbuntu.Version, os.Version)
}

func TestParseSysRelease(t *testing.T) {
	file, err := ioutil.TempFile("", "system-release")
	assert.Nil(t, err)
	_, err = file.WriteString(mockCentosSystemFile)
	assert.Nil(t, err)

	defer os.Remove(file.Name())

	os, err := openSystemReleaseFile(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, mockCentos.Name, os.Name)
	assert.Equal(t, mockCentos.Version, os.Version)
}

var (
	mockCentos              = OS{Name: "centos", Version: "6.10"}
	mockUbuntu              = OS{Name: "ubuntu", Version: "18.04"}
	mockCentosSystemFile    = "CentOS release 6.10 (Final)"
	mockUbuntuOSReleaseFile = `NAME="Ubuntu"
VERSION="18.04.5 LTS (Bionic Beaver)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 18.04.5 LTS"
VERSION_ID="18.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=bionic
UBUNTU_CODENAME=bionic
`
)
