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

// GetAwsUsGovCfg gets a single AwsUsGovCfg integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAwsUsGovCfg(guid string) (
	response AwsUsGovCfgIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsUsGovCfg updates a single AwsUsGovCfg integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsUsGovCfg(data CloudAccount) (
	response AwsUsGovCfgIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsUsGovCfgIntegrationResponse struct {
	Data AwsUsGovCfg `json:"data"`
}

type AwsUsGovCfg struct {
	v2CommonIntegrationData
	Data AwsUsGovCfgData `json:"data"`
}

type AwsUsGovCfgData struct {
	Credentials AwsUsGovCfgCredentials `json:"accessKeyCredentials"`
}

type AwsUsGovCfgCredentials struct {
	AwsAccountID    string `json:"accountId"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}
