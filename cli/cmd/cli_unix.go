// +build !windows
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
	"os"

	"github.com/AlecAivazis/survey/v2"
)

// used by configure.go
var configureListCmdSetProfileEnv = `export LW_PROFILE="my-profile"`

// promptIconsFuncs configures the prompt icons for Unix systems
var promptIconsFunc = func(icons *survey.IconSet) {
	icons.Question.Text = "▸"
}

// Env variables found in GCP, AWS and Azure cloudshell
var (
	gcpCloudEnv   = "CLOUD_SHELL"
	awsCloudEnv   = "AWS_EXECUTION_ENV"
	AzureCloudEnv = "POWERSHELL_DISTRIBUTION_CHANNEL"
)

// UpdateCommand returns the command that a user should run to update the cli
// to the latest available version (unix specific command)
func (c *cliState) UpdateCommand() string {
	if os.Getenv(HomebrewInstall) != "" {
		return `
  brew upgrade lacework-cli
`
	}

	if isCloudshell() {
		return `
  curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash -s -- -d $HOME/bin
`
	}
	return `
  curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash
`
}

func isCloudshell() bool {
	return os.Getenv(gcpCloudEnv) == "true" || os.Getenv(awsCloudEnv) == "Cloudshell" || os.Getenv(AzureCloudEnv) == "CloudShell"
}
