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

package api_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestNewClientWithLogLevel(t *testing.T) {
	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	logOutput := captureOutput(func() {
		c, err := api.NewClient("test",
			api.WithURL(fakeServer.URL()),
			api.WithLogLevel("debug"),
		)
		if assert.Nil(t, err) {
			assert.Equal(t, "v1", c.ApiVersion(), "modified API version should be v2")
		}
	})

	// we are asserting a log message similar to:
	//
	// [timestamp] DEBUG api/client.go:92 api client created {
	//    "url": "http://127.0.0.1:55317",
	//    "log_level": "debug",
	//    "timeout": 1800
	// }
	assert.Contains(t, logOutput, "DEBUG")
	assert.Contains(t, logOutput, "api/client.go")
	assert.Contains(t, logOutput, "api client created")
	assert.Contains(t, logOutput, "\"url\"")
	assert.Contains(t, logOutput, "http://127.0.0.1:")
	assert.Contains(t, logOutput, "log_level")
	assert.Contains(t, logOutput, "debug")
	assert.Contains(t, logOutput, "timeout")
	assert.Contains(t, logOutput, "3600")
}

// captureOutput executes a function and captures the STDOUT and STDERR,
// useful to test logging messages
func captureOutput(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	stderr := os.Stderr
	os.Stderr = w
	defer func() {
		os.Stderr = stderr
	}()

	f()
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}
