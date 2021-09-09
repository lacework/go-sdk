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

import "github.com/lacework/go-sdk/lwdomain"

// AccountService is a service that interacts with Account related
// endpoints from the Lacework Server
type AccountService struct {
	client *Client
}

func (svc *AccountService) GetOrganizationInfo() (
	response accountOrganizationInfoResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET",
		apiAccountOrganizationInfo,
		nil,
		&response,
	)
	return
}

type accountOrganizationInfoResponse struct {
	OrgAccount    bool   `json:"orgAccount"`
	OrgAccountURL string `json:"orgAccountUrl,omitempty"`
}

func (r accountOrganizationInfoResponse) AccountName() string {
	d, _ := lwdomain.New(r.OrgAccountURL)
	return d.String()
}
