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

import "github.com/lacework/go-sdk/lwdomain"

// OrganizationInfoService is the service that interacts with
// the OrganizationInfo schema from the Lacework APIv2 Server
type OrganizationInfoService struct {
	client *Client
}

func (svc *OrganizationInfoService) Get() (
	response OrganizationInfoResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET",
		apiV2OrganizationInfo,
		nil,
		&response,
	)
	return
}

type OrganizationInfoResponse struct {
	Data []OrganizationInfo `json:"data"`
}

type OrganizationInfo struct {
	OrgAccount    bool   `json:"orgAccount"`
	OrgAccountURL string `json:"orgAccountUrl,omitempty"`
}

func (r OrganizationInfo) AccountName() string {
	d, _ := lwdomain.New(r.OrgAccountURL)
	return d.String()
}
