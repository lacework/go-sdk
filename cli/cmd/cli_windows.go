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
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// used by configure.go
var configureListCmdSetProfileEnv = `$env:LW_PROFILE = 'my-profile'`

// promptIconsFuncs configures the prompt icons for Windows systems
var promptIconsFunc = func(icons *survey.IconSet) {
	icons.Question.Text = ">"
}

// customPromptIconsFunc configures the prompt icons with custom string for Windows systems
var customPromptIconsFunc = func(s string) func(icons *survey.IconSet) {
	return func(icons *survey.IconSet) {
		icons.Question.Text = fmt.Sprintf("> %s", s)
	}
}

// A variety of colorized icons used throughout the code
var (
	successIcon = color.HiGreenString("√")
	failureIcon = color.HiRedString("×") //nolint
)

// UpdateCommand returns the command that a user should run to update the cli
// to the latest available version (windows specific command)
func (c *cliState) UpdateCommand() string {
	if os.Getenv(ChocolateyInstall) != "" {
		return `
  choco upgrade lacework-cli
`
	}

	return `
  Set-ExecutionPolicy Bypass -Scope Process -Force;
  iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.ps1'))
`
}
