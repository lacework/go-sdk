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
	CloudAccounts       *CloudAccountsService
	ContainerRegistries *ContainerRegistriesService
	AgentAccessTokens   *AgentAccessTokensService
	Query               *QueryService
	Policy              *PolicyService
}

func NewV2Endpoints(c *Client) *V2Endpoints {
	return &V2Endpoints{c,
		&UserProfileService{c},
		&AlertChannelsService{c},
		&CloudAccountsService{c},
		&ContainerRegistriesService{c},
		&AgentAccessTokensService{c},
		&QueryService{c},
		&PolicyService{c},
	}
}
