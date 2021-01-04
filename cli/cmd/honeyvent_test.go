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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHoneyventDefaultParameters(t *testing.T) {
	assert.NotNil(t, cli.Event)
	assert.Equal(t, Version, cli.Event.Version)
	assert.Equal(t, runtime.GOOS, cli.Event.Os)
	assert.Equal(t, runtime.GOARCH, cli.Event.Arch)
	assert.Equal(t, cli.Profile, cli.Event.Profile)
	assert.Equal(t, cli.Account, cli.Event.Account)
	assert.Equal(t, cli.KeyID, cli.Event.ApiKey)
	assert.NotEmpty(t, cli.Event.TraceID)
	assert.Empty(t, cli.Event.SpanID)
	assert.Empty(t, cli.Event.ParentID)
}

func TestSendHoneyventTracingFields(t *testing.T) {
	// by default, the span_id must be empty
	assert.Empty(t, cli.Event.SpanID)
	assert.Empty(t, cli.Event.ParentID)

	// mocking sending first honeyvent
	cli.SendHoneyvent()

	// the span_id should be set to the cli id
	// but the parent_id must continue to be empty
	assert.Equal(t, cli.id, cli.Event.SpanID)
	assert.Empty(t, cli.Event.ParentID)

	// mocking sending second honeyvent
	cli.SendHoneyvent()

	// any further event should set the parent_id as the cli id
	// and generate a new id for the span_id
	assert.NotEmpty(t, cli.Event.SpanID)
	assert.NotEqual(t, cli.id, cli.Event.SpanID)
	assert.Equal(t, cli.id, cli.Event.ParentID)
}

func TestSendHoneyventFeatureFields(t *testing.T) {
	// by default, the feature, feature.data and duration_ms should be empty
	assert.Empty(t, cli.Event.DurationMs)
	assert.Empty(t, cli.Event.Error)
	assert.Empty(t, cli.Event.Feature)
	assert.Empty(t, cli.Event.FeatureData)

	// a new feature will need to set at least the feature field
	cli.Event.Feature = featPollCtrScan

	// additionally, features could define the feature.data and duration_ms
	cli.Event.FeatureData = map[string]interface{}{"key": "value"}
	cli.Event.DurationMs = 639023
	cli.Event.Error = "something happened"

	// mocking sending honeyvent
	cli.SendHoneyvent()

	// after submitting the honeyvent, the global
	// event struct should be resetted
	assert.Empty(t, cli.Event.DurationMs)
	assert.Empty(t, cli.Event.Error)
	assert.Empty(t, cli.Event.Feature)
	assert.Empty(t, cli.Event.FeatureData)
}

func TestSendHoneyventDisableTelemetry(t *testing.T) {
	// testing that the func SendHoneyvent won't run when the
	// environment variable 'DisableTelemetry'  is set
	os.Setenv(DisableTelemetry, "1")
	defer os.Setenv(DisableTelemetry, "")

	// setting up a test feature
	cli.Event.Feature = "test_feature"

	// mocking sending honeyvent
	cli.SendHoneyvent()

	// all these fields should not be empty after sending the event
	assert.Equal(t, "test_feature", cli.Event.Feature)
	// this validates that the environment variable is not sending
	// events when it is set (disabled)
}
