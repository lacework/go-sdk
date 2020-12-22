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

// A runner package that executes commands on remote hosts.
package lwrunner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Runner struct {
	Hostname string
	Port     int
	*ssh.ClientConfig
}

func New(user, host string) *Runner {
	// @afiune notify the user?
	hostKeyCallback := ssh.InsecureIgnoreHostKey()

	// try to use the known_hosts file as a host key callback
	home, err := homedir.Dir()
	if err == nil {
		hostKeyFromKnownHosts, err := knownhosts.New(path.Join(home, ".ssh", "known_hosts"))
		if err == nil {
			hostKeyCallback = hostKeyFromKnownHosts
		}
	}

	if os.Getenv("LW_SSH_USER") != "" {
		user = os.Getenv("LW_SSH_USER")
	}

	return &Runner{
		host,
		22,
		&ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{defaultAuthMethod()},
			HostKeyCallback: hostKeyCallback,
		},
	}
}

func (run Runner) UseIdentityFile(file string) error {
	signer, err := newSignerFromFile(file)
	if err != nil {
		return err
	}
	run.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	return nil
}

func (run Runner) UsePassword(secret string) {
	run.Auth = []ssh.AuthMethod{ssh.Password(secret)}
}

func (run *Runner) Address() string {
	return fmt.Sprintf("%s:%d", run.Hostname, run.Port)
}

// Exec executes a command on the configured remote host
func (run *Runner) Exec(cmd string) (stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	// @afiune make the port configurable
	conn, err := ssh.Dial("tcp", run.Address(), run.ClientConfig)
	if err != nil {
		return
	}

	session, err := conn.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(cmd)
	return
}

func newSignerFromFile(keyname string) (ssh.Signer, error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(buf)
}

func defaultAuthMethod() ssh.AuthMethod {
	var (
		signers = []ssh.Signer{}
		keys    = []string{}
	)
	home, err := homedir.Dir()
	if err == nil {
		keys = append(keys, path.Join(home, ".ssh", "id_rsa"))
	}

	if os.Getenv("LW_SSH_IDENTITY_FILE") != "" {
		keys = append(keys, os.Getenv("LW_SSH_IDENTITY_FILE"))
	}

	for _, keyname := range keys {
		signer, err := newSignerFromFile(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
	}

	return ssh.PublicKeys(signers...)
}
