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

// v2ConfigService is a service that interacts with the Configs
// endpoints from the Lacework APIv2 Server
type v2ConfigService struct {
	client *Client
	Azure  *v2AzureConfigService
	Gcp    *v2GcpConfigService
}

func NewV2ConfigService(c *Client) *v2ConfigService {
	return &v2ConfigService{c,
		&v2AzureConfigService{c},
		&v2GcpConfigService{c},
	}
}
