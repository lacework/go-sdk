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
	"fmt"
	"io/ioutil"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// the representation of the ~/.lacework.toml
type config struct {
	Account   string `toml:"account"`
	ApiKey    string `toml:"api_key"`
	ApiSecret string `toml:"api_secret"`
}

var (
	// configureCmd represents the configure command
	configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Set up your Lacework cli",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			var (
				promptAccount = promptui.Prompt{
					Label:   "Account",
					Default: cli.Account,
					Validate: func(input string) error {
						if len(input) == 0 {
							return errors.New(
								"The account subdomain of URL is required. (i.e. <ACCOUNT>.lacework.net)")
						}
						return nil
					},
				}
				promptApiKey = promptui.Prompt{
					Label:   "Access Key ID",
					Default: cli.KeyID,
					Validate: func(input string) error {
						if len(input) < 55 {
							return errors.New(
								"The API access key id must have more than 55 characters.")
						}
						return nil
					},
				}
				promptApiSecret = promptui.Prompt{
					Label:   "Secret Access Key",
					Default: cli.Secret,
					Validate: func(input string) error {
						if len(input) < 30 {
							return errors.New(
								"The API secret access key must have more than 30 characters.")
						}
						return nil
					},
				}
			)
			account, err := promptAccount.Run()
			if err != nil {
				return err
			}
			apiKey, err := promptApiKey.Run()
			if err != nil {
				return err
			}
			apiSecret, err := promptApiSecret.Run()
			if err != nil {
				return err
			}

			var (
				c        = config{account, apiKey, apiSecret}
				buf      = new(bytes.Buffer)
				confPath = viper.ConfigFileUsed()
			)

			if err := toml.NewEncoder(buf).Encode(c); err != nil {
				return err
			}

			if confPath == "" {
				home, err := homedir.Dir()
				if err != nil {
					return err
				}
				confPath = path.Join(home, ".lacework.toml")
			}

			err = ioutil.WriteFile(confPath, buf.Bytes(), 0600)
			if err != nil {
				return err
			}
			fmt.Println("\nYou are all set!")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(configureCmd)
}
