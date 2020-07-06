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

package integration

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
)

func TestConfigureCommand(t *testing.T) {
	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("test-account")
			c.ExpectString("Access Key ID:")
			c.SendLine("INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00")
			c.ExpectString("Secret Access Key:")
			c.SendLine("_00000000000000000000000000000000")
			c.ExpectString("You are all set!")
		},
		"configure",
	)

	assert.Equal(t, `[default]
  account = "test-account"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_00000000000000000000000000000000"
`, laceworkTOML, "there is a problem with the generated config")
}

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
  api_key = "my-key"
  api_secret = "my-secret"
`, string(laceworkTOML), "there is a problem with the generated config")
}

func TestConfigureCommandWithProfileFlag(t *testing.T) {
	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("test-account")
			c.ExpectString("Access Key ID:")
			c.SendLine("INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00")
			c.ExpectString("Secret Access Key:")
			c.SendLine("_00000000000000000000000000000000")
			c.ExpectString("You are all set!")
		},
		"configure", "--profile", "my-profile",
	)

	assert.Equal(t, `[my-profile]
  account = "test-account"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_00000000000000000000000000000000"
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithJSONFileFlag(t *testing.T) {
	// create a JSON file similar to what the Lacework Web UI would provide
	s := createJSONFileLikeWebUI(`{"keyId": "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00","secret": "_cccccccccccccccccccccccccccccccc"}`)
	defer os.Remove(s)

	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("") // using the default, which should be auto-populated from the provided --profile flag
			c.ExpectString("Access Key ID:")
			c.SendLine("") // using the default, which should be loaded from the JSON file
			c.ExpectString("Secret Access Key:")
			c.SendLine("") // using the default, which should be loaded from the JSON file
			c.ExpectString("You are all set!")
		},
		"configure", "--json_file", s, "--profile", "web-ui-test",
	)

	assert.Equal(t, `[web-ui-test]
  account = "web-ui-test"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_cccccccccccccccccccccccccccccccc"
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithJSONFileFlagError(t *testing.T) {
	out, err, exitcode := LaceworkCLI("configure", "--json_file", "foo")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Contains(t,
		err.String(),
		"ERROR unable to load keys from the provided json file: open foo: no such file or directory",
		"STDERR should be empty")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestConfigureCommandWithEnvironmentVariables(t *testing.T) {
	os.Setenv("LW_ACCOUNT", "env-vars")
	os.Setenv("LW_API_KEY", "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00")
	os.Setenv("LW_API_SECRET", "_cccccccccccccccccccccccccccccccc")
	defer os.Setenv("LW_ACCOUNT", "")
	defer os.Setenv("LW_API_KEY", "")
	defer os.Setenv("LW_API_SECRET", "")

	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("") // using the default, which should be loaded from the environment variables
			c.ExpectString("Access Key ID:")
			c.SendLine("") // using the default, which should be loaded from the environment variables
			c.ExpectString("Secret Access Key:")
			c.SendLine("") // using the default, which should be loaded from the environment variables
			c.ExpectString("You are all set!")
		},
		"configure",
	)

	assert.Equal(t, `[default]
  account = "env-vars"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_cccccccccccccccccccccccccccccccc"
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithAPIkeysFromFlags(t *testing.T) {
	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("") // using the default, which should be loaded from the provided flags
			c.ExpectString("Access Key ID:")
			c.SendLine("") // using the default, which should be loaded from the provided flags
			c.ExpectString("Secret Access Key:")
			c.SendLine("") // using the default, which should be loaded from the provided flags
			c.ExpectString("You are all set!")
		},
		"configure",
		"--account", "from-flags",
		"--api_key", "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00",
		"--api_secret", "_cccccccccccccccccccccccccccccccc",
	)

	assert.Equal(t, `[default]
  account = "from-flags"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_cccccccccccccccccccccccccccccccc"
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithExistingConfigAndMultiProfile(t *testing.T) {
	dir := createTOMLConfig()
	defer os.RemoveAll(dir)

	_, laceworkTOML := runConfigureTestFromDir(t, dir,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("super-cool-profile")
			c.ExpectString("Access Key ID:")
			c.SendLine("TEST_ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
			c.ExpectString("Secret Access Key:")
			c.SendLine("_uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu")
			c.ExpectString("You are all set!")
		},
		"configure", "--profile", "new-profile",
	)

	assert.Equal(t, `[default]
  account = "test.account"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_00000000000000000000000000000000"

[dev]
  account = "dev.example"
  api_key = "DEVDEV_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC000"
  api_secret = "_11111111111111111111111111111111"

[integration]
  account = "integration"
  api_key = "INTEGRATION_3DF1234AABBCCDD5678XXYYZZ1234ABC8BEC6500DC70001"
  api_secret = "_1234abdc00ff11vv22zz33xyz1234abc"

[new-profile]
  account = "super-cool-profile"
  api_key = "TEST_ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"
  api_secret = "_uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu"
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandErrors(t *testing.T) {
	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("")
			c.ExpectString("The account subdomain of URL is required")
			c.SendLine("my-account")
			c.ExpectString("Access Key ID:")
			c.SendLine("")
			c.ExpectString("The API access key id must have more than 55 characters")
			c.SendLine("INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00")
			c.ExpectString("Secret Access Key:")
			c.SendLine("")
			c.ExpectString("The API secret access key must have more than 30 characters")
			c.SendLine("_00000000000000000000000000000000")
			c.ExpectString("You are all set!")
		},
		"configure",
	)

	assert.Equal(t, `[default]
  account = "my-account"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_00000000000000000000000000000000"
`, laceworkTOML, "there is a problem with the generated config")
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

func runConfigureTest(t *testing.T, conditions func(*expect.Console), args ...string) (string, string) {
	// create a temporal directory where we will check that the
	// configuration file is deployed (.lacework.toml)
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	return runConfigureTestFromDir(t, dir, conditions, args...)
}

func runConfigureTestFromDir(t *testing.T, dir string, conditions func(*expect.Console), args ...string) (string, string) {
	console, state, err := vt10x.NewVT10XConsole()
	if err != nil {
		panic(err)
	}
	defer console.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		conditions(console)
	}()

	// spawn a new `lacework configure' command
	cmd := NewLaceworkCLI(dir, args...)
	cmd.Stdin = console.Tty()
	cmd.Stdout = console.Tty()
	cmd.Stderr = console.Tty()
	err = cmd.Start()
	assert.Nil(t, err)

	// read the remaining bytes
	console.Tty().Close()
	<-donec

	configPath := path.Join(dir, ".lacework.toml")
	assert.Contains(t, state.String(), "You are all set!", "you are not all set, check configure cmd")
	assert.FileExists(t, configPath, "the configuration file is missing")
	laceworkTOML, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	return state.String(), string(laceworkTOML)
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

[integration]
account = 'integration'
api_key = 'INTEGRATION_3DF1234AABBCCDD5678XXYYZZ1234ABC8BEC6500DC70001'
api_secret = '_1234abdc00ff11vv22zz33xyz1234abc'

[dev]
account = 'dev.example'
api_key = 'DEVDEV_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC000'
api_secret = '_11111111111111111111111111111111'
`)
	err = ioutil.WriteFile(configFile, c, 0644)
	if err != nil {
		panic(err)
	}
	return dir
}
