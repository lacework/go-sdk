//
// Author:: Vatasha White (<vatasha.white@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

import "github.com/pkg/errors"

const (
	// The list of valid inputs for DatadogSite field
	DatadogSiteEu  datadogSite = "eu"
	DatadogSiteCom datadogSite = "com"

	// The list of valid inputs for DatadogService field
	DatadogServiceLogsDetails   datadogService = "Logs Detail"
	DatadogServiceEventsSummary datadogService = "Events Summary"
	DatadogServiceLogsSummary   datadogService = "Logs Summary"
)

// GetDatadog gets a single instance of a Datadog alert channel
// with the corresponding integration guid
func (svc *AlertChannelsService) GetDatadog(guid string) (response DatadogAlertChannelResponseV2, err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateDatadog updates a single instance of a Datadog integration on the Lacework server
func (svc *AlertChannelsService) UpdateDatadog(data AlertChannel) (response DatadogAlertChannelResponseV2, err error) {
	err = svc.update(data.ID(), data, &response)
	return
}

// DatadogSite returns the datadogSite type for the corresponding string input
func DatadogSite(site string) (datadogSite, error) {
	if val, ok := datadogSites[site]; ok {
		return val, nil
	}
	return "", errors.Errorf("%v is not a valid Datadog Site", site)
}

// DatadogService returns the datadogService type for the corresponding string input
func DatadogService(service string) (datadogService, error) {
	if val, ok := datadogServices[service]; ok {
		return val, nil
	}
	return "", errors.Errorf("%v is not a valid Datadog Service", service)
}

type datadogSite string
type datadogService string

type DatadogDataV2 struct {
	ApiKey      string         `json:"apiKey"`
	DatadogSite datadogSite    `json:"datadogSite,omitempty"`
	DatadogType datadogService `json:"datadogType,omitempty"`
}

type DatadogAlertChannelV2 struct {
	v2CommonIntegrationData
	Data DatadogDataV2 `json:"data"`
}

type DatadogAlertChannelResponseV2 struct {
	Data DatadogAlertChannelV2 `json:"data"`
}

var datadogSites = map[string]datadogSite{
	string(DatadogSiteEu):  DatadogSiteEu,
	string(DatadogSiteCom): DatadogSiteCom,
}

var datadogServices = map[string]datadogService{
	string(DatadogServiceLogsDetails):   DatadogServiceLogsDetails,
	string(DatadogServiceEventsSummary): DatadogServiceEventsSummary,
	string(DatadogServiceLogsSummary):   DatadogServiceLogsSummary,
}
