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

	"github.com/fatih/color"
	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// cliState holds the state of the entire Lacework CLI
type cliState struct {
	Profile  string
	Account  string
	KeyID    string
	Secret   string
	Token    string
	LogLevel string

	JsonF *prettyjson.Formatter
	Log   *zap.SugaredLogger

	profileDetails map[string]interface{}
}

// NewDefaultState creates a new cliState with some defaults
func NewDefaultState() cliState {
	return cliState{
		Profile: "default",
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

	c.Log.Debugw("state loaded",
		"profile", c.Profile,
		"account", c.Account,
		"api_key", c.KeyID,
		"api_secret", c.Secret,
	)

	c.loadStateFromViper()
	return nil
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
