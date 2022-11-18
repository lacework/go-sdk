//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

// MapPolicy enable and apply policyguids or disable
func (svc *ContainerRegistriesService) MapPolicy(guid string, data MapPolicyRequest) (
	response MapPolicyResponse,
	err error,
) {
	err = svc.mapPolicy(guid, data, &response)
	return
}

type MapPolicyRequest struct {
	Evaluate    bool     `json:"evaluate"`
	PolicyGuids []string `json:"policyGuids,omitempty"`
}

type MapPolicyResponse struct {
	Data MapPolicyIntegration `json:"data"`
}

type MapPolicyIntegration struct {
	v2CommonIntegrationData
	Props       interface{}   `json:"props"`
	State       interface{}   `json:"state"`
	Data        interface{}   `json:"data"`
	ServerToken V2ServerToken `json:"server_token"`
}
