//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

// GetAwsCfg gets a single AwsCfg integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAwsCfg(guid string) (
	response AwsCfgIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsCfg updates a single AwsCfg integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsCfg(data CloudAccount) (
	response AwsCfgIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsCfgIntegrationResponse struct {
	Data AwsCfg `json:"data"`
}

type AwsCfg struct {
	v2CommonIntegrationData
	Data AwsCfgData `json:"data"`
}

type AwsCfgData struct {
	Credentials  AwsCfgCredentials `json:"crossAccountCredentials"`
	AwsAccountID string            `json:"awsAccountId,omitempty"`
}

type AwsCfgCredentials struct {
	RoleArn    string `json:"roleArn"`
	ExternalID string `json:"externalId"`
}
