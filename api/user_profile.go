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
	"regexp"
	"strings"
)

// UserProfileService is the service that interacts with the UserProfile
// schema from the Lacework APIv2 Server
type UserProfileService struct {
	client *Client
}

func (svc *UserProfileService) Get() (response UserProfileResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2UserProfile, nil, &response)
	return
}

type UserProfileResponse struct {
	Data []UserProfile `json:"data"`
}

type UserProfile struct {
	Username   string    `json:"username"`
	OrgAccount bool      `json:"orgAccount"`
	URL        string    `json:"url"`
	OrgAdmin   bool      `json:"orgAdmin"`
	OrgUser    bool      `json:"orgUser"`
	Accounts   []Account `json:"accounts"`
}

func (p *UserProfile) OrgAccountName() string {
	// TODO @afiune should we handle other datacenters?
	rx, err := regexp.Compile(`\.lacework\.net.*`)
	if err == nil {
		if split := rx.Split(p.URL, -1); len(split) != 0 {
			return strings.ToLower(split[0])
		}
	}

	return p.URL
}

func (p *UserProfile) SubAccountNames() []string {
	names := make([]string, 0)
	orgAccountName := p.OrgAccountName()
	for _, acc := range p.Accounts {
		accName := strings.ToLower(acc.AccountName)
		if accName == orgAccountName {
			continue
		}
		if acc.Enabled() {
			names = append(names, accName)
		}
	}
	return names
}

type Account struct {
	Admin       bool   `json:"admin"`
	AccountName string `json:"accountName"`
	CustGUID    string `json:"custGuid"`
	UserGUID    string `json:"userGuid"`
	UserEnabled int    `json:"userEnabled"`
}

func (a *Account) Enabled() bool {
	return a.UserEnabled == 1
}
