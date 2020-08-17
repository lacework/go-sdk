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

// Go package that allows a program to easily pipe it's
// standard output through a pager program
package lwpager

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

type Pager struct {
	Cmd *exec.Cmd
	Out io.WriteCloser
}

// The environment variables to check for the name of (and arguments to)
// the pager to run
var PagerEnv = []string{"LW_PAGER", "PAGER"}

// List of pager commands to look for inside the $PATH environment variables
var PagerCmds = []string{"less", "more"}

func pagerExecPath() (path string, args []string, err error) {
	for _, testVar := range PagerEnv {
		path = os.Getenv(testVar)
		if path != "" {
			args = strings.Fields(path)
			return args[0], args[1:], nil
		}
	}

	// by default look for pager commands only if PagerCmds is empty
	err = exec.ErrNotFound
	for _, testPath := range PagerCmds {
		path, err = exec.LookPath(testPath)
		if err == nil {
			return path, nil, nil
		}
	}
	return "", nil, err
}

// New returns a new io.WriteCloser connected to a pager.
// The returned out can be used as a replacement to os.Stdout,
// everything written to it is piped to a pager.
// To determine what pager to run, the environment variables listed
// in PagerEnv are checked.
// If all are empty/unset then the commands listed in PagerCmds
// are looked for in $PATH.
func New() (*Pager, error) {
	pager := new(Pager)
	path, args, err := pagerExecPath()
	if err != nil {
		return nil, err
	}

	pager.Cmd = exec.Command(path, args...)
	pager.Cmd.Stdout = os.Stdout
	pager.Cmd.Stderr = os.Stderr
	w, err := pager.Cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	pager.Out = w
	return pager, nil
}

func Start() (*Pager, error) {
	pager, err := New()
	if err != nil {
		return nil, err
	}

	if err = pager.Cmd.Start(); err != nil {
		return nil, err
	}

	return pager, nil
}

// Stdout sets the global variable os.Stdout to the result of New()
// and returns the old os.Stdout value.
//func Stdout() (io.WriteCloser, error) {
//p, err := New()
//if err != nil {
//return nil, err
//}
//io.Copy(os.Stdout, p)
//return p, nil
//}

// Wait closes the pipe to the pager setup with New() or Stdout() and waits
// for it to exit.
//
// This should normally be called before the program exists,
// typically via a defer call in main().
func (p *Pager) Wait() {
	if p.Cmd == nil {
		return
	}
	p.Out.Close()
	_ = p.Cmd.Wait()
}
