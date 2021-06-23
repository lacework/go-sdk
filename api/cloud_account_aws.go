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

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// GetAwsCtSqs gets a single AwsCtSqs integration matching the
// provided integration guid
func (svc *CloudAccountsService) GetAwsCtSqs(guid string) (
	response AwsCtSqsIntegrationResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateAwsCtSqs updates a single AwsCtSqs integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsCtSqs(data AwsCtSqsIntegration) (
	response AwsCtSqsIntegrationResponse,
	err error,
) {
	err = svc.update(data.IntgGuid, data, &response)
	return
}

type AwsCtSqsIntegrationResponse struct {
	Data AwsCtSqsIntegration `json:"data"`
}

type AwsCtSqsIntegration struct {
	v2CommonIntegrationData
	Data AwsCtSqsData `json:"data"`
}

type AwsCtSqsData struct {
	Credentials  AwsCtSqsCredentials `json:"crossAccountCredentials"`
	QueueUrl     string              `json:"queueUrl"`
	AwsAccountID string              `json:"awsAccountId,omitempty"`

	// This field must be a base64 encode with the following format:
	//
	// "data:application/json;name=i.json;base64,[ENCODING]"
	//
	// [ENCODING] is the the base64 encode, use EncodeAccountMappingFile() to encode a JSON mapping file
	AccountMappingFile string `json:"accountMappingFile,omitempty"`
}

type AwsCtSqsCredentials struct {
	RoleArn    string `json:"roleArn"`
	ExternalID string `json:"externalId"`
}

func (aws *AwsCtSqsData) EncodeAccountMappingFile(mapping []byte) {
	encodedMappings := base64.StdEncoding.EncodeToString(mapping)
	aws.AccountMappingFile = fmt.Sprintf("data:application/json;name=i.json;base64,%s", encodedMappings)
}

func (aws *AwsCtSqsData) DecodeAccountMappingFile() ([]byte, error) {
	if len(aws.AccountMappingFile) == 0 {
		return []byte{}, nil
	}

	var (
		b64      = strings.Split(aws.AccountMappingFile, ",")
		raw, err = base64.StdEncoding.DecodeString(b64[1])
	)
	if err != nil {
		return []byte{}, err
	}

	return raw, nil
}
