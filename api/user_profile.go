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
	Data []struct {
		Username   string `json:"username"`
		OrgAccount bool   `json:"orgAccount"`
		URL        string `json:"url"`
		OrgAdmin   bool   `json:"orgAdmin"`
		OrgUser    bool   `json:"orgUser"`
		Accounts   []struct {
			Admin       bool   `json:"admin"`
			AccountName string `json:"accountName"`
			CustGUID    string `json:"custGuid"`
			UserGUID    string `json:"userGuid"`
			UserEnabled int    `json:"userEnabled"`
		} `json:"accounts"`
	} `json:"data"`
}
