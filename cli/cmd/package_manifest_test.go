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
				Pkg: "kernel", PkgVer: "4.14.203-156.332.amzn2", // without EPOCH
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
