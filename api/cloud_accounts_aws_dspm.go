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

// GetAwsDspm gets a single AwsDspm integration matching the provided integration guid
func (svc *CloudAccountsService) GetAwsDspm(guid string) (
	response AwsDspmResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// CreateAwsDspm creates an AwsDspm Cloud Account integration
func (svc *CloudAccountsService) CreateAwsDspm(data CloudAccount) (
	response AwsDspmResponse,
	err error,
) {
	err = svc.create(data, &response)
	return
}

// UpdateAwsDspm updates a single AwsDspm integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsDspm(data CloudAccount) (
	response AwsDspmResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsDspmResponse struct {
	Data AwsDspm `json:"data"`
}

type AwsDspm struct {
	v2CommonIntegrationData
	Data AwsDspmData `json:"data"`
}

// AwsDspmData contains the data needed by Lacework platform services.
type AwsDspmData struct {
	AccountID         string                         `json:"awsAccountId,omitempty"`
	BucketArn         string                         `json:"bucketArn,omitempty"`
	CrossAccountCreds AwsDspmCrossAccountCredentials `json:"crossAccountCredentials"`
}

type AwsDspmCrossAccountCredentials struct {
	ExternalID string `json:"externalId,omitempty"`
	RoleArn    string `json:"roleArn,omitempty"`
}
