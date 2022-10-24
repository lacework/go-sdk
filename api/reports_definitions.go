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
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// ReportDefinitionsService is a service that interacts with the ReportDefinitions
// endpoints from the Lacework APIv2 Server
type ReportDefinitionsService struct {
	client *Client
}

// List returns a ReportDefinitionResponse
func (svc *ReportDefinitionsService) List() (response ReportDefinitionsResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2ReportDefinitions, nil, &response)
	return
}

// Get returns a ReportDefinitionResponse
func (svc *ReportDefinitionsService) Get(reportDefinitionGuid string) (response ReportDefinitionResponse, err error) {
	if reportDefinitionGuid == "" {
		return ReportDefinitionResponse{}, errors.New("specify a report definition guid")
	}
	apiPath := fmt.Sprintf(apiV2ReportDefinitionsFromGUID, reportDefinitionGuid)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

type ReportDefinitionsResponse struct {
	Data []ReportDefinition `json:"data"`
}

type ReportDefinitionResponse struct {
	Data ReportDefinition `json:"data"`
}

type ReportDefinition struct {
	ReportDefinitionGuid    string                  `json:"reportDefinitionGuid"`
	ReportName              string                  `json:"reportName"`
	DisplayName             string                  `json:"displayName"`
	ReportType              string                  `json:"reportType"`
	ReportNotificationType  string                  `json:"reportNotificationType"`
	SubReportType           string                  `json:"subReportType"`
	ReportDefinitionDetails ReportDefinitionDetails `json:"reportDefinition"`
	Props                   ReportDefinitionProps   `json:"props"`
	Version                 int                     `json:"version"`
	CreatedBy               string                  `json:"createdBy"`
	CreatedTime             time.Time               `json:"createdTime"`
	Enabled                 int                     `json:"enabled"`
}

type ReportDefinitionDetails struct {
	Sections  []ReportDefinitionSection `json:"sections"`
	Overrides []ReportDefinitionSection `json:"overrides"`
}

type ReportDefinitionOverrides struct {
	Policy string `json:"policy"`
	Title  string `json:"title"`
}

type ReportDefinitionSection struct {
	Category string   `json:"category"`
	Title    string   `json:"title"`
	Policies []string `json:"policies"`
}

type ReportDefinitionProps struct {
	Engine       string `json:"engine"`
	ReleaseLabel string `json:"releaseLabel"`
}
