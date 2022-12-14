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

// A Logger wrapper for Lacework based of zap logger Go package.
package lwlogger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// LogLevelEnv represents the level that the logger is configured
	LogLevelEnv        = "LW_LOG"
	SupportedLogLevels = [4]string{"", "ERROR", "INFO", "DEBUG"}

	// LogFormatEnv controls the format of the logger
	LogFormatEnv        = "LW_LOG_FORMAT"
	DefaultLogFormat    = "JSON"
	SupportedLogFormats = [2]string{"JSON", "CONSOLE"}

	// LogDevelopmentModeEnv switches the logger to development mode
	LogDevelopmentModeEnv = "LW_LOG_DEV"

	// LogToNativeLoggerEnv is used for those consumers like terraform that control
	// the logs that are presented to the user, when this environment is turned
	// on, the logger implementation will use the native Go logger 'log.Writer()'
	LogToNativeLoggerEnv = "LW_LOG_NATIVE"
)

// New initialize a new logger with the provided level and options
func New(level string, options ...zap.Option) *zap.Logger {
	// give priority to the environment variable
	if envLevel := LogLevelFromEnvironment(); envLevel != "" {
		level = envLevel
	}

	zapConfig := zap.Config{
		Level: zapLogLevel(level),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Development:      inDevelopmentMode(),
		Encoding:         logFormatFromEnv(),
		EncoderConfig:    laceworkEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := zapConfig.Build(options...)
	if err != nil {
		fmt.Printf("Error: unable to initialize logger: %v\n", err)
		return zap.NewExample(options...)
	}

	return l
}

// NewWithWriter initialize a new logger with the provided level and options
// but redirecting the logs to the provider io.Writer
func NewWithWriter(level string, out io.Writer, options ...zap.Option) *zap.Logger {
	// give priority to the environment variable
	if envLevel := LogLevelFromEnvironment(); envLevel != "" {
		level = envLevel
	}

	var (
		writeSyncer = zapcore.AddSync(out)
		core        = zapcore.NewCore(
			zapEncoderFromFormat(logFormatFromEnv()),
			writeSyncer,
			zapLogLevel(level),
		)
		localOpts = []zap.Option{
			zap.ErrorOutput(writeSyncer),
			zap.AddCaller(),
			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)
			}),
		}
	)

	return zap.New(core, options...).WithOptions(localOpts...)
}

func ValidLevel(level string) bool {
	for _, l := range SupportedLogLevels {
		if l == level {
			return true
		}
	}
	return false
}

// LogLevelFromEnvironment checks the environment variable 'LW_LOG'
func LogLevelFromEnvironment() string {
	switch os.Getenv(LogLevelEnv) {
	case "info", "INFO":
		return "INFO"
	case "debug", "DEBUG":
		return "DEBUG"
	case "error", "ERROR":
		return "ERROR"
	default:
		return ""
	}
}

func zapLogLevel(level string) zap.AtomicLevel {
	switch level {
	case "INFO":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "DEBUG":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	}
}

func inDevelopmentMode() bool {
	return os.Getenv(LogDevelopmentModeEnv) == "true"
}

func logFormatFromEnv() string {
	switch os.Getenv(LogFormatEnv) {
	case "console", "CONSOLE":
		return "console"
	case "json", "JSON":
		return "json"
	}
	// @afiune the library require the format to be lowercase
	return strings.ToLower(DefaultLogFormat)
}

func zapEncoderFromFormat(format string) zapcore.Encoder {
	switch format {
	case "console":
		return zapcore.NewConsoleEncoder(laceworkEncoderConfig())
	case "json":
		return zapcore.NewJSONEncoder(laceworkEncoderConfig())
	default:
		// @afiune we should never land here but just in case ;)
		return zapcore.NewJSONEncoder(laceworkEncoderConfig())
	}
}

func laceworkEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
