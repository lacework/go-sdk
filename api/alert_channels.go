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

	"github.com/pkg/errors"
)

// AlertChannelsService is the service that interacts with
// the AlertChannels schema from the Lacework APIv2 Server
type AlertChannelsService struct {
	client *Client
}

// NewAlertChannel returns an instance of the AlertChannelRaw struct with the
// provided Alert Channel integration type, name and raw data as an interface{}.
//
// NOTE: This function must be used by any Alert Channel type.
//
// Basic usage: Initialize a new EmailUserAlertChannel struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   emailAlertChan := api.NewAlertChannel("foo",
//     api.EmailUserAlertChannelType,
//     api.EmailUserData{
//       ChannelProps: api.EmailUserChannelProps{
//         Recipients: []string{"name@example.com"},
//       },
//     },
//   )
//
//   client.V2.AlertChannels.Create(emailAlertChan)
//
func NewAlertChannel(name string, iType alertChannelType, data interface{}) AlertChannelRaw {
	return AlertChannelRaw{
		v2CommonIntegrationData: v2CommonIntegrationData{
			Name:    name,
			Type:    iType.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// AlertChannel is an interface that helps us implement a few functions
// that any Alert Channel might use, there are some cases, like during
// Update, where we need to get the ID of the Alert Channel and its type,
// this will allow users to pass any Alert Channel that implements these
// methods
type AlertChannel interface {
	ID() string
	AlertChannelType() alertChannelType
}

type alertChannelType int

const (
	// type that defines a non-existing Alert Channel integration
	NoneAlertChannelType alertChannelType = iota
	EmailUserAlertChannelType
	SlackChannelAlertChannelType
	AwsS3AlertChannelType
	CloudwatchEbAlertChannelType
	DatadogAlertChannelType
	WebhookAlertChannelType
	VictorOpsAlertChannelType
	CiscoSparkWebhookAlertChannelType
	MicrosoftTeamsAlertChannelType
	SplunkHecAlertChannelType
)

// AlertChannelTypeTypes is the list of available Alert Channel integration types
var AlertChannelTypes = map[alertChannelType]string{
	NoneAlertChannelType:              "None",
	EmailUserAlertChannelType:         "EmailUser",
	SlackChannelAlertChannelType:      "SlackChannel",
	AwsS3AlertChannelType:             "AwsS3",
	CloudwatchEbAlertChannelType:      "CloudwatchEb",
	DatadogAlertChannelType:           "Datadog",
	WebhookAlertChannelType:           "Webhook",
	VictorOpsAlertChannelType:         "VictorOps",
	CiscoSparkWebhookAlertChannelType: "CiscoSparkWebhook",
	MicrosoftTeamsAlertChannelType:    "MicrosoftTeams",
	SplunkHecAlertChannelType:         "SplunkHec",
}

// String returns the string representation of a Alert Channel integration type
func (i alertChannelType) String() string {
	return AlertChannelTypes[i]
}

// FindAlertChannelType looks up inside the list of available alert channel types
// the matching type from the provided string, if none, returns NoneAlertChannelType
func FindAlertChannelType(alertChannel string) (alertChannelType, bool) {
	for cType, cStr := range AlertChannelTypes {
		if cStr == alertChannel {
			return cType, true
		}
	}
	return NoneAlertChannelType, false
}

// List returns a list of Alert Channel integrations
func (svc *AlertChannelsService) List() (response AlertChannelsResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2AlertChannels, nil, &response)
	return
}

// Create creates a single Alert Channel integration
func (svc *AlertChannelsService) Create(integration AlertChannelRaw) (
	response AlertChannelResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// Delete deletes a Alert Channel integration that matches the provided guid
func (svc *AlertChannelsService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2AlertChannelFromGUID, guid),
		nil,
		nil,
	)
}

// Test tests an Alert Channel integration that matches the provided guid
func (svc *AlertChannelsService) Test(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	apiPath := fmt.Sprintf(apiV2AlertChannelTest, guid)
	return svc.client.RequestDecoder("POST", apiPath, nil, nil)
}

// Get returns a raw response of the Alert Channel with the matching integration guid.
//
// To return a more specific Go struct of a Alert Channel integration, use the proper
// method such as GetEmailUser() where the function name is composed by:
//
//  Get<Type>(guid)
//
//    Where <Type> is the Alert Channel integration type.
func (svc *AlertChannelsService) Get(guid string, response interface{}) error {
	return svc.get(guid, &response)
}

type AlertChannelRaw struct {
	v2CommonIntegrationData
	Data interface{} `json:"data,omitempty"`
}

func (alert AlertChannelRaw) AlertChannelType() alertChannelType {
	t, _ := FindAlertChannelType(alert.Type)
	return t
}

type AlertChannelResponse struct {
	Data AlertChannelRaw `json:"data"`
}

type AlertChannelsResponse struct {
	Data []AlertChannelRaw `json:"data"`
}

func (svc *AlertChannelsService) create(data interface{}, response interface{}) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2AlertChannels, data, response)
}

func (svc *AlertChannelsService) get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}
	apiPath := fmt.Sprintf(apiV2AlertChannelFromGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, response)
}

func (svc *AlertChannelsService) update(guid string, data interface{}, response interface{}) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}
	apiPath := fmt.Sprintf(apiV2AlertChannelFromGUID, guid)
	return svc.client.RequestEncoderDecoder("PATCH", apiPath, data, response)
}
