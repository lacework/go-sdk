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

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// Exec executes a command on a remote host
func Exec(host, cmd string) (string, error) {
	user := "root"
	if os.Getenv("LW_SSH_USER") != "" {
		user = os.Getenv("LW_SSH_USER")
	}
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{createKeyRing()},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return executeCmd(cmd, host, config)
}

func createSigner(keyname string) (ssh.Signer, error) {
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

func createKeyRing() ssh.AuthMethod {
	signers := []ssh.Signer{}
	keys := []string{os.Getenv("HOME") + "/.ssh/id_rsa"}

	if os.Getenv("LW_SSH_IDENTITY_FILE") != "" {
		keys = append(keys, os.Getenv("LW_SSH_IDENTITY_FILE"))
	}

	for _, keyname := range keys {
		signer, err := createSigner(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
	}

	return ssh.PublicKeys(signers...)
}

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) (string, error) {
	// @afiune make the port configurable
	conn, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		return "", err
	}

	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var (
		stdoutBuf bytes.Buffer
		stderrBuf bytes.Buffer
	)
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	err = session.Run(cmd)
	if err != nil {
		combinedOutput := ""
		if stdoutBuf.String() != "" {
			combinedOutput = fmt.Sprintf("%s\nSTDOUT:\n%s\n", combinedOutput, stdoutBuf.String())
		}

		if stderrBuf.String() != "" {
			combinedOutput = fmt.Sprintf("%s\nSTDERR:\n%s\n", combinedOutput, stderrBuf.String())
		}
		return "", errors.Wrap(err, combinedOutput)
	}

	return stdoutBuf.String(), nil
}
