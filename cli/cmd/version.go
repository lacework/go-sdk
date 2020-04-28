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
			cli.StartProgress(" Checking available updates...")
			sdk, err := lwupdater.Check("go-sdk", fmt.Sprintf("v%s", Version))
			cli.StopProgress()
			if err != nil {
				exitwithCode(errors.Wrap(err, "unable to check for updates"), 4)
			}
			if sdk.Outdated {
				cli.OutputHuman(fmt.Sprintf(
					"\nA newer version of the Lacework CLI is available! The latest version is %s,\n"+
						"to update execute the following command:\n%s\n",
					sdk.Latest, cli.UpdateCommand()))
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
