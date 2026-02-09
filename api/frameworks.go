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

import (
	"time"
)

// FrameworksService is a service that interacts with the Frameworks
// endpoints from the Lacework APIv2 Server
type FrameworksService struct {
	client *Client
}

// List returns a ReportDefinitionResponse
func (svc *FrameworksService) List() (response FrameworksResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2Frameworks, nil, &response)
	return
}

type FrameworksResponse struct {
	Data []Framework `json:"data"`
}

type FrameworkResponse struct {
	Data Framework `json:"data"`
}

type Framework struct {
	Guid           string     `json:"guid" yaml:"guid"`
	Name           string     `json:"name" yaml:"name"`
	Domains        []string   `json:"domains" yaml:"domains"`
	Sections       []Section  `json:"sections" yaml:"sections"`
	Owner          string     `json:"owner" yaml:"owner"`
	Revision       int        `json:"revision" yaml:"revision"`
	LastUpdateTime *time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	LastUpdateUser string     `json:"lastUpdateUser" yaml:"lastUpdateUser"`
	CreatedTime    *time.Time `json:"createdTime" yaml:"createdTime"`
	CreatedBy      string     `json:"createdBy" yaml:"createdBy"`
}

type Section struct {
	Name     string   `json:"name" yaml:"name"`
	Policies []Policy `json:"policies" yaml:"policies"`
}
