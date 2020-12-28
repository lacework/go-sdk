//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

// NewAwsS3AlertChannel returns an instance of AwsS3AlertChannel
// with the provided name and data.
//
// Basic usage: Initialize a new AwsS3AlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   awsS3Channel := api.NewAwsS3AlertChannel("foo",
//     api.AwsS3ChannelData{
//       Credentials: api.AwsS3Creds{
//         ExternalID:  "1234",
//         RoleArn:     "arn:aws:iam::account-id:role/role-name-with-path",
//         BucketArn:   "arn:aws:s3:::bucket_name/key_name",
//       },
//     },
//   )
//
//   client.Integrations.CreateAwsS3AlertChannel(awsS3Channel)
//
func NewAwsS3AlertChannel(name string, data AwsS3ChannelData) AwsS3AlertChannel {
	return AwsS3AlertChannel{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    AwsS3ChannelIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateAwsS3AlertChannel creates an AWS_S3 alert channel integration on the Lacework Server
func (svc *IntegrationsService) CreateAwsS3AlertChannel(integration AwsS3AlertChannel) (
	response AwsS3AlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetAwsS3AlertChannel gets an AWS_S3 alert channel integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetAwsS3AlertChannel(guid string) (
	response AwsS3AlertChannelResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsS3AlertChannel updates a single AWS_S3 alert channel integration
func (svc *IntegrationsService) UpdateAwsS3AlertChannel(data AwsS3AlertChannel) (
	response AwsS3AlertChannelResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

// ListAwsS3AlertChannel lists the AWS_S3 external integrations available on the Lacework Server
func (svc *IntegrationsService) ListAwsS3AlertChannel() (response AwsS3AlertChannelResponse, err error) {
	err = svc.listByType(AwsS3ChannelIntegration, &response)
	return
}

type AwsS3AlertChannelResponse struct {
	Data    []AwsS3AlertChannel `json:"data"`
	Ok      bool                `json:"ok"`
	Message string              `json:"message"`
}

type AwsS3AlertChannel struct {
	commonIntegrationData
	Data AwsS3ChannelData `json:"DATA"`
}

type AwsS3ChannelData struct {
	Credentials AwsS3Creds `json:"S3_CROSS_ACCOUNT_CREDENTIALS" mapstructure:"S3_CROSS_ACCOUNT_CREDENTIALS"`
}

type AwsS3Creds struct {
	ExternalID string `json:"EXTERNAL_ID" mapstructure:"EXTERNAL_ID"`
	RoleArn    string `json:"ROLE_ARN" mapstructure:"ROLE_ARN"`
	BucketArn  string `json:"BUCKET_ARN" mapstructure:"BUCKET_ARN"`
}
