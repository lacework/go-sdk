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

// GetAwsGovCfg gets a single AwsGovCfg integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAwsGovCfg(guid string) (
	response AwsGovCfgIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsGovCfg updates a single AwsGovCfg integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsGovCfg(data CloudAccount) (
	response AwsGovCfgIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsGovCfgIntegrationResponse struct {
	Data AwsGovCfg `json:"data"`
}

type AwsGovCfg struct {
	v2CommonIntegrationData
	Data AwsGovCfgData `json:"data"`
}

type AwsGovCfgData struct {
	Credentials AwsGovCfgCredentials `json:"accessKeyCredentials"`
}

type AwsGovCfgCredentials struct {
	AwsAccountID    string `json:"awsAccountId"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}
