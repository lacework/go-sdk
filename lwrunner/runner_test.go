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

package lwrunner_test

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"

	"github.com/lacework/go-sdk/lwrunner"
)

func TestLwRunnerNew(t *testing.T) {
	// we use the default know host file inside the HOME directory
	// of the current user, that is why we need to mock it
	mockHome, err := ioutil.TempDir("", "lwrunner")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(mockHome)

	err = os.Mkdir(path.Join(mockHome, ".ssh"), 0755)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(mockHome, ".ssh", "known_hosts"), []byte(""), 0600)
	if err != nil {
		panic(err)
	}
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", mockHome)
	defer os.Setenv("HOME", homeCache)

	subject, err := lwrunner.New("root", "192.1.1.2", nil)
	if assert.Nil(t, err) {
		assert.Equal(t, 22, subject.Port)
		assert.Equal(t, "root", subject.User)
		assert.Equal(t, "192.1.1.2", subject.Hostname)
	}
}

func TestLwRunnerNewIgnoreHostKey(t *testing.T) {
	subject, err := lwrunner.New("ubuntu", "my-test-host", ssh.InsecureIgnoreHostKey())
	if assert.Nil(t, err) {
		assert.Equal(t, 22, subject.Port)
		assert.Equal(t, "ubuntu", subject.User)
		assert.Equal(t, "my-test-host", subject.Hostname)
	}
}

func TestLwRunnerNewCustomCallback(t *testing.T) {
	subject, err := lwrunner.New("ec2-user", "host.example.com", customHostCallback)
	if assert.Nil(t, err) {
		assert.Equal(t, 22, subject.Port)
		assert.Equal(t, subject.User, "ec2-user")
		assert.Equal(t, subject.Hostname, "host.example.com")
	}
}

// test function to mock host callback
func customHostCallback(host string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func TestLwRunnerNewUserEnvVariable(t *testing.T) {
	os.Setenv("LW_SSH_USER", "root")
	defer os.Setenv("LW_SSH_USER", "")

	subject, err := lwrunner.New("ubuntu", "a-test-host", ssh.InsecureIgnoreHostKey())
	if assert.Nil(t, err) {
		assert.Equal(t, subject.User, "root")
		assert.Equal(t, subject.Hostname, "a-test-host")
	}
}

func TestLwRunnerUsePassword(t *testing.T) {
	subject, err := lwrunner.New("ec2-user", "host.example.com", customHostCallback)
	if assert.Nil(t, err) {
		assert.Equal(t, "ec2-user", subject.User)
		assert.Equal(t, "host.example.com:22", subject.Address())
	}

	subject.UsePassword("secret123")

	assert.Equal(t, 1, len(subject.Auth))
}

func TestLwRunnerUseIdentityFile(t *testing.T) {
	subject, err := lwrunner.New("ec2-user", "host.example.com", customHostCallback)
	if assert.Nil(t, err) {
		assert.Equal(t, "ec2-user", subject.User)
		assert.Equal(t, "host.example.com:22", subject.Address())
	}

	err = subject.UseIdentityFile("file-not-found")
	assert.NotNil(t, err)
}

func TestLwRunnerDefaultKnownHostsPath(t *testing.T) {
	subject, err := lwrunner.DefaultKnownHostsPath()
	if assert.Nil(t, err) {
		assert.Contains(t, subject, ".ssh/known_hosts")
	}
}
