//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

// SchemasService is the service that retrieves schemas for v2
type SchemasService struct {
	client   *Client
	Services map[integrationSchema]V2Service
}

type integrationSchema int

const (
	None integrationSchema = iota
	AlertChannels
	AlertProfiles
	AlertRules
	ContainerRegistries
	CloudAccounts
	ResourceGroups
	ReportRules
	TeamMembers
	VulnerabilityExceptions
)

func (svc *SchemasService) GetService(schemaName integrationSchema) V2Service {
	return svc.Services[schemaName]
}
