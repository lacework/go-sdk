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

func TestComplianceCommandHelp(t *testing.T) {
	out, err, exitcode := LaceworkCLI("help", "compliance")
	assert.Equal(t,
		`Manage compliance reports for GCP, Azure, or AWS cloud providers.

To start sending data about your environment to Lacework for compliance reporting
analysis, configure one or more cloud integration using the following command:

  $ lacework integration create

Or, if you prefer to do it via the WebUI, log in to your account at:

    https://<ACCOUNT>.lacework.net

Then navigate to Settings > Integrations > Cloud Accounts.

Use the following command to list all available integrations in your account:

  $ lacework integrations list

Usage:
  lacework compliance [command]

Aliases:
  compliance, comp

Available Commands:
  aws         compliance for AWS
  azure       compliance for Microsoft Azure
  gcp         compliance for Google Cloud

Flags:
  -h, --help   help for compliance

Global Flags:
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml

Use "lacework compliance [command] --help" for more information about a command.
`,
		out.String(),
		"the compliance help message changed, please update")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}
