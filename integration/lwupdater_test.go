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

package integration_test

import (
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/lwupdater"
)

func TestCheckErrorEmptyProject(t *testing.T) {
	info, err := lwupdater.Check("", "")
	assert.Empty(t, info)
	if assert.NotNil(t, err) {
		assert.Equal(t, "specify a valid project", err.Error())
	}
}

// @afiune this test requires to actually have internet access,
// I wonder if this will cause problems in the future, if so,
// we should disable it.
func TestCheck(t *testing.T) {
	info, err := lwupdater.Check("go-sdk", "0.1.6")
	if assert.Nil(t, err) {
		assert.Equal(t, "go-sdk", info.Project)
		assert.Equal(t, "0.1.6", info.CurrentVersion)
		assert.Regexp(t, `^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)$`, info.LatestVersion)
		assert.True(t, info.Outdated)
	}
}

func TestCheckSkipDevVersions(t *testing.T) {
	info, err := lwupdater.Check("go-sdk", "1.30.9-dev")
	if assert.Nil(t, err) {
		assert.Equal(t, "1.30.9-dev", info.CurrentVersion)
		assert.False(t, info.Outdated)
	}
}

func TestCheckDisabled(t *testing.T) {
	// disable the updater
	os.Setenv(lwupdater.DisableEnv, "1")
	defer os.Setenv(lwupdater.DisableEnv, "")
	info, err := lwupdater.Check("go-sdk", "v0.1.0")
	assert.Nil(t, err)
	assert.Empty(t, info)
}

func TestVersionLoadCacheError(t *testing.T) {
	dir, err := os.MkdirTemp("", "t")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	cacheFile := path.Join(dir, "version_cache")
	v, err := lwupdater.LoadCache(cacheFile)
	assert.Empty(t, v)
	if assert.NotNil(t, err) {
		if runtime.GOOS == "windows" {
			assert.Contains(t, err.Error(), "The system cannot find the file specified")
		} else {
			assert.Contains(t, err.Error(), "no such file or directory")
		}
	}
}

func TestVersionStoreLoadCache(t *testing.T) {
	dir, err := os.MkdirTemp("", "t")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	cacheFile := path.Join(dir, "version_cache")

	mockVersion := lwupdater.Version{
		Project:        "lwupdater",
		CurrentVersion: "0.1.0",
		LastCheckTime:  time.Now(),
	}
	err = mockVersion.StoreCache(cacheFile)
	assert.Nil(t, err)

	subjectVersion, err := lwupdater.LoadCache(cacheFile)
	if assert.Nil(t, err) {
		assert.Equal(t, mockVersion.Project, subjectVersion.Project)
		assert.Equal(t, mockVersion.CurrentVersion, subjectVersion.CurrentVersion)
		assert.Equal(t,
			mockVersion.LastCheckTime.Format(time.RFC3339),
			subjectVersion.LastCheckTime.Format(time.RFC3339),
		)
	}
}
