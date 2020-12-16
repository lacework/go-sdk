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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/lwupdater"
)

var (
	// All the following "unknown" variables are being injected at
	// build time via the cross-platform directive inside the Makefile
	//
	// Version is the semver coming from the VERSION file
	Version = "unknown"

	// GitSHA is the git ref that the cli was built from
	GitSHA = "unknown"

	// BuildTime is a human-readable time when the cli was built at
	BuildTime = "unknown"

	// versionCmd represents the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print the Lacework CLI version",
		Long: `
Prints out the installed version of the Lacework CLI and checks for newer
versions available for update.

Set the environment variable 'LW_UPDATES_DISABLE=1' to avoid checking for updates.`,
		Args: cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			if cli.JSONOutput() {
				errcheckEXIT(
					cli.OutputJSONString(
						fmt.Sprintf(
							`{"version":"v%s","git_sha":"%s","build_time":"%s"}`,
							Version, GitSHA, BuildTime,
						),
					),
				)
				return
			}
			cli.OutputHuman("lacework v%s (sha:%s) (time:%s)\n", Version, GitSHA, BuildTime)

			// check the latest version of the cli
			if err := versionCheck(); err != nil {
				exitwithCode(err, 4)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// VersionCache is the representation of the file named 'version_cache' stored
// at the directory ~/.config/lacework
type VersionCache struct {
	Project        string    `json:"project"`
	CurrentVersion string    `json:"current_version"`
	LastCheckTime  time.Time `json:"last_check_time"`
}

// versionCheck checks if the user is running the latest version of the cli,
// if not, displays a friendly message about the new version available
func versionCheck() error {
	cli.Log.Debugw("check version of the lacework-cli version", "repository", "go-sdk")
	sdk, err := lwupdater.Check("go-sdk", fmt.Sprintf("v%s", Version))
	if err != nil {
		return errors.Wrap(err, "unable to check updates")
	}

	if sdk.Outdated {
		cli.OutputHuman(fmt.Sprintf(
			"\nA newer version of the Lacework CLI is available! The latest version is %s,\n"+
				"to update execute the following command:\n%s\n",
			sdk.Latest, cli.UpdateCommand()))
	}

	return nil
}

// dailyVersionCheck will execute a version check on a daily basis, the function uses
// the file ~/.config/lacework/version_cache to track the time of last check
func dailyVersionCheck() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	var (
		configDir   = path.Join(home, ".config")
		lwConfigDir = path.Join(configDir, "lacework")
		cacheFile   = path.Join(lwConfigDir, "version_cache")
	)

	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		// first time running the daily version check, create directory
		if err := os.MkdirAll(lwConfigDir, 0755); err != nil {
			return err
		}

		if err := updateVersionCache(cacheFile); err != nil {
			return err
		}

		return versionCheck()
	}

	cli.Log.Debugw("verify cached version", "cache_file", cacheFile)
	cacheJSON, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return err
	}

	var versionCache VersionCache
	if err := json.Unmarshal(cacheJSON, &versionCache); err != nil {
		return err
	}

	cli.Log.Debugw("version cache", "content", versionCache)
	checkTime := time.Now().AddDate(0, 0, -1)
	if versionCache.LastCheckTime.Before(checkTime) {
		cli.Log.Debugw("triggering daily version check")
		if err := updateVersionCache(cacheFile); err != nil {
			return err
		}

		return versionCheck()
	}

	cli.Log.Debugw("threshold not yet met. skipping daily version check",
		"threshold", "1d",
		"last_check_time", versionCache.LastCheckTime,
		"next_check_time", versionCache.LastCheckTime.AddDate(0, 0, 1))

	return nil
}

func updateVersionCache(cacheFile string) error {
	var (
		versionCache = VersionCache{
			Project:        "lacework-cli",
			CurrentVersion: Version,
			LastCheckTime:  time.Now(),
		}
		buf = new(bytes.Buffer)
	)
	cli.Log.Debugw("storing version cache", "content", versionCache)
	if err := json.NewEncoder(buf).Encode(versionCache); err != nil {
		return err
	}

	err := ioutil.WriteFile(cacheFile, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
