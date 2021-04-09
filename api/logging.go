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

package api

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/lacework/go-sdk/lwlogger"
)

// WithLogLevel sets the log level of the client, available: info or debug
func WithLogLevel(level string) Option {
	return clientFunc(func(c *Client) error {
		// do not re initialize our logger if the log level
		// is the same as the desired one
		if level == c.logLevel {
			return nil
		}

		if !lwlogger.ValidLevel(level) {
			return fmt.Errorf("invalid log level '%s'", level)
		}

		c.log.Debug("setting up client", zap.String("log_level", level))
		c.logLevel = level
		c.initLogger()
		return nil
	})
}

// WithLogLevelAndWriter sets the log level of the client
// and writes the log messages to the provided io.Writer
func WithLogLevelAndWriter(level string, w io.Writer) Option {
	return clientFunc(func(c *Client) error {
		if !lwlogger.ValidLevel(level) {
			return fmt.Errorf("invalid log level '%s'", level)
		}

		c.log.Debug("setting up client", zap.String("log_level", level))
		c.logLevel = level
		c.initLoggerWithWriter(w)
		return nil
	})
}

// WithLogWriter configures the client to log messages to the provided io.Writer
func WithLogWriter(w io.Writer) Option {
	return clientFunc(func(c *Client) error {
		c.initLoggerWithWriter(w)
		return nil
	})
}

// WithLogLevelAndFile sets the log level of the client
// and writes the log messages to the provided file
func WithLogLevelAndFile(level, filename string) Option {
	return clientFunc(func(c *Client) error {
		if !lwlogger.ValidLevel(level) {
			return fmt.Errorf("invalid log level '%s'", level)
		}

		c.logLevel = level

		logWriter, err := os.OpenFile(filename, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
		if err != nil {
			return errors.Wrap(err, "unable to open file to initialize api logger ")
		}

		c.initLoggerWithWriter(logWriter)
		return nil
	})
}

// WithLogFile configures the client to write messages to the provided file
func WithLogFile(filename string) Option {
	return clientFunc(func(c *Client) error {
		logWriter, err := os.OpenFile(filename, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
		if err != nil {
			return errors.Wrap(err, "unable to open file to initialize api logger ")
		}

		c.log.Debug("setting up client redirect logger", zap.String("file", filename))
		c.initLoggerWithWriter(logWriter)
		return nil
	})
}

// initLogger initializes the logger with a set of default fields
func (c *Client) initLogger() {
	if c.log != nil {
		_ = c.log.Sync()
	}
	c.log = lwlogger.New(c.logLevel,
		zap.Fields(
			zap.Field(zap.String("id", c.id)),
			zap.Field(zap.String("account", c.account)),
		),
	)

	// verify if the log level has been configure through environment variable
	if envLevel := lwlogger.LogLevelFromEnvironment(); envLevel != "" {
		c.log.Debug("setting up client, override log level",
			zap.String("before", c.logLevel),
			zap.String("after", envLevel),
		)
		c.logLevel = envLevel
	}
}

// initLoggerWithWriter initializes a new logger with a set
// of default fields and configues the provided io.Writer
func (c *Client) initLoggerWithWriter(w io.Writer) {
	if c.log != nil {
		_ = c.log.Sync()
	}
	c.log = lwlogger.NewWithWriter(c.logLevel, w,
		zap.Fields(
			zap.Field(zap.String("id", c.id)),
			zap.Field(zap.String("account", c.account)),
		),
	)

	// verify if the log level has been configure through environment variable
	if envLevel := lwlogger.LogLevelFromEnvironment(); envLevel != "" {
		c.log.Debug("setting up client, override log level",
			zap.String("before", c.logLevel),
			zap.String("after", envLevel),
		)
		c.logLevel = envLevel
	}
}

// debugMode returns true if the client is configured to display debug level logs
func (c *Client) debugMode() bool {
	return c.logLevel == "DEBUG"
}
