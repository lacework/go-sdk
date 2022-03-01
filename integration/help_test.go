//go:build help

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
	"embed"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	availableCommandsBlobRE = regexp.MustCompile(`(?ims)^Available Commands:.*?^\s*$`)
	availableCommandRE      = regexp.MustCompile(`(?im)^\s+([\w-]+)`)
	//go:embed test_resources/help/*
	helpCanon embed.FS
)

func getAllCommands(in string, commands [][]string, working []string) [][]string {
	availableCommandBlob := availableCommandsBlobRE.FindString(in)
	availableCommands := availableCommandRE.FindAllStringSubmatch(availableCommandBlob, -1)

	for _, match := range availableCommands {
		cmd := match[1]

		// push + add item
		this_working := append(working, cmd)
		commands = append(commands, this_working)

		// get output
		out, _, _ := LaceworkCLI(append([]string{"help"}, this_working...)...)

		// recurse
		commands = getAllCommands(out.String(), commands, this_working)
	}

	return commands
}

func TestHelpAll(t *testing.T) {
	out, _, exitcode := LaceworkCLI("help")
	if exitcode != 0 {
		assert.FailNow(t, "Something went terribly wrong")
	}

	commands := getAllCommands(out.String(), [][]string{}, []string{})

	for _, cmd := range commands {
		cmdStr := strings.Join(cmd, "_")

		t.Run(cmdStr, func(t *testing.T) {
			filePath := fmt.Sprintf("test_resources/help/%s", cmdStr)
			windowsFilePath := fmt.Sprintf("test_resources/help/windows/%s", cmdStr)

			// run command
			out, err, exitcode := LaceworkCLI(append([]string{"help"}, cmd...)...)

			// validate proper execution
			assert.Empty(t, err.String(), "STDERR should be empty")
			assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

			// validate expected output
			if runtime.GOOS == "windows" {
				canon, err := helpCanon.ReadFile(windowsFilePath)
				if err != nil {
					assert.Equal(t, out.String(), string(canon))
					return
				}
			}
			canon, _ := helpCanon.ReadFile(filePath)
			assert.Equal(t, out.String(), string(canon))
		})
	}
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
	canon, _ := helpCanon.ReadFile("test_resources/help/no-command-provided")
	assert.Equal(t,
		string(canon),
		out.String(),
		"the main help message changed, please update")
	assert.Empty(t,
		err.String(),
		"STDERR message doesn't match")
	assert.Equal(t, 127, exitcode,
		"EXITCODE is not the expected one")
}
