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
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"
	"net"
	"os"
	"path"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"

	"github.com/lacework/go-sdk/lwrunner"
)

func TestLwRunnerNew(t *testing.T) {
	subject := lwrunner.New("root", "192.1.1.2", nil)
	assert.Equal(t, 22, subject.Port)
	assert.Equal(t, "root", subject.User)
	assert.Equal(t, "192.1.1.2", subject.Hostname)
}

func TestLwRunnerNewIgnoreHostKey(t *testing.T) {
	subject := lwrunner.New("ubuntu", "my-test-host", ssh.InsecureIgnoreHostKey())
	assert.Equal(t, 22, subject.Port)
	assert.Equal(t, "ubuntu", subject.User)
	assert.Equal(t, "my-test-host", subject.Hostname)
}

func TestLwRunnerNewCustomCallback(t *testing.T) {
	subject := lwrunner.New("ec2-user", "host.example.com", customHostCallback)
	assert.Equal(t, 22, subject.Port)
	assert.Equal(t, subject.User, "ec2-user")
	assert.Equal(t, subject.Hostname, "host.example.com")
}

// test function to mock host callback
func customHostCallback(host string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func TestLwRunnerNewUserEnvVariable(t *testing.T) {
	os.Setenv("LW_SSH_USER", "root")
	defer os.Setenv("LW_SSH_USER", "")

	subject := lwrunner.New("ubuntu", "a-test-host", ssh.InsecureIgnoreHostKey())
	assert.Equal(t, subject.User, "root")
	assert.Equal(t, subject.Hostname, "a-test-host")
}

func TestLwRunnerUsePassword(t *testing.T) {
	subject := lwrunner.New("ec2-user", "host.example.com", customHostCallback)
	assert.Equal(t, "ec2-user", subject.User)
	assert.Equal(t, "host.example.com:22", subject.Address())

	subject.UsePassword("secret123")
	assert.Equal(t, 1, len(subject.Auth))
}

func TestLwRunnerUseIdentityFile(t *testing.T) {
	subject := lwrunner.New("ec2-user", "host.example.com", customHostCallback)
	assert.Equal(t, "ec2-user", subject.User)
	assert.Equal(t, "host.example.com:22", subject.Address())

	err := subject.UseIdentityFile("file-not-found")
	assert.NotNil(t, err)
}

func TestLwRunnerDefaultKnownHostsPath(t *testing.T) {
	subject, err := lwrunner.DefaultKnownHostsPath()
	if assert.Nil(t, err) {
		assert.Contains(t, subject, ".ssh/known_hosts")
	}
}

func TestDefaultIdentityFilePath(t *testing.T) {
	subject, err := lwrunner.DefaultIdentityFilePath()
	if assert.Nil(t, err) {
		assert.Contains(t, subject, ".ssh")
		assert.Contains(t, subject, "id_rsa")
	}
}

func TestDefaultIdentityFilePathEnvVariable(t *testing.T) {
	expected := "/pat/to/key"
	os.Setenv("LW_SSH_IDENTITY_FILE", expected)
	defer os.Setenv("LW_SSH_IDENTITY_FILE", "")

	subject, err := lwrunner.DefaultIdentityFilePath()
	if assert.Nil(t, err) {
		assert.Equal(t, subject, expected)
	}
}

func TestLwRunnerAddKnownHostNoSSHDir(t *testing.T) {
	mockHome, err := ioutil.TempDir("", "lwrunner")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(mockHome)

	knownFile := path.Join(mockHome, ".ssh", "known_hosts")
	netAddr := mockNetAddr{}
	// generate test RSA keypair in SSH format
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	rsaPub := priv.PublicKey
	sshPub, err := ssh.NewPublicKey(&rsaPub)
	assert.NoError(t, err)

	// Add known host to mocked home directory
	subject := lwrunner.AddKnownHost("mock-test", netAddr, sshPub, knownFile)
	assert.NoError(t, subject)

	// Check the known host file
	content, err := ioutil.ReadFile(knownFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "mock-test")

	// Try again, it should work
	// Add known host to mocked home directory
	subject = lwrunner.AddKnownHost("second-time", netAddr, sshPub, knownFile)
	assert.NoError(t, subject)

	// Check the known host file
	content, err = ioutil.ReadFile(knownFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "second-time")
}

func TestLwRunnerAddKnownWithSSHDir(t *testing.T) {
	mockHome, err := ioutil.TempDir("", "lwrunner")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(mockHome)

	// Mock that the ~/.ssh dir exists
	err = os.Mkdir(path.Join(mockHome, ".ssh"), 0700)
	if err != nil {
		panic(err)
	}

	knownFile := path.Join(mockHome, ".ssh", "known_hosts")
	netAddr := mockNetAddr{}
	// generate test RSA keypair in SSH format
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	rsaPub := priv.PublicKey
	sshPub, err := ssh.NewPublicKey(&rsaPub)
	assert.NoError(t, err)

	// Add known host to mocked home directory
	subject := lwrunner.AddKnownHost("mock-test", netAddr, sshPub, knownFile)
	assert.NoError(t, subject)

	// Check the known host file
	content, err := ioutil.ReadFile(knownFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "mock-test")
}

type mockNetAddr struct{}

func (m mockNetAddr) Network() string {
	return "tcp"
}
func (m mockNetAddr) String() string {
	return "1.1.1.1"
}

func TestLwRunnerDefaultKnownHosts(t *testing.T) {
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

	homedir.DisableCache = true
	subject, err := lwrunner.DefaultKnownHosts()
	assert.NotNil(t, subject)
	if assert.Nil(t, err) {
		assert.NotNil(t, subject("mock.hostname.example.com:22", mockAddr{}, mockPublicKey{}))
	}
}

type mockAddr struct{}

func (m mockAddr) Network() string {
	return "tcp"
}
func (m mockAddr) String() string {
	return "mock.hostname.example.com:22"
}

type mockPublicKey struct{}

func (m mockPublicKey) Type() string {
	return "ssh-rsa"
}
func (m mockPublicKey) Marshal() []byte {
	return []byte{}
}
func (m mockPublicKey) Verify(_ []byte, _ *ssh.Signature) error {
	return nil
}
