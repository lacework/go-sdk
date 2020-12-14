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
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwlogger"
)

// cliState holds the state of the entire Lacework CLI
type cliState struct {
	Profile    string
	Account    string
	Subaccount string
	KeyID      string
	Secret     string
	Token      string
	LogLevel   string

	LwApi *api.Client
	JsonF *prettyjson.Formatter
	Log   *zap.SugaredLogger
	Event *Honeyvent

	id             string
	workers        sync.WaitGroup
	spinner        *spinner.Spinner
	jsonOutput     bool
	nonInteractive bool
	profileDetails map[string]interface{}
}

// NewDefaultState creates a new cliState with some defaults
func NewDefaultState() *cliState {
	c := &cliState{
		id:      newID(),
		Profile: "default",
		Log:     lwlogger.New("").Sugar(),
		JsonF: &prettyjson.Formatter{
			KeyColor:    color.New(color.FgCyan, color.Bold),
			StringColor: color.New(color.FgGreen, color.Bold),
			BoolColor:   color.New(color.FgYellow, color.Bold),
			NumberColor: color.New(color.FgRed, color.Bold),
			NullColor:   color.New(color.FgWhite, color.Bold),
			Indent:      2,
			Newline:     "\n",
		},
	}

	// initialize honeycomb library and honeyvent
	c.InitHoneyvent()

	return c
}

// SetProfile sets the provided profile into the cliState and loads the entire
// state of the Lacework CLI by calling 'LoadState()'
func (c *cliState) SetProfile(profile string) error {
	if profile == "" {
		return errors.New("Specify a profile.")
	}

	c.Profile = profile
	c.Log.Debugw("custom profile", "profile", profile)
	return c.LoadState()
}

// LoadState loads the state of the cli in the following order, loads the
// configured profile out from the viper loaded config, if the profile is
// set to the default and it is not found, we assume that the user is running
// the CLI with parameters or environment variables, so we proceed to load
// those. Though, if the profile is NOT the default, we error out with some
// breadcrumbs to help the user configure the CLI. After loading the profile,
// this function verifies parameters and env variables coming from viper
func (c *cliState) LoadState() error {
	defer func() {
		// update global honeyvent with loaded state
		c.Event.Account = c.Account
		c.Event.Profile = c.Profile
		c.Event.ApiKey = c.KeyID
	}()

	c.profileDetails = viper.GetStringMap(c.Profile)
	if len(c.profileDetails) == 0 {
		if c.Profile != "default" {
			return fmt.Errorf(
				"The profile '%s' could not be found.\n\nTry running 'lacework configure --profile %s'.",
				c.Profile, c.Profile,
			)
		} else {
			c.Log.Debugw("unable to load state from config")
			c.loadStateFromViper()
			return nil
		}
	}

	c.KeyID = c.extractValueString("api_key")
	c.Secret = c.extractValueString("api_secret")
	c.Account = c.extractValueString("account")
	c.Subaccount = c.extractValueString("subaccount")

	c.Log.Debugw("state loaded",
		"profile", c.Profile,
		"account", c.Account,
		"subaccount", c.Subaccount,
		"api_key", c.KeyID,
		"api_secret", c.Secret,
	)

	c.loadStateFromViper()
	return nil
}

// LoadProfiles loads all the profiles from the configuration file
func (c *cliState) LoadProfiles() (Profiles, error) {
	var (
		profiles = Profiles{}
		confPath = viper.ConfigFileUsed()
	)

	if confPath == "" {
		return profiles, errors.New("unable to load profiles. No configuration file found.")
	}

	cli.Log.Debugw("decoding config", "path", confPath)
	if _, err := toml.DecodeFile(confPath, &profiles); err != nil {
		return profiles, errors.Wrap(err, "unable to decode profiles from config")
	}

	cli.Log.Debugw("profiles loaded from config", "profiles", profiles)
	return profiles, nil
}

// VerifySettings checks if the CLI state has the neccessary settings to run,
// if not, it throws an error with breadcrumbs to help the user configure the CLI
func (c *cliState) VerifySettings() error {
	if c.Profile == "" ||
		c.Account == "" ||
		c.Secret == "" ||
		c.KeyID == "" {
		return fmt.Errorf(
			"there is one or more settings missing.\n\nTry running 'lacework configure'.",
		)
	}

	return nil
}

// NewClient creates and stores a new Lacework API client to be used by the CLI
func (c *cliState) NewClient() error {
	err := c.VerifySettings()
	if err != nil {
		return err
	}

	apiOpts := []api.Option{
		api.WithLogLevel(c.LogLevel),
		api.WithApiV2(),
		api.WithApiKeys(c.KeyID, c.Secret),
		api.WithHeader("User-Agent", fmt.Sprintf("Command-Line/%s", Version)),
	}

	if c.Subaccount != "" {
		apiOpts = append(apiOpts, api.WithHeader("Account-Name", c.Subaccount))
	}

	client, err := api.NewClient(c.Account, apiOpts...)
	if err != nil {
		return errors.Wrap(err, "unable to generate api client")
	}

	c.LwApi = client
	return nil
}

// InteractiveMode returns true if the cli is running in interactive mode
func (c *cliState) InteractiveMode() bool {
	return !c.nonInteractive
}

// NonInteractive turns off interactive mode, that is, no progress bars and spinners
func (c *cliState) NonInteractive() {
	c.Log.Info("turning off interactive mode")
	c.nonInteractive = true
}

// StartProgress starts a new progress spinner with the provider suffix and stores it
// into the cli state, make sure to run StopSpinner when you are done processing
func (c *cliState) StartProgress(suffix string) {
	if c.nonInteractive {
		c.Log.Debugw("skipping spinner",
			"noninteractive", c.nonInteractive,
			"action", "start_progress",
		)
		return
	}

	// humans like spinners (^.^)
	if c.HumanOutput() {
		// make sure there is not a spinner already running
		c.StopProgress()

		c.Log.Debug("starting spinner")
		c.spinner = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		c.spinner.Suffix = suffix
		c.spinner.Start()
	}
}

// StopProgress stops the running progress spinner, if any
func (c *cliState) StopProgress() {
	if c.nonInteractive {
		c.Log.Debugw("skipping spinner",
			"noninteractive", c.nonInteractive,
			"action", "stop_progress",
		)
		return
	}

	// humans like spinners (^.^)
	if c.HumanOutput() {
		if c.spinner != nil {
			c.Log.Debug("stopping spinner")
			c.spinner.Stop()
			c.spinner = nil
		}
	}
}

// EnableJSONOutput enables the cli to display JSON output
func (c *cliState) EnableJSONOutput() {
	c.Log.Info("switch output to json format")
	c.jsonOutput = true
}

// EnableJSONOutput enables the cli to display human readable output
func (c *cliState) EnableHumanOutput() {
	c.Log.Info("switch output to human format")
	c.jsonOutput = false
}

// JSONOutput returns true if the cli is configured to display JSON output
func (c *cliState) JSONOutput() bool {
	return c.jsonOutput
}

// HumanOutput returns true if the cli is configured to siplay human readable output
func (c *cliState) HumanOutput() bool {
	return !c.jsonOutput
}

// loadStateFromViper loads parameters and environment variables
// coming from viper into the CLI state
func (c *cliState) loadStateFromViper() {
	if v := viper.GetString("api_key"); v != "" {
		c.KeyID = v
		c.Log.Debugw("state updated", "api_key", c.KeyID)
	}

	if v := viper.GetString("api_secret"); v != "" {
		c.Secret = v
		c.Log.Debugw("state updated", "api_secret", c.Secret)
	}

	if v := viper.GetString("account"); v != "" {
		c.Account = v
		c.Log.Debugw("state updated", "account", c.Account)
	}

	if v := viper.GetString("subaccount"); v != "" {
		c.Subaccount = v
		c.Log.Debugw("state updated", "subaccount", c.Subaccount)
	}
}

func (c *cliState) extractValueString(key string) string {
	if val, ok := c.profileDetails[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		c.Log.Warnw("config value type mismatch",
			"expected_type", "string",
			"file", viper.ConfigFileUsed(),
			"profile", c.Profile,
			"key", key,
			"value", val,
		)
		return ""
	}
	c.Log.Warnw("unable to find key from config",
		"file", viper.ConfigFileUsed(),
		"profile", c.Profile,
		"key", key,
	)
	return ""
}

// newID generates a new client id, this id is useful for logging purposes
// when there are more than one client running on the same machine
// TODO @afiune move this into its own go package (look at api/client.go)
func newID() string {
	now := time.Now().UTC().UnixNano()
	seed := rand.New(rand.NewSource(now))
	return strconv.FormatInt(seed.Int63(), 16)
}
