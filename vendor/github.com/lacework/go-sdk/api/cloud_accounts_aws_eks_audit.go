//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

// GetAwsEksAudit gets a single AwsEksAudit integration matching the provided integration guid
func (svc *CloudAccountsService) GetAwsEksAudit(guid string) (
	response AwsEksAuditIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsEksAudit updates a single AwsEksAudit integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsEksAudit(data CloudAccount) (
	response AwsEksAuditIntegrationResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsEksAuditIntegrationResponse struct {
	Data AwsEksAuditIntegration `json:"data"`
}

type AwsEksAuditIntegration struct {
	v2CommonIntegrationData
	Data AwsEksAuditData `json:"data"`
}

type AwsEksAuditData struct {
	Credentials AwsEksAuditCredentials `json:"crossAccountCredentials"`
	SnsArn      string                 `json:"snsArn"`
	S3BucketArn string                 `json:"s3BucketArn,omitempty"`
}

type AwsEksAuditCredentials struct {
	RoleArn    string `json:"roleArn"`
	ExternalID string `json:"externalId"`
}
