//
// Author:: Teddy Reed (<teddy.reed@lacework.net>)
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

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// GetAwsSidekickOrg gets a single AwsSidekickOrg integration matching the provided integration guid
func (svc *CloudAccountsService) GetAwsSidekickOrg(guid string) (
	response AwsSidekickOrgResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// CreateAwsSidekickOrg creates an AwsSidekickOrg Cloud Account integration
func (svc *CloudAccountsService) CreateAwsSidekickOrg(data CloudAccount) (
	response AwsSidekickOrgResponse,
	err error,
) {
	err = svc.create(data, &response)
	return
}

// UpdateAwsSidekickOrg updates a single AwsSidekickOrg integration on the Lacework Server
func (svc *CloudAccountsService) UpdateAwsSidekickOrg(data CloudAccount) (
	response AwsSidekickOrgResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type AwsSidekickOrgResponse struct {
	Data AwsSidekickOrg `json:"data"`
}

type AwsSidekickOrg struct {
	v2CommonIntegrationData
	awsSidekickToken `json:"serverToken"`
	Data             AwsSidekickOrgData `json:"data"`
}

type AwsSidekickOrgData struct {
	//QueryText represents an lql json string
	QueryText string `json:"queryText,omitempty"`

	//ScanFrequency in hours, 24 == 24 hours
	ScanFrequency int `json:"scanFrequency"`

	ScanContainers          bool `json:"scanContainers"`
	ScanHostVulnerabilities bool `json:"scanHostVulnerabilities"`

	//Properties specific to the AWS organization integration type
	ScanningAccount   string `json:"scanningAccount"`
	ManagementAccount string `json:"managementAccount,omitempty"`
	MonitoredAccounts string `json:"monitoredAccounts"`

	AccountID          string                             `json:"awsAccountId,omitempty"`
	BucketArn          string                             `json:"bucketArn,omitempty"`
	CrossAccountCreds  AwsSidekickCrossAccountCredentials `json:"crossAccountCredentials"`
	AccountMappingFile string                             `json:"accountMappingFile,omitempty"`
}

func (aws *AwsSidekickOrgData) EncodeAccountMappingFile(mapping []byte) {
	encodedMappings := base64.StdEncoding.EncodeToString(mapping)
	aws.AccountMappingFile = fmt.Sprintf("data:application/json;name=i.json;base64,%s", encodedMappings)
}

func (aws *AwsSidekickOrgData) DecodeAccountMappingFile() ([]byte, error) {
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
