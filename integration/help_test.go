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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpCommand(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help")
	assert.Contains(t,
		out.String(),
		"Use \"lacework [command] --help\" for more information about a command.",
		"STDOUT bottom message doesn't match")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestHelpCommandForConfigureCommand(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "configure")
	assert.Equal(t,
		`
Configure settings that the Lacework CLI uses to interact with the Lacework
platform. These include your Lacework account, API access key and secret.

To create a set of API keys, log in to your Lacework account via WebUI and
navigate to Settings > API Keys and click + Create New. Enter a name for
the key and an optional description, then click Save. To get the secret key,
download the generated API key file.

Use the argument --json_file to preload the downloaded API key file.

If this command is run with no arguments, the Lacework CLI will store all
settings under the default profile. The information in the default profile
is used any time you run a Lacework CLI command that doesn't explicitly
specify a profile to use.

You can configure multiple profiles by using the --profile argument. If a
config file does not exist (the default location is ~/.lacework.toml), the
Lacework CLI will create it for you.

Usage:
  lacework configure [flags]

Flags:
  -h, --help               help for configure
  -j, --json_file string   loads the generated API key JSON file from the WebUI

Global Flags:
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml
`,
		out.String(),
		"the configure help message changed, please update")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestHelpCommandDisplayHelpFromUnknownCommand(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "foo")
	// this is an unknown command, we should display the help message via STDERR
	assert.Contains(t,
		err.String(),
		"Use \"lacework [command] --help\" for more information about a command.",
		"STDERR bottom message doesn't match")
	assert.Contains(t,
		err.String(),
		"Unknown help topic [`foo`]",
		"missing unknown message")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestCommandDoesNotExist(t *testing.T) {
	out, err, exitcode := LaceworkCLI("foo")
	assert.Empty(t,
		out.String(),
		"STDOUT should be empty")
	assert.Contains(t,
		err.String(),
		"ERROR unknown command \"foo\" for \"lacework\"",
		"STDERR message doesn't match")
	assert.Equal(t, 1, exitcode,
		"EXITCODE is not the expected one")
}

func TestNoCommandProvided(t *testing.T) {
	out, err, exitcode := LaceworkCLI()
	assert.Equal(t,
		`
The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. Use it to manage compliance reports,
external integrations, vulnerability scans, and other operations.

Start by configuring the Lacework CLI with the command:

    $ lacework configure

This will prompt you for your Lacework account and a set of API access keys.

Usage:
  lacework [command]

Available Commands:
  access-token  generate temporary access tokens
  api           helper to call Lacework's RestfulAPI
  compliance    manage compliance reports
  configure     configure the Lacework CLI
  event         inspect Lacework events
  integration   manage external integrations
  version       print the Lacework CLI version
  vulnerability view vulnerability reports and run on-demand scans

Flags:
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml

Use "lacework [command] --help" for more information about a command.
`,
		out.String(),
		"the main help message changed, please update")
	assert.Empty(t,
		err.String(),
		"STDERR message doesn't match")
	assert.Equal(t, 127, exitcode,
		"EXITCODE is not the expected one")
}
