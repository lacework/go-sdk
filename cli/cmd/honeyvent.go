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
	"os"
	"runtime"
	"sync"

	"github.com/honeycombio/libhoney-go"
)

var (
	// HoneyApiKey is a variable that is injected at build time via
	// the cross-platform directive inside the Makefile, this key is
	// used to send events to Honeycomb so that we can understand how
	// our customers use the Lacework CLI
	HoneyApiKey = "unknown"

	// HoneyDataset is the dataset in Honeycomb that we send tracing
	// data this variable will be set depending on the environment we
	// are running on. During development, we send all events and
	// tracing data to a default dataset.
	HoneyDataset = "lacework-cli-dev"
)

const (
	// DisableTelemetry is an environment variable that can be used to
	// disable telemetry sent to Honeycomb
	DisableTelemetry = "LW_TELEMETRY_DISABLE"

	// List of Features
	//
	// A feature within the Lacework CLI is any functionality that
	// can't be traced or tracked by the default event sent to Honeycomb,
	// it is a behavior that we, Lacework engineers, would like to
	// trace and understand its usage and adoption.
	//
	// By default the Feature field within the Honeyvent is empty,
	// define a new feature below and set it before sending a new
	// Honeyvent. Additionally, there is a FeatureData field that
	// any feature can use to inject any specific information
	// related to that feature.
	//
	// Example:
	//
	// ```go
	// cli.Event.Feature = featPollCtrScan
	// cli.Event.FeatureData = map[string]interface{"key", "value"}
	// cli.SendHoneyvent()
	// ```
	//
	// Polling mechanism feature
	featPollCtrScan = "poll_ctr_scan"

	// Daily version check feature
	featDailyVerCheck = "daily_check"
)

// Honeyvent defines what a Honeycomb event looks like for the Lacework CLI
type Honeyvent struct {
	Version     string      `json:"version"`
	Os          string      `json:"os"`
	Arch        string      `json:"arch"`
	Command     string      `json:"command,omitempty"`
	Args        []string    `json:"args,omitempty"`
	Account     string      `json:"account,omitempty"`
	Profile     string      `json:"profile,omitempty"`
	ApiKey      string      `json:"api_key,omitempty"`
	Feature     string      `json:"feature,omitempty"`
	FeatureData interface{} `json:"feature.data,omitempty"`
	DurationMs  int64       `json:"duration_ms,omitempty"`
	Error       string      `json:"error,omitempty"`

	// tracing data for multiple events, this is useful for specific features
	// within the Lacework CLI such as daily version check, polling mechanism, etc.
	TraceID  string `json:"trace.trace_id,omitempty"`
	SpanID   string `json:"trace.span_id,omitempty"`
	ParentID string `json:"trace.parent_id,omitempty"`
}

// InitHoneyvent initialize honeycomb library and main Honeyvent, such event
// could be modified during a command execution to add extra parameters such
// as error message, feature data, etc.
func (c *cliState) InitHoneyvent() {
	hc := libhoney.Config{
		WriteKey: HoneyApiKey,
		Dataset:  HoneyDataset,
	}
	_ = libhoney.Init(hc)

	c.Event = &Honeyvent{
		Os:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		Version: Version,
		Profile: c.Profile,
		Account: c.Account,
		ApiKey:  c.KeyID,
		TraceID: newID(),
	}
}

// Wait should be called before finishing the execution of any CLI command,
// it waits for pending workers (a.k.a. honeyvents) to be transmitted
func (c *cliState) Wait() {
	// wait for any missing worker
	c.workers.Wait()

	// flush any pending calls to Honeycomb
	libhoney.Close()
}

// SendHoneyvent is used throughout the CLI to send Honeyvents, these events
// have tracing data to understand how the commands are being executed, what
// features are used and the overall command flow. This function sends the
// events via goroutines so that we don't block the execution of the main process
//
// NOTE: the CLI will send at least one event per command execution
func (c *cliState) SendHoneyvent() {
	if disabled := os.Getenv(DisableTelemetry); disabled != "" {
		return
	}

	if c.Event.SpanID == "" {
		// root span of a trace which is defined by having its parent_id omitted
		c.Event.SpanID = c.id
	} else {
		// parent_id is set always to the root span since this is a command-line
		c.Event.ParentID = c.id
		c.Event.SpanID = newID()
	}

	c.Log.Debugw("new honeyvent", "dataset", HoneyDataset,
		"trace_id", c.Event.TraceID,
		"span_id", c.Event.SpanID,
		"parent_id", c.Event.ParentID,
	)
	honeyvent := libhoney.NewEvent()
	_ = honeyvent.Add(c.Event)

	c.workers.Add(1)
	go func(wg *sync.WaitGroup, event *libhoney.Event) {
		defer wg.Done()

		c.Log.Debugw("sending honeyvent", "dataset", HoneyDataset)
		err := event.Send()
		if err != nil {
			c.Log.Debugw("unable to send honeyvent", "error", err)
		}

	}(&c.workers, honeyvent)

	// after adding a worker to submit a honeyvent, we remove
	// all temporal fields such as feature, feature.data, error
	c.Event.DurationMs = 0
	c.Event.Error = ""
	c.Event.Feature = ""
	c.Event.FeatureData = nil
}
