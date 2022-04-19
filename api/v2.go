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

import (
	"net/url"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// V2Endpoints groups all APIv2 endpoints available, they are grouped by
// schema which matches with our service architecture
type V2Endpoints struct {
	client *Client

	// Every schema must have its own service
	UserProfile             *UserProfileService
	AlertChannels           *AlertChannelsService
	AlertProfiles           *AlertProfilesService
	AlertRules              *AlertRulesService
	ReportRules             *ReportRulesService
	CloudAccounts           *CloudAccountsService
	ContainerRegistries     *ContainerRegistriesService
	ResourceGroups          *ResourceGroupsService
	AgentAccessTokens       *AgentAccessTokensService
	Query                   *QueryService
	Policy                  *PolicyService
	Entities                *EntitiesService
	Schemas                 *SchemasService
	Datasources             *DatasourcesService
	TeamMembers             *TeamMembersService
	VulnerabilityExceptions *VulnerabilityExceptionsService
	Vulnerabilities         *v2VulnerabilitiesService
}

func NewV2Endpoints(c *Client) *V2Endpoints {
	v2 := &V2Endpoints{c,
		&UserProfileService{c},
		&AlertChannelsService{c},
		&AlertProfilesService{c, AlertTemplatesService{c}},
		&AlertRulesService{c},
		&ReportRulesService{c},
		&CloudAccountsService{c},
		&ContainerRegistriesService{c},
		&ResourceGroupsService{c},
		&AgentAccessTokensService{c},
		&QueryService{c},
		&PolicyService{c},
		&EntitiesService{c},
		&SchemasService{c, map[integrationSchema]V2Service{}},
		&DatasourcesService{c},
		&TeamMembersService{c},
		&VulnerabilityExceptionsService{c},
		NewV2VulnerabilitiesService(c),
	}

	v2.Schemas.Services = map[integrationSchema]V2Service{
		AlertChannels:           &AlertChannelsService{c},
		AlertProfiles:           &AlertProfilesService{c, AlertTemplatesService{c}},
		AlertRules:              &AlertRulesService{c},
		CloudAccounts:           &CloudAccountsService{c},
		ContainerRegistries:     &ContainerRegistriesService{c},
		ResourceGroups:          &ResourceGroupsService{c},
		TeamMembers:             &TeamMembersService{c},
		ReportRules:             &ReportRulesService{c},
		VulnerabilityExceptions: &VulnerabilityExceptionsService{c},
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

type V2Pagination struct {
	Rows      int `json:"rows"`
	TotalRows int `json:"totalRows"`
	Urls      struct {
		NextPage string `json:"nextPage"`
	} `json:"urls"`
}

// Pageable is the interface that structs should implement to become
// pageable and be able to use the client.NextPage() function
type Pageable interface {
	PageInfo() *V2Pagination
	ResetPaging()
}

// NextPage
//
// Use this function to access the next page from an API v2 endpoint, the provided
// response must implement the Pageable interface and when it is passed, it will
// be overwritten, if the response doesn't have paging information this function
// returns false and not error
//
// Usage: To iterate over all pages
//
// ```go
// var (
// 		response = api.MachineDetailEntityResponse{}
// 		err      = client.V2.Entities.Search(&response, api.SearchFilter{})
// )
//
// for {
// 		// Use information from response.Data
// 		fmt.Printf("Data from page: %d\n", len(response.Data))
//
// 		pageOk, err := client.NextPage(&response)
// 		if err != nil {
// 			fmt.Printf("Unable to access next page, error '%s'", err.Error())
// 			break
// 		}
//
// 		if pageOk {
// 			continue
// 		}
// 		break
// }
// ```
func (c *Client) NextPage(p Pageable) (bool, error) {
	if p == nil {
		return false, nil
	}
	pagination := p.PageInfo()
	if pagination == nil {
		c.log.Info("pagination information not found")
		return false, nil
	}

	c.log.Info("pagination", zap.Int("rows", pagination.Rows),
		zap.Int("total_rows", pagination.TotalRows),
		zap.String("next_page", pagination.Urls.NextPage),
	)
	if pagination.Urls.NextPage == "" {
		return false, nil
	}

	pageURL, err := url.Parse(pagination.Urls.NextPage)
	if err != nil {
		return false, errors.Wrap(err, "unable to part next page url")
	}

	p.ResetPaging()
	c.log.Info("pagination reset")
	err = c.RequestDecoder("GET", pageURL.Path, nil, p)
	return true, err
}
