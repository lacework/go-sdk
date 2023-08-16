//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

	"github.com/pkg/errors"
)

// ReportDistributionsService is a service that interacts with the ReportDistributions
// endpoints from the Lacework APIv2 Server
type ReportDistributionsService struct {
	client *Client
}

// The report distribution frequency type
type reportDistributionFrequency int

const (
	ReportDistributionFrequencyDaily reportDistributionFrequency = iota
	ReportDistributionFrequencyWeekly
	ReportDistributionFrequencyBiweekly
	ReportDistributionFrequencyMonthly
)

func (frequency reportDistributionFrequency) String() string {
	return reportDistributionTypes[frequency]
}

var reportDistributionTypes = map[reportDistributionFrequency]string{
	ReportDistributionFrequencyDaily:    "daily",
	ReportDistributionFrequencyWeekly:   "weekly",
	ReportDistributionFrequencyBiweekly: "biweekly",
	ReportDistributionFrequencyMonthly:  "monthly",
}

func ReportDistributionFrequencies() (frequencies []string) {
	for _, v := range reportDistributionTypes {
		frequencies = append(frequencies, v)
	}
	return
}

// The report distribution violation type
type reportDistributionViolation int

const (
	ReportDistributionViolationCompliant reportDistributionViolation = iota
	ReportDistributionViolationNonCompliant
	ReportDistributionViolationSuppressed
	ReportDistributionViolationCouldNotAssess
	ReportDistributionViolationManual
)

func (subType reportDistributionViolation) String() string {
	return reportDistributionSubTypes[subType]
}

func ReportDistributionViolations() (values []string) {
	for _, v := range reportDistributionSubTypes {
		values = append(values, v)
	}
	return
}

var reportDistributionSubTypes = map[reportDistributionViolation]string{
	ReportDistributionViolationCompliant:      "Compliant",
	ReportDistributionViolationNonCompliant:   "NonCompliant",
	ReportDistributionViolationSuppressed:     "Suppressed",
	ReportDistributionViolationCouldNotAssess: "CouldNotAssess",
	ReportDistributionViolationManual:         "Manual",
}

// The report distribution scope type
type reportDistributionScope int

const (
	ReportDistributionScopeResourceGroup reportDistributionScope = iota
	ReportDistributionScopeCloudIntegration
)

func (scope reportDistributionScope) String() string {
	return reportDistributionScopeTypes[scope]
}

func ReportDistributionScopes() (values []string) {
	for _, v := range reportDistributionScopeTypes {
		values = append(values, v)
	}
	return
}

var reportDistributionScopeTypes = map[reportDistributionScope]string{
	ReportDistributionScopeResourceGroup:    "Resource Group",
	ReportDistributionScopeCloudIntegration: "Cloud Account Integration",
}

// List returns a ReportDistributionResponse
func (svc *ReportDistributionsService) List() (response ReportDistributionsResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2ReportDistributions, nil, &response)
	return
}

// Get returns a ReportDistributionResponse
func (svc *ReportDistributionsService) Get(guid string) (response ReportDistributionResponse, err error) {
	if guid == "" {
		return ReportDistributionResponse{}, errors.New("specify a report distribution guid")
	}
	apiPath := fmt.Sprintf(apiV2ReportDistributionsFromGUID, guid)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

// Delete a ReportDistribution
func (svc *ReportDistributionsService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify a report distribution guid")
	}

	return svc.client.RequestDecoder("DELETE", fmt.Sprintf(apiV2ReportDistributionsFromGUID, guid), nil, nil)
}

func (svc *ReportDistributionsService) Create(report ReportDistribution) (
	response ReportDistributionResponse, err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2ReportDistributions, report, &response)
	return
}

func (svc *ReportDistributionsService) Update(guid string, report ReportDistributionUpdate) (
	response ReportDistributionResponse, err error,
) {
	if guid == "" {
		return response, errors.New("specify a report distribution guid")
	}

	err = svc.client.RequestEncoderDecoder("PATCH",
		fmt.Sprintf(apiV2ReportDistributionsFromGUID, guid), report, &response)
	return
}

func (distribution *ReportDistribution) UpdateConfig() ReportDistributionUpdate {
	return ReportDistributionUpdate{
		DistributionName: distribution.DistributionName,
		Data:             distribution.Data,
		AlertChannels:    distribution.AlertChannels,
		Frequency:        distribution.Frequency,
	}
}

type ReportDistributionsResponse struct {
	Data []ReportDistribution `json:"data"`
}

type ReportDistributionResponse struct {
	Data ReportDistribution `json:"data"`
}

type ReportDistribution struct {
	ReportDistributionGuid string                 `json:"reportDistributionGuid,omitempty"`
	ReportDefinitionGuid   string                 `json:"reportDefinitionGuid"`
	DistributionName       string                 `json:"distributionName"`
	Data                   ReportDistributionData `json:"data"`
	AlertChannels          []string               `json:"alertChannels"`
	Frequency              string                 `json:"frequency"`
}

type ReportDistributionData struct {
	Severities     []string                        `json:"severities"`
	Violations     []string                        `json:"violations"`
	ResourceGroups []string                        `json:"resourceGroups"`
	Integrations   []ReportDistributionIntegration `json:"integrations"`
}

type ReportDistributionIntegration struct {
	TenantID       string `json:"tenantId,omitempty"`
	SubscriptionID string `json:"subscriptionId,omitempty"`
	AccountID      string `json:"accountId,omitempty"`
	OrganizationID string `json:"organizationId,omitempty"`
	ProjectID      string `json:"projectId,omitempty"`
}

type ReportDistributionUpdate struct {
	DistributionName string                 `json:"distributionName,omitempty"`
	Data             ReportDistributionData `json:"data,omitempty"`
	AlertChannels    []string               `json:"alertChannels,omitempty"`
	Frequency        string                 `json:"frequency,omitempty"`
}
