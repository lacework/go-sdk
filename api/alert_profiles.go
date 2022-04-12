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

import (
	"fmt"

	"github.com/pkg/errors"
)

// AlertProfilesService is the service that interacts with
// the AlertProfiles schema from the Lacework APIv2 Server
type AlertProfilesService struct {
	client    *Client
	Templates AlertTemplatesService
}

// NewAlertProfile returns an instance of the AlertProfileConfig struct
//
// Basic usage: Initialize a new AlertProfileConfig struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   alertProfile := api.NewAlertProfile(
//		"Foo",
//		api.AlertProfileConfig{
// 		...
//     },
//   )
//
//   client.V2.AlertProfiles.Create(AlertProfile)
//
func NewAlertProfile(id string, extends string, alerts []AlertProfileAlert) AlertProfileConfig {
	profile := AlertProfileConfig{
		Guid:    id,
		Extends: extends,
		Alerts:  alerts,
	}
	return profile
}

// List returns a list of Alert Profiles
func (svc *AlertProfilesService) List() (response AlertProfilesResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2AlertProfiles, nil, &response)
	return
}

// Create creates a single Alert Profile
func (svc *AlertProfilesService) Create(profile AlertProfileConfig) (
	response AlertProfileResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2AlertProfiles, profile, &response)
	return
}

// Delete deletes a Alert Profile that matches the provided guid
func (svc *AlertProfilesService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2AlertProfileFromGUID, guid),
		nil,
		nil,
	)
}

// Update updates a single Alert Profile of the provided guid.
func (svc *AlertProfilesService) Update(guid string, data AlertProfileUpdate) (
	response AlertProfileResponse,
	err error,
) {
	if guid == "" {
		err = errors.New("specify a Guid")
		return
	}
	apiPath := fmt.Sprintf(apiV2AlertProfileFromGUID, guid)
	err = svc.client.RequestEncoderDecoder("PATCH", apiPath, data, &response)
	return
}

// Get returns a raw response of the Alert Profile with the matching guid.
func (svc *AlertProfilesService) Get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify a Guid")
	}
	apiPath := fmt.Sprintf(apiV2AlertProfileFromGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, &response)
}

type AlertProfile struct {
	Guid            string                        `json:"alertProfileId,omitempty"`
	Extends         string                        `json:"extends"`
	Fields          []AlertProfileField           `json:"fields,omitempty"`
	DescriptionKeys []AlertProfileDescriptionKeys `json:"descriptionKeys,omitempty"`
	Alerts          []AlertProfileAlert           `json:"alerts"`
}

type AlertProfileConfig struct {
	Guid    string              `json:"alertProfileId"`
	Extends string              `json:"extends"`
	Alerts  []AlertProfileAlert `json:"alerts"`
}

type AlertProfileField struct {
	Name string `json:"name"`
}

type AlertProfileDescriptionKeys struct {
	Name string `json:"name"`
	Spec string `json:"spec"`
}

type AlertProfileAlert struct {
	Name        string `json:"name"`
	EventName   string `json:"eventName"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
}

type AlertProfileUpdate struct {
	Alerts []AlertProfileAlert `json:"alerts"`
}

type AlertProfileResponse struct {
	Data AlertProfile `json:"data"`
}

type AlertProfilesResponse struct {
	Data []AlertProfile `json:"data"`
}
