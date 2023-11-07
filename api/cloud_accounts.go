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
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/lwtime"
)

// CloudAccountsService is the service that interacts with
// the CloudAccounts schema from the Lacework APIv2 Server
type CloudAccountsService struct {
	client *Client
}

// NewCloudAccount returns an instance of the CloudAccountRaw struct with the
// provided Cloud Account integration type, name and raw data as an interface{}.
//
// NOTE: This function must be used by any Cloud Account type.
//
// Basic usage: Initialize a new AwsIntegration struct, then use the new
// instance to do CRUD operations
//
//	client, err := api.NewClient("account")
//	if err != nil {
//	  return err
//	}
//
//	awsCtSqs := api.NewCloudAccount("foo",
//	  api.AwsCtSqsCloudAccount,
//	  api.AwsCtSqsData{
//	    QueueUrl: "https://sqs.us-west-2.amazonaws.com/123456789000/lw",
//	    Credentials: &api.AwsCtSqsCredentials {
//	      RoleArn: "arn:aws:XYZ",
//	      ExternalID: "1",
//	    },
//	  },
//	)
//
//	client.V2.CloudAccount.Create(awsCtSqs)
func NewCloudAccount(name string, iType cloudAccountType, data interface{}) CloudAccountRaw {
	return CloudAccountRaw{
		v2CommonIntegrationData: v2CommonIntegrationData{
			Name:    name,
			Type:    iType.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CloudAccount is an interface that helps us implement a few functions
// that any Cloud Account might use, there are some cases, like during
// Update, where we need to get the ID of the Cloud Account and its type,
// this will allow users to pass any Cloud Account that implements these
// methods
type CloudAccount interface {
	ID() string
	CloudAccountType() cloudAccountType
}

type cloudAccountType int

const (
	// type that defines a non-existing Cloud Account integration
	NoneCloudAccount cloudAccountType = iota
	AwsCfgCloudAccount
	AwsCtSqsCloudAccount
	AwsEksAuditCloudAccount
	AwsSidekickCloudAccount
	AwsSidekickOrgCloudAccount
	AwsUsGovCfgCloudAccount
	AwsUsGovCtSqsCloudAccount
	AzureAlSeqCloudAccount
	AzureCfgCloudAccount
	GcpAtSesCloudAccount
	GcpCfgCloudAccount
	GcpGkeAuditCloudAccount
	GcpSidekickCloudAccount
	AzureSidekickCloudAccount
	GcpAlPubSubCloudAccount
	OciCfgCloudAccount
)

// CloudAccountTypes is the list of available Cloud Account integration types
var CloudAccountTypes = map[cloudAccountType]string{
	NoneCloudAccount:           "None",
	AwsCfgCloudAccount:         "AwsCfg",
	AwsCtSqsCloudAccount:       "AwsCtSqs",
	AwsEksAuditCloudAccount:    "AwsEksAudit",
	AwsSidekickCloudAccount:    "AwsSidekick",
	AwsSidekickOrgCloudAccount: "AwsSidekickOrg",
	AwsUsGovCfgCloudAccount:    "AwsUsGovCfg",
	AwsUsGovCtSqsCloudAccount:  "AwsUsGovCtSqs",
	AzureAlSeqCloudAccount:     "AzureAlSeq",
	AzureCfgCloudAccount:       "AzureCfg",
	GcpAtSesCloudAccount:       "GcpAtSes",
	GcpCfgCloudAccount:         "GcpCfg",
	GcpGkeAuditCloudAccount:    "GcpGkeAudit",
	GcpSidekickCloudAccount:    "GcpSidekick",
	AzureSidekickCloudAccount:  "AzureSidekick",
	GcpAlPubSubCloudAccount:    "GcpAlPubSub",
	OciCfgCloudAccount:         "OciCfg",
}

// String returns the string representation of a Cloud Account integration type
func (i cloudAccountType) String() string {
	return CloudAccountTypes[i]
}

// FindCloudAccountType looks up inside the list of available cloud account types
// the matching type from the provided string, if none, returns NoneCloudAccount
func FindCloudAccountType(cloudAccount string) (cloudAccountType, bool) {
	for cType, cStr := range CloudAccountTypes {
		if cStr == cloudAccount {
			return cType, true
		}
	}
	return NoneCloudAccount, false
}

// List returns a list of Cloud Account integrations
func (svc *CloudAccountsService) List() (response CloudAccountsResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2CloudAccounts, nil, &response)
	return
}

// ListByType lists the cloud accounts from the provided type that are available
// on the Lacework Server
func (svc *CloudAccountsService) ListByType(caType cloudAccountType) (response CloudAccountsResponse, err error) {
	err = svc.get(caType.String(), &response)
	return
}

// Create creates a single Cloud Account integration
func (svc *CloudAccountsService) Create(integration CloudAccountRaw) (
	response CloudAccountResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// Delete deletes a Cloud Account integration that matches the provided guid
func (svc *CloudAccountsService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2CloudAccountsWithParam, guid),
		nil,
		nil,
	)
}

// Migrate marks a Cloud Account integration that matches the provided guid for migration
func (svc *CloudAccountsService) Migrate(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	data := MigrateRequestData{
		MigrateData{
			IntgGuid: guid,
			Props: Props{
				Migrate:            true,
				MigrationTimestamp: time.Now(),
			},
		},
	}

	return svc.client.RequestEncoderDecoder(
		"PATCH",
		apiV2MigrateGcpAtSes,
		data,
		nil,
	)
}

// Get returns a raw response of the Cloud Account with the matching integration guid.
//
// To return a more specific Go struct of a Cloud Account integration, use the proper
// method such as GetAwsCtSqs() where the function name is composed by:
//
//	Get<Type>(guid)
//
//	  Where <Type> is the Cloud Account integration type.
func (svc *CloudAccountsService) Get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}
	return svc.get(guid, &response)
}

type CloudAccountRaw struct {
	v2CommonIntegrationData
	Data interface{} `json:"data,omitempty"`
}

func (cloud CloudAccountRaw) GetData() any {
	return cloud.Data
}

func (cloud CloudAccountRaw) GetCommon() v2CommonIntegrationData {
	return cloud.v2CommonIntegrationData
}

func (cloud CloudAccountRaw) CloudAccountType() cloudAccountType {
	t, _ := FindCloudAccountType(cloud.Type)
	return t
}

type CloudAccountResponse struct {
	Data CloudAccountRaw `json:"data"`
}

type CloudAccountsResponse struct {
	Data []CloudAccountRaw `json:"data"`
}

type v2CommonIntegrationData struct {
	IntgGuid             string              `json:"intgGuid,omitempty"`
	Name                 string              `json:"name"`
	CreatedOrUpdatedTime string              `json:"createdOrUpdatedTime,omitempty"`
	CreatedOrUpdatedBy   string              `json:"createdOrUpdatedBy,omitempty"`
	Type                 string              `json:"type"`
	Enabled              int                 `json:"enabled"`
	IsOrg                int                 `json:"isOrg,omitempty"`
	State                *V2IntegrationState `json:"state,omitempty"`
}

func (c v2CommonIntegrationData) ID() string {
	return c.IntgGuid
}

func (c v2CommonIntegrationData) Status() string {
	if c.Enabled == 1 {
		return "Enabled"
	}
	return "Disabled"
}

func (c v2CommonIntegrationData) StateString() string {
	if c.State != nil && c.State.Ok {
		return "Ok"
	}
	return "Pending"
}

type V2IntegrationState struct {
	Ok                 bool                   `json:"ok"`
	Details            map[string]interface{} `json:"details"`
	LastUpdatedTime    lwtime.Epoch           `json:"lastUpdatedTime"`
	LastSuccessfulTime lwtime.Epoch           `json:"lastSuccessfulTime"`
}

func (svc *CloudAccountsService) create(data interface{}, response interface{}) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2CloudAccounts, data, response)
}

func (svc *CloudAccountsService) get(param string, response interface{}) error {
	apiPath := fmt.Sprintf(apiV2CloudAccountsWithParam, param)
	return svc.client.RequestDecoder("GET", apiPath, nil, response)
}

func (svc *CloudAccountsService) update(guid string, data interface{}, response interface{}) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}
	apiPath := fmt.Sprintf(apiV2CloudAccountsWithParam, guid)
	return svc.client.RequestEncoderDecoder("PATCH", apiPath, data, response)
}
