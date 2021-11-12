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

// V2Endpoints groups all APIv2 endpoints available, they are grouped by
// schema which matches with our service architecture
type V2Endpoints struct {
	client *Client

	// Every schema must have its own service
	UserProfile         *UserProfileService
	AlertChannels       *AlertChannelsService
	AlertRules          *AlertRulesService
	ReportRules         *ReportRulesService
	CloudAccounts       *CloudAccountsService
	ContainerRegistries *ContainerRegistriesService
	ResourceGroups      *ResourceGroupsService
	AgentAccessTokens   *AgentAccessTokensService
	Query               *QueryService
	Policy              *PolicyService
	Schemas             *SchemasService
	Datasources         *DatasourcesService
}

func NewV2Endpoints(c *Client) *V2Endpoints {
	v2 := &V2Endpoints{c,
		&UserProfileService{c},
		&AlertChannelsService{c},
		&AlertRulesService{c},
		&ReportRulesService{c},
		&CloudAccountsService{c},
		&ContainerRegistriesService{c},
		&ResourceGroupsService{c},
		&AgentAccessTokensService{c},
		&QueryService{c},
		&PolicyService{c},
		&SchemasService{c, map[integrationSchema]V2Service{}},
		&DatasourcesService{c},
	}

	v2.Schemas.Services = map[integrationSchema]V2Service{
		AlertChannels:       &AlertChannelsService{c},
		AlertRules:          &AlertRulesService{c},
		CloudAccounts:       &CloudAccountsService{c},
		ContainerRegistries: &ContainerRegistriesService{c},
		ResourceGroups:      &ResourceGroupsService{c},
		ReportRules:         &ReportRulesService{c},
	}
	return v2
}

type V2Service interface {
	Get(string, interface{}) error
	Delete(string) error
}

type V2CommonIntegration struct {
	Data v2CommonIntegrationData `json:"data"`
}
