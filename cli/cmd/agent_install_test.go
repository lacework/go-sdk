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

package cmd

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLatestAgentInstallDownloadUrl(t *testing.T) {
	downloadUrl, err := latestAgentInstallDownloadUrl()
	if assert.Nil(t, err) {
		assert.Contains(t, downloadUrl, "https://s3-us-west-2.amazonaws.com/www.lacework.net/download")
		assert.Contains(t, downloadUrl, "install.sh")
	}
}

func TestLatestAgentVersionSHA(t *testing.T) {
	sha, err := latestAgentVersionSHA()
	if assert.Nil(t, err) {
		// Example: 3.3.5_2020-11-16_master_ac0e65055f11f4f59bab6ea4dfa61dcafaa9a3f1
		assert.Regexpf(t,
			regexp.MustCompile("[0-9]*\\.[0-9]*\\.[0-9]*_[0-9]*-[0-9]*-[0-9]*_\\w*_*"), sha,
			"agent version SHA doesn't match regex AGENT_VERSION_YYYY-MM-DD_BRANCHNAME_GITHASH")
	}
}

func TestDefaultIdentityFile(t *testing.T) {
	subject, err := defaultIdentityFile()
	if assert.Nil(t, err) {
		assert.Contains(t, subject, ".ssh")
		assert.Contains(t, subject, "id_rsa")
	}
}

func TestFormatRunnerError(t *testing.T) {
	cases := []struct {
		expected error
		stdout   *bytes.Buffer
		stderr   *bytes.Buffer
		err      error
	}{
		{expected: nil,
			stdout: bytes.NewBufferString(""),
			stderr: bytes.NewBufferString(""),
			err:    nil},
		{expected: errors.New("something happened without stdout and stderr"),
			stdout: bytes.NewBufferString(""),
			stderr: bytes.NewBufferString(""),
			err:    errors.New("something happened without stdout and stderr")},
		{expected: errors.New("\n\nSTDOUT:\nonly something in stdout"),
			stdout: bytes.NewBufferString("only something in stdout"),
			stderr: bytes.NewBufferString(""),
			err:    nil},
		{expected: errors.New("\n\nSTDERR:\nonly something in stderr"),
			stdout: bytes.NewBufferString(""),
			stderr: bytes.NewBufferString("only something in stderr"),
			err:    nil},
		{expected: errors.New("\n\nSTDOUT:\nsomething in stdout\n\nSTDERR:\nand something in stderr"),
			stdout: bytes.NewBufferString("something in stdout"),
			stderr: bytes.NewBufferString("and something in stderr"),
			err:    nil},
		{expected: errors.New("\n\nSTDOUT:\nsomething here\n\nSTDERR:\nand something here: and here"),
			stdout: bytes.NewBufferString("something here"),
			stderr: bytes.NewBufferString("and something here"),
			err:    errors.New("and here")},
	}
	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			subject := formatRunnerError(*kase.stdout, *kase.stderr, kase.err)
			if kase.expected == nil {
				assert.Nil(t, subject)
			} else {
				assert.Equal(t, kase.expected.Error(), subject.Error())
			}
		})
	}
}
