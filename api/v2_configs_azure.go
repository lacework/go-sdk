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

import "fmt"

// v2AzureConfigService is a service that interacts with the APIv2
// vulnerabilities endpoints for hosts
type v2AzureConfigService struct {
	client *Client
}

// List returns a list of Azure tenants and subscriptions
func (svc *v2AzureConfigService) List() (response AzureConfigsResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2ConfigsAzure, nil, &response)
	if err != nil {
		return
	}

	return
}

// ListSubscriptions returns a list of Azure subscriptions for a given tenant
func (svc *v2AzureConfigService) ListSubscriptions(tenantID string) (response AzureConfigsResponse, err error) {
	err = svc.client.RequestDecoder("GET", fmt.Sprintf(apiV2ConfigsAzureSubscriptions, tenantID), nil, &response)
	if err != nil {
		return
	}

	return
}

type AzureConfigsResponse struct {
	Data []AzureConfigData `json:"data"`
}

type AzureConfigData struct {
	Tenant        string   `json:"tenant"`
	Subscriptions []string `json:"subscriptions"`
}
