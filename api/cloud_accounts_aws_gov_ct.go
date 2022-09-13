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

// GetAwsUsGovCtSqs gets a single AwsUsGovCtSqs integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAwsUsGovCtSqs(guid string) (
	response AwsUsGovCtSqsIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsUsGovCtSqs updates a single AwsUsGovCtSqs integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsUsGovCtSqs(data CloudAccount) (
	response AwsUsGovCtSqsIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsUsGovCtSqsIntegrationResponse struct {
	Data AwsUsGovCtSqs `json:"data"`
}

type AwsUsGovCtSqs struct {
	v2CommonIntegrationData
	Data AwsUsGovCtSqsData `json:"data"`
}

type AwsUsGovCtSqsData struct {
	Credentials AwsUsGovCtSqsCredentials `json:"accessKeyCredentials"`
	QueueUrl    string                   `json:"queueUrl"`
}

type AwsUsGovCtSqsCredentials struct {
	AwsAccountID    string `json:"accountId"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}
