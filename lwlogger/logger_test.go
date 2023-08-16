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

package lwlogger_test

import (
	"bytes"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/lacework/go-sdk/internal/capturer"
	"github.com/lacework/go-sdk/lwlogger"
)

func TestLoggerNew(t *testing.T) {
	// by default we do not log an INFO or DEBUG, only ERROR
	logOutput := capturer.CaptureOutput(func() {
		lwL := lwlogger.New("")
		lwL.Info("interesting info")
		lwL.Debug("we are debugging")
		lwL.Error("we have errors")
	})

	// we are asserting a log message similar to:
	//
	// {
	//   "level":      "error",
	//   "msg":        "we have errors",
	//   "caller":     "lwlogger/logger_test.go:37",
	//   "ts":         "2020-04-03T18:22:17-06:00",
	//   "stacktrace": "<...long long stacktrace...>"
	// }
	assert.Contains(t, logOutput, "\"level\"")
	assert.Contains(t, logOutput, "\"error\"")
	assert.Contains(t, logOutput, "\"caller\"")
	assert.Contains(t, logOutput, "lwlogger/logger_test.go:")
	assert.Contains(t, logOutput, "\"msg\"")
	assert.Contains(t, logOutput, "we have errors")
	assert.Contains(t, logOutput, "\"stacktrace\"")
	assert.Contains(t, logOutput, "\"ts\"")

	assert.NotContains(t, logOutput, "we are debugging",
		"we are in ERROR mode, no DEBUG messages should be displayed",
	)
	assert.NotContains(t, logOutput, "interesting info",
		"we are in ERROR mode, no INFO messages should be displayed",
	)
}

func TestLoggerNewINFO(t *testing.T) {
	// switching to INFO will show ERRORS and INFO logs
	logOutput := capturer.CaptureOutput(func() {
		lwL := lwlogger.New("INFO")
		lwL.Info("interesting info")
		lwL.Debug("we are debugging")
		lwL.Error("we have errors")
	})

	// we are asserting a log messages similar to:
	//
	// {
	//   "level":      "error",
	//   "msg":        "we have errors",
	//   "caller":     "lwlogger/logger_test.go:37",
	//   "ts":         "2020-04-03T18:22:17-06:00",
	//   "stacktrace": "<...long long stacktrace...>"
	// }
	// {
	//   "level":      "info",
	//   "msg":        "interesting info",
	//   "caller":     "lwlogger/logger_test.go:74",
	//   "ts":         "2020-04-03T18:22:17-06:00",
	//   "stacktrace": "<...long long stacktrace...>"
	// }
	assert.Contains(t, logOutput, "\"level\"")
	assert.Contains(t, logOutput, "\"error\"")
	assert.Contains(t, logOutput, "\"info\"")
	assert.Contains(t, logOutput, "\"caller\"")
	assert.Contains(t, logOutput, "lwlogger/logger_test.go:")
	assert.Contains(t, logOutput, "\"msg\"")
	assert.Contains(t, logOutput, "we have errors")   // ERROR
	assert.Contains(t, logOutput, "interesting info") // INFO
	assert.Contains(t, logOutput, "\"stacktrace\"")
	assert.Contains(t, logOutput, "\"ts\"")

	assert.NotContains(t, logOutput, "we are debugging",
		"we are in INFO mode, no DEBUG messages should be displayed",
	)
}

func TestLoggerNewDEBUG(t *testing.T) {
	// switching to DEBUG will show INFO, DEBUG and ERRORS messages
	logOutput := capturer.CaptureOutput(func() {
		lwL := lwlogger.New("DEBUG")
		lwL.Info("interesting info")
		lwL.Debug("we are debugging")
		lwL.Error("we have errors")
	})

	// we are asserting a log message similar to:
	//
	// {
	//   "level":      "error",
	//   "msg":        "we have errors",
	//   "caller":     "lwlogger/logger_test.go:37",
	//   "ts":         "2020-04-03T18:22:17-06:00",
	//   "stacktrace": "<...long long stacktrace...>"
	// }
	// {
	//   "level":      "info",
	//   "msg":        "interesting info",
	//   "caller":     "lwlogger/logger_test.go:105",
	//   "ts":         "2020-04-03T18:22:17-06:00",
	//   "stacktrace": "<...long long stacktrace...>"
	// }
	// {
	//   "level":      "debug",
	//   "msg":        "we are debugging",
	//   "caller":     "lwlogger/logger_test.go:106",
	//   "ts":         "2020-04-03T18:22:17-06:00",
	//   "stacktrace": "<...long long stacktrace...>"
	// }
	assert.Contains(t, logOutput, "\"level\"")
	assert.Contains(t, logOutput, "\"error\"")
	assert.Contains(t, logOutput, "\"info\"")
	assert.Contains(t, logOutput, "\"debug\"")
	assert.Contains(t, logOutput, "\"caller\"")
	assert.Contains(t, logOutput, "lwlogger/logger_test.go:")
	assert.Contains(t, logOutput, "\"msg\"")
	assert.Contains(t, logOutput, "we have errors")   // ERROR
	assert.Contains(t, logOutput, "interesting info") // INFO
	assert.Contains(t, logOutput, "we are debugging") // DEBUG
	assert.Contains(t, logOutput, "\"stacktrace\"")
	assert.Contains(t, logOutput, "\"ts\"")
}

func TestLoggerNewLogEnv(t *testing.T) {
	// we tell the logger to start with ERROR log but we pass
	// DEBUG via the environment variable
	os.Setenv(lwlogger.LogLevelEnv, "DEBUG")
	defer os.Setenv(lwlogger.LogLevelEnv, "")

	logOutput := capturer.CaptureOutput(func() {
		lwL := lwlogger.New("")
		lwL.Info("interesting info")
		lwL.Debug("we are debugging")
		lwL.Error("we have errors")
	})

	// we should see all log messages
	assert.Contains(t, logOutput, "we have errors")   // ERROR
	assert.Contains(t, logOutput, "interesting info") // INFO
	assert.Contains(t, logOutput, "we are debugging") // DEBUG
}

func TestLoggerNewLogFormatEnv(t *testing.T) {
	// we tell the logger to start with ERROR log but we pass
	// DEBUG via the environment variable
	os.Setenv(lwlogger.LogFormatEnv, "CONSOLE")
	defer os.Setenv(lwlogger.LogFormatEnv, "")

	logOutput := capturer.CaptureOutput(func() {
		lwL := lwlogger.New("DEBUG")
		lwL.Info("interesting info")
		lwL.Debug("we are debugging")
	})

	assert.Contains(t, logOutput, "interesting info") // INFO
	assert.Contains(t, logOutput, "we are debugging") // DEBUG
}

func TestLoggerNewWithWriter(t *testing.T) {
	// create a temporal file to write the logs
	tmpfile, err := os.CreateTemp("", "logger")
	assert.Nil(t, err)
	defer os.Remove(tmpfile.Name())
	logOutput, err := os.OpenFile(tmpfile.Name(), syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
	assert.Nil(t, err)

	// create a new logger with a io.Writer from the temporal file
	lwlogger.NewWithWriter("DEBUG", logOutput).Debug("we are debugging")

	// after writing the log message, read the file and assert its content
	logContentB, err := os.ReadFile(tmpfile.Name())
	assert.Nil(t, err)
	logContent := string(logContentB)
	assert.Contains(t, logContent, "we are debugging")
}

func TestLoggerNewWithOptions(t *testing.T) {
	logOutput := capturer.CaptureOutput(func() {
		lwL := lwlogger.New("DEBUG",
			zap.Fields(zap.Field(zap.String("my_field", "awesome"))),
		)
		lwL.Debug("we are debugging")
	})

	assert.Contains(t, logOutput, "we are debugging")
	assert.Contains(t, logOutput, "\"my_field\"")
	assert.Contains(t, logOutput, "\"awesome\"")
}

func TestLoggerMerge(t *testing.T) {
	var bufOne bytes.Buffer
	var bufTwo bytes.Buffer

	logOne := lwlogger.NewWithWriter("INFO", &bufOne)
	logTwo := lwlogger.NewWithWriter("DEBUG", &bufTwo)

	mergedLog := lwlogger.Merge(logOne, logTwo)
	mergedLog.Info("ABCD")
	mergedLog.Debug("XYZ")
	_ = mergedLog.Sync()

	assert.Contains(t, bufOne.String(), "ABCD")
	assert.NotContains(t, bufOne.String(), "XYZ")
	assert.Contains(t, bufTwo.String(), "ABCD")
	assert.Contains(t, bufTwo.String(), "XYZ")
}

func TestValidLevel(t *testing.T) {
	assert.True(t, lwlogger.ValidLevel("INFO"))
	assert.True(t, lwlogger.ValidLevel("DEBUG"))
	assert.True(t, lwlogger.ValidLevel(""))
	assert.False(t, lwlogger.ValidLevel("FOO"))
}
