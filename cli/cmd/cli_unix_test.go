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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliStateUpdateCommand(t *testing.T) {
	assert.Contains(t,
		cli.UpdateCommand(),
		"curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash",
	)

	t.Run("Homebrew installation", func(t *testing.T) {
		os.Setenv("LW_HOMEBREW_INSTALL", "1")
		defer os.Setenv("LW_HOMEBREW_INSTALL", "")
		assert.Contains(t, cli.UpdateCommand(), "brew upgrade lacework-cli")
	})

	t.Run("Gcp CloudShell Installation", func(t *testing.T) {
		os.Setenv("CLOUD_SHELL", "true")
		defer os.Setenv("CLOUD_SHELL", "")
		assert.Contains(t, cli.UpdateCommand(), "curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash -s -- -d $HOME/bin")
	})

	t.Run("Aws CloudShell Installation", func(t *testing.T) {
		os.Setenv("AWS_EXECUTION_ENV", "Cloudshell")
		defer os.Setenv("AWS_EXECUTION_ENV", "")
		assert.Contains(t, cli.UpdateCommand(), "curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash -s -- -d $HOME/bin")
	})

	t.Run("Azure CloudShell Installation", func(t *testing.T) {
		os.Setenv("POWERSHELL_DISTRIBUTION_CHANNEL", "CloudShell")
		defer os.Setenv("POWERSHELL_DISTRIBUTION_CHANNEL", "")
		assert.Contains(t, cli.UpdateCommand(), "curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash -s -- -d $HOME/bin")
	})
}
