//
// Author:: Darren Murray(<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

// GetAwsAgentlessScanning gets a single AwsAgentlessScanning integration matching the provided integration guid
func (svc *CloudAccountsService) GetAwsAgentlessScanning(guid string) (
	response AwsAgentlessScanningResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsAgentlessScanning updates a single AwsAgentlessScanning integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsAgentlessScanning(data CloudAccount) (
	response AwsAgentlessScanningResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsAgentlessScanningResponse struct {
	Data AwsAgentlessScanning `json:"data"`
}

type AwsAgentlessScanning struct {
	v2CommonIntegrationData
	Data AwsAgentlessScanningData `json:"data"`
}

type AwsAgentlessScanningData struct {
	//QueryText represents an lql json string
	QueryText string `json:"queryText"`

	//ScanFrequency in hours, 24 == 24 hours
	ScanFrequency int `json:"scanFrequency"`

	ScanContainers          bool `json:"scanContainers"`
	ScanHostVulnerabilities bool `json:"scanHostVulnerabilities"`
}
