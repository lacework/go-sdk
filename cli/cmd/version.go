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

	// The name of the version cache file needed for daily version checks
	VersionCacheFile = "version_cache"

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
			if _, err := versionCheck(); err != nil {
				exitwithCode(err, 4)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// versionCheck checks if the user is running the latest version of the cli,
// if not, displays a friendly message about the new version available
func versionCheck() (*lwupdater.Version, error) {
	cli.Log.Debugw("check version of the lacework-cli", "repository", "github.com/lacework/go-sdk")
	sdk, err := lwupdater.Check("go-sdk", fmt.Sprintf("v%s", Version))
	if err != nil {
		return nil, errors.Wrap(err, "unable to check updates")
	}

	if sdk.Outdated {
		cli.OutputHuman(fmt.Sprintf(
			"\nA newer version of the Lacework CLI is available! The latest version is %s,\n"+
				"to update execute the following command:\n%s\n",
			sdk.LatestVersion, cli.UpdateCommand()))
	}

	return sdk, nil
}

func versionCacheDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".config", "lacework"), nil
}

// dailyVersionCheck will execute a version check on a daily basis, the function uses
// the file ~/.config/lacework/version_cache to track the last check time
func dailyVersionCheck() error {
	cacheDir, err := versionCacheDir()
	if err != nil {
		return err
	}

	cacheFile := path.Join(cacheDir, VersionCacheFile)
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		// first time running the daily version check, create directory
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return err
		}

		currentVersion, err := versionCheck()
		if err != nil {
			return err
		}

		cli.Log.Debugw("storing version cache", "content", currentVersion)
		if err := currentVersion.StoreCache(cacheFile); err != nil {
			return err
		}

	}

	cli.Log.Debugw("verifying cached version", "cache_file", cacheFile)
	versionCache, err := lwupdater.LoadCache(cacheFile)
	if err != nil {
		return err
	}

	cli.Log.Debugw("version cache", "content", versionCache)

	// since our check is daily, substract one day from now and compare it
	checkTime := time.Now().AddDate(0, 0, -1)
	if versionCache.LastCheckTime.Before(checkTime) {
		versionCache.LastCheckTime = time.Now()

		cli.Log.Debugw("storing new version cache", "content", versionCache)
		err := versionCache.StoreCache(cacheFile)
		if err != nil {
			return err
		}

		_, err = versionCheck()
		return err
	}

	cli.Log.Debugw("threshold not yet met. skipping daily version check",
		"threshold", "1d",
		"current_version", versionCache.CurrentVersion,
		"latest_version", versionCache.LatestVersion,
		"version_outdated", versionCache.Outdated,
		"last_check_time", versionCache.LastCheckTime,
		"next_check_time", versionCache.LastCheckTime.AddDate(0, 0, 1))

	return nil
}
