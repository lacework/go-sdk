// +build configure

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

package integration

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureCommandNonInteractive(t *testing.T) {
	// create a temporal directory where we will check that the
	// configuration file is deployed (.lacework.toml)
	home, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(home)
	out, errB, exitcode := LaceworkCLIWithHome(home, "configure",
		"--noninteractive",
		"-a", "my-account",
		"-k", "my-key",
		"-s", "my-secret",
		"--subaccount", "my-sub-account",
	)

	assert.Empty(t, errB.String())
	assert.Equal(t, 0, exitcode)
	assert.Equal(t, "You are all set!\n", out.String(),
		"you are not all set, check configure cmd")

	configPath := path.Join(home, ".lacework.toml")
	assert.FileExists(t, configPath, "the configuration file is missing")
	laceworkTOML, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, `[default]
  account = "my-account"
  subaccount = "my-sub-account"
  api_key = "my-key"
  api_secret = "my-secret"
  version = 2
`, string(laceworkTOML), "there is a problem with the generated config")
}

func createJSONFileLikeWebUI(content string) string {
	contentBytes := []byte(content)
	tmpfile, err := ioutil.TempFile("", "json_file")
	if err != nil {
		panic(err)
	}

	if _, err := tmpfile.Write(contentBytes); err != nil {
		panic(err)
	}
	return tmpfile.Name()
}

func createTOMLConfig() string {
	dir, err := ioutil.TempDir("", "lacework-toml")
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(dir, ".lacework.toml")
	c := []byte(`[default]
account = 'test.account'
api_key = 'INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00'
api_secret = '_00000000000000000000000000000000'
version = 2

[integration]
account = 'integration'
api_key = 'INTEGRATION_3DF1234AABBCCDD5678XXYYZZ1234ABC8BEC6500DC70001'
api_secret = '_1234abdc00ff11vv22zz33xyz1234abc'
version = 2

[dev]
account = 'dev.example'
api_key = 'DEVDEV_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC000'
api_secret = '_11111111111111111111111111111111'
version = 2

[v1]
account = 'v1.example'
api_key = 'V1CONFIG_KEY'
api_secret = '_secret'

[v2]
account = 'v2.example'
api_key = 'V2CONFIG_KEY'
api_secret = '_secret'
subaccount = 'sub-account'
`)
	err = ioutil.WriteFile(configFile, c, 0644)
	if err != nil {
		panic(err)
	}
	return dir
}
