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

// The report definition type. At present "COMPLIANCE" is the only supported type for custom report definitions
type reportDefinitionType int

const (
	ReportDefinitionTypeCompliance reportDefinitionType = iota
)

func (reportType reportDefinitionType) String() string {
	return reportDefinitionTypes[reportType]
}

var reportDefinitionTypes = map[reportDefinitionType]string{
	ReportDefinitionTypeCompliance: "COMPLIANCE",
}

// The report definition subtype. Supported values are "AWS", "Azure", "GCP"
type reportDefinitionSubType int

const (
	ReportDefinitionSubTypeAws reportDefinitionSubType = iota
	ReportDefinitionSubTypeGcp
	ReportDefinitionSubTypeAzure
)

func (subType reportDefinitionSubType) String() string {
	return reportDefinitionSubTypes[subType]
}

func ReportDefinitionSubTypes() (values []string) {
	for _, v := range reportDefinitionSubTypes {
		values = append(values, v)
	}
	return
}

var reportDefinitionSubTypes = map[reportDefinitionSubType]string{
	ReportDefinitionSubTypeAws:   "AWS",
	ReportDefinitionSubTypeGcp:   "GCP",
	ReportDefinitionSubTypeAzure: "Azure",
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

// GetVersions returns a list of all versions of a reportDefinition
func (svc *ReportDefinitionsService) GetVersions(reportDefinitionGuid string) (response ReportDefinitionsResponse, err error) {
	if reportDefinitionGuid == "" {
		return ReportDefinitionsResponse{}, errors.New("specify a report definition guid")
	}
	apiPath := fmt.Sprintf(apiV2ReportDefinitionsVersions, reportDefinitionGuid)
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

func (svc *ReportDefinitionsService) Create(report ReportDefinition) (response ReportDefinitionResponse, err error) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2ReportDefinitions, report, &response)
	return
}

func (svc *ReportDefinitionsService) Update(guid string, report ReportDefinitionUpdate) (response ReportDefinitionResponse, err error) {
	if guid == "" {
		return response, errors.New("specify a report definition guid")
	}

	err = svc.client.RequestEncoderDecoder("PATCH", fmt.Sprintf(apiV2ReportDefinitionsFromGUID, guid), report, &response)
	return
}

func (svc *ReportDefinitionsService) Revert(guid string, version int) (response ReportDefinitionResponse, err error) {
	if guid == "" {
		return response, errors.New("specify a report definition guid")
	}

	apiPath := fmt.Sprintf(apiV2ReportDefinitionsRevert, guid, version)
	err = svc.client.RequestEncoderDecoder("PATCH", apiPath, "", &response)
	return
}

// NewReportDefinition creates a new report definition for Create function
func NewReportDefinition(cfg ReportDefinitionConfig) ReportDefinition {
	return ReportDefinition{
		ReportName:              cfg.ReportName,
		DisplayName:             cfg.DisplayName,
		ReportType:              cfg.ReportType,
		SubReportType:           cfg.SubReportType,
		ReportDefinitionDetails: ReportDefinitionDetails{Sections: cfg.Sections},
	}
}

// NewReportDefinitionUpdate creates a new report definition for Update function
func NewReportDefinitionUpdate(cfg ReportDefinitionConfig) ReportDefinitionUpdate {
	return ReportDefinitionUpdate{
		ReportName:              cfg.ReportName,
		DisplayName:             cfg.DisplayName,
		ReportDefinitionDetails: &ReportDefinitionDetails{Sections: cfg.Sections},
	}
}

var ReportDefinitionSubtypes = []string{"AWS", "Azure", "GCP"}

type ReportDefinitionConfig struct {
	ReportName    string                    `json:"reportName" yaml:"reportName"`
	DisplayName   string                    `json:"displayName" yaml:"displayName"`
	ReportType    string                    `json:"reportType" yaml:"reportType"`
	SubReportType string                    `json:"subReportType" yaml:"subReportType"`
	Sections      []ReportDefinitionSection `json:"sections,omitempty" yaml:"sections,omitempty"`
}

// ReportDefinitionUpdate represents fields allowed for update request
type ReportDefinitionUpdate struct {
	ReportName              string                   `json:"reportName,omitempty" yaml:"reportName,omitempty"`
	DisplayName             string                   `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	ReportDefinitionDetails *ReportDefinitionDetails `json:"reportDefinition,omitempty" yaml:"reportDefinition,omitempty"`
}

type ReportDefinitionsResponse struct {
	Data []ReportDefinition `json:"data"`
}

type ReportDefinitionResponse struct {
	Data ReportDefinition `json:"data"`
}

type ReportDefinition struct {
	ReportDefinitionGuid    string                  `json:"reportDefinitionGuid,omitempty" yaml:"reportDefinitionGuid,omitempty"`
	ReportName              string                  `json:"reportName" yaml:"reportName"`
	DisplayName             string                  `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	ReportType              string                  `json:"reportType" yaml:"reportType"`
	ReportNotificationType  string                  `json:"reportNotificationType,omitempty" yaml:"reportNotificationType,omitempty"`
	SubReportType           string                  `json:"subReportType" yaml:"subReportType"`
	ReportDefinitionDetails ReportDefinitionDetails `json:"reportDefinition" yaml:"reportDefinition"`
	Props                   *ReportDefinitionProps  `json:"props,omitempty" yaml:"props,omitempty"`
	DistributionType        string                  `json:"distributionType,omitempty" yaml:"distributionType,omitempty"`
	AlertChannels           []string                `json:"alertChannels,omitempty" yaml:"alertChannels,omitempty"`
	Frequency               string                  `json:"frequency,omitempty" yaml:"frequency,omitempty"`
	Version                 int                     `json:"version,omitempty" yaml:"version,omitempty"`
	UpdateType              string                  `json:"updateType,omitempty" yaml:"updateType,omitempty"`
	CreatedBy               string                  `json:"createdBy,omitempty" yaml:"createdBy,omitempty"`
	CreatedTime             *time.Time              `json:"createdTime,omitempty" yaml:"createdTime,omitempty"`
	Enabled                 int                     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

// IsCustom returns true if report definition is user created, not created by SYSTEM
func (report ReportDefinition) IsCustom() bool {
	return report.CreatedBy != "SYSTEM"
}

func (report ReportDefinition) Config() ReportDefinitionConfig {
	return ReportDefinitionConfig{
		ReportName:    report.ReportName,
		ReportType:    report.ReportType,
		DisplayName:   report.DisplayName,
		SubReportType: report.SubReportType,
		Sections:      report.ReportDefinitionDetails.Sections,
	}
}

type ReportDefinitionDetails struct {
	Sections  []ReportDefinitionSection   `json:"sections"`
	Overrides []ReportDefinitionOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
}

type ReportDefinitionOverrides struct {
	Policy string `json:"policy" yaml:"policy"`
	Title  string `json:"title" yaml:"title"`
}

type ReportDefinitionSection struct {
	Category string   `json:"category" yaml:"category"`
	Title    string   `json:"title" yaml:"title"`
	Policies []string `json:"policies" yaml:"policies"`
}

type ReportDefinitionProps struct {
	Engine         string   `json:"engine,omitempty" yaml:"engine,omitempty"`
	ReleaseLabel   string   `json:"releaseLabel,omitempty" yaml:"releaseLabel,omitempty"`
	ResourceGroups []string `json:"resourceGroups,omitempty" yaml:"resourceGroups,omitempty"`
	Integrations   []string `json:"integrations,omitempty" yaml:"integrations,omitempty"`
}
