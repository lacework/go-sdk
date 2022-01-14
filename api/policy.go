//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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
	"net/url"

	"github.com/pkg/errors"
)

// PolicyService is a service that interacts with the Custom Policies
// endpoints from the Lacework Server
type PolicyService struct {
	client *Client
}

// ValidPolicySeverities is a list of all valid policy severities
var ValidPolicySeverities = []string{"critical", "high", "medium", "low", "info"}

type NewPolicy struct {
	EvaluatorID   string `json:"evaluatorId,omitempty" yaml:"evaluatorId"`
	PolicyID      string `json:"policyId" yaml:"policyId" `
	PolicyType    string `json:"policyType" yaml:"policyType"`
	QueryID       string `json:"queryId" yaml:"queryId"`
	Title         string `json:"title" yaml:"title"`
	Enabled       bool   `json:"enabled" yaml:"enabled"`
	Description   string `json:"description" yaml:"description"`
	Remediation   string `json:"remediation" yaml:"remediation"`
	Severity      string `json:"severity" yaml:"severity"`
	Limit         int    `json:"limit,omitempty" yaml:"limit,omitempty"`
	EvalFrequency string `json:"evalFrequency,omitempty" yaml:"evalFrequency,omitempty"`
	AlertEnabled  bool   `json:"alertEnabled" yaml:"alertEnabled"`
	AlertProfile  string `json:"alertProfile" yaml:"alertProfile"`
}

/* In order to properly PATCH we need to omit items that aren't specified.
For booleans and integers Golang will omit zero values false and 0 respectively.
This would prevent someone from toggling something to disabled or 0 respectively.
As such we are using pointers instead of primitives for booleans and integers in this struct
*/
type UpdatePolicy struct {
	EvaluatorID   string `json:"evaluatorId,omitempty" yaml:"evaluatorId,omitempty"`
	PolicyID      string `json:"policyId,omitempty" yaml:"policyId,omitempty"`
	PolicyType    string `json:"policyType,omitempty" yaml:"policyType,omitempty"`
	QueryID       string `json:"queryId,omitempty" yaml:"queryId,omitempty"`
	Title         string `json:"title,omitempty" yaml:"title,omitempty"`
	Enabled       *bool  `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Description   string `json:"description,omitempty" yaml:"description,omitempty"`
	Remediation   string `json:"remediation,omitempty" yaml:"remediation,omitempty"`
	Severity      string `json:"severity,omitempty" yaml:"severity,omitempty"`
	Limit         *int   `json:"limit,omitempty" yaml:"limit,omitempty"`
	EvalFrequency string `json:"evalFrequency,omitempty" yaml:"evalFrequency,omitempty"`
	AlertEnabled  *bool  `json:"alertEnabled,omitempty" yaml:"alertEnabled,omitempty"`
	AlertProfile  string `json:"alertProfile,omitempty" yaml:"alertProfile,omitempty"`
}

type Policy struct {
	EvaluatorID    string `json:"evaluatorId"`
	PolicyID       string `json:"policyId"`
	PolicyType     string `json:"policyType"`
	QueryID        string `json:"queryId"`
	Title          string `json:"title"`
	Enabled        bool   `json:"enabled"`
	Description    string `json:"description"`
	Remediation    string `json:"remediation"`
	Severity       string `json:"severity"`
	Limit          int    `json:"limit"`
	EvalFrequency  string `json:"evalFrequency"`
	AlertEnabled   bool   `json:"alertEnabled"`
	AlertProfile   string `json:"alertProfile"`
	Owner          string `json:"owner"`
	LastUpdateTime string `json:"lastUpdateTime"`
	LastUpdateUser string `json:"lastUpdateUser"`
}

type PolicyResponse struct {
	Data    Policy `json:"data"`
	Message string `json:"message"`
}

type PoliciesResponse struct {
	Data    []Policy `json:"data"`
	Message string   `json:"message"`
}

func (svc *PolicyService) Create(np NewPolicy) (
	response PolicyResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2Policies, np, &response)
	return
}

func (svc *PolicyService) List() (
	response PoliciesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiV2Policies, nil, &response)
	return
}

func (svc *PolicyService) Get(policyID string) (
	response PolicyResponse,
	err error,
) {
	if policyID == "" {
		err = errors.New("policy ID must be provided")
		return
	}
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf("%s/%s", apiV2Policies, url.QueryEscape(policyID)),
		nil,
		&response,
	)
	return
}

func (svc *PolicyService) Update(up UpdatePolicy) (
	response PolicyResponse,
	err error,
) {
	if up.PolicyID == "" {
		err = errors.New("policy ID must be provided")
		return
	}
	var policyID = up.PolicyID
	up.PolicyID = "" // omit this for PATCH

	err = svc.client.RequestEncoderDecoder(
		"PATCH",
		fmt.Sprintf("%s/%s", apiV2Policies, url.QueryEscape(policyID)),
		up,
		&response,
	)
	return
}

func (svc *PolicyService) Delete(policyID string) (
	response PolicyResponse,
	err error,
) {
	if policyID == "" {
		err = errors.New("policy ID must be provided")
		return
	}
	err = svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf("%s/%s", apiV2Policies, url.QueryEscape(policyID)),
		nil,
		&response,
	)
	return
}
