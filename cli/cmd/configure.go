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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lacework/go-sdk/internal/domain"
	"github.com/lacework/go-sdk/internal/format"
	"github.com/lacework/go-sdk/lwconfig"
)

var (
	// configureJsonFile is the API key file downloaded form the Lacework WebUI
	configureJsonFile string

	// configureCmd represents the configure command
	configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "configure the Lacework CLI",
		Args:  cobra.NoArgs,
		Long: `Configure settings that the Lacework CLI uses to interact with the Lacework
platform. These include your Lacework account, API access key and secret.

To create a set of API keys, log in to your Lacework account via WebUI and
navigate to Settings > API Keys and click + Create New. Enter a name for
the key and an optional description, then click Save. To get the secret key,
download the generated API key file.

Use the flag --json_file to preload the downloaded API key file.

If this command is run with no flags, the Lacework CLI will store all
settings under the default profile. The information in the default profile
is used any time you run a Lacework CLI command that doesn't explicitly
specify a profile to use.

You can configure multiple profiles by using the --profile flag. If a
config file does not exist (the default location is ~/.lacework.toml),
the Lacework CLI will create it for you.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runConfigureSetup()
		},
	}

	configureListCmd = &cobra.Command{
		Use:   "list",
		Short: "list all configured profiles at ~/.lacework.toml",
		Args:  cobra.NoArgs,
		Long: `List all profiles configured into the config file ~/.lacework.toml

To switch to a different profile permanently in your current terminal,
export the environment variable:

    ` + configureListCmdSetProfileEnv,
		RunE: func(_ *cobra.Command, _ []string) error {
			profiles, err := cli.LoadProfiles()
			if err != nil {
				return err
			}

			cli.OutputHuman(
				renderSimpleTable(
					[]string{"Profile", "Account", "Subaccount", "API Key", "API Secret", "V"},
					buildProfilesTableContent(cli.Profile, profiles),
				),
			)
			return nil
		},
	}

	configureGetCmd = &cobra.Command{
		Use:   "show <config_key>",
		Short: "show current configuration data",
		Args:  cobra.ExactArgs(1),
		Long: `Prints the current computed configuration data from the specified configuration
key. The order of precedence to compute the configuration is flags, environment
variables, and the configuration file ~/.lacework.toml. 

The available configuration keys are:
* profile
* account
* api_secret
* api_key

To show the configuration from a different profile, use the flag --profile.

    $ lacework configure show account --profile my-profile`,
		RunE: func(_ *cobra.Command, args []string) error {
			data, ok := showConfigurationDataFromKey(args[0])
			if !ok {
				// TODO change this to be dynamic
				return errors.New("unknown configuration key. (available: profile, account, subaccount, api_secret, api_key, version)")
			}

			if data != "" {
				cli.OutputHuman(data)
				cli.OutputHuman("\n")
			}

			return nil
		},
	}
)

func showConfigurationDataFromKey(key string) (string, bool) {
	switch key {
	case "profile":
		return cli.Profile, true
	case "account":
		return cli.Account, true
	case "subaccount":
		return cli.Subaccount, true
	case "api_secret":
		return cli.Secret, true
	case "api_key":
		return cli.KeyID, true
	case "version":
		return fmt.Sprintf("%d", cli.CfgVersion), true
	default:
		return "", false
	}
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.AddCommand(configureListCmd)
	configureCmd.AddCommand(configureGetCmd)

	configureCmd.Flags().StringVarP(&configureJsonFile,
		"json_file", "j", "", "loads the API key JSON file downloaded from the WebUI",
	)
}

func runConfigureSetup() error {
	cli.Log.Debugw("configuring cli", "profile", cli.Profile)

	// make sure that the state is loaded to use during configuration
	cli.loadStateFromViper()

	// if the Lacework account is empty, and the profile that is being configured is
	// not the 'default' profile, auto-populate the account with the provided profile
	if cli.Account == "" && cli.Profile != "default" {
		cli.Account = cli.Profile
	}

	if len(configureJsonFile) != 0 {
		err := loadUIJsonFile(configureJsonFile)
		if err != nil {
			return errors.Wrap(err, "unable to load keys from the provided json file")
		}
	}

	// all new configurations should default to version 2
	cli.CfgVersion = 2

	newProfile := lwconfig.Profile{
		Version:    cli.CfgVersion,
		Subaccount: cli.Subaccount,
		Account:    cli.Account,
		ApiKey:     cli.KeyID,
		ApiSecret:  cli.Secret,
	}
	if cli.InteractiveMode() {
		if err := promptConfigureSetup(&newProfile); err != nil {
			return err
		}

		// before trying to detect if the account is organizational or not, and to
		// check if there are sub-accounts, we need to update the CLI settings
		cli.Log.Debug("storing interactive information into the cli state")
		cli.Account = newProfile.Account
		cli.Subaccount = newProfile.Subaccount
		cli.Secret = newProfile.ApiSecret
		cli.KeyID = newProfile.ApiKey

		// generate a new API client to connect and check for sub-accounts
		if err := cli.NewClient(); err != nil {
			return err
		}

		// get sub-accounts from organizational accounts
		subaccount, err := getSubAccountForOrgAdmins()
		if err != nil {
			return err
		}

		// only configure the subaccount if it is not empty
		if subaccount != "" {
			newProfile.Subaccount = subaccount
		}

		cli.OutputHuman("\n")
	}

	if err := newProfile.Verify(); err != nil {
		return errors.Wrap(err, "unable to configure the command-line")
	}

	if err := lwconfig.StoreProfileAt(viper.ConfigFileUsed(), cli.Profile, newProfile); err != nil {
		return errors.Wrap(err, "unable to configure the command-line")
	}

	cli.OutputHuman("You are all set!\n")
	return nil
}

func promptConfigureSetup(newProfile *lwconfig.Profile) error {
	questions := []*survey.Question{
		{
			Name: "account",
			Prompt: &survey.Input{
				Message: "Account:",
				Default: cli.Account,
			},
			Validate: promptRequiredStringLen(1,
				"The account subdomain of URL is required. (i.e. <ACCOUNT>.lacework.net)",
			),
			Transform: func(ans interface{}) interface{} {
				answer, ok := ans.(string)
				if ok && strings.Contains(answer, ".lacework.net") {

					d, err := domain.New(answer)
					if err != nil {
						cli.Log.Warn(err)
						return answer
					}
					cli.OutputHuman("\nPassing full 'lacework.net' domain not required. Using '%s'\n", d.String())
					return d.String()
				}

				return ans
			},
		},
		{
			Name: "api_key",
			Prompt: &survey.Input{
				Message: "Access Key ID:",
				Default: cli.KeyID,
			},
			Validate: promptRequiredStringLen(55,
				"The API access key id must have more than 55 characters.",
			),
		},
	}

	secretQuest := &survey.Question{
		Name: "api_secret",
		Validate: func(input interface{}) error {
			str, ok := input.(string)
			if !ok || len(str) < 30 {
				if len(str) == 0 && len(cli.Secret) != 0 {
					return nil
				}
				return errors.New("The API secret access key must have more than 30 characters.")
			}
			return nil
		},
	}

	secretMessage := "Secret Access Key:"
	if len(cli.Secret) != 0 {
		secretMessage = fmt.Sprintf("Secret Access Key: (%s)", format.Secret(4, cli.Secret))
	}
	secretQuest.Prompt = &survey.Password{
		Message: secretMessage,
	}

	err := survey.Ask(append(questions, secretQuest), newProfile,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	if len(newProfile.ApiSecret) == 0 {
		newProfile.ApiSecret = cli.Secret
	}

	return nil
}

func getSubAccountForOrgAdmins() (string, error) {
	cli.StartProgress(" Verifying credentials ...")
	user, err := cli.LwApi.V2.UserProfile.Get()
	cli.StopProgress()
	if err != nil {
		cli.Log.Warnw("unable to access UserProfile endpoint",
			"error", err,
		)
		// We do NOT error here since API v2 is sending 500 errors
		// for mortal users, we need to fix this on the server side
		return "", nil
	}

	// We only ask for the sub-account if the account is an organizational account
	// and it has at least one sub-account other than the primary account
	if len(user.Data) != 0 &&
		user.Data[0].OrgAccount &&
		len(user.Data[0].SubAccountNames()) > 0 {

		var (
			subaccount string
			primary    = fmt.Sprintf("PRIMARY (%s)", cli.Account)
		)
		err := survey.AskOne(&survey.Select{
			Message: "(Org Admins) Managing a sub-account?",
			Default: strings.ToLower(cli.Subaccount),
			Options: append([]string{primary}, user.Data[0].SubAccountNames()...),
		}, &subaccount, survey.WithIcons(promptIconsFunc))
		if err != nil {
			return "", err
		}

		if subaccount != primary {
			return subaccount, nil
		}
	}

	return "", nil
}

// apiKeyDetails represents the details of an API key, we use this struct
// internally to unmarshal the JSON file provided by the Lacework WebUI
type apiKeyDetails struct {
	Account    string `json:"account,omitempty"`
	SubAccount string `json:"subAccount,omitempty"`
	KeyID      string `json:"keyId"`
	Secret     string `json:"secret"`
}

func loadUIJsonFile(file string) error {
	cli.Log.Debugw("loading API key JSON file", "path", file)
	jsonData, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	cli.Log.Debugw("JSON file", "raw", string(jsonData))
	var auth apiKeyDetails
	err = json.Unmarshal(jsonData, &auth)
	if err != nil {
		return err
	}

	cli.KeyID = auth.KeyID
	cli.Secret = auth.Secret
	cli.Subaccount = strings.ToLower(auth.SubAccount)

	if auth.Account != "" {
		d, err := domain.New(auth.Account)
		if err != nil {
			return err
		}
		cli.Account = d.String()
	}

	return nil
}

func buildProfilesTableContent(current string, profiles lwconfig.Profiles) [][]string {
	out := [][]string{}
	for profile, creds := range profiles {
		out = append(out, []string{
			profile,
			creds.Account,
			creds.Subaccount,
			creds.ApiKey,
			format.Secret(4, creds.ApiSecret),
			fmt.Sprintf("%d", creds.Version),
		})
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return out[i][0] < out[j][0]
	})

	for i := range out {
		if out[i][0] == current {
			out[i][0] = fmt.Sprintf("> %s", out[i][0])
		} else {
			out[i][0] = fmt.Sprintf("  %s", out[i][0])
		}
	}

	return out
}
