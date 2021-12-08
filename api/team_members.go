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
	"fmt"

	"github.com/pkg/errors"
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

// NewTeamMemberOrg returns an instance of the team member org struct
//
// Basic usage: Initialize a new TeamMemberOrg struct and then use the new instance to perform CRUD operations.
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   teamMember := api.NewTeamMemberOrg(
//		"FooBar",
//		api.TeamMemberProps{
//		Company: "ACME Inc",
//		FirstName: "Foo",
//		LastName: "Bar"
//      },
//   },
// )
//
// client.V2.TeamMembers.CreateOrg(teamMember)
//
func NewTeamMemberOrg(username string, props TeamMemberProps) TeamMemberOrg {
	return TeamMemberOrg{
		Props:             props,
		UserEnabled:       1,
		UserName:          username,
		OrgAdmin:          false,
		OrgUser:           true,
		AdminRoleAccounts: []string{},
		UserRoleAccounts:  []string{},
	}
}

// List returns a list of team members
func (svc *TeamMembersService) List() (res TeamMembersResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2TeamMembers, nil, &res)
	return
}

// Create creates a single team member
func (svc *TeamMembersService) Create(tm TeamMember) (res TeamMemberResponse, err error) {
	if svc.client.OrgAccess() {
		return res, errors.New("client configured to manage org-level datasets, use CreateOrg()")
	}
	err = svc.client.RequestEncoderDecoder("POST", apiV2TeamMembers, tm, &res)
	return
}

// CreateOrg creates a single team member at the org level
//TODO Move all ORG stuff into a different file
func (svc *TeamMembersService) CreateOrg(tm TeamMemberOrg) (res TeamMemberOrgResponse, err error) {
	if !svc.client.OrgAccess() {
		return res, errors.New("client configured to manage account-level datasets, use Create()")
	}
	err = svc.client.RequestEncoderDecoder("POST", apiV2TeamMembers, tm, &res)
	return
}

// Delete deletes a single team member at the account level with the corresponding guid
func (svc *TeamMembersService) Delete(guid string) error {
	if svc.client.OrgAccess() {
		return errors.New("client configured to manage org-level datasets, use DeleteOrg()")
	}
	if guid == "" {
		return errors.New("please specify a guid")
	}
	return svc.client.RequestDecoder("DELETE", fmt.Sprintf(apiV2TeamMembersFromGUID, guid), nil, nil)
}

// DeleteOrg deletes a single team member at the org level with the corresponding guid
func (svc *TeamMembersService) DeleteOrg(guid string) error {
	if !svc.client.OrgAccess() {
		return errors.New("client configured to manage account-level datasets, use Delete()")
	}
	if guid == "" {
		return errors.New("please specify a guid")
	}
	return svc.client.RequestDecoder("DELETE", fmt.Sprintf(apiV2TeamMembersFromGUID, guid), nil, nil)
}

// Update updates a single team member at the account-level with the corresponding guid
func (svc *TeamMembersService) Update(tm TeamMember) (res TeamMemberResponse, err error) {
	if svc.client.OrgAccess() {
		return res, errors.New("client configured to manage org-level datasets, use UpdateOrg()")
	}
	if tm.UserGuid == "" {
		err = errors.New("please specify a guid")
		return
	}
	userGuid := tm.UserGuid
	// Omit userGuid for patch requests
	tm.UserGuid = ""
	// Omit userName for patch requests as it cannot be modified
	tm.UserName = ""
	err = svc.client.RequestEncoderDecoder("PATCH", fmt.Sprintf(apiV2TeamMembersFromGUID, userGuid), tm, &res)
	return
}

// UpdateOrg updates a single team member at the org-level with the corresponding username
func (svc *TeamMembersService) UpdateOrg(tm TeamMemberOrg) (res TeamMemberOrgResponse, err error) {
	if !svc.client.OrgAccess() {
		return res, errors.New("client configured to manage account-level datasets, use Update()")
	}
	if tm.UserName == "" {
		err = errors.New("please specify a username")
		return
	}
	tms, errSearch := svc.SearchUsername(tm.UserName)
	if errSearch != nil || len(tms.Data) == 0 {
		err = errors.Wrap(err, "unable to find user with specified username")
		return
	}
	return svc.UpdateOrgById(tm)
}

// UpdateOrgById updates a single team member at the org-level with the corresponding guid
func (svc *TeamMembersService) UpdateOrgById(tm TeamMemberOrg) (res TeamMemberOrgResponse, err error) {
	if !svc.client.OrgAccess() {
		return res, errors.New("client configured to manage account-level datasets, use Update()")
	}
	if tm.UserGuid == "" {
		err = errors.New("please specify a user guid")
		return
	}
	userGuid := tm.UserGuid
	// Omit UserGuid from the patch body as it cannot be modified
	tm.UserGuid = ""
	// Omit userEnabled field from the patch body as it cannot be modified
	tm.UserEnabled = 0
	// Omit userName field from the patch body as it cannot be modified
	tm.UserName = ""
	// Omit Company field from the patch body as it cannot be modified
	tm.Props.Company = ""

	err = svc.client.RequestEncoderDecoder("PATCH", fmt.Sprintf(apiV2TeamMembersFromGUID, userGuid), tm, &res)
	return
}

// Get returns a response of the team member
func (svc *TeamMembersService) Get(guid string, res interface{}) error {
	if guid == "" {
		return errors.New("please specify a guid")
	}
	return svc.client.RequestDecoder("GET", fmt.Sprintf(apiV2TeamMembersFromGUID, guid), nil, &res)

}

func (svc *TeamMembersService) SearchUsername(username string) (res TeamMembersResponse, err error) {
	if username == "" {
		err = errors.New("specify a username to search for a team member")
		return
	}
	err = svc.client.RequestEncoderDecoder("POST",
		apiV2TeamMembersSearch,
		SearchFilter{
			Filters: []Filter{
				Filter{
					Field:      "userName",
					Expression: "eq",
					Value:      username,
				},
			},
		},
		&res,
	)
	return
}

type TeamMemberProps struct {
	AccountAdmin bool `json:"accountAdmin,omitempty"`
	//Company is empty for patch requests on updateOrg as it cannot be modified
	Company                string      `json:"company,omitempty"`
	CreatedTime            string      `json:"createdTime,omitempty"`
	FirstName              string      `json:"firstName"`
	JitCreated             bool        `json:"jitCreated,omitempty"`
	LastLoginTime          interface{} `json:"lastLoginTime,omitempty"`
	LastName               string      `json:"lastName"`
	LastSessionCreatedTime interface{} `json:"lastSessionCreatedTime,omitempty"`
	OrgAdmin               bool        `json:"orgAdmin,omitempty"`
	OrgUser                bool        `json:"orgUser,omitempty"`
	UpdatedBy              string      `json:"updatedBy,omitempty"`
	UpdatedTime            interface{} `json:"updatedTime,omitempty"`
}

// TeamMember is for a standalone team member without org access
type TeamMember struct {
	CustGuid    string          `json:"custGuid,omitempty"`
	Props       TeamMemberProps `json:"props"`
	UserEnabled int             `json:"userEnabled"`
	UserGuid    string          `json:"userGuid,omitempty"`
	UserName    string          `json:"userName,omitempty"`
}

type TeamMemberResponse struct {
	Data TeamMember `json:"data"`
}

type TeamMembersResponse struct {
	Data []TeamMember `json:"data"`
}

// TeamMemberOrg is for an organizational team member
type TeamMemberOrg struct {
	AdminRoleAccounts []string        `json:"adminRoleAccounts"`
	OrgAdmin          bool            `json:"orgAdmin"`
	OrgUser           bool            `json:"orgUser"`
	Props             TeamMemberProps `json:"props"`
	UserEnabled       int             `json:"userEnabled,omitempty"`
	UserGuid          string          `json:"userGuid,omitempty"`
	UserName          string          `json:"userName,omitempty"`
	UserRoleAccounts  []string        `json:"userRoleAccounts"`
}

type TeamMemberAccount struct {
	AccountName string `json:"accountName"`
	Admin       bool   `json:"admin"`
	CustGuid    string `json:"custGuid"`
	UserEnabled int    `json:"userEnabled"`
	UserGuid    string `json:"userGuid"`
}

type TeamMemberOrgData struct {
	Accounts   []TeamMemberAccount `json:"accounts"`
	OrgAccount bool                `json:"orgAccount"`
	OrgAdmin   bool                `json:"orgAdmin"`
	OrgUser    bool                `json:"orgUser"`
	Url        string              `json:"url"`
	UserName   string              `json:"userName"`
}

type TeamMemberOrgResponse struct {
	Data TeamMemberOrgData `json:"data"`
}
