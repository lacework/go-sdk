//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

import "encoding/json"

// PolicyService is a service that interacts with the Custom Policies
// endpoints from the Lacework Server
type PolicyService struct {
	client *Client
}

type PolicyCreateResponse struct {
	Data    []Policy `json:"data"`
	Ok      bool     `json:"ok"`
	Message string   `json:"message"`
}

type Policy struct {
	PolicyID string `json:"policy_id"`
}

func (svc *PolicyService) Create(policy string) (
	response PolicyCreateResponse,
	err error,
) {
	var p map[string]interface{}
	if err = json.Unmarshal([]byte(policy), &p); err != nil {
		return
	}
	err = svc.client.RequestEncoderDecoder("POST", ApiPolicy, p, &response)
	return
}
