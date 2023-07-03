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
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestRemoveInactivePackagesFromManifest(t *testing.T) {
	manifest := new(api.VulnerabilitiesPackageManifest)
	subject := cli.removeInactivePackagesFromManifest(manifest, "rpm")
	assert.Equal(t, manifest, subject)
}

func TestRemoveInactivePackagesFromManifestRemoveKernelRPM(t *testing.T) {
	manifest := &api.VulnerabilitiesPackageManifest{
		OsPkgInfoList: []api.VulnerabilitiesOsPkgInfo{
			{
				Os: "amzn", OsVer: "2",
				Pkg: "kernel", PkgVer: "0:4.14.209-160.339.amzn2", // with EPOCH
			},
			{
				Os: "amzn", OsVer: "2",
				Pkg: "kernel", PkgVer: "4.14.203-156.331.amzn2", // without EPOCH
			},
		},
	}
	subject := cli.removeInactivePackagesFromManifest(manifest, "rpm")
	assert.Empty(t, subject)
}

func TestRemoveInactivePackagesFromManifestRemoveKernelDPKG(t *testing.T) {
	manifest := &api.VulnerabilitiesPackageManifest{
		OsPkgInfoList: []api.VulnerabilitiesOsPkgInfo{
			{
				Os: "ubuntu", OsVer: "18.04",
				Pkg: "linux-image-5.3.0-1035-aws", PkgVer: "5.3.0-1035.37",
			},
			{
				Os: "ubuntu", OsVer: "18.04",
				Pkg: "sudo", PkgVer: "1.8.21p2-3ubuntu1.2", // not a kernel pkg
			},
		},
	}
	subject := cli.removeInactivePackagesFromManifest(manifest, "dpkg-query")
	assert.NotEmpty(t, subject)
	assert.Equal(t, &api.VulnerabilitiesPackageManifest{
		OsPkgInfoList: []api.VulnerabilitiesOsPkgInfo{
			{
				Os: "ubuntu", OsVer: "18.04",
				Pkg: "sudo", PkgVer: "1.8.21p2-3ubuntu1.2", // this pkg should persist
			},
		},
	}, subject)
}

func TestRemoveInactivePackagesFromManifestUnknownManager(t *testing.T) {
	manifest := new(api.VulnerabilitiesPackageManifest)
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
			manifest := &api.VulnerabilitiesPackageManifest{
				OsPkgInfoList: make([]api.VulnerabilitiesOsPkgInfo, kase.size),
			}
			subject := splitPackageManifest(manifest, kase.chunks)
			assert.Equal(t, kase.expectedSize, len(subject))
		})
	}
}

func TestFanOutHostScans(t *testing.T) {
	// mock the api client
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI("Vulnerabilities/SoftwarePackages/scan",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(403)
			fmt.Fprintf(w, "User not authorized.")
		},
	)

	client, err := api.NewClient("test",
		api.WithURL(fakeServer.URL()),
		api.WithToken("TOKEN"))
	assert.Nil(t, err)
	cli.LwApi = client
	defer func() {
		cli.LwApi = nil
	}()

	subject, err := fanOutHostScans()
	assert.Nil(t, err)
	assert.Equal(t, api.VulnerabilitySoftwarePackagesResponse{}, subject)

	subject, err = fanOutHostScans(nil)
	assert.Nil(t, err)
	assert.Equal(t, api.VulnerabilitySoftwarePackagesResponse{}, subject)

	// more than 10 morkers should return an error
	multiManifests := make([]*api.VulnerabilitiesPackageManifest, 11)
	subject, err = fanOutHostScans(multiManifests...)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(),
			"limit of packages exceeded",
		)
	}
	assert.Equal(t, api.VulnerabilitySoftwarePackagesResponse{}, subject)

	subject, err = fanOutHostScans(&api.VulnerabilitiesPackageManifest{})
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(),
			"[403] User not authorized.", // intentional error since we are mocking the api token
		)
	}
	assert.Equal(t, api.VulnerabilitySoftwarePackagesResponse{}, subject)
}

func TestMergeVulnerabilitySoftwarePackagesResponses(t *testing.T) {
	cases := []struct {
		expected api.VulnerabilitySoftwarePackagesResponse
		from     api.VulnerabilitySoftwarePackagesResponse
		to       api.VulnerabilitySoftwarePackagesResponse
	}{
		// merge two responses into one single response 1 + 1 = 2
		{
			expected: api.VulnerabilitySoftwarePackagesResponse{
				Data: []api.VulnerabilitySoftwarePackage{
					api.VulnerabilitySoftwarePackage{}, api.VulnerabilitySoftwarePackage{},
				},
			},
			from: api.VulnerabilitySoftwarePackagesResponse{
				Data: []api.VulnerabilitySoftwarePackage{api.VulnerabilitySoftwarePackage{}}},
			to: api.VulnerabilitySoftwarePackagesResponse{
				Data: []api.VulnerabilitySoftwarePackage{api.VulnerabilitySoftwarePackage{}}},
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			mergeHostVulnScanPkgManifestResponses(&c.to, &c.from)
			assert.Equal(t, c.expected, c.to)
		})
	}
}

func TestIsEsmEnabledWhenIsNotEnabled(t *testing.T) {
	assert.False(t, cli.IsEsmEnabled())
}

func TestIsEsmEnabledWhenIsActuallyEnabled(t *testing.T) {
	file, err := os.CreateTemp("", "ubuntu-advantage-status-json")
	assert.Nil(t, err)
	_, err = file.WriteString(mockUbuntuUAStatusFile)
	assert.Nil(t, err)

	procUAStatusFile = file.Name()

	defer func() {
		os.Remove(file.Name())
		// rollback
		procUAStatusFile = "/proc/1/root/var/lib/ubuntu-advantage/status.json"
	}()

	assert.True(t, cli.IsEsmEnabled())
}

func TestParseOsRelease(t *testing.T) {
	file, err := os.CreateTemp("", "os-release")
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
	file, err := os.CreateTemp("", "system-release")
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
	mockUbuntuUAStatusFile = `{
  "execution_details": "No Ubuntu Advantage operations are running",
  "machine_id": "abc123",
  "origin": "free",
  "_doc": "Content provided in json response is currently considered Experimental and may change",
  "account": {
    "name": "test@lacework.net",
    "id": "abc123",
    "external_account_ids": [],
    "created_at": "2022-01-25T23:06:47+00:00"
  },
  "expires": "9999-12-31T00:00:00+00:00",
  "execution_status": "inactive",
  "version": "27.2.2~16.04.1",
  "notices": [],
  "_schema_version": "0.1",
  "config_path": "/etc/ubuntu-advantage/uaclient.conf",
  "attached": true,
  "services": [
    {
      "entitled": "yes",
      "description_override": null,
      "name": "cc-eal",
      "status_details": "CC EAL2 is not configured",
      "description": "Common Criteria EAL2 Provisioning Packages",
      "status": "disabled",
      "available": "yes"
    },
    {
      "entitled": "yes",
      "description_override": null,
      "name": "cis",
      "status_details": "CIS Audit is not configured",
      "description": "Center for Internet Security Audit Tools",
      "status": "disabled",
      "available": "yes"
    },
    {
      "entitled": "no",
      "description_override": null,
      "name": "esm-apps",
      "status_details": "",
      "description": "UA Apps: Extended Security Maintenance (ESM)",
      "status": "—",
      "available": "yes"
    },
    {
      "entitled": "yes",
      "description_override": null,
      "name": "esm-infra",
      "status_details": "UA Infra: ESM is active",
      "description": "UA Infra: Extended Security Maintenance (ESM)",
      "status": "enabled",
      "available": "yes"
    },
    {
      "entitled": "yes",
      "description_override": null,
      "name": "fips",
      "status_details": "FIPS is not configured",
      "description": "NIST-certified core packages",
      "status": "disabled",
      "available": "yes"
    },
    {
      "entitled": "yes",
      "description_override": null,
      "name": "fips-updates",
      "status_details": "FIPS Updates is not configured",
      "description": "NIST-certified core packages with priority security updates",
      "status": "disabled",
      "available": "yes"
    },
    {
      "entitled": "yes",
      "description_override": null,
      "name": "livepatch",
      "status_details": "",
      "description": "Canonical Livepatch service",
      "status": "enabled",
      "available": "yes"
    }
  ],
  "contract": {
    "products": [
      "free"
    ],
    "name": "test@lacework.net",
    "id": "abc123",
    "tech_support_level": "n/a",
    "created_at": "2022-01-25T23:06:47+00:00"
  },
  "config": {
    "security_url": "https://ubuntu.com/security",
    "data_dir": "/var/lib/ubuntu-advantage",
    "ua_config": {
      "apt_http_proxy": null,
      "https_proxy": null,
      "http_proxy": null,
      "apt_https_proxy": null
    },
    "log_level": "debug",
    "log_file": "/var/log/ubuntu-advantage.log",
    "contract_url": "https://contracts.canonical.com"
  },
  "effective": "2022-01-25T23:06:47+00:00"
}`
)
