//go:build alert_profile

// Author:: Darren Murray (<darren.murray@lacework.net>)
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
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
)

func TestAlertProfileUpdateEditor(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("ap", "update")
	assert.Contains(t, out.String(), "Select an alert profile to update:")
	assert.Contains(t, err.String(), "ERROR unable to update alert profile:")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

// CUSTOM_CUSTOMER_DEMO_GCP
func TestAlertProfileUpdate(t *testing.T) {
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	tmResult, err := runAlertProfileTest(t,
		func(c *expect.Console) {
			c.ExpectString("? Update alert templates for profile CUSTOM_CUSTOMER_DEMO_GCP [Enter to launch editor] ")
			c.SendLine("")
			time.Sleep(time.Millisecond)
			c.SendLine(":wq!") // save and close
			time.Sleep(time.Millisecond)
			c.Close()
		},
		"ap", "update", "CUSTOM_CUSTOMER_DEMO_GCP")

	assert.Contains(t, "The alert profile CUSTOM_CUSTOMER_DEMO_GCP was updated", strings.TrimSpace(tmResult))
}

func runAlertProfileTest(t *testing.T, conditions func(*expect.Console), args ...string) (string, error) {
	dir := createTOMLConfigFromCIvars()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	os.Setenv("api_token", "test")
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	return runAlertProfileTestFromDir(t, dir, conditions, args...)
}

func runAlertProfileTestFromDir(t *testing.T, dir string, conditions func(*expect.Console), args ...string) (string, error) {
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
