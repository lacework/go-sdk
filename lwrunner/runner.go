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
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"

	"github.com/lacework/go-sdk/internal/file"
	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Runner struct {
	Hostname string
	Port     int
	*ssh.ClientConfig
}

func New(user, host string, callback ssh.HostKeyCallback) *Runner {
	if os.Getenv("LW_SSH_USER") != "" {
		user = os.Getenv("LW_SSH_USER")
	}

	defaultCallback, err := DefaultKnownHosts()
	if err == nil && callback == nil {
		callback = defaultCallback
	}

	return &Runner{
		host,
		22,
		&ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{},
			HostKeyCallback: callback,
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

// DefaultKnownHosts returns a host key callback from default known hosts path
func DefaultKnownHosts() (ssh.HostKeyCallback, error) {
	path, err := DefaultKnownHostsPath()
	if err != nil {
		return nil, err
	}

	return knownhosts.New(path)
}

// DefaultKnownHostsPath returns default user ~/.ssh/known_hosts file
func DefaultKnownHostsPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".ssh", "known_hosts"), nil
}

// AddKnownHost adds a host to the provided known hosts file, if no known hosts
// file is provided, it will fallback to default known_hosts file
func AddKnownHost(host string, remote net.Addr, key ssh.PublicKey, knownFile string) (err error) {
	if knownFile == "" {
		path, err := DefaultKnownHostsPath()
		if err != nil {
			return err
		}

		knownFile = path
	}

	if !file.FileExists(knownFile) {
		if err := os.MkdirAll(path.Dir(knownFile), 0700); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(knownFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	var (
		remoteNormalized = knownhosts.Normalize(remote.String())
		hostNormalized   = knownhosts.Normalize(host)
		addresses        = []string{remoteNormalized}
	)

	if hostNormalized != remoteNormalized {
		addresses = append(addresses, hostNormalized)
	}

	_, err = f.WriteString(knownhosts.Line(addresses, key) + "\n")
	return err
}

// CheckKnownHost checks if a host is in known hosts file, if no known hosts
// file is provided, it will fallback to default known_hosts file
func CheckKnownHost(host string, remote net.Addr, key ssh.PublicKey, knownFile string) (found bool, err error) {
	var keyErr *knownhosts.KeyError

	// Fallback to default known_hosts file
	if knownFile == "" {
		path, err := DefaultKnownHostsPath()
		if err != nil {
			return false, err
		}

		knownFile = path
	}

	// get host key callback
	callback, err := knownhosts.New(knownFile)
	if err != nil {
		return false, err
	}

	// check if host already exists
	err = callback(host, remote, key)
	if err == nil {
		// host is known (already exists)
		return true, nil
	}

	// if keyErr.Want is greater than 0 length, that means host is in file with different key
	if errors.As(err, &keyErr) && len(keyErr.Want) > 0 {
		return true, keyErr
	}

	// if not, pass it back to the user
	if err != nil {
		return false, err
	}

	// key is not trusted because it is not in the known hosts file
	return false, nil
}

func DefaultIdentityFilePath() (string, error) {
	if os.Getenv("LW_SSH_IDENTITY_FILE") != "" {
		return os.Getenv("LW_SSH_IDENTITY_FILE"), nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".ssh", "id_rsa"), nil
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
