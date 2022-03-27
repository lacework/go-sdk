//go:build version

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

package integration

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/lwupdater"
)

// To test the daily version check we need to set the environment
// variable CI_TEST_LWUPDATER=1 so that the framework does not disable it
func TestDailyVersionCheckEndToEnd(t *testing.T) {
	enableTestingUpdaterEnv()
	defer disableTestingUpdaterEnv()

	// create a temporal directory to use it as our home directory to test
	// the version_cache mechanism
	home := createTOMLConfigFromCIvars()
	defer os.RemoveAll(home)

	out, errB, exitcode := LaceworkCLIWithHome(home, "configure", "list")
	assert.Empty(t, errB.String())
	assert.Equal(t, 0, exitcode)
	assert.Contains(t, out.String(),
		"A newer version of the Lacework CLI is available! The latest version is v",
		"version check message changed?",
	)

	versionCacheFile := path.Join(home, ".config", "lacework", "version_cache")
	assert.FileExists(t, versionCacheFile, "the version_cache file is missing")

	var actualCache lwupdater.Version
	cacheJSON, err := ioutil.ReadFile(versionCacheFile)
	if assert.Nil(t, err) {
		err := json.Unmarshal(cacheJSON, &actualCache)
		if assert.Nil(t, err) {
			assert.Equal(t, actualCache.Project, "go-sdk")
			assert.True(t, time.Now().After(actualCache.LastCheckTime))
			assert.True(t, time.Now().AddDate(0, 0, -1).Before(actualCache.LastCheckTime))
		}
	}

	// re-running the same command should not check and display the version update
	out, errB, exitcode = LaceworkCLIWithHome(home, "configure", "list")
	assert.Empty(t, errB.String())
	assert.Equal(t, 0, exitcode)
	assert.NotContains(t, out.String(),
		"A newer version of the Lacework CLI is available! The latest version is v",
		"version update message should not be displayed",
	)

	// version cache file should continue to exist
	assert.FileExists(t, versionCacheFile, "the version_cache file is missing")

	var nextCache lwupdater.Version
	cacheJSON, err = ioutil.ReadFile(versionCacheFile)
	if assert.Nil(t, err) {
		err := json.Unmarshal(cacheJSON, &nextCache)
		if assert.Nil(t, err) {
			assert.Equal(t, actualCache.Project, nextCache.Project)
			assert.Equal(t, actualCache.LastCheckTime, nextCache.LastCheckTime)
			assert.Equal(t, actualCache.CurrentVersion, nextCache.CurrentVersion)
		}
	}

	// manipulate the version cache file, set the last check to two days ago
	actualCache.LastCheckTime = time.Now().AddDate(0, 0, -2)

	err = actualCache.StoreCache(versionCacheFile)
	assert.Nil(t, err)

	out, errB, exitcode = LaceworkCLIWithHome(home, "configure", "list")
	assert.Empty(t, errB.String())
	assert.Equal(t, 0, exitcode)
	assert.Contains(t, out.String(),
		"A newer version of the Lacework CLI is available! The latest version is v",
		"version check message should be there",
	)

	assert.FileExists(t, versionCacheFile, "the version_cache file is missing")

	var lastCache lwupdater.Version
	cacheJSON, err = ioutil.ReadFile(versionCacheFile)
	if assert.Nil(t, err) {
		err := json.Unmarshal(cacheJSON, &lastCache)
		if assert.Nil(t, err) {
			assert.Equal(t, lastCache.Project, "go-sdk")
			assert.True(t, time.Now().After(lastCache.LastCheckTime))
			assert.True(t, time.Now().AddDate(0, 0, -1).Before(lastCache.LastCheckTime))
		}
	}
}

func TestVersionCommand(t *testing.T) {
	enableTestingUpdaterEnv()
	defer disableTestingUpdaterEnv()

	out, errB, exitcode := LaceworkCLIWithTOMLConfig("version")
	assert.Empty(t, errB.String())
	assert.Equal(t, 0, exitcode)
	assert.Contains(t, out.String(),
		"A newer version of the Lacework CLI is available! The latest version is v",
		"version update message should be displayed",
	)
}

func TestDailyVersionCheckShouldNotRunWhenInNonInteractiveMode(t *testing.T) {
	enableTestingUpdaterEnv()
	defer disableTestingUpdaterEnv()

	home := createTOMLConfigFromCIvars()
	defer os.RemoveAll(home)

	out, err, exitcode := LaceworkCLIWithHome(home, "configure", "list", "--noninteractive")
	assert.Empty(t, err.String())
	assert.Equal(t, 0, exitcode)
	assert.NotContains(t, out.String(),
		"A newer version of the Lacework CLI is available! The latest version is v",
		"we shouldn't see this daily version check message",
	)

	versionCacheFile := path.Join(home, ".config", "lacework", "version_cache")
	assert.NoFileExists(t, versionCacheFile, "the version_cache file is missing")
}
