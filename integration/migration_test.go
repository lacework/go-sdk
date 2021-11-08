//go:build migration

// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestV2MigrationWithConfigFile(t *testing.T) {
	// create a temporal directory to use it as our home directory to test
	// the version_cache mechanism
	home := createTOMLConfigFromCIvars()
	defer os.RemoveAll(home)

	out, err, exitcode := LaceworkCLIWithHome(home, "agent", "token", "list")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.NotEmpty(t,
		out.String(),
		"STDOUT should not be empty")

	// @afiune the command is not really important but the actual result in the config

	t.Run("check the backup dir", func(t *testing.T) {
		backupsDir := path.Join(home, ".config", "lacework", "cfg_backups")
		if assert.DirExists(t, backupsDir, "the config backups directory is missing") {

			f, _ := os.Open(backupsDir)
			defer f.Close()

			fileNames, _ := f.Readdirnames(0)
			if assert.Equal(t, 1, len(fileNames), "a single backup should exist") {
				assert.Regexpf(t, regexp.MustCompile(".lacework.toml.*.*.bkp"), fileNames[0], "backup config file mismatch")
			}
		}
	})

	t.Run("check migrated config file", func(t *testing.T) {
		configPath := path.Join(home, ".lacework.toml")
		if assert.FileExists(t, configPath, "the configuration file is missing") {
			laceworkTOML, err := ioutil.ReadFile(configPath)
			if assert.Nil(t, err) {
				laceworkTOMLString := string(laceworkTOML)
				assert.Contains(t, laceworkTOMLString, "account", "there is a problem with the v2 migrated config")
				assert.Contains(t, laceworkTOMLString, "subaccount", "there is a problem with the v2 migrated config") // only for our tech-ally account
				assert.Contains(t, laceworkTOMLString, "api_key", "there is a problem with the v2 migrated config")
				assert.Contains(t, laceworkTOMLString, "api_secret", "there is a problem with the v2 migrated config")
				assert.Contains(t, laceworkTOMLString, "version = 2", "there is a problem with the v2 migrated config")
			}
		}
	})
}

func TestV2MigrationWithFlagsOrEnvVariables(t *testing.T) {
	home, errDir := ioutil.TempDir("", "lacework-cli")
	if errDir != nil {
		panic(errDir)
	}
	defer os.RemoveAll(home)

	out, err, exitcode := LaceworkCLIWithHome(home, "agent", "token", "list",
		"-a", os.Getenv("CI_ACCOUNT"), "-k", os.Getenv("CI_API_KEY"), "-s", os.Getenv("CI_API_SECRET"),
	)
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.NotEmpty(t,
		out.String(),
		"STDOUT should not be empty")

	// @afiune the command is not really important but the actual result in the config

	t.Run("check the backup dir", func(t *testing.T) {
		backupsDir := path.Join(home, ".config", "lacework", "cfg_backups")
		assert.NoDirExists(t, backupsDir, "the config backups directory should not exist")
	})

	t.Run("check migrated config file", func(t *testing.T) {
		configPath := path.Join(home, ".lacework.toml")
		assert.NoFileExists(t, configPath, "the configuration file should not exist")
	})
}
