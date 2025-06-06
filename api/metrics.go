//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2024, Lacework Inc.
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
	"runtime"
)

const DisableTelemetry = "LW_TELEMETRY_DISABLE"

// MetricsService is a service that sends events to Lacework APIv2 Server metrics endpoint
type MetricsService struct {
	client *Client
}

func (svc *MetricsService) Send(event MetricEvent) (err error) {
	if disabled := os.Getenv(DisableTelemetry); disabled != "" {
		return
	}

	event.setAccountDetails(*svc.client)
	event.SampleRate100 = true
	event.TelemetrySource = "external"
	event.TelemetryType = "customer"

	err = svc.client.RequestEncoderDecoder("POST", fmt.Sprintf(apiV2OtelMetrics, event.Dataset), event, nil)
	return
}

func NewMetricEvent(version, feature, dataset string) MetricEvent {
	event := MetricEvent{
		Os:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		TraceID: newID(),
		Version: version,
		Dataset: dataset,
		Feature: feature,
	}

	return event
}

func (h *MetricEvent) setAccountDetails(client Client) {
	if h.Account == "" {
		h.Account = client.account
	}
	if h.Subaccount == "" {
		h.Subaccount = client.subaccount
	}
}

// MetricEvent defines what a metric event looks like for the Lacework CLI
type MetricEvent struct {
	Version       string      `json:"version"`
	CfgVersion    int         `json:"config_version"`
	Os            string      `json:"os"`
	Arch          string      `json:"arch"`
	Command       string      `json:"command,omitempty"`
	Args          []string    `json:"args,omitempty"`
	Flags         []string    `json:"flags,omitempty"`
	Account       string      `json:"account,omitempty"`
	Subaccount    string      `json:"subaccount,omitempty"`
	Profile       string      `json:"profile,omitempty"`
	ApiKey        string      `json:"api_key,omitempty"`
	Feature       string      `json:"feature,omitempty"`
	FeatureData   interface{} `json:"feature.data,omitempty"`
	DurationMs    int64       `json:"duration_ms,omitempty"`
	Error         string      `json:"error,omitempty"`
	InstallMethod string      `json:"install_method,omitempty"`
	Component     string      `json:"component,omitempty"`
	Dataset       string      `json:"dataset,omitempty"`

	// tracing data for multiple events, this is useful for specific features
	// within the Lacework CLI such as daily version check, polling mechanism, etc.
	TraceID   string `json:"trace.trace_id,omitempty"`
	SpanID    string `json:"trace.span_id,omitempty"`
	ParentID  string `json:"trace.parent_id,omitempty"`
	ContextID string `json:"trace.context_id,omitempty"`

	SampleRate100   bool   `json:"sample_rate_100,omitempty"`
	TelemetrySource string `json:"telemetry_source,omitempty"`
	TelemetryType   string `json:"telemetry_type,omitempty"`
}

func (e *MetricEvent) AddFeatureField(key string, value interface{}) {
	if e.FeatureData == nil {
		e.FeatureData = map[string]interface{}{key: value}
		return
	}

	if v, ok := e.FeatureData.(map[string]interface{}); ok {
		v[key] = value
		e.FeatureData = v
	}
}
