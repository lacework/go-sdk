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

// GetAwsS3 gets a single AwsS3 alert channel matching the
// provided integration guid
func (svc *AlertChannelsService) GetAwsS3(guid string) (
	response AwsS3AlertChannelResponseV2,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsS3 updates a single AwsS3 integration on the Lacework Server
func (svc *AlertChannelsService) UpdateAwsS3(data AlertChannel) (
	response AwsS3AlertChannelResponseV2,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsS3AlertChannelResponseV2 struct {
	Data AwsS3AlertChannelV2 `json:"data"`
}

type AwsS3AlertChannelV2 struct {
	v2CommonIntegrationData
	Data AwsS3DataV2 `json:"data"`
}

type AwsS3DataV2 struct {
	Credentials AwsS3Credentials `json:"s3CrossAccountCredentials"`
}

type AwsS3Credentials struct {
	RoleArn    string `json:"roleArn"`
	ExternalID string `json:"externalId"`
	BucketArn  string `json:"bucketArn"`
}
