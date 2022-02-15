//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	configureSwitchProfileCmd = &cobra.Command{
		Use:     "switch-profile <profile>",
		Aliases: []string{"switch", "use"},
		Args:    cobra.ExactArgs(1),
		Short:   "Switch between configured profiles",
		Long: `Switch between profiles configured into the config file ~/.lacework.toml

An alternative to temporarily switch to a different profile in your current terminal
is to export the environment variable:

    ` + configureListCmdSetProfileEnv,
		RunE: func(_ *cobra.Command, args []string) error {
			if args[0] == "default" {
				cli.Log.Debug("removing global profile cache, going back to default")
				if err := cli.Cache.Erase("global/profile"); err != nil {
					return errors.Wrap(err, "unable to switch profile")
				}
				cli.OutputHuman("Profile switched back to default.\n")
				return nil
			}

			profiles, err := cli.LoadProfiles()
			if err != nil {
				return err
			}

			if _, ok := profiles[args[0]]; ok {
				cli.Log.Debugw("storing global profile cache",
					"current_profile", cli.Profile,
					"new_profile", args[0],
				)
				if err := cli.Cache.Write("global/profile", []byte(args[0])); err != nil {
					return errors.Wrap(err, "unable to switch profile")
				}
				cli.OutputHuman("Profile switched to '%s'.\n", args[0])
				return nil
			}

			return errors.Errorf(
				"Profile '%s' not found. Try 'lacework configure list' to see all configured profiles.",
				args[0],
			)
		},
	}
)

func init() {
	configureCmd.AddCommand(configureSwitchProfileCmd)
}
