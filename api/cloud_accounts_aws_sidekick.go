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

// GetAwsSidekick gets a single AwsSidekick integration matching the provided integration guid
func (svc *CloudAccountsService) GetAwsSidekick(guid string) (
	response AwsSidekickResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// CreateAwsSidekick creates an AwsSidekick Cloud Account integration
func (svc *CloudAccountsService) CreateAwsSidekick(data CloudAccount) (
	response AwsSidekickResponse,
	err error,
) {
	err = svc.create(data, &response)
	return
}

// UpdateAwsSidekick updates a single AwsSidekick integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsSidekick(data CloudAccount) (
	response AwsSidekickResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsSidekickResponse struct {
	Data AwsSidekick `json:"data"`
}

type AwsSidekick struct {
	v2CommonIntegrationData
	awsSidekickToken `json:"serverToken"`
	Data             AwsSidekickData `json:"data"`
}

type awsSidekickToken struct {
	ServerToken string `json:"serverToken"`
	Uri         string `json:"uri"`
}

type AwsSidekickData struct {
	//QueryText represents an lql json string
	QueryText string `json:"queryText,omitempty"`

	//ScanFrequency in hours, 24 == 24 hours
	ScanFrequency int `json:"scanFrequency"`

	ScanContainers          bool `json:"scanContainers"`
	ScanHostVulnerabilities bool `json:"scanHostVulnerabilities"`
	ScanMultiVolume         bool `json:"scanMultiVolume"`
	ScanStoppedInstances    bool `json:"scanStoppedInstances"`

	AccountID         string                             `json:"awsAccountId,omitempty"`
	BucketArn         string                             `json:"bucketArn,omitempty"`
	CrossAccountCreds AwsSidekickCrossAccountCredentials `json:"crossAccountCredentials"`
}

type AwsSidekickCrossAccountCredentials struct {
	RoleArn    string `json:"roleArn,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
}
