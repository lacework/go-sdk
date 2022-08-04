//go:build !windows && team_member

// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
)

func TestCreateTeamMember(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	tmResult, err := runTeamMembersTest(t,
		func(c *expect.Console) {
			c.ExpectString("Email:")
			c.SendLine("test.user@email.com")
			c.ExpectString("First Name:")
			c.SendLine("Test")
			c.ExpectString("Last Name:")
			c.SendLine("User")
			c.ExpectString("Company:")
			c.SendLine("Lacework")
			c.ExpectString("Create at Organization Level?")
			c.SendLine("N")
			c.ExpectString("Account Admin?")
			c.Close()
		},
		"tm",
		"create")

	assert.Contains(t, expectedOutput, strings.TrimSpace(tmResult))
}

func TestTeamMemberValidateEmail(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	tmResult, err := runTeamMembersTest(t,
		func(c *expect.Console) {
			c.ExpectString("Email:")
			c.SendLine("invalid")
			c.ExpectString("X Sorry, your reply was invalid: not a valid email invalid")
			c.Close()
		},
		"tm",
		"create")

	assert.Contains(t, "X Sorry, your reply was invalid: not a valid email invalid", strings.TrimSpace(tmResult))
}

func runTeamMembersTest(t *testing.T, conditions func(*expect.Console), args ...string) (string, error) {
	dir := createDummyTOMLConfig()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	os.Setenv("api_token", "test")
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	return runTeamMemberTestFromDir(t, dir, conditions, args...)
}

func runTeamMemberTestFromDir(t *testing.T, dir string, conditions func(*expect.Console), args ...string) (string, error) {
	console, state, err := vt10x.NewVT10XConsole()
	if err != nil {
		return "", err
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

	cmd := NewLaceworkCLI(dir, nil, args...)
	cmd.Stdin = console.Tty()
	cmd.Stdout = console.Tty()
	cmd.Stderr = console.Tty()
	err = cmd.Start()
	assert.Nil(t, err)

	console.Tty().Close()
	<-donec

	return state.String(), err
}

var expectedOutput = `▸ Email:  test.user@email.com                                                   
▸ First Name:  Test                                                             
▸ Last Name:  User                                                              
▸ Company:  Lacework                                                            
▸ Create at Organization Level? No                                              
? Account Admin?`
