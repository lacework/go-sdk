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
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/lacework/go-sdk/v2/api"
	"github.com/lacework/go-sdk/v2/lwdomain"
)

var (
	// MetricDataset is the dataset in Elastic that we send tracing
	// data this variable will be set depending on the environment we
	// are running on. During development, we send all events and
	// tracing data to a default dataset.
	MetricDataset = "lacework-cli-dev"
)

const (
	// DisableTelemetry is an environment variable that can be used to
	// disable telemetry sent to Elastic
	DisableTelemetry = "LW_TELEMETRY_DISABLE"

	// HomebrewInstall is an environment variable that denotes the
	// install method was via homebrew package manager
	HomebrewInstall = "LW_HOMEBREW_INSTALL"

	// ChocolateyInstall is an environment variable that denotes the
	// install method was via chocolatey package manager
	ChocolateyInstall = "LW_CHOCOLATEY_INSTALL"

	// List of Features
	//
	// A feature within the Lacework CLI is any functionality that
	// can't be traced or tracked by the default event sent to Elastic,
	// it is a behavior that we, Lacework engineers, would like to
	// trace and understand its usage and adoption.
	//
	// By default the Feature field within the MetricEvent is empty,
	// define a new feature below and set it before sending a new
	// MetricEvent. Additionally, there is a FeatureData field that
	// any feature can use to inject any specific information
	// related to that feature.
	//
	// Example:
	//
	// ```go
	// cli.Event.Feature = featPollCtrScan
	// cli.Event.AddFeatureField("key", "value")
	// cli.SendMetricEvent()
	// ```
	//
	// Polling mechanism feature
	featPollCtrScan = "poll_ctr_scan"

	// Daily version check feature
	featDailyVerCheck = "daily_check"

	// Daily Component version check feature
	featDailyCompVerCheck = "daily_comp_check"

	// Generate package manifest feature
	featGenPkgManifest = "gen_pkg_manifest"

	// Split package manifest feature
	featSplitPkgManifest = "split_pkg_manifest"

	// Migration API v1 -> v2 feature
	featMigrateConfigV2 = "migrate_config_v2"
)

// InitMetricEvent initialize Elastic library and main MetricEvent, such event
// could be modified during a command execution to add extra parameters such
// as error message, feature data, etc.
func (c *cliState) InitMetricEvent() {
	c.Event = &api.MetricEvent{
		Os:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		Version:       Version,
		Profile:       c.Profile,
		Account:       c.Account,
		Subaccount:    c.Subaccount,
		ApiKey:        c.KeyID,
		CfgVersion:    c.CfgVersion,
		TraceID:       newID(),
		InstallMethod: installMethod(),
		Dataset:       MetricDataset,
	}
}

// Wait should be called before finishing the execution of any CLI command,
// it waits for pending workers (a.k.a. MetricEvents) to be transmitted
func (c *cliState) Wait() {
	// wait for any missing worker
	c.workers.Wait()

	// stop gRPC server gracefully
	c.Stop()
}

// SendMetricEvent is used throughout the CLI to send metric events, these events
// have tracing data to understand how the commands are being executed, what
// features are used and the overall command flow. This function sends the
// events via goroutines so that we don't block the execution of the main process
//
// NOTE: the CLI will send at least one event per command execution
func (c *cliState) SendMetricEvent() {
	if disabled := os.Getenv(DisableTelemetry); disabled != "" {
		return
	}

	if c.LwApi == nil {
		c.Log.Debug("unable to send MetricEvent", "error")
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

	if c.Event.ContextID == "" {
		c.Event.ContextID = os.Getenv("LACEWORK_CONTEXT_ID")
	}

	// Lacework accounts are NOT case-sensitive but some users configure them
	// in uppercase and others in lowercase, therefore we will normalize all
	// account to be lowercase so that we don't see different accounts in metrics.
	c.Event.Account = strings.ToLower(c.Event.Account)

	// Detect if the account has the full domain, if so, subtract the account
	if match, _ := regexp.MatchString(".lacework.net", c.Account); match {
		d, err := lwdomain.New(c.Account)
		if err == nil {
			c.Event.Account = strings.ToLower(d.String())
		}
	}

	c.Log.Debugw("new metric event", "dataset", MetricDataset,
		"trace_id", c.Event.TraceID,
		"span_id", c.Event.SpanID,
		"parent_id", c.Event.ParentID,
		"context_id", c.Event.ContextID,
	)

	event := *c.Event
	event.Dataset = c.Event.Dataset

	c.workers.Add(1)
	go func(wg *sync.WaitGroup, event api.MetricEvent) {
		defer wg.Done()

		c.Log.Debugw("sending MetricEvent", "dataset", MetricDataset)

		err := c.LwApi.V2.Metrics.Send(event)
		if err != nil {
			c.Log.Debugw("unable to send MetricEvent", "error", err)
		}

	}(&c.workers, event)

	// after adding a worker to submit a metric event, we remove
	// all temporal fields such as feature, feature.data, error
	c.Event.DurationMs = 0
	c.Event.Error = ""
	c.Event.Feature = ""
	c.Event.FeatureData = nil
}

func installMethod() string {
	if os.Getenv(HomebrewInstall) != "" {
		return "HOMEBREW"
	}

	if os.Getenv(ChocolateyInstall) != "" {
		return "CHOCOLATEY"
	}
	return ""
}

// parseFlags is a helper used to parse all the flags that the user provided
func parseFlags(args []string) (flags []string) {
	for len(args) > 0 {
		arg := args[0]
		args = args[1:]
		if len(arg) <= 1 || arg[0] != '-' {
			// not a flag
			continue
		}

		flags = append(flags, strings.Split(arg, "=")[0])
	}
	return
}
