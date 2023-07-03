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
	"log"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/capturer"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestNewClientWithLogLevel(t *testing.T) {
	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	logOutput := capturer.CaptureOutput(func() {
		c, err := api.NewClient("test",
			api.WithURL(fakeServer.URL()),
			api.WithLogLevel("INFO"),
		)
		if assert.Nil(t, err) {
			assert.Equal(t, "v2", c.ApiVersion(), "API version should be v2")
		}
	})

	testNewClientLogOutput(t, logOutput)
	assert.Contains(t, logOutput, "INFO")
}

func TestClientWithLogLevelAndWriter(t *testing.T) {
	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	// there could be cases where cunsumers, like terraform, control
	// the logs that are presented to the user. Here is a test that
	// will use the native Go logger to write the API logs to its Writer
	//
	// configuring the native Go logger to use a temporal file
	tmpfile := configureNativeGoLoggerAsConsumers(t)
	defer os.Remove(tmpfile)

	// use the configured native Go logger as the io.Writer
	c, err := api.NewClient("test",
		api.WithURL(fakeServer.URL()),
		api.WithToken("TOKEN"),
		api.WithLogLevelAndWriter("DEBUG", log.Writer()),
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "v2", c.ApiVersion(), "API version should be v2")
	}

	// generating a DEBUG log by creating a new request
	_, err = c.NewRequest("GET", "foo", nil)
	assert.Nil(t, err)

	logContentB, err := os.ReadFile(tmpfile)
	assert.Nil(t, err)
	logContent := string(logContentB)

	// this tests the info log message from the new client
	testNewClientLogOutput(t, logContent)

	// we are asserting a log message similar to:
	// {
	//   "level":"debug",
	//   "ts":"2020-04-04T02:31:55-06:00",
	//   "caller":"api/http.go:78",
	//   "msg":"request",
	//   "id":"323870a20f48f018",
	//   "account":"test",
	//   "method":"GET",
	//   "url":"http://127.0.0.1:58753",
	//   "endpoint":"/api/v2/foo",
	//   "headers":{
	//     "Accept":"application/json",
	//     "Authorization":"TOKEN",
	//     "Method":"GET"
	//   },
	//   "body":""
	// }
	assert.Contains(t, logContent, "debug")
	assert.Contains(t, logContent, "DEBUG")
	assert.Contains(t, logContent, "api/client.go")
	assert.Contains(t, logContent, "\"request\"")
	assert.Contains(t, logContent, "\"endpoint\"")
	assert.Contains(t, logContent, "/api/v2/foo")
	assert.Contains(t, logContent, "\"method\"")
	assert.Contains(t, logContent, "GET")
	assert.Contains(t, logContent, "\"id\"")
	assert.Contains(t, logContent, "\"ts\"")
	assert.Contains(t, logContent, "\"headers\"")
	assert.Contains(t, logContent, "\"body\"")
	assert.Contains(t, logContent, "\"Authorization\"")
	assert.Contains(t, logContent, "TOKEN")
}

func TestClientWithLogLevelAndFile(t *testing.T) {
	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	tmpfile, err := os.CreateTemp("", "logger")
	assert.Nil(t, err)
	defer os.Remove(tmpfile.Name())

	c, err := api.NewClient("test",
		api.WithURL(fakeServer.URL()),
		api.WithLogLevelAndFile("INFO", tmpfile.Name()),
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "v2", c.ApiVersion(), "API version should be v2")
	}

	logContentB, err := os.ReadFile(tmpfile.Name())
	assert.Nil(t, err)
	logContent := string(logContentB)

	// this tests the info log message from the new client
	testNewClientLogOutput(t, logContent)
	assert.Contains(t, logContent, "INFO")
}

func testNewClientLogOutput(t *testing.T, out string) {
	// we are asserting a log message similar to:
	//
	// {
	//   "level":     "info",
	//   "ts":        "2020-04-03T17:12:37-06:00",
	//   "caller":    "api/client.go:98",
	//   "msg":       "api client created",
	//   "id":        "12ba597e8b7b2379",
	//   "account":   "test",
	//   "url":       "http://127.0.0.1:52544",
	//   "version":   "v2",
	//   "log_level": "INFO",
	//   "timeout":   3600
	// }
	assert.Contains(t, out, "\"level\"")
	assert.Contains(t, out, "info")
	assert.Contains(t, out, "\"caller\"")
	assert.Contains(t, out, "api/client.go")
	assert.Contains(t, out, "\"msg\"")
	assert.Contains(t, out, "api client created")
	assert.Contains(t, out, "\"url\"")
	assert.Contains(t, out, "http://127.0.0.1:")
	assert.Contains(t, out, "\"timeout\"")
	assert.Contains(t, out, "3600")
	assert.Contains(t, out, "\"account\"")
	assert.Contains(t, out, "test")
	assert.Contains(t, out, "\"version\"")
	assert.Contains(t, out, "v2")
	assert.Contains(t, out, "\"log_level\"")
	assert.Contains(t, out, "\"id\"")
	assert.Contains(t, out, "\"ts\"")
}

func configureNativeGoLoggerAsConsumers(t *testing.T) string {
	tmpfile, err := os.CreateTemp("", "logger")
	assert.Nil(t, err)
	logOutput, err := os.OpenFile(tmpfile.Name(), syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
	assert.Nil(t, err)
	log.SetOutput(logOutput)
	return tmpfile.Name()
}
