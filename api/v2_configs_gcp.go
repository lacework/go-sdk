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

// v2GcpConfigService is a service that interacts with the APIv2
// vulnerabilities endpoints for hosts
type v2GcpConfigService struct {
	client *Client
}

// List returns a list of Gcp organizations and projects
func (svc *v2GcpConfigService) List() (response GcpConfigsResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2ConfigsGcp, nil, &response)
	if err != nil {
		return
	}

	return
}

type GcpConfigsResponse struct {
	Data []GcpConfigData `json:"data"`
}

type GcpConfigData struct {
	Tenant        string   `json:"tenant"`
	Subscriptions []string `json:"subscriptions"`
}
