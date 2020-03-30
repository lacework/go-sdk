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
	"os"

	"go.uber.org/zap"
)

// WithLogLevel sets the log level of the client, available: info or debug
func WithLogLevel(level string) Option {
	return clientFunc(func(c *Client) error {
		switch level {
		case "info", "debug":
			c.logLevel = level
		default:
			c.logLevel = "info"
		}

		c.initializeLogger()
		return nil
	})
}

// initializeLogger initializes the logger, by default we assume production,
// but if debug mode is turned on, we switch to development
func (c *Client) initializeLogger() {
	// give priority to the environment variable: LW_DEBUG
	debug := os.Getenv("LW_DEBUG")
	if debug == "true" {
		c.logLevel = "debug"
	} else if debug == "false" {
		c.logLevel = "info"
	}

	var err error
	if c.logLevel == "debug" {
		c.log, err = zap.NewDevelopment()
	} else {
		c.log, err = zap.NewProduction()
	}

	// if we find any error initializing zap, default to a standard logger
	if err != nil {
		fmt.Printf("Error: unable to initialize zap logger: %v\n", err)
		c.log = zap.NewExample()
	}
}
