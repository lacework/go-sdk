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

// GetAwsGovCt gets a single AwsGovCt integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAwsGovCt(guid string) (
	response AwsGovCtIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsGovCt updates a single AwsGovCt integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsGovCt(data CloudAccount) (
	response AwsGovCtIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsGovCtIntegrationResponse struct {
	Data AwsGovCt `json:"data"`
}

type AwsGovCt struct {
	v2CommonIntegrationData
	Data AwsGovCtData `json:"data"`
}

type AwsGovCtData struct {
	Credentials AwsGovCtCredentials `json:"accessKeyCredentials"`
	QueueUrl    string              `json:"queueUrl"`
}

type AwsGovCtCredentials struct {
	AwsAccountID    string `json:"awsAccountId"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}
