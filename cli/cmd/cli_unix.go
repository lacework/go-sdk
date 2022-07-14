//go:build !windows

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
	"github.com/fatih/color"
)

// used by configure.go
var configureListCmdSetProfileEnv = `export LW_PROFILE="my-profile"`

// promptIconsFuncs configures the prompt icons for Unix systems
var promptIconsFunc = func(icons *survey.IconSet) {
	icons.Question.Text = "▸"
}

// A variety of colorized icons used throughout the code
var (
	successIcon = color.HiGreenString("✓")
	failureIcon = color.HiRedString("✖") //nolint
)

// Env variables found in GCP, AWS and Azure cloudshell.
// Used to determine if cli is running on cloudshell.
const (
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

	if isCloudShell() {
		return `
  curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash -s -- -d $HOME/bin
`
	}
	return `
  curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash
`
}

// isCloudShell uses env variables specific to GCP, AWS and Azure
// to determine if the Lacework CLI is running on cloudshell
func isCloudShell() bool {
	return isAwsCloudShell() || isGcpCloudShell() || isAzureCloudShell()
}

// isAzureCloudShell uses the native env variable POWERSHELL_DISTRIBUTION_CHANNEL="CloudShell"
// to determine if the Lacework CLI is running on Azure cloudshell
func isAzureCloudShell() bool {
	return os.Getenv(AzureCloudEnv) == "CloudShell"
}

// isGcpCloudShell uses the native env variable CLOUD_SHELL=true
// to determine if the Lacework CLI is running on GCP cloudshell
func isGcpCloudShell() bool {
	return os.Getenv(gcpCloudEnv) == "true"
}

// isAwsCloudShell uses the native env variable AWS_EXECUTION_ENV="Cloudshell"
// to determine if the Lacework CLI is running on AWS cloudshell
func isAwsCloudShell() bool {
	return os.Getenv(awsCloudEnv) == "Cloudshell"
}
