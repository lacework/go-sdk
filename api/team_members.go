//
// Author:: Vatasha White (<vatasha.white@lacework.net>)
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
	"errors"
	"fmt"
)

type TeamMembersService struct {
	client *Client
}

// NewTeamMember returns an instance of the Team Member struct
//
// Basic usage: Initialize a new TeamMember struct and then use the new instance to perform CRUD operations.
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   teamMember := api.NewTeamMember(
//		"FooBar",
//		api.TeamMemberProps{
//		Company: "ACME Inc",
//		FirstName: "Foo",
//		LastName: "Bar"
//      },
//   },
// )
//
// client.V2.TeamMembers.Create(teamMember)
//
func NewTeamMember(username string, props TeamMemberProps) TeamMember {
	return TeamMember{
		Props:       props,
		UserEnabled: 1,
		UserName:    username,
	}
}

// List returns a list of team members
func (svc *TeamMembersService) List() (res TeamMembersResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2TeamMembers, nil, &res)
	return
}

// Create creates a single team member
func (svc *TeamMembersService) Create(tm TeamMember) (res TeamMemberResponse, err error) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2TeamMembers, tm, &res)
	return
}

// Delete deletes a single team member with the corresponding guid
func (svc *TeamMembersService) Delete(guid string) error {
	if guid == "" {
		return errors.New("please specify a guid")
	}

	return svc.client.RequestDecoder("DELETE", fmt.Sprintf(apiV2TeamMembersFromGUID, guid), nil, nil)
}

// Update updates a single team member with the corresponding guid
func (svc *TeamMembersService) Update(tm TeamMember) (res TeamMemberResponse, err error) {
	if tm.UserGuid == "" {
		err = errors.New("please specify a guid")
		return
	}
	err = svc.client.RequestEncoderDecoder("PATCH", fmt.Sprintf(apiV2TeamMembersFromGUID, tm.UserGuid), tm, &res)
	return
}

// Get returns a response of the team member
func (svc *TeamMembersService) Get(guid string, res interface{}) error {
	if guid == "" {
		return errors.New("please specify a guid")
	}
	return svc.client.RequestDecoder("GET", fmt.Sprintf(apiV2TeamMembersFromGUID, guid), nil, &res)

}

type TeamMemberProps struct {
	AccountAdmin           bool   `json:"accountAdmin,omitempty"`
	Company                string `json:"company"`
	CreatedTime            string `json:"createdTime,omitempty"`
	FirstName              string `json:"firstName"`
	JitCreated             bool   `json:"jitCreated,omitempty"`
	LastLoginTime          string `json:"lastLoginTime,omitempty"`
	LastName               string `json:"lastName"`
	LastSessionCreatedTime string `json:"lastSessionCreatedTime,omitempty"`
	OrgAdmin               bool   `json:"orgAdmin,omitempty"`
	OrgUser                bool   `json:"orgUser,omitempty"`
	UpdatedBy              string `json:"updatedBy,omitempty"`
	UpdatedTime            string `json:"UpdatedTime,omitempty"`
}

type TeamMember struct {
	CustGuid    string          `json:"custGuid,omitempty"`
	Props       TeamMemberProps `json:"props"`
	UserEnabled int             `json:"userEnabled"`
	UserGuid    string          `json:"userGuid,omitempty"`
	UserName    string          `json:"userName"`
}

type TeamMemberResponse struct {
	Data TeamMember `json:"data"`
}

type TeamMembersResponse struct {
	Data []TeamMember `json:"data"`
}
