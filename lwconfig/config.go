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
	"io/ioutil"
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

// Verify will return an error is there is one required setting missing
func (p *Profile) Verify() error {
	if p.Account == "" {
		return errors.New("account missing")
	}
	if p.ApiKey == "" {
		return errors.New("api_key missing")
	}
	if p.ApiSecret == "" {
		return errors.New("api_secret missing")
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
	path, err = DefaultConfigPath()
	if err != nil {
		return Profiles{}, err
	}

	return LoadProfilesFrom(path)
}

// LoadProfilesFrom loads all the profiles from the provided location
func LoadProfilesFrom(path string) (Profiles, error) {
	if path == "" {
		return Profiles{}, errors.New("unable to load profiles. Specify a configuration file.")
	}

	var profiles Profiles
	if _, err := toml.DecodeFile(path, &profiles); err != nil {
		return profiles, errors.Wrap(err, "unable to decode profiles from config")
	}

	return profiles, nil
}

// StoreProfileAt updates a single profile from the provided configuration file
func StoreProfileAt(path, name string, profile Profile) error {
	var (
		profiles = Profiles{}
		buf      = new(bytes.Buffer)
		err      error
	)

	if path == "" {
		path, err = DefaultConfigPath()
		if err != nil {
			return err
		}
	}

	profiles, err = LoadProfilesFrom(path)
	if err != nil {
		return err
	}

	profiles[name] = profile
	if err = toml.NewEncoder(buf).Encode(profiles); err != nil {
		return err
	}

	return ioutil.WriteFile(path, buf.Bytes(), 0600)
}
