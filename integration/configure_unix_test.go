// +build configure configure,!windows

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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
)

func TestConfigureCommand(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

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
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandForFrankfurtDatacenter(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			// if the full URL was provided we transform it and inform the user
			c.SendLine("my-account-in.fra.lacework.net")
			c.ExpectString("Passing full 'lacework.net' domain not required. Using 'my-account-in.fra'")
			c.ExpectString("Access Key ID:")
			c.SendLine("FRANK_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC0011")
			c.ExpectString("Secret Access Key:")
			c.SendLine("_00000000000000000000000000000000")
			c.ExpectString("You are all set!")
		},
		"configure",
	)
	assert.Equal(t, `[default]
  account = "my-account-in.fra"
  api_key = "FRANK_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC0011"
  api_secret = "_00000000000000000000000000000000"
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandForOrgAdmins(t *testing.T) {
	if os.Getenv("CI_STANDALONE_ACCOUNT") != "" {
		t.Skip("skipping organizational account test")
	}
	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine(os.Getenv("CI_V2_ACCOUNT"))
			c.ExpectString("Access Key ID:")
			c.SendLine(os.Getenv("CI_API_KEY"))
			c.ExpectString("Secret Access Key:")
			c.SendLine(os.Getenv("CI_API_SECRET"))
			c.ExpectString("Verifying credentials ...")
			c.ExpectString("(Org Admins) Managing a sub-account?")
			// @afiune this is needed just because we have two accounts that start exactly the same
			// and so, we need to key in ARROW DOWN to chose the right one.
			c.SendLine(fmt.Sprintf("%s\x1B[B", os.Getenv("CI_ACCOUNT")))
			c.ExpectString("You are all set!")
		},
		"configure",
	)

	assert.Equal(t, `[default]
  account = "`+os.Getenv("CI_V2_ACCOUNT")+`"
  subaccount = "`+os.Getenv("CI_ACCOUNT")+`"
  api_key = "`+os.Getenv("CI_API_KEY")+`"
  api_secret = "`+os.Getenv("CI_API_SECRET")+`"
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithProfileFlag(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

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
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithNewJSONFileFlagForStandaloneAccounts(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// create a New JSON file similar to what the Lacework Web UI would provide
	s := createJSONFileLikeWebUI(`
{
  "account": "standalone.lacework.net",
  "keyId": "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00",
  "secret": "_cccccccccccccccccccccccccccccccc"
}
`)
	defer os.Remove(s)

	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("") // using the default, which should be auto-populated from the new JSON file
			c.ExpectString("Access Key ID:")
			c.SendLine("") // using the default, which should be loaded from the JSON file
			c.ExpectString("Secret Access Key:")
			c.SendLine("") // using the default, which should be loaded from the JSON file
			c.ExpectString("You are all set!")
		},
		"configure", "--json_file", s,
	)

	assert.Equal(t, `[default]
  account = "standalone"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_cccccccccccccccccccccccccccccccc"
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithNewJSONFileFlagForOrganizationalAccounts(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// create a New JSON file similar to what the Lacework Web UI would provide
	s := createJSONFileLikeWebUI(`
{
  "subAccount": "sub-account-name",
  "account": "organization.lacework.net",
  "keyId": "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00",
  "secret": "_cccccccccccccccccccccccccccccccc"
}
`)
	defer os.Remove(s)

	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("") // using the default, which should be auto-populated from the new JSON file
			c.ExpectString("Access Key ID:")
			c.SendLine("") // using the default, which should be loaded from the JSON file
			c.ExpectString("Secret Access Key:")
			c.SendLine("") // using the default, which should be loaded from the JSON file
			c.ExpectString("You are all set!")
		},
		"configure", "--json_file", s,
	)

	assert.Equal(t, `[default]
  account = "organization"
  subaccount = "sub-account-name"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_cccccccccccccccccccccccccccccccc"
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithOldJSONFileFlag(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// create a JSON file similar to what the Lacework Web UI would provide
	s := createJSONFileLikeWebUI(`
{
  "keyId": "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00",
  "secret": "_cccccccccccccccccccccccccccccccc"
}
`)
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
		"configure", "--json_file", s, "--profile", "v1-web-ui-test",
	)

	assert.Equal(t, `[v1-web-ui-test]
  account = "v1-web-ui-test"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_cccccccccccccccccccccccccccccccc"
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithEnvironmentVariables(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	os.Setenv("LW_ACCOUNT", "env-vars")
	os.Setenv("LW_SUBACCOUNT", "sublime")
	os.Setenv("LW_API_KEY", "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00")
	os.Setenv("LW_API_SECRET", "_cccccccccccccccccccccccccccccccc")
	defer os.Setenv("LW_NOCACHE", "")
	defer os.Setenv("LW_ACCOUNT", "")
	defer os.Setenv("LW_SUBACCOUNT", "")
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
  subaccount = "sublime"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_cccccccccccccccccccccccccccccccc"
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithAPIkeysFromFlagsWithoutSubaccount(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

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
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithAPIkeysFromFlagsWithSubaccount(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

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
		"--subaccount", "sublime",
		"--api_key", "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00",
		"--api_secret", "_cccccccccccccccccccccccccccccccc",
	)

	assert.Equal(t, `[default]
  account = "from-flags"
  subaccount = "sublime"
  api_key = "INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00"
  api_secret = "_cccccccccccccccccccccccccccccccc"
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandWithExistingConfigAndMultiProfile(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

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
  version = 2

[dev]
  account = "dev.example"
  api_key = "DEVDEV_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC000"
  api_secret = "_11111111111111111111111111111111"
  version = 2

[integration]
  account = "integration"
  api_key = "INTEGRATION_3DF1234AABBCCDD5678XXYYZZ1234ABC8BEC6500DC70001"
  api_secret = "_1234abdc00ff11vv22zz33xyz1234abc"
  version = 2

[new-profile]
  account = "super-cool-profile"
  api_key = "TEST_ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"
  api_secret = "_uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu"
  version = 2

[v1]
  account = "v1.example"
  api_key = "V1CONFIG_KEY"
  api_secret = "_secret"
`, laceworkTOML, "there is a problem with the generated config")
}

func TestConfigureCommandErrors(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	_, laceworkTOML := runConfigureTest(t,
		func(c *expect.Console) {
			c.ExpectString("Account:")
			c.SendLine("")
			c.ExpectString("The account subdomain of URL is required")
			// if the full URL was provided we transform it and inform the user
			c.SendLine("https://my-account.lacework.net")
			c.ExpectString("Passing full 'lacework.net' domain not required. Using 'my-account'")
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
  version = 2
`, laceworkTOML, "there is a problem with the generated config")
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

	if os.Getenv("DEBUG") != "" {
		state.DebugLogger = log.Default()
	}

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
	assert.Nil(t, err)
	return state.String(), string(laceworkTOML)
}

func TestConfigureCommandWithJSONFileFlagError(t *testing.T) {
	out, err, exitcode := LaceworkCLI("configure", "--json_file", "foo")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Contains(t,
		err.String(),
		"ERROR unable to load keys from the provided json file: open foo: no such file or directory",
		"STDERR error message changed, please check")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestConfigureListHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("configure", "list", "--help")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Contains(t,
		out.String(),
		"$ export LW_PROFILE=\"my-profile\"",
		"STDOUT the environment variable in the help message is not correct")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}
