//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

package lwcomponent

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// RunAndOutput runs the command and outputs to os.Stdout and os.Stderr,
// the provided environment variables will be accessible by the component
func (c Component) RunAndOutput(args []string, envs ...string) error {
	loc, err := c.Path()
	if err != nil {
		return errors.Wrap(err, baseRunErr)
	}

	cmd := exec.Command(loc, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, envs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return c.run(cmd)
}

// RunAndReturn runs the command and returns its standard output and standard error,
// the provided environment variables will be accessible by the component
func (c Component) RunAndReturn(args []string, stdin io.Reader, envs ...string) (
	stdout string,
	stderr string,
	err error,
) {
	var outBuff, errBuff bytes.Buffer

	loc, err := c.Path()
	if err != nil {
		err = errors.Wrap(err, baseRunErr)
		return
	}

	cmd := exec.Command(loc, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, envs...)
	cmd.Stdin = stdin
	cmd.Stdout = &outBuff
	cmd.Stderr = &errBuff

	err = c.run(cmd)

	stdout, stderr = outBuff.String(), errBuff.String()
	return
}

func (c Component) run(cmd *exec.Cmd) error {
	if c.IsExecutable() {

		// verify component
		if err := c.isVerified(); err != nil {
			return errors.Wrap(err, baseRunErr)
		}

		if err := cmd.Run(); err != nil {
			// default to -1 in case we can't get the actual exit code, which
			// is better than returning an error with exit code 0 (default int)
			exitCode := -1

			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			}

			return &RunError{Err: err, Message: baseRunErr, ExitCode: exitCode}
		}

		return nil
	}

	return errors.Errorf("%s: component %s is not a binary", baseRunErr, c.Name)
}
