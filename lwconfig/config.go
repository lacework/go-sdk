//
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

// Helps you manage the Lacework configuration file ($HOME/.lacework.toml)
package lwconfig

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// Profiles is the representation of the $HOME/.lacework.toml
//
// Example:
//
// [default]
// account = "example"
// api_key = "EXAMPLE_0123456789"
// api_secret = "_0123456789"
//
// [dev]
// account = "dev"
// api_key = "DEV_0123456789"
// api_secret = "_0123456789"
//
// [prod]
// account = "coolcorp"
// subaccount = "prod-business"
// api_key = "PROD_0123456789"
// api_secret = "_0123456789"
// version = 2
type Profiles map[string]Profile

// Profile represents a single profile within a configuration file
type Profile struct {
	Account    string `toml:"account"`
	Subaccount string `toml:"subaccount,omitempty"`
	ApiKey     string `toml:"api_key" survey:"api_key"`
	ApiSecret  string `toml:"api_secret" survey:"api_secret"`
	Version    int    `toml:"version,omitzero"`
}

const (
	ApiKeyMinLength    = 55
	ApiSecretMinLength = 30
)

// Verify will return an error is there is one required setting missing
func (p *Profile) Verify() error {
	if p.Account == "" {
		return errors.New("account missing")
	}

	if err := p.verifyApiKey(); err != nil {
		return err
	}

	if err := p.verifyApiSecret(); err != nil {
		return err
	}

	return nil
}

func (p *Profile) verifyApiKey() error {
	if p.ApiKey == "" {
		return errors.New("api_key missing")
	}

	if len(p.ApiKey) < ApiKeyMinLength {
		return errors.New(fmt.Sprintf("api_key must have more than %d characters", ApiKeyMinLength))
	}

	return nil
}

func (p *Profile) verifyApiSecret() error {
	if p.ApiSecret == "" {
		return errors.New("api_secret missing")
	}

	if len(p.ApiSecret) < ApiSecretMinLength {
		return errors.New(fmt.Sprintf("api_secret must have more than %d characters", ApiSecretMinLength))
	}

	return nil
}

// DefaultConfigPath returns the default path where the Lacework config file
// is located, which is at $HOME/.lacework.toml
func DefaultConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".lacework.toml"), nil
}

// LoadProfiles loads all the profiles from the default location ($HOME/.lacework.toml)
func LoadProfiles() (Profiles, error) {
	configPath, err := DefaultConfigPath()
	if err != nil {
		return Profiles{}, err
	}

	return LoadProfilesFrom(configPath)
}

// LoadProfilesFrom loads all the profiles from the provided location
func LoadProfilesFrom(configPath string) (Profiles, error) {
	if configPath == "" {
		return Profiles{}, errors.New("unable to load profiles. Specify a configuration file.")
	}

	var profiles Profiles
	if _, err := toml.DecodeFile(configPath, &profiles); err != nil {
		return profiles, errors.Wrap(err, "unable to decode profiles from config")
	}

	return profiles, nil
}

// StoreProfileAt updates a single profile from the provided configuration file
func StoreProfileAt(configPath, name string, profile Profile) error {
	if configPath == "" {
		defaultPath, err := DefaultConfigPath()
		if err != nil {
			return err
		}
		configPath = defaultPath
	}

	var (
		profiles = Profiles{}
		err      error
	)
	if _, err = os.Stat(configPath); err == nil {
		if profiles, err = LoadProfilesFrom(configPath); err != nil {
			return err
		}
	}

	profiles[name] = profile
	return StoreAt(configPath, profiles)
}

// StoreAt stores the provided profiles into the selected configuration file
func StoreAt(configPath string, profiles Profiles) error {
	if configPath == "" {
		defaultPath, err := DefaultConfigPath()
		if err != nil {
			return err
		}
		configPath = defaultPath
	}

	var buf = new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(profiles); err != nil {
		return err
	}

	return os.WriteFile(configPath, buf.Bytes(), 0600)
}
