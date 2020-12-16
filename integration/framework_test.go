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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/lacework/go-sdk/lwupdater"
)

// Use this function to execute a real lacework CLI command, under the hood the function
// will detect the correct binary depending on the running OS and architecture, if you
// need to override the binary to use at runtime, set the `LW_CLI_BIN` environment
// variable to the path of the binary you wish to use.
//
// example:
//
//  func TestHelpCommand(t *testing.T) {
//    out, err, exitcode := LaceworkCLI("help")
//
//    assert.Contains(t,
//      out.String(),
//      "Use \"lacework [command] --help\" for more information about a command.",
//      "STDOUT doesn't match")
//    assert.Empty(t,
//      err.String(),
//      "STDERR should be empty")
//    assert.Equal(t, 0, exitcode,
//      "EXITCODE is not the expected one")
//  }
//
func LaceworkCLI(args ...string) (bytes.Buffer, bytes.Buffer, int) {
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	return runLaceworkCLI(dir, args...)
}

func LaceworkCLIWithTOMLConfig(args ...string) (bytes.Buffer, bytes.Buffer, int) {
	dir := createTOMLConfigFromCIvars()
	defer os.RemoveAll(dir)

	return runLaceworkCLI(dir, args...)
}

func LaceworkCLIWithDummyConfig(args ...string) (bytes.Buffer, bytes.Buffer, int) {
	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)

	return runLaceworkCLI(dir, args...)
}

func LaceworkCLIWithHome(dir string, args ...string) (bytes.Buffer, bytes.Buffer, int) {
	return runLaceworkCLI(dir, args...)
}

func NewLaceworkCLI(workingDir string, args ...string) *exec.Cmd {
	cmd := exec.Command(findLaceworkCLIBinary(), args...)
	cmd.Env = os.Environ()
	if len(workingDir) != 0 {
		cmd.Dir = workingDir
		env := append(os.Environ(), fmt.Sprintf("HOME=%s", workingDir))

		// by default, we disable all lwupdater requests, unless we are testing it
		// to test it, set the environment variable CI_TEST_LWUPDATER
		if os.Getenv(ciTestingUpdaterEnv) == "" {
			env = append(env, fmt.Sprintf("%s=1", lwupdater.DisableEnv))
		}
		cmd.Env = env
	}
	return cmd
}

// By default, we disable all lwupdater requests, unless we are testing it
// to test it, set the environment variable CI_TEST_LWUPDATER=1
//
// Example:
//
// func TestUpdaterExample(t *testing.T) {
//   enableTestingUpdaterEnv()
//   defer disableTestingUpdaterEnv()
//
//   // exacute an updater test
// }
var ciTestingUpdaterEnv = "CI_TEST_LWUPDATER"

func enableTestingUpdaterEnv() {
	os.Setenv(ciTestingUpdaterEnv, "1")
}

func disableTestingUpdaterEnv() {
	os.Setenv(ciTestingUpdaterEnv, "")
}

func runLaceworkCLI(workingDir string, args ...string) (stdout bytes.Buffer, stderr bytes.Buffer, exitcode int) {
	cmd := NewLaceworkCLI(workingDir, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitcode = 0
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitcode = exitError.ExitCode()
		} else {
			exitcode = 999
			fmt.Println(stderr)
			if _, err := stderr.WriteString(err.Error()); err != nil {
				// @afiune we should never get here but if we do, lets print the error
				fmt.Println(err)
			}
		}
	}
	return
}

func findLaceworkCLIBinary() string {
	if bin := os.Getenv("LW_CLI_BIN"); bin != "" {
		return bin
	}

	// TODO @afiune add ext for windows support
	if runtime.GOOS != "" && runtime.GOARCH != "" {
		return fmt.Sprintf("lacework-cli-%s-%s", runtime.GOOS, runtime.GOARCH)
	}

	return "lacework"
}

func createTOMLConfigFromCIvars() string {
	if os.Getenv("CI_ACCOUNT") == "" ||
		os.Getenv("CI_API_KEY") == "" ||
		os.Getenv("CI_API_SECRET") == "" {
		log.Fatal(missingCIEnvironmentVariables())
	}

	dir, err := ioutil.TempDir("", "lacework-toml")
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(dir, ".lacework.toml")
	c := []byte(`[default]
account = '` + os.Getenv("CI_ACCOUNT") + `'
api_key = '` + os.Getenv("CI_API_KEY") + `'
api_secret = '` + os.Getenv("CI_API_SECRET") + `'
`)
	err = ioutil.WriteFile(configFile, c, 0644)
	if err != nil {
		panic(err)
	}
	return dir
}

func missingCIEnvironmentVariables() string {
	return `
ERROR
  Missing CI environment variables.

  To run the integration tests you need to setup a few environment variables, look
  at https://github.com/lacework/go-sdk/tree/master/cli#integration-tests for
  more information.

`
}

func createDummyTOMLConfig() string {
	dir, err := ioutil.TempDir("", "lacework-toml")
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(dir, ".lacework.toml")
	c := []byte(`[default]
account = 'dummy'
api_key = 'DUMMY_1234567890abcdefg'
api_secret = '_superdummysecret'

[test]
account = 'test.account'
api_key = 'INTTEST_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC00'
api_secret = '_00000000000000000000000000000000'

[integration]
account = 'integration'
api_key = 'INTEGRATION_3DF1234AABBCCDD5678XXYYZZ1234ABC8BEC6500DC70'
api_secret = '_1234abdc00ff11vv22zz33xyz1234abc'

[dev]
account = 'dev.example'
api_key = 'DEVDEV_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890AAABBBCCC000'
api_secret = '_11111111111111111111111111111111'
`)
	err = ioutil.WriteFile(configFile, c, 0644)
	if err != nil {
		panic(err)
	}
	return dir
}

// store a file in Circle CI Working directory, only if we are running on CircleCI
func storeFileInCircleCI(f string) {
	if jobDir := os.Getenv("CIRCLE_WORKING_DIRECTORY"); jobDir != "" {
		var (
			file      = filepath.Base(f)
			artifacts = path.Join(jobDir, "circleci-artifacts")
			err       = os.Mkdir(artifacts, 0755)
		)
		if err != nil {
			fmt.Println(err)
		}

		err = os.Rename(f, path.Join(artifacts, file))
		if err != nil {
			fmt.Println(err)
		}
	}
}
