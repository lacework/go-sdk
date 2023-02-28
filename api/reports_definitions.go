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

const ReportDefinitionNotificationTypeAws = "AWS"
const ReportDefinitionNotificationTypeGcp = "GCP"
const ReportDefinitionNotificationTypeAzure = "Azure"

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

// Delete a ReportDefinition
func (svc *ReportDefinitionsService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify a report definition guid")
	}

	return svc.client.RequestDecoder("DELETE", fmt.Sprintf(apiV2ReportDefinitionsFromGUID, guid), nil, nil)
}

// NewReportDefinition creates a new report definition for Create function
func NewReportDefinition(cfg ReportDefinitionConfig) ReportDefinition {
	return ReportDefinition{
		ReportName:              cfg.ReportName,
		ReportType:              cfg.ReportType,
		SubReportType:           cfg.SubReportType,
		ReportDefinitionDetails: ReportDefinitionDetails{cfg.Sections, cfg.Overrides},
		Props:                   cfg.Props,
		DistributionType:        cfg.DistributionType,
		AlertChannels:           cfg.AlertChannels,
		Frequency:               cfg.Frequency,
	}
}

var ReportDefinitionSubtypes = []string{"AWS", "Azure", "GCP"}

type ReportDefinitionConfig struct {
	ReportName       string                      `json:"reportName"`
	ReportType       string                      `json:"reportType"`
	SubReportType    string                      `json:"subReportType"`
	Sections         []ReportDefinitionSection   `json:"sections"`
	Overrides        []ReportDefinitionOverrides `json:"overrides"`
	Props            ReportDefinitionProps       `json:"props"`
	AlertChannels    []string                    `json:"alertChannels"`
	DistributionType string                      `json:"distributionType"`
	Frequency        string                      `json:"frequency"`
	UpdateType       string                      `json:"updateType,omitempty"`
}

type ReportDefinitionsResponse struct {
	Data []ReportDefinition `json:"data"`
}

type ReportDefinitionResponse struct {
	Data ReportDefinition `json:"data"`
}

type ReportDefinition struct {
	ReportDefinitionGuid    string                  `json:"reportDefinitionGuid,omitempty"`
	ReportName              string                  `json:"reportName"`
	DisplayName             string                  `json:"displayName,omitempty"`
	ReportType              string                  `json:"reportType"`
	ReportNotificationType  string                  `json:"reportNotificationType,omitempty"`
	SubReportType           string                  `json:"subReportType"`
	ReportDefinitionDetails ReportDefinitionDetails `json:"reportDefinition"`
	Props                   ReportDefinitionProps   `json:"props"`
	DistributionType        string                  `json:"distributionType"`
	AlertChannels           []string                `json:"alertChannels,omitempty"`
	Frequency               string                  `json:"frequency,omitempty"`
	Version                 int                     `json:"version,omitempty"`
	UpdateType              string                  `json:"updateType,omitempty"`
	CreatedBy               string                  `json:"createdBy,omitempty"`
	CreatedTime             *time.Time              `json:"createdTime,omitempty"`
	Enabled                 int                     `json:"enabled,omitempty"`
}

type ReportDefinitionDetails struct {
	Sections  []ReportDefinitionSection   `json:"sections"`
	Overrides []ReportDefinitionOverrides `json:"overrides,omitempty"`
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
	Engine         string   `json:"engine,omitempty"`
	ReleaseLabel   string   `json:"releaseLabel,omitempty"`
	ResourceGroups []string `json:"resourceGroups,omitempty"`
	Integrations   []string `json:"integrations,omitempty"`
}
